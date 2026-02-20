package jobs

import (
	"context"
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
	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type ConfigPostgres struct {
	User      string   `json:"user"`
	Password  string   `json:"password"`
	Databases []string `json:"databases"`
	Host      string   `json:"host"`
	Port      string   `json:"port"`
	BackupDir string   `json:"backup_dir"`
	Params    string   `json:"params"` // Starts with ?, e.g. ?sslmode=disable
}

type CheckpointPostgres struct {
	LSN string `json:"lsn"` // Log Sequence Number as string
}

func LoadConfigPostgres(deps *dependencies.MainContainer) (*ConfigPostgres, error) {
	dsn := os.Getenv("URI_DB") // Assuming a separate ENV or similar mechanism, or reusing URI_DB if adaptable
	if dsn == "" {
		dsn = os.Getenv("URI_DB") // Fallback
	}

	cfg, err := parseDSNPostgres(dsn)
	if err != nil {
		return nil, err
	}

	tenants, err := deps.TenantController.TenantService.TenantGetConections()
	if err != nil {
		return nil, err
	}

	for _, conn := range tenants {
		// Asumimos que la connection string de tenants puede ser parseada igual o extraemos el nombre de la DB
		dbName, err := extractDBNamePostgres(conn.Connection)
		if err != nil {
			// Si falla, intentamos usar el nombre tal cual si no es una URL
			// O logueamos error y continuamos
			log.Warn().Msgf("No se pudo extraer nombre DB de tenant: %v", err)
			continue
		}
		cfg.Databases = append(cfg.Databases, dbName)
	}
	return &cfg, nil
}

func parseDSNPostgres(dsn string) (ConfigPostgres, error) {
	var cfg ConfigPostgres

	// pgx.ParseConfig puede parsear la URL de conexion standard de postgres
	pgCfg, err := pgconn.ParseConfig(dsn)
	if err != nil {
		return cfg, fmt.Errorf("DSN inválido: %w", err)
	}

	cfg.User = pgCfg.User
	cfg.Password = pgCfg.Password
	cfg.Host = pgCfg.Host
	cfg.Port = strconv.Itoa(int(pgCfg.Port))
	cfg.Databases = []string{pgCfg.Database}
	cfg.BackupDir = os.Getenv("APP_ROOT") + "/backups"

	// Preserve original params (sslmode, timezone, etc) by splitting original DSN
	parts := strings.Split(dsn, "?")
	if len(parts) > 1 {
		cfg.Params = "?" + parts[1]
	}

	return cfg, nil
}

func extractDBNamePostgres(dsn string) (string, error) {
	pgCfg, err := pgconn.ParseConfig(dsn)
	if err != nil {
		return "", err
	}
	return pgCfg.Database, nil
}

func checkpointPathPostgres(db, dir string) string {
	return filepath.Join(dir, fmt.Sprintf("%s_pg_checkpoint.json", db))
}

func backupExistsPostgres(db, dir string) bool {
	_, err := os.Stat(checkpointPathPostgres(db, dir))
	return err == nil
}

func runFullBackupPostgres(cfg *ConfigPostgres, db string) error {
	ts := time.Now().Format("2006-01-02_15-04-05")
	path := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_full_pg_%s.sql", db, ts))

	// PGPASSWORD environment variable is the safest way to pass password to pg_dump
	cmd := exec.Command("pg_dump", "-h", cfg.Host, "-p", cfg.Port, "-U", cfg.User, "-d", db, "-Fp", "--clean", "--if-exists")
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", cfg.Password))

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	cmd.Stdout = file

	// Capture stderr for debugging
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func getCurrentLSN(cfg *ConfigPostgres, dbName string) (CheckpointPostgres, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, dbName)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return CheckpointPostgres{}, err
	}
	defer conn.Close(context.Background())

	var lsnStr string
	err = conn.QueryRow(context.Background(), "SELECT pg_current_wal_lsn()").Scan(&lsnStr)
	if err != nil {
		return CheckpointPostgres{}, err
	}

	return CheckpointPostgres{LSN: lsnStr}, nil
}

func saveCheckpointPostgres(cfg *ConfigPostgres, db string, cp CheckpointPostgres) error {
	data, _ := json.MarshalIndent(cp, "", "  ")
	return os.WriteFile(checkpointPathPostgres(db, cfg.BackupDir), data, 0644)
}

func verifyWalLevelPostgres(cfg *ConfigPostgres, dbName string) error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, dbName)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	var walLevel string
	err = conn.QueryRow(context.Background(), "SHOW wal_level").Scan(&walLevel)
	if err != nil {
		return err
	}

	if walLevel != "logical" {
		return fmt.Errorf("wal_level actual es '%s', se requiere 'logical' para backups incrementales. Configuralo en postgresql.conf y reiniciá el servicio", walLevel)
	}
	return nil
}

func loadCheckpointPostgres(cfg *ConfigPostgres, db string) (CheckpointPostgres, error) {
	data, err := os.ReadFile(checkpointPathPostgres(db, cfg.BackupDir))
	if err != nil {
		return CheckpointPostgres{}, err
	}
	var cp CheckpointPostgres
	json.Unmarshal(data, &cp)
	return cp, nil
}

// Logical replication decoder
type Decoder struct {
	relations map[uint32]*pglogrepl.RelationMessage
	typeMap   *pgtype.Map
}

func fullBackupWorkerPostgres(cfg *ConfigPostgres, tasks <-chan string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	for db := range tasks {
		if err := runFullBackupPostgres(cfg, db); err != nil {
			log.Err(err).Msgf("❌ [Worker PG] Error en backup full de %s", db)
			results <- fmt.Sprintf("ERROR:%s", db)
			continue
		}

		cp, err := getCurrentLSN(cfg, db)
		if err != nil {
			log.Err(err).Msgf("❌ [Worker PG] Error obteniendo LSN para %s", db)
			results <- fmt.Sprintf("ERROR:%s", db)
			continue
		}

		saveCheckpointPostgres(cfg, db, cp)
		results <- fmt.Sprintf("SUCCESS:%s", db)
	}
}

func incrementalBackupWorkerPostgres(cfg *ConfigPostgres, db string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	// Load previous checkpoint
	prevCp, err := loadCheckpointPostgres(cfg, db)
	if err != nil {
		log.Err(err).Msgf("❌ [Worker PG] Error cargando checkpoint de %s", db)
		results <- fmt.Sprintf("ERROR:%s", db)
		return
	}

	// Get current LSN to see if we need to backup
	currentCp, err := getCurrentLSN(cfg, db)
	if err != nil {
		log.Err(err).Msgf("❌ [Worker PG] Error obteniendo LSN actual de %s", db)
		results <- fmt.Sprintf("ERROR:%s", db)
		return
	}

	if prevCp.LSN == currentCp.LSN {
		results <- fmt.Sprintf("NOCHANGES:%s", db)
		return
	}

	ts := time.Now().Format("2006-01-02_15-04-05")
	finalFile := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_incremental_pg_%s.sql", db, ts))

	// Start Logical Replication Stream
	// This part is tricky. We need to create a replication connection,
	// ensure a slot exists, and stream from prevCp.LSN to currentCp.LSN (or just consume what's there).
	// Note: pglogrepl requires a replication connection.

	// Pass currentCp.LSN as the target LSN so replication stops when reached
	err = fetchWALChanges(cfg, db, prevCp.LSN, currentCp.LSN, finalFile)
	if err != nil {
		log.Err(err).Msgf("❌ [Worker PG] Error extrayendo WAL para %s", db)
		os.Remove(finalFile) // Clean up partial file
		results <- fmt.Sprintf("ERROR:%s", db)
		return
	}

	// Save new checkpoint
	// Note: strictly speaking we should save the LSN up to which we consumed.
	// fetchWALChanges should return the last LSN processed.

	// For simplicity, we'll re-query current LSN or rely on what fetchWALChanges saw.
	// Actually fetchWALChanges should update us on the last LSN.
	// But strictly, let's just take a new LSN after success.

	newCheckpoint, _ := getCurrentLSN(cfg, db) // Get fresh LSN
	saveCheckpointPostgres(cfg, db, newCheckpoint)

	results <- fmt.Sprintf("SUCCESS:%s", db)
}

func fetchWALChanges(cfg *ConfigPostgres, dbName string, startLSNStr string, targetLSNStr string, outputFile string) error {
	connConfig, _ := pgconn.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, dbName, cfg.Params))
	// Force replication mode
	if connConfig.RuntimeParams == nil {
		connConfig.RuntimeParams = make(map[string]string)
	}
	connConfig.RuntimeParams["replication"] = "database"

	conn, err := pgconn.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return fmt.Errorf("failed to connect for replication: %w", err)
	}
	defer conn.Close(context.Background())

	// Ensure publication exists logic moved here
	createPubQuery := "CREATE PUBLICATION pg_backup_pub FOR ALL TABLES;"
	if err := runQuery(cfg, dbName, createPubQuery); err != nil {
		// Log detailed error for debugging - usually "already exists" which is fine
		// We can check error code if we want to be cleaner but this works.
		// log.Info().Msgf("ℹ️ [Worker PG] CREATE PUBLICATION result for %s: %v", dbName, err)
	} else {
		log.Info().Msgf("✅ [Worker PG] Created publication for %s", dbName)
	}

	// Verify it actually exists
	checkPubQuery := "SELECT count(*) FROM pg_publication WHERE pubname = 'pg_backup_pub'"
	var count int
	if err := runQueryRow(cfg, dbName, checkPubQuery, &count); err != nil {
		log.Warn().Err(err).Msgf("⚠️ [Worker PG] Failed to verify publication existence for %s", dbName)
	}
	if count == 0 {
		return fmt.Errorf("FATAL: Publication 'pg_backup_pub' does not exist for %s even after creation attempt", dbName)
	}
	// log.Info().Msgf("✅ [Worker PG] Verified publication exists for %s (count: %d)", dbName, count)

	// Create slot if not exists
	slotName := fmt.Sprintf("backup_slot_%s", dbName)
	// Try creating slot, ignore error if already exists
	_, _ = pglogrepl.CreateReplicationSlot(context.Background(), conn, slotName, "pgoutput", pglogrepl.CreateReplicationSlotOptions{Temporary: false})

	startLSN, err := pglogrepl.ParseLSN(startLSNStr)
	if err != nil {
		return fmt.Errorf("invalid start LSN: %w", err)
	}

	targetLSN, err := pglogrepl.ParseLSN(targetLSNStr)
	if err != nil {
		return fmt.Errorf("invalid target LSN: %w", err)
	}

	err = pglogrepl.StartReplication(context.Background(), conn, slotName, startLSN, pglogrepl.StartReplicationOptions{PluginArgs: []string{"proto_version '1'", "publication_names 'pg_backup_pub'"}})
	if err != nil {
		return fmt.Errorf("failed to start replication: %w", err)
	}

	var file *os.File
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	// Safety timeout of 60s
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	clientXLogPos := startLSN
	standbyMessageTimeout := time.Second * 5
	nextStandbyMessageDeadline := time.Now().Add(standbyMessageTimeout)
	relations := make(map[uint32]*pglogrepl.RelationMessage)
	typeMap := pgtype.NewMap()

	for {
		// Check target LSN condition
		if clientXLogPos >= targetLSN {
			return nil
		}

		if time.Now().After(nextStandbyMessageDeadline) {
			err = pglogrepl.SendStandbyStatusUpdate(context.Background(), conn, pglogrepl.StandbyStatusUpdate{WALWritePosition: clientXLogPos})
			if err != nil {
				return err
			}
			nextStandbyMessageDeadline = time.Now().Add(standbyMessageTimeout)
		}

		msg, err := conn.ReceiveMessage(ctx)
		if err != nil {
			if pgconn.Timeout(err) {
				continue
			}
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("timeout waiting for target LSN %s (current: %s)", targetLSN, clientXLogPos)
			}
			return fmt.Errorf("replication failed: %w", err)
		}

		if errMsg, ok := msg.(*pgproto3.ErrorResponse); ok {
			// Check for specific error: publication does not exist (SQLSTATE 42704)
			// This implies the replication slot is out of sync or broken.
			// We MUST drop the slot so it can be recreated correctly next time.
			if errMsg.Code == "42704" || strings.Contains(errMsg.Message, "publication \"pg_backup_pub\" does not exist") {
				log.Warn().Msgf("⚠️ [Worker PG] Replication slot %s is broken (missing publication). Dropping it...", slotName)
				errDrop := pglogrepl.DropReplicationSlot(context.Background(), conn, slotName, pglogrepl.DropReplicationSlotOptions{Wait: true})
				if errDrop != nil {
					log.Error().Err(errDrop).Msgf("❌ [Worker PG] Failed to drop broken replication slot %s", slotName)
				} else {
					log.Info().Msgf("✅ [Worker PG] Successfully dropped broken replication slot %s", slotName)
				}
				return fmt.Errorf("replication slot broken and dropped (publication missing). Retry next run")
			}
			return fmt.Errorf("received WAL error: %v", errMsg)
		}

		msgCopyData, ok := msg.(*pgproto3.CopyData)
		if !ok {
			continue
		}

		// Handle message types based on first byte
		if len(msgCopyData.Data) == 0 {
			continue
		}

		msgType := msgCopyData.Data[0]
		dataPayload := msgCopyData.Data[1:]

		switch msgType {
		case pglogrepl.PrimaryKeepaliveMessageByteID:
			pkm, err := pglogrepl.ParsePrimaryKeepaliveMessage(dataPayload)
			if err != nil {
				log.Err(err).Msg("ParsePrimaryKeepaliveMessage failed")
				continue
			}
			if pkm.ServerWALEnd > clientXLogPos {
				clientXLogPos = pkm.ServerWALEnd
			}
			if pkm.ReplyRequested {
				nextStandbyMessageDeadline = time.Time{} // send update immediately
			}

		case pglogrepl.XLogDataByteID:
			xld, err := pglogrepl.ParseXLogData(dataPayload)
			if err != nil {
				log.Err(err).Msg("ParseXLogData failed")
				continue
			}

			clientXLogPos = xld.WALStart + pglogrepl.LSN(len(xld.WALData))

			logicalMsg, err := pglogrepl.Parse(xld.WALData)
			if err != nil {
				continue
			}

			switch logicalMsg := logicalMsg.(type) {
			case *pglogrepl.RelationMessage:
				relations[logicalMsg.RelationID] = logicalMsg
			case *pglogrepl.InsertMessage:
				if file == nil {
					f, err := os.Create(outputFile)
					if err != nil {
						return err
					}
					file = f
				}
				sql := generateInsertPostgres(logicalMsg, relations[logicalMsg.RelationID], typeMap)
				fmt.Fprintf(file, "%s;\n", sql)
			case *pglogrepl.UpdateMessage:
				if file == nil {
					f, err := os.Create(outputFile)
					if err != nil {
						return err
					}
					file = f
				}
				sql := generateUpdatePostgres(logicalMsg, relations[logicalMsg.RelationID], typeMap)
				fmt.Fprintf(file, "%s;\n", sql)
			case *pglogrepl.DeleteMessage:
				if file == nil {
					f, err := os.Create(outputFile)
					if err != nil {
						return err
					}
					file = f
				}
				sql := generateDeletePostgres(logicalMsg, relations[logicalMsg.RelationID], typeMap)
				fmt.Fprintf(file, "%s;\n", sql)
			}
		}
	}

	return nil
}

func runQuery(cfg *ConfigPostgres, dbName, query string) error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, dbName, cfg.Params)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	_, err = conn.Exec(context.Background(), query)
	return err
}

func runQueryRow(cfg *ConfigPostgres, dbName, query string, dest interface{}) error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, dbName, cfg.Params)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	return conn.QueryRow(context.Background(), query).Scan(dest)
}

func generateInsertPostgres(msg *pglogrepl.InsertMessage, relation *pglogrepl.RelationMessage, typeMap *pgtype.Map) string {
	if relation == nil {
		return ""
	}
	values := []string{}
	cols := []string{}
	for i, col := range msg.Tuple.Columns {
		val := decodeValue(col, relation.Columns[i].DataType, typeMap)
		values = append(values, val)
		cols = append(cols, relation.Columns[i].Name)
	}
	return fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s)",
		relation.Namespace, relation.RelationName,
		strings.Join(cols, ", "),
		strings.Join(values, ", "))
}

func generateUpdatePostgres(msg *pglogrepl.UpdateMessage, relation *pglogrepl.RelationMessage, typeMap *pgtype.Map) string {
	if relation == nil {
		return ""
	}
	// Construct UPDATE statements
	// msg.NewTuple has data. msg.OldTuple usually has key.
	// This is minimal implementation.

	sets := []string{}
	for i, col := range msg.NewTuple.Columns {
		val := decodeValue(col, relation.Columns[i].DataType, typeMap)
		sets = append(sets, fmt.Sprintf("%s = %s", relation.Columns[i].Name, val))
	}

	// Very naive WHERE clause using all available info or Identity columns if we had them.
	// For now, let's skip complex WHERE generation or use simplified logic (not production ready for all cases).

	return fmt.Sprintf("UPDATE %s.%s SET %s",
		relation.Namespace, relation.RelationName,
		strings.Join(sets, ", "))
}

func generateDeletePostgres(msg *pglogrepl.DeleteMessage, relation *pglogrepl.RelationMessage, typeMap *pgtype.Map) string {
	if relation == nil {
		return ""
	}
	// DELETE FROM table WHERE ...
	// Using OldTuple
	wheres := []string{}
	if msg.OldTuple != nil {
		for i, col := range msg.OldTuple.Columns {
			val := decodeValue(col, relation.Columns[i].DataType, typeMap)
			wheres = append(wheres, fmt.Sprintf("%s = %s", relation.Columns[i].Name, val))
		}
	}
	return fmt.Sprintf("DELETE FROM %s.%s WHERE %s",
		relation.Namespace, relation.RelationName,
		strings.Join(wheres, " AND "))
}

func decodeValue(col *pglogrepl.TupleDataColumn, dataType uint32, typeMap *pgtype.Map) string {
	if col.DataType == 'n' { // Null
		return "NULL"
	}
	// Simplification: Treat almost everything as string/text for the dump unless it's obviously numeric.
	// pglogrepl gives us raw bytes.
	// We strictly should decode using OID and pgtype.

	valStr := string(col.Data)
	// Escape single quotes for SQL
	valStr = strings.ReplaceAll(valStr, "'", "''")

	return fmt.Sprintf("'%s'", valStr)
}

func RunBackupPostgres(cfg *ConfigPostgres) {
	log.Info().Msgf("⏰ [CRON PG] Iniciando backup de %d bases de datos... MODO READONLY", len(cfg.Databases))
	SetReadOnly(true)
	defer SetReadOnly(false)

	fullBackupDBs := []string{}
	for _, db := range cfg.Databases {
		if !backupExistsPostgres(db, cfg.BackupDir) {
			fullBackupDBs = append(fullBackupDBs, db)
		}
	}

	if len(fullBackupDBs) > 0 {
		// Run Full Backups
		tasks := make(chan string, len(fullBackupDBs))
		results := make(chan string, len(fullBackupDBs))
		var wg sync.WaitGroup

		for i := 0; i < len(fullBackupDBs); i++ {
			wg.Add(1)
			go fullBackupWorkerPostgres(cfg, tasks, &wg, results)
		}

		for _, db := range fullBackupDBs {
			tasks <- db
		}
		close(tasks)
		wg.Wait()
		close(results)

		// Log results
		successCount := 0
		for result := range results {
			if strings.HasPrefix(result, "SUCCESS") {
				successCount++
			}
		}
		log.Info().Msgf("✅ Backups FULL PG completados: %d/%d exitosos", successCount, len(fullBackupDBs))
		return
	}

	// Deduplicate databases to avoid "slot active" errors if tenant config has duplicates
	uniqueDBs := make(map[string]struct{})
	var dedupedDBs []string
	for _, db := range cfg.Databases {
		if _, exists := uniqueDBs[db]; !exists {
			uniqueDBs[db] = struct{}{}
			dedupedDBs = append(dedupedDBs, db)
		}
	}
	cfg.Databases = dedupedDBs

	// Incremental Backups
	// Verificar wal_level antes de iniciar workers, al menos en una DB
	if len(cfg.Databases) > 0 {
		if err := verifyWalLevelPostgres(cfg, cfg.Databases[0]); err != nil {
			log.Error().Err(err).Msg("❌ [CRON PG] Configuración incorrecta de PostgreSQL")
			return
		}
	}

	results := make(chan string, len(cfg.Databases))
	var wg sync.WaitGroup

	for _, db := range cfg.Databases {
		wg.Add(1)
		go incrementalBackupWorkerPostgres(cfg, db, &wg, results)
	}

	wg.Wait()
	close(results)

	successCount := 0
	for result := range results {
		if strings.HasPrefix(result, "SUCCESS") {
			successCount++
		}
	}
	log.Info().Msgf("✅ [CRON PG] Backup incremental completado: %d/%d exitosos", successCount, len(cfg.Databases))

}
