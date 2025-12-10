package jobs

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

	log.Printf("✓ Binlog status: File=%s, Position=%d (columnas detectadas: %d)", cp.BinlogFile, cp.Position, len(columns))

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

func resetCheckpoint(cfg *Config, db string) error {
	cpPath := checkpointPath(db, cfg.BackupDir)
	if err := os.Remove(cpPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	log.Printf("⚠️  Checkpoint eliminado para %s, se hará backup full", db)
	return nil
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

	log.Printf("🔄 [%s] Extrayendo eventos desde %d hasta %d...", database, cp.Position, currentPos)

	rotateCount := 0
	consecutiveEmptyEvents := 0

	for {
		ev, err := streamer.GetEvent(ctx)
		if err != nil {
			if sqlCount > 0 {
				log.Printf("✅ [%s] Finalizado con %d SQL statements", database, sqlCount)
				break
			}
			return fmt.Errorf("error leyendo evento: %w", err)
		}

		eventCount++

		if ev.Header.LogPos >= currentPos {
			log.Printf("✅ [%s] Alcanzada posición objetivo %d", database, currentPos)
			break
		}

		if ev.Header.EventType == replication.ROTATE_EVENT {
			rotateCount++
			log.Printf("🔄 [%s] ROTATE_EVENT #%d detectado", database, rotateCount)
			if rotateCount >= 2 {
				log.Printf("✅ [%s] Múltiples ROTATE detectados, finalizando", database)
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
				log.Printf("✅ [%s] Sin más eventos relevantes, finalizando", database)
				break
			}
		}
	}

	log.Printf("✅ [%s] Procesados %d eventos, %d sentencias SQL", database, eventCount, sqlCount)

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
				log.Printf("⚠️  Error obteniendo schema de %s.%s: %v", schema, table, err)
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
		log.Printf("📦 [Worker] Ejecutando backup full para: %s", db)
		if err := runFullBackup(cfg, db); err != nil {
			log.Printf("❌ [Worker] Error en backup full de %s: %v", db, err)
			results <- fmt.Sprintf("ERROR:%s", db)
			continue
		}

		cp, err := getBinlogStatus(cfg)
		if err != nil {
			log.Printf("❌ [Worker] Error obteniendo binlog status para %s: %v", db, err)
			results <- fmt.Sprintf("ERROR:%s", db)
			continue
		}

		log.Printf("✅ [Worker] Checkpoint inicial para %s: File=%s, Position=%d", db, cp.BinlogFile, cp.Position)
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
			log.Printf("❌ [Worker] Error cargando checkpoint de %s: %v", db, err)
			results <- fmt.Sprintf("ERROR:%s", db)
			continue
		}

		if cp.BinlogFile == "" {
			log.Printf("⚠️  [Worker] Checkpoint inválido para %s, skip", db)
			results <- fmt.Sprintf("SKIP:%s", db)
			continue
		}

		if cp.Position >= currentPos {
			log.Printf("ℹ️  [Worker] %s: Sin cambios nuevos", db)
			results <- fmt.Sprintf("NOCHANGES:%s", db)
			continue
		}

		log.Printf("📋 [Worker] %s: Procesando desde position %d hasta %d", db, cp.Position, currentPos)

		ts := time.Now().Format("2006-01-02_15-04-05")
		finalFile := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_incremental_%s.sql", db, ts))

		if err := extractBinlogForDatabase(cfg, cp, currentPos, db, finalFile); err != nil {
			log.Printf("❌ [Worker] Error extrayendo binlog para %s: %v", db, err)
			results <- fmt.Sprintf("ERROR:%s", db)
			continue
		}

		content, err := os.ReadFile(finalFile)
		if err != nil {
			log.Printf("❌ [Worker] Error leyendo binlog de %s: %v", db, err)
			os.Remove(finalFile)
			results <- fmt.Sprintf("ERROR:%s", db)
			continue
		}

		if len(content) == 0 {
			log.Printf("ℹ️  [Worker] %s: Sin cambios SQL ejecutables", db)
			os.Remove(finalFile)
		} else {
			log.Printf("✅ [Worker] %s: Backup incremental generado (%d bytes)", db, len(content))
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
	log.Printf("⏰ [CRON] Iniciando backup de %d bases de datos... MODO READONLY", len(cfg.Databases))
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
		log.Printf("🚀 Procesando %d backups FULL con hasta %d workers paralelos", len(fullBackupDBs), maxWorkers)

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

		log.Printf("✅ Backups FULL completados: %d/%d exitosos", successCount, len(fullBackupDBs))
		return
	}

	// Fase 2: Backups incrementales
	currentBinlogStatus, err := getBinlogStatus(cfg)
	if err != nil {
		log.Printf("❌ Error obteniendo binlog status actual: %v", err)
		return
	}

	log.Printf("📊 Estado actual del binlog: File=%s, Position=%d", currentBinlogStatus.BinlogFile, currentBinlogStatus.Position)

	// Filtrar DBs que necesitan incremental
	var incrementalDBs []string
	for _, db := range cfg.Databases {
		cp, err := loadCheckpoint(cfg, db)
		if err == nil && cp.Position < currentBinlogStatus.Position {
			incrementalDBs = append(incrementalDBs, db)
		}
	}

	if len(incrementalDBs) == 0 {
		log.Printf("ℹ️  No hay cambios para procesar en ninguna base de datos")
		return
	}

	log.Printf("🚀 Procesando %d backups INCREMENTALES con hasta %d workers paralelos", len(incrementalDBs), maxWorkers)

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

	log.Printf("✅ [CRON] Backup completado: %d/%d exitosos", successCount, len(incrementalDBs))
}

// package jobs

// import (
// 	"context"
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/go-mysql-org/go-mysql/mysql"
// 	"github.com/go-mysql-org/go-mysql/replication"
// )

// type Config struct {
// 	User      string   `json:"user"`
// 	Password  string   `json:"password"`
// 	Databases []string `json:"databases"`
// 	Host      string   `json:"host"`
// 	Port      string   `json:"port"`
// 	BackupDir string   `json:"backup_dir"`
// }

// type Checkpoint struct {
// 	BinlogFile string `json:"binlog_file"`
// 	Position   uint32 `json:"position"`
// }

// type TableSchema struct {
// 	Columns []string
// }

// func LoadConfig(deps *dependencies.MainContainer) (*Config, error) {
// 	dsn := os.Getenv("URI_DB")
// 	cfg, err := parseDSN(dsn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	tenants, err := deps.TenantController.TenantService.TenantGetConections()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, conn := range tenants {
// 		dbName, err := extractDBName(conn.Connection)
// 		if err != nil {
// 			return nil, err
// 		}
// 		cfg.Databases = append(cfg.Databases, dbName)
// 	}
// 	return &cfg, nil
// }

// func parseDSN(dsn string) (Config, error) {
// 	var cfg Config
// 	parts := strings.SplitN(dsn, "@tcp(", 2)
// 	if len(parts) != 2 {
// 		return cfg, fmt.Errorf("DSN inválido")
// 	}
// 	up := strings.SplitN(parts[0], ":", 2)
// 	if len(up) != 2 {
// 		return cfg, fmt.Errorf("falta ':' en usuario:contraseña")
// 	}
// 	cfg.User, cfg.Password = up[0], up[1]

// 	hostEnd := strings.Index(parts[1], ")/")
// 	if hostEnd == -1 {
// 		return cfg, fmt.Errorf("DSN mal formado")
// 	}
// 	hostPort := strings.SplitN(parts[1][:hostEnd], ":", 2)
// 	if len(hostPort) != 2 {
// 		return cfg, fmt.Errorf("Host:Port inválido")
// 	}
// 	cfg.Host, cfg.Port = hostPort[0], hostPort[1]

// 	dbPart := parts[1][hostEnd+2:]
// 	if i := strings.Index(dbPart, "?"); i != -1 {
// 		dbPart = dbPart[:i]
// 	}
// 	cfg.Databases = []string{dbPart}
// 	cfg.BackupDir = os.Getenv("APP_ROOT") + "/backups"

// 	return cfg, nil
// }

// func extractDBName(dsn string) (string, error) {
// 	beforeParams := strings.SplitN(dsn, "?", 2)[0]
// 	parts := strings.Split(beforeParams, "/")
// 	if len(parts) < 2 {
// 		return "", fmt.Errorf("DB no encontrada en DSN")
// 	}
// 	return parts[len(parts)-1], nil
// }

// func checkpointPath(db, dir string) string {
// 	return filepath.Join(dir, fmt.Sprintf("%s_checkpoint.json", db))
// }

// func backupExists(db, dir string) bool {
// 	_, err := os.Stat(checkpointPath(db, dir))
// 	return err == nil
// }

// // Crear conexión SQL estándar
// func getDBConnection(cfg *Config) (*sql.DB, error) {
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := db.Ping(); err != nil {
// 		db.Close()
// 		return nil, err
// 	}

// 	return db, nil
// }

// // Full backup sigue usando mariadb-dump (más confiable)
// func runFullBackup(cfg *Config, db string) error {
// 	ts := time.Now().Format("2006-01-02_15-04-05")
// 	path := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_full_%s.sql", db, ts))
// 	args := []string{
// 		"-u", cfg.User,
// 		"-p" + cfg.Password,
// 		"-h", cfg.Host,
// 		"-P", cfg.Port,
// 		"--databases", db,
// 		"--routines", "--events", "--single-transaction",
// 	}

// 	cmd := exec.Command("mariadb-dump", args...)
// 	file, err := os.Create(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()
// 	cmd.Stdout = file

// 	return cmd.Run()
// }

// // Obtener binlog status usando database/sql
// func getBinlogStatus(cfg *Config) (Checkpoint, error) {
// 	db, err := getDBConnection(cfg)
// 	if err != nil {
// 		return Checkpoint{}, fmt.Errorf("error conectando: %w", err)
// 	}
// 	defer db.Close()

// 	rows, err := db.Query("SHOW MASTER STATUS")
// 	if err != nil {
// 		return Checkpoint{}, fmt.Errorf("error ejecutando SHOW MASTER STATUS: %w", err)
// 	}
// 	defer rows.Close()

// 	if !rows.Next() {
// 		return Checkpoint{}, fmt.Errorf("no hay binlog activo en el servidor")
// 	}

// 	// Obtener columnas para saber cuántas hay
// 	columns, err := rows.Columns()
// 	if err != nil {
// 		return Checkpoint{}, fmt.Errorf("error obteniendo columnas: %w", err)
// 	}

// 	// Crear slice de interfaces para escanear todas las columnas
// 	values := make([]interface{}, len(columns))
// 	valuePtrs := make([]interface{}, len(columns))
// 	for i := range values {
// 		valuePtrs[i] = &values[i]
// 	}

// 	if err := rows.Scan(valuePtrs...); err != nil {
// 		return Checkpoint{}, fmt.Errorf("error escaneando resultado: %w", err)
// 	}

// 	// Las primeras dos columnas siempre son File y Position
// 	var file string
// 	var position uint32

// 	if values[0] != nil {
// 		if b, ok := values[0].([]byte); ok {
// 			file = string(b)
// 		} else if s, ok := values[0].(string); ok {
// 			file = s
// 		}
// 	}

// 	if values[1] != nil {
// 		switch v := values[1].(type) {
// 		case int64:
// 			position = uint32(v)
// 		case uint64:
// 			position = uint32(v)
// 		case []byte:
// 			fmt.Sscanf(string(v), "%d", &position)
// 		case string:
// 			fmt.Sscanf(v, "%d", &position)
// 		}
// 	}

// 	if file == "" {
// 		return Checkpoint{}, fmt.Errorf("no hay binlog activo en el servidor")
// 	}

// 	cp := Checkpoint{
// 		BinlogFile: file,
// 		Position:   position,
// 	}

// 	log.Printf("✓ Binlog status: File=%s, Position=%d (columnas detectadas: %d)", cp.BinlogFile, cp.Position, len(columns))

// 	return cp, nil
// }

// // Obtener schema de una tabla
// func getTableSchema(cfg *Config, database, table string) (*TableSchema, error) {
// 	db, err := getDBConnection(cfg)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer db.Close()

// 	query := fmt.Sprintf("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s' ORDER BY ORDINAL_POSITION", database, table)

// 	rows, err := db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var columns []string
// 	for rows.Next() {
// 		var col string
// 		if err := rows.Scan(&col); err != nil {
// 			return nil, err
// 		}
// 		columns = append(columns, col)
// 	}

// 	return &TableSchema{Columns: columns}, nil
// }

// func saveCheckpoint(cfg *Config, db string, cp Checkpoint) error {
// 	data, _ := json.MarshalIndent(cp, "", "  ")
// 	return os.WriteFile(checkpointPath(db, cfg.BackupDir), data, 0644)
// }

// func loadCheckpoint(cfg *Config, db string) (Checkpoint, error) {
// 	data, err := os.ReadFile(checkpointPath(db, cfg.BackupDir))
// 	if err != nil {
// 		return Checkpoint{}, err
// 	}
// 	var cp Checkpoint
// 	json.Unmarshal(data, &cp)
// 	return cp, nil
// }

// func resetCheckpoint(cfg *Config, db string) error {
// 	cpPath := checkpointPath(db, cfg.BackupDir)
// 	if err := os.Remove(cpPath); err != nil && !os.IsNotExist(err) {
// 		return err
// 	}
// 	log.Printf("⚠️  Checkpoint eliminado para %s, se hará backup full", db)
// 	return nil
// }

// // Extraer binlog incremental usando go-mysql replication
// func extractBinlogForDatabase(cfg *Config, cp Checkpoint, database string, outputFile string) error {
// 	port, _ := strconv.Atoi(cfg.Port)

// 	syncCfg := replication.BinlogSyncerConfig{
// 		ServerID: 100,
// 		Flavor:   "mariadb",
// 		Host:     cfg.Host,
// 		Port:     uint16(port),
// 		User:     cfg.User,
// 		Password: cfg.Password,
// 	}

// 	syncer := replication.NewBinlogSyncer(syncCfg)
// 	defer syncer.Close()

// 	streamer, err := syncer.StartSync(mysql.Position{
// 		Name: cp.BinlogFile,
// 		Pos:  cp.Position,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("error iniciando sync: %w", err)
// 	}

// 	file, err := os.Create(outputFile)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// Timeout más corto: 3 segundos sin eventos = fin
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	eventCount := 0
// 	sqlCount := 0
// 	schemaCache := make(map[string]*TableSchema)
// 	lastEventTime := time.Now()

// 	log.Printf("🔄 [%s] Extrayendo eventos desde posición %d...", database, cp.Position)

// 	rotateCount := 0
// 	noEventCount := 0

// 	for {
// 		// Si llevamos más de 2 segundos sin eventos y ya tenemos algo, terminamos
// 		if time.Since(lastEventTime) > 2*time.Second && sqlCount > 0 {
// 			log.Printf("✅ [%s] 2s sin eventos, finalizando con %d SQL", database, sqlCount)
// 			break
// 		}

// 		select {
// 		case <-ctx.Done():
// 			log.Printf("✅ [%s] Timeout alcanzado, finalizando extracción", database)
// 			goto done
// 		default:
// 		}

// 		// Timeout más agresivo por evento: 1 segundo
// 		eventCtx, eventCancel := context.WithTimeout(context.Background(), 1*time.Second)
// 		ev, err := streamer.GetEvent(eventCtx)
// 		eventCancel()

// 		if err != nil {
// 			if err == context.DeadlineExceeded {
// 				// Si timeout y ya tenemos SQL, es porque terminamos
// 				if sqlCount > 0 {
// 					log.Printf("✅ [%s] No hay más eventos (timeout 1s), finalizando", database)
// 					break
// 				}
// 				goto done
// 			}
// 			if sqlCount > 0 {
// 				log.Printf("⚠️  [%s] Error leyendo más eventos, finalizando con %d SQL", database, sqlCount)
// 				goto done
// 			}
// 			return fmt.Errorf("error leyendo evento: %w", err)
// 		}

// 		eventCount++
// 		lastEventTime = time.Now()

// 		if ev.Header.EventType == replication.ROTATE_EVENT {
// 			rotateCount++
// 			log.Printf("🔄 [%s] ROTATE_EVENT #%d detectado, continuando lectura...", database, rotateCount)
// 			if rotateCount >= 2 {
// 				log.Printf("✅ [%s] Múltiples ROTATE detectados, finalizando", database)
// 				break
// 			}
// 			continue
// 		}

// 		sql := eventToSQL(cfg, ev, database, schemaCache)
// 		if sql != "" {
// 			sqlCount++
// 			noEventCount = 0
// 			fmt.Fprintf(file, "%s;\n\n", sql)
// 		} else {
// 			noEventCount++
// 			// Si llevamos 3 eventos vacíos y tenemos SQL, terminamos
// 			if noEventCount >= 3 && sqlCount > 0 {
// 				log.Printf("✅ [%s] Sin más eventos relevantes, finalizando", database)
// 				break
// 			}
// 		}
// 	}

// done:
// 	log.Printf("✅ [%s] Procesados %d eventos, %d sentencias SQL", database, eventCount, sqlCount)

// 	return nil
// }

// func eventToSQL(cfg *Config, ev *replication.BinlogEvent, targetDB string, schemaCache map[string]*TableSchema) string {
// 	switch e := ev.Event.(type) {
// 	case *replication.QueryEvent:
// 		schema := string(e.Schema)
// 		query := string(e.Query)

// 		if schema != targetDB {
// 			return ""
// 		}

// 		if strings.HasPrefix(query, "BEGIN") ||
// 			strings.HasPrefix(query, "COMMIT") ||
// 			strings.HasPrefix(query, "ROLLBACK") ||
// 			strings.HasPrefix(query, "SET ") ||
// 			strings.HasPrefix(query, "DELIMITER") {
// 			return ""
// 		}

// 		return query

// 	case *replication.RowsEvent:
// 		schema := string(e.Table.Schema)
// 		table := string(e.Table.Table)

// 		if schema != targetDB {
// 			return ""
// 		}

// 		// Obtener schema de la tabla
// 		cacheKey := fmt.Sprintf("%s.%s", schema, table)
// 		if _, exists := schemaCache[cacheKey]; !exists {
// 			tableSchema, err := getTableSchema(cfg, schema, table)
// 			if err != nil {
// 				log.Printf("⚠️  Error obteniendo schema de %s.%s: %v", schema, table, err)
// 				return ""
// 			}
// 			schemaCache[cacheKey] = tableSchema
// 		}

// 		tableSchema := schemaCache[cacheKey]

// 		switch ev.Header.EventType {
// 		case replication.WRITE_ROWS_EVENTv0, replication.WRITE_ROWS_EVENTv1, replication.WRITE_ROWS_EVENTv2:
// 			return generateInsert(schema, table, e, tableSchema)
// 		case replication.UPDATE_ROWS_EVENTv0, replication.UPDATE_ROWS_EVENTv1, replication.UPDATE_ROWS_EVENTv2:
// 			return generateUpdate(schema, table, e, tableSchema)
// 		case replication.DELETE_ROWS_EVENTv0, replication.DELETE_ROWS_EVENTv1, replication.DELETE_ROWS_EVENTv2:
// 			return generateDelete(schema, table, e, tableSchema)
// 		}
// 	}

// 	return ""
// }

// func generateInsert(schema, table string, e *replication.RowsEvent, tableSchema *TableSchema) string {
// 	if len(e.Rows) == 0 {
// 		return ""
// 	}

// 	var sql strings.Builder

// 	// Construir lista de columnas
// 	cols := strings.Join(tableSchema.Columns, "`, `")
// 	sql.WriteString(fmt.Sprintf("INSERT INTO `%s`.`%s` (`%s`) VALUES ", schema, table, cols))

// 	for i, row := range e.Rows {
// 		if i > 0 {
// 			sql.WriteString(", ")
// 		}
// 		sql.WriteString("(")
// 		for j, val := range row {
// 			if j > 0 {
// 				sql.WriteString(", ")
// 			}
// 			sql.WriteString(formatValue(val))
// 		}
// 		sql.WriteString(")")
// 	}

// 	return sql.String()
// }

// func generateUpdate(schema, table string, e *replication.RowsEvent, tableSchema *TableSchema) string {
// 	var queries []string

// 	for i := 0; i < len(e.Rows); i += 2 {
// 		if i+1 >= len(e.Rows) {
// 			break
// 		}

// 		before := e.Rows[i]
// 		after := e.Rows[i+1]

// 		var sql strings.Builder
// 		sql.WriteString(fmt.Sprintf("UPDATE `%s`.`%s` SET ", schema, table))

// 		// SET clause
// 		first := true
// 		for j, val := range after {
// 			if j < len(tableSchema.Columns) {
// 				if !first {
// 					sql.WriteString(", ")
// 				}
// 				sql.WriteString(fmt.Sprintf("`%s`=%s", tableSchema.Columns[j], formatValue(val)))
// 				first = false
// 			}
// 		}

// 		// WHERE clause
// 		sql.WriteString(" WHERE ")
// 		first = true
// 		for j, val := range before {
// 			if j < len(tableSchema.Columns) {
// 				if !first {
// 					sql.WriteString(" AND ")
// 				}
// 				sql.WriteString(fmt.Sprintf("`%s`=%s", tableSchema.Columns[j], formatValue(val)))
// 				first = false
// 			}
// 		}

// 		queries = append(queries, sql.String())
// 	}

// 	return strings.Join(queries, ";\n")
// }

// func generateDelete(schema, table string, e *replication.RowsEvent, tableSchema *TableSchema) string {
// 	var queries []string

// 	for _, row := range e.Rows {
// 		var sql strings.Builder
// 		sql.WriteString(fmt.Sprintf("DELETE FROM `%s`.`%s` WHERE ", schema, table))

// 		first := true
// 		for j, val := range row {
// 			if j < len(tableSchema.Columns) {
// 				if !first {
// 					sql.WriteString(" AND ")
// 				}
// 				sql.WriteString(fmt.Sprintf("`%s`=%s", tableSchema.Columns[j], formatValue(val)))
// 				first = false
// 			}
// 		}

// 		queries = append(queries, sql.String())
// 	}

// 	return strings.Join(queries, ";\n")
// }

// func formatValue(val interface{}) string {
// 	if val == nil {
// 		return "NULL"
// 	}

// 	switch v := val.(type) {
// 	case string:
// 		escaped := strings.ReplaceAll(v, "\\", "\\\\")
// 		escaped = strings.ReplaceAll(escaped, "'", "\\'")
// 		return fmt.Sprintf("'%s'", escaped)
// 	case []byte:
// 		escaped := strings.ReplaceAll(string(v), "\\", "\\\\")
// 		escaped = strings.ReplaceAll(escaped, "'", "\\'")
// 		return fmt.Sprintf("'%s'", escaped)
// 	case int, int8, int16, int32, int64:
// 		return fmt.Sprintf("%d", v)
// 	case uint, uint8, uint16, uint32, uint64:
// 		return fmt.Sprintf("%d", v)
// 	case float32, float64:
// 		return fmt.Sprintf("%v", v)
// 	case bool:
// 		if v {
// 			return "1"
// 		}
// 		return "0"
// 	case time.Time:
// 		return fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
// 	default:
// 		return fmt.Sprintf("'%v'", v)
// 	}
// }

// func RunBackup(cfg *Config) {
// log.Printf("⏰ [CRON] Iniciando backup de %d bases de datos... MODO READONLY", len(cfg.Databases))
// 	SetReadOnly(true)
// 	defer SetReadOnly(false)

// 	needsFullBackup := false
// 	for _, db := range cfg.Databases {
// 		if !backupExists(db, cfg.BackupDir) {
// 			log.Printf("DB '%s' necesita backup full", db)
// 			needsFullBackup = true
// 		}
// 	}

// 	if needsFullBackup {
// 		for _, db := range cfg.Databases {
// 			if !backupExists(db, cfg.BackupDir) {
// 				log.Printf("📦 Ejecutando backup full para: %s", db)
// 				if err := runFullBackup(cfg, db); err != nil {
// 					log.Printf("❌ Error en backup full de %s: %v", db, err)
// 					continue
// 				}

// 				cp, err := getBinlogStatus(cfg)
// 				if err != nil {
// 					log.Printf("❌ Error obteniendo binlog status: %v", err)
// 					continue
// 				}

// 				log.Printf("✅ Checkpoint inicial para %s: File=%s, Position=%d", db, cp.BinlogFile, cp.Position)
// 				saveCheckpoint(cfg, db, cp)
// 			}
// 		}
// 		return
// 	}

// 	currentBinlogStatus, err := getBinlogStatus(cfg)
// 	if err != nil {
// 		log.Printf("❌ Error obteniendo binlog status actual: %v", err)
// 		return
// 	}

// 	log.Printf("📊 Estado actual del binlog: File=%s, Position=%d",
// 		currentBinlogStatus.BinlogFile, currentBinlogStatus.Position)

// 	for _, db := range cfg.Databases {
// 		cp, err := loadCheckpoint(cfg, db)
// 		if err != nil {
// 			log.Printf("❌ Error cargando checkpoint de %s: %v", db, err)
// 			continue
// 		}

// 		if cp.BinlogFile == "" {
// 			log.Printf("⚠️  Checkpoint inválido para %s, regenerando...", db)
// 			resetCheckpoint(cfg, db)
// 			continue
// 		}

// 		if cp.Position >= currentBinlogStatus.Position {
// 			log.Printf("ℹ️  %s: Sin cambios nuevos (checkpoint: %d >= actual: %d)",
// 				db, cp.Position, currentBinlogStatus.Position)
// 			continue
// 		}

// 		log.Printf("📋 %s: Procesando desde position %d hasta %d",
// 			db, cp.Position, currentBinlogStatus.Position)

// 		ts := time.Now().Format("2006-01-02_15-04-05")
// 		finalFile := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_incremental_%s.sql", db, ts))

// 		if err := extractBinlogForDatabase(cfg, cp, db, finalFile); err != nil {
// 			log.Printf("❌ Error extrayendo binlog para %s: %v", db, err)
// 			continue
// 		}

// 		content, err := os.ReadFile(finalFile)
// 		if err != nil {
// 			log.Printf("❌ Error leyendo binlog de %s: %v", db, err)
// 			os.Remove(finalFile)
// 			continue
// 		}

// 		if len(content) == 0 {
// 			log.Printf("ℹ️  %s: Sin cambios SQL ejecutables", db)
// 			os.Remove(finalFile)
// 		} else {
// 			log.Printf("✅ %s: Backup incremental generado (%d bytes)", db, len(content))
// 		}

// 		saveCheckpoint(cfg, db, currentBinlogStatus)
// 	}

// 	log.Printf("✅ [CRON] Backup completado")
// }
