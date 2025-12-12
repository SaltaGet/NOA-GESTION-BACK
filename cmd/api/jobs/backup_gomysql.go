package jobs

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

type Config struct {
	User      string   `json:"user"`
	Password  string   `json:"password"`
	Databases []string `json:"databases"`
	Host      string   `json:"host"`
	Port      string   `json:"port"`
	BackupDir string   `json:"backup_dir"`
}

type Checkpoint struct {
	BinlogFile string `json:"binlog_file"`
	Position   uint32 `json:"position"`
}

type TableSchema struct {
	Columns []string
}

type BackupTask struct {
	Database string
	Type     string // "full" o "incremental"
}

func LoadConfig(deps *dependencies.MainContainer) (*Config, error) {
	dsn := os.Getenv("URI_DB")
	cfg, err := parseDSN(dsn)
	if err != nil {
		return nil, err
	}

	tenants, err := deps.TenantController.TenantService.TenantGetConections()
	if err != nil {
		return nil, err
	}

	for _, conn := range tenants {
		dbName, err := extractDBName(conn.Connection)
		if err != nil {
			return nil, err
		}
		cfg.Databases = append(cfg.Databases, dbName)
	}
	return &cfg, nil
}

func parseDSN(dsn string) (Config, error) {
	var cfg Config
	parts := strings.SplitN(dsn, "@tcp(", 2)
	if len(parts) != 2 {
		return cfg, fmt.Errorf("DSN inválido")
	}
	up := strings.SplitN(parts[0], ":", 2)
	if len(up) != 2 {
		return cfg, fmt.Errorf("falta ':' en usuario:contraseña")
	}
	cfg.User, cfg.Password = up[0], up[1]

	hostEnd := strings.Index(parts[1], ")/")
	if hostEnd == -1 {
		return cfg, fmt.Errorf("DSN mal formado")
	}
	hostPort := strings.SplitN(parts[1][:hostEnd], ":", 2)
	if len(hostPort) != 2 {
		return cfg, fmt.Errorf("Host:Port inválido")
	}
	cfg.Host, cfg.Port = hostPort[0], hostPort[1]

	dbPart := parts[1][hostEnd+2:]
	if i := strings.Index(dbPart, "?"); i != -1 {
		dbPart = dbPart[:i]
	}
	cfg.Databases = []string{dbPart}
	cfg.BackupDir = os.Getenv("APP_ROOT") + "/backups"

	return cfg, nil
}

func extractDBName(dsn string) (string, error) {
	beforeParams := strings.SplitN(dsn, "?", 2)[0]
	parts := strings.Split(beforeParams, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("DB no encontrada en DSN")
	}
	return parts[len(parts)-1], nil
}

func checkpointPath(db, dir string) string {
	return filepath.Join(dir, fmt.Sprintf("%s_checkpoint.json", db))
}

func backupExists(db, dir string) bool {
	_, err := os.Stat(checkpointPath(db, dir))
	return err == nil
}

func getDBConnection(cfg *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func runFullBackup(cfg *Config, db string) error {
	ts := time.Now().Format("2006-01-02_15-04-05")
	path := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_full_%s.sql", db, ts))
	args := []string{
		"-u", cfg.User,
		"-p" + cfg.Password,
		"-h", cfg.Host,
		"-P", cfg.Port,
		"--databases", db,
		"--routines", "--events", "--single-transaction",
	}

	cmd := exec.Command("mariadb-dump", args...)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	cmd.Stdout = file

	return cmd.Run()
}

func getBinlogStatus(cfg *Config) (Checkpoint, error) {
	db, err := getDBConnection(cfg)
	if err != nil {
		return Checkpoint{}, fmt.Errorf("error conectando: %w", err)
	}
	defer db.Close()

	rows, err := db.Query("SHOW MASTER STATUS")
	if err != nil {
		return Checkpoint{}, fmt.Errorf("error ejecutando SHOW MASTER STATUS: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return Checkpoint{}, fmt.Errorf("no hay binlog activo en el servidor")
	}

	columns, err := rows.Columns()
	if err != nil {
		return Checkpoint{}, fmt.Errorf("error obteniendo columnas: %w", err)
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return Checkpoint{}, fmt.Errorf("error escaneando resultado: %w", err)
	}

	var file string
	var position uint32

	if values[0] != nil {
		if b, ok := values[0].([]byte); ok {
			file = string(b)
		} else if s, ok := values[0].(string); ok {
			file = s
		}
	}

	if values[1] != nil {
		switch v := values[1].(type) {
		case int64:
			position = uint32(v)
		case uint64:
			position = uint32(v)
		case []byte:
			fmt.Sscanf(string(v), "%d", &position)
		case string:
			fmt.Sscanf(v, "%d", &position)
		}
	}

	if file == "" {
		return Checkpoint{}, fmt.Errorf("no hay binlog activo en el servidor")
	}

	cp := Checkpoint{
		BinlogFile: file,
		Position:   position,
	}

	return cp, nil
}

func getTableSchema(cfg *Config, database, table string) (*TableSchema, error) {
	db, err := getDBConnection(cfg)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s' ORDER BY ORDINAL_POSITION", database, table)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	return &TableSchema{Columns: columns}, nil
}

func saveCheckpoint(cfg *Config, db string, cp Checkpoint) error {
	data, _ := json.MarshalIndent(cp, "", "  ")
	return os.WriteFile(checkpointPath(db, cfg.BackupDir), data, 0644)
}

func loadCheckpoint(cfg *Config, db string) (Checkpoint, error) {
	data, err := os.ReadFile(checkpointPath(db, cfg.BackupDir))
	if err != nil {
		return Checkpoint{}, err
	}
	var cp Checkpoint
	json.Unmarshal(data, &cp)
	return cp, nil
}

func extractBinlogForDatabase(cfg *Config, cp Checkpoint, currentPos uint32, database string, outputFile string) error {
	port, _ := strconv.Atoi(cfg.Port)

	syncCfg := replication.BinlogSyncerConfig{
		ServerID: 100 + uint32(time.Now().UnixNano()%1000), // ServerID único por goroutine
		Flavor:   "mariadb",
		Host:     cfg.Host,
		Port:     uint16(port),
		User:     cfg.User,
		Password: cfg.Password,
	}

	syncer := replication.NewBinlogSyncer(syncCfg)
	defer syncer.Close()

	streamer, err := syncer.StartSync(mysql.Position{
		Name: cp.BinlogFile,
		Pos:  cp.Position,
	})
	if err != nil {
		return fmt.Errorf("error iniciando sync: %w", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	ctx := context.Background()

	eventCount := 0
	sqlCount := 0
	schemaCache := make(map[string]*TableSchema)

	rotateCount := 0
	consecutiveEmptyEvents := 0

	for {
		ev, err := streamer.GetEvent(ctx)
		if err != nil {
			if sqlCount > 0 {
				break
			}
			return fmt.Errorf("error leyendo evento: %w", err)
		}

		eventCount++

		if ev.Header.LogPos >= currentPos {
			break
		}

		if ev.Header.EventType == replication.ROTATE_EVENT {
			rotateCount++
			if rotateCount >= 2 {
				break
			}
			continue
		}

		sql := eventToSQL(cfg, ev, database, schemaCache)
		if sql != "" {
			sqlCount++
			consecutiveEmptyEvents = 0
			fmt.Fprintf(file, "%s;\n\n", sql)
		} else {
			consecutiveEmptyEvents++
			if consecutiveEmptyEvents >= 10 && sqlCount > 0 {
				break
			}
		}
	}

	log.Info().Msgf("✅ [%s] Procesados %d eventos, %d sentencias SQL", database, eventCount, sqlCount)

	return nil
}

func eventToSQL(cfg *Config, ev *replication.BinlogEvent, targetDB string, schemaCache map[string]*TableSchema) string {
	switch e := ev.Event.(type) {
	case *replication.QueryEvent:
		schema := string(e.Schema)
		query := string(e.Query)

		if schema != targetDB {
			return ""
		}

		if strings.HasPrefix(query, "BEGIN") ||
			strings.HasPrefix(query, "COMMIT") ||
			strings.HasPrefix(query, "ROLLBACK") ||
			strings.HasPrefix(query, "SET ") ||
			strings.HasPrefix(query, "DELIMITER") {
			return ""
		}

		return query

	case *replication.RowsEvent:
		schema := string(e.Table.Schema)
		table := string(e.Table.Table)

		if schema != targetDB {
			return ""
		}

		cacheKey := fmt.Sprintf("%s.%s", schema, table)
		if _, exists := schemaCache[cacheKey]; !exists {
			tableSchema, err := getTableSchema(cfg, schema, table)
			if err != nil {
				log.Warn().Msgf("⚠️  Error obteniendo schema de %s.%s: %v", schema, table, err)
				return ""
			}
			schemaCache[cacheKey] = tableSchema
		}

		tableSchema := schemaCache[cacheKey]

		switch ev.Header.EventType {
		case replication.WRITE_ROWS_EVENTv0, replication.WRITE_ROWS_EVENTv1, replication.WRITE_ROWS_EVENTv2:
			return generateInsert(schema, table, e, tableSchema)
		case replication.UPDATE_ROWS_EVENTv0, replication.UPDATE_ROWS_EVENTv1, replication.UPDATE_ROWS_EVENTv2:
			return generateUpdate(schema, table, e, tableSchema)
		case replication.DELETE_ROWS_EVENTv0, replication.DELETE_ROWS_EVENTv1, replication.DELETE_ROWS_EVENTv2:
			return generateDelete(schema, table, e, tableSchema)
		}
	}

	return ""
}

func generateInsert(schema, table string, e *replication.RowsEvent, tableSchema *TableSchema) string {
	if len(e.Rows) == 0 {
		return ""
	}

	var sql strings.Builder

	cols := strings.Join(tableSchema.Columns, "`, `")
	sql.WriteString(fmt.Sprintf("INSERT INTO `%s`.`%s` (`%s`) VALUES ", schema, table, cols))

	for i, row := range e.Rows {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString("(")
		for j, val := range row {
			if j > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(formatValue(val))
		}
		sql.WriteString(")")
	}

	return sql.String()
}

func generateUpdate(schema, table string, e *replication.RowsEvent, tableSchema *TableSchema) string {
	var queries []string

	for i := 0; i < len(e.Rows); i += 2 {
		if i+1 >= len(e.Rows) {
			break
		}

		before := e.Rows[i]
		after := e.Rows[i+1]

		var sql strings.Builder
		sql.WriteString(fmt.Sprintf("UPDATE `%s`.`%s` SET ", schema, table))

		first := true
		for j, val := range after {
			if j < len(tableSchema.Columns) {
				if !first {
					sql.WriteString(", ")
				}
				sql.WriteString(fmt.Sprintf("`%s`=%s", tableSchema.Columns[j], formatValue(val)))
				first = false
			}
		}

		sql.WriteString(" WHERE ")
		first = true
		for j, val := range before {
			if j < len(tableSchema.Columns) {
				if !first {
					sql.WriteString(" AND ")
				}
				sql.WriteString(fmt.Sprintf("`%s`=%s", tableSchema.Columns[j], formatValue(val)))
				first = false
			}
		}

		queries = append(queries, sql.String())
	}

	return strings.Join(queries, ";\n")
}

func generateDelete(schema, table string, e *replication.RowsEvent, tableSchema *TableSchema) string {
	var queries []string

	for _, row := range e.Rows {
		var sql strings.Builder
		sql.WriteString(fmt.Sprintf("DELETE FROM `%s`.`%s` WHERE ", schema, table))

		first := true
		for j, val := range row {
			if j < len(tableSchema.Columns) {
				if !first {
					sql.WriteString(" AND ")
				}
				sql.WriteString(fmt.Sprintf("`%s`=%s", tableSchema.Columns[j], formatValue(val)))
				first = false
			}
		}

		queries = append(queries, sql.String())
	}

	return strings.Join(queries, ";\n")
}

func formatValue(val interface{}) string {
	if val == nil {
		return "NULL"
	}

	switch v := val.(type) {
	case string:
		escaped := strings.ReplaceAll(v, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "'", "\\'")
		return fmt.Sprintf("'%s'", escaped)
	case []byte:
		escaped := strings.ReplaceAll(string(v), "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "'", "\\'")
		return fmt.Sprintf("'%s'", escaped)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		if v {
			return "1"
		}
		return "0"
	case time.Time:
		return fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

// Worker para procesar backups full
func fullBackupWorker(cfg *Config, tasks <-chan string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	for db := range tasks {
		if err := runFullBackup(cfg, db); err != nil {
			log.Err(err).Msgf("❌ [Worker] Error en backup full de %s", db)
			results <- fmt.Sprintf("ERROR:%s", db)
			continue
		}

		cp, err := getBinlogStatus(cfg)
		if err != nil {
			log.Err(err).Msgf("❌ [Worker] Error obteniendo binlog status para %s", db)
			results <- fmt.Sprintf("ERROR:%s", db)
			continue
		}

		saveCheckpoint(cfg, db, cp)
		results <- fmt.Sprintf("SUCCESS:%s", db)
	}
}

// Worker para procesar backups incrementales
func incrementalBackupWorker(cfg *Config, currentPos uint32, tasks <-chan string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	for db := range tasks {
		cp, err := loadCheckpoint(cfg, db)
		if err != nil {
			log.Err(err).Msgf("❌ [Worker] Error cargando checkpoint de %s", db)
			results <- fmt.Sprintf("ERROR:%s", db)
			continue
		}

		if cp.BinlogFile == "" {
			log.Warn().Msgf("⚠️  [Worker] Checkpoint inválido para %s, skip", db)
			results <- fmt.Sprintf("SKIP:%s", db)
			continue
		}

		if cp.Position >= currentPos {
			results <- fmt.Sprintf("NOCHANGES:%s", db)
			continue
		}

		ts := time.Now().Format("2006-01-02_15-04-05")
		finalFile := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_incremental_%s.sql", db, ts))

		if err := extractBinlogForDatabase(cfg, cp, currentPos, db, finalFile); err != nil {
			log.Err(err).Msgf("❌ [Worker] Error extrayendo binlog para %s", db)
			results <- fmt.Sprintf("ERROR:%s", db)
			continue
		}

		content, err := os.ReadFile(finalFile)
		if err != nil {
			log.Err(err).Msgf("❌ [Worker] Error leyendo binlog de %s", db)
			os.Remove(finalFile)
			results <- fmt.Sprintf("ERROR:%s", db)
			continue
		}

		if len(content) == 0 {
			os.Remove(finalFile)
		} else {
			log.Info().Msgf("✅ [Worker] %s: Backup incremental generado (%d bytes)", db, len(content))
		}
	
		binlog := ""
		if currentPos > 0 {
			binlog = cp.BinlogFile
		}

		saveCheckpoint(cfg, db, Checkpoint{
			BinlogFile: binlog,
			Position:   currentPos,
		})

		results <- fmt.Sprintf("SUCCESS:%s", db)
	}
}

func RunBackup(cfg *Config) {
	log.Info().Msgf("⏰ [CRON] Iniciando backup de %d bases de datos... MODO READONLY", len(cfg.Databases))
	SetReadOnly(true)
	defer SetReadOnly(false)

	const maxWorkers = 10

	// Fase 1: Identificar qué backups full necesitamos
	var fullBackupDBs []string
	for _, db := range cfg.Databases {
		if !backupExists(db, cfg.BackupDir) {
			fullBackupDBs = append(fullBackupDBs, db)
		}
	}

	// Si hay backups full, procesarlos en paralelo
	if len(fullBackupDBs) > 0 {
		tasks := make(chan string, len(fullBackupDBs))
		results := make(chan string, len(fullBackupDBs))
		var wg sync.WaitGroup

		// Crear workers (máximo 10)
		numWorkers := maxWorkers
		if len(fullBackupDBs) < maxWorkers {
			numWorkers = len(fullBackupDBs)
		}

		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go fullBackupWorker(cfg, tasks, &wg, results)
		}

		// Enviar tareas
		for _, db := range fullBackupDBs {
			tasks <- db
		}
		close(tasks)

		// Esperar a que terminen
		wg.Wait()
		close(results)

		// Recolectar resultados
		successCount := 0
		for result := range results {
			if strings.HasPrefix(result, "SUCCESS") {
				successCount++
			}
		}

		log.Info().Msgf("✅ Backups FULL completados: %d/%d exitosos", successCount, len(fullBackupDBs))
		return
	}

	// Fase 2: Backups incrementales
	currentBinlogStatus, err := getBinlogStatus(cfg)
	if err != nil {
		log.Err(err).Msgf("❌ Error obteniendo binlog status actual")
		return
	}

	// Filtrar DBs que necesitan incremental
	var incrementalDBs []string
	for _, db := range cfg.Databases {
		cp, err := loadCheckpoint(cfg, db)
		if err == nil && cp.Position < currentBinlogStatus.Position {
			incrementalDBs = append(incrementalDBs, db)
		}
	}

	if len(incrementalDBs) == 0 {
		return
	}

	tasks := make(chan string, len(incrementalDBs))
	results := make(chan string, len(incrementalDBs))
	var wg sync.WaitGroup

	numWorkers := maxWorkers
	if len(incrementalDBs) < maxWorkers {
		numWorkers = len(incrementalDBs)
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go incrementalBackupWorker(cfg, currentBinlogStatus.Position, tasks, &wg, results)
	}

	for _, db := range incrementalDBs {
		tasks <- db
	}
	close(tasks)

	wg.Wait()
	close(results)

	successCount := 0
	for result := range results {
		if strings.HasPrefix(result, "SUCCESS") {
			successCount++
		}
	}

	log.Info().Msgf("✅ [CRON] Backup completado: %d/%d exitosos", successCount, len(incrementalDBs))
}
