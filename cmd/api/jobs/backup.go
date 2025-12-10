package jobs

// import (
// 	"bytes"
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
// 	Position   int    `json:"position"`
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

// func runCommand(name string, args []string, outputPath string) error {
// 	cmd := exec.Command(name, args...)
// 	file, err := os.Create(outputPath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()
// 	cmd.Stdout = file
// 	return cmd.Run()
// }

// func runFullBackup(cfg *Config, db string) error {
// 	ts := time.Now().Format("2006-01-02_15-04-05")
// 	path := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_full_%s.sql", db, ts))
// 	args := []string{
// 		"-u", cfg.User,
// 		"-p" + cfg.Password,
// 		"-h", cfg.Host,
// 		"-P", cfg.Port,
// 		"--databases", db,
// 		"--routines", "--events", "--single-transaction"}
// 	return runCommand("mariadb-dump", args, path)
// }

// func getBinlogStatus(cfg *Config) (Checkpoint, error) {
// 	commands := []string{"mariadb", "mysql"}
	
// 	var output []byte
// 	var err error
// 	var successCmd string
	
// 	args := []string{
// 		"-u" + cfg.User,
// 		"-p" + cfg.Password,
// 		"-h", cfg.Host,
// 		"-P", cfg.Port,
// 		"-e", "SHOW MASTER STATUS\\G",
// 	}
	
// 	for _, cmd := range commands {
// 		output, err = exec.Command(cmd, args...).CombinedOutput()
// 		if err == nil {
// 			successCmd = cmd
// 			break
// 		}
// 	}
	
// 	if err != nil {
// 		return Checkpoint{}, fmt.Errorf("failed to execute mysql/mariadb: %w, output: %s", err, string(output))
// 	}
	
// 	log.Printf("===== SHOW MASTER STATUS OUTPUT =====")
// 	log.Printf("Comando usado: %s", successCmd)
// 	log.Printf("Output completo:\n%s", string(output))
// 	log.Printf("=====================================")
	
// 	lines := strings.Split(string(output), "\n")
// 	var cp Checkpoint
	
// 	for i, l := range lines {
// 		log.Printf("Línea %d: '%s'", i, l)
		
// 		if strings.Contains(l, "File:") {
// 			parts := strings.SplitN(l, ":", 2)
// 			if len(parts) == 2 {
// 				cp.BinlogFile = strings.TrimSpace(parts[1])
// 				log.Printf("✓ File encontrado: '%s'", cp.BinlogFile)
// 			}
// 		} else if strings.Contains(l, "Position:") {
// 			parts := strings.SplitN(l, ":", 2)
// 			if len(parts) == 2 {
// 				posStr := strings.TrimSpace(parts[1])
// 				cp.Position, _ = strconv.Atoi(posStr)
// 				log.Printf("✓ Position encontrado: %d (string: '%s')", cp.Position, posStr)
// 			}
// 		}
// 	}
	
// 	log.Printf("===== CHECKPOINT FINAL =====")
// 	log.Printf("BinlogFile: '%s' (empty=%v)", cp.BinlogFile, cp.BinlogFile == "")
// 	log.Printf("Position: %d", cp.Position)
// 	log.Printf("============================")
	
// 	if cp.BinlogFile == "" || cp.Position == 0 {
// 		return cp, fmt.Errorf("binlog inválido: File='%s', Position=%d", cp.BinlogFile, cp.Position)
// 	}
	
// 	return cp, nil
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

// func RunBackup(cfg *Config) {
// 	log.Printf("⏰ [CRON] Iniciando backup de %d bases de datos... MODO READONLY", len(cfg.Databases))
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
		
// 		log.Printf("📋 %s: Procesando desde position %d hasta %d (%d bytes)", 
// 			db, cp.Position, currentBinlogStatus.Position, 
// 			currentBinlogStatus.Position - cp.Position)
		
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
		
// 		if !hasExecutableSQL(string(content)) {
// 			log.Printf("ℹ️  %s: Sin cambios SQL ejecutables", db)
// 			os.Remove(finalFile) // Eliminar archivo vacío
// 			saveCheckpoint(cfg, db, currentBinlogStatus)
// 			continue
// 		}
		
// 		log.Printf("✅ %s: Backup incremental generado (%d bytes)", db, len(content))
		
// 		saveCheckpoint(cfg, db, currentBinlogStatus)
// 	}
	
// 	log.Printf("✅ [CRON] Backup completado")
// }

// func extractBinlogForDatabase(cfg *Config, cp Checkpoint, database string, outputFile string) error {
// 	args := []string{
// 		"--read-from-remote-server",
// 		fmt.Sprintf("--host=%s", cfg.Host),
// 		fmt.Sprintf("--port=%s", cfg.Port),
// 		fmt.Sprintf("--user=%s", cfg.User),
// 		fmt.Sprintf("--password=%s", cfg.Password),
// 		fmt.Sprintf("--start-position=%d", cp.Position),
// 		fmt.Sprintf("--database=%s", database),
// 		cp.BinlogFile,
// 	}
	
// 	log.Printf("🔄 [%s] Extrayendo binlog desde posición %d...", database, cp.Position)
	
// 	cmd := exec.Command("mariadb-binlog", args...)
	
// 	var stdout, stderr bytes.Buffer
// 	cmd.Stdout = &stdout
// 	cmd.Stderr = &stderr
	
// 	if err := cmd.Run(); err != nil {
// 		return fmt.Errorf("mariadb-binlog failed: %w\nSTDERR: %s", err, stderr.String())
// 	}
	
// 	log.Printf("📊 [%s] Binlog extraído: %d bytes", database, stdout.Len())
	
// 	lines := strings.Split(stdout.String(), "\n")
// 	log.Printf("🔍 [%s] Primeras 40 líneas del binlog:", database)
// 	for i, line := range lines {
// 		if i >= 40 {
// 			break
// 		}
// 		log.Printf("  L%d: %s", i+1, line)
// 	}
	
// 	return os.WriteFile(outputFile, stdout.Bytes(), 0644)
// }

// func hasExecutableSQL(content string) bool {
// 	lines := strings.Split(content, "\n")
	
// 	sqlKeywords := []string{
// 		"INSERT INTO",
// 		"UPDATE ",
// 		"DELETE FROM",
// 		"CREATE TABLE",
// 		"ALTER TABLE",
// 		"DROP TABLE",
// 		"TRUNCATE",
// 		"REPLACE INTO",
// 	}
	
// 	for _, line := range lines {
// 		trimmed := strings.TrimSpace(line)
		
// 		if trimmed == "" ||
// 		   strings.HasPrefix(trimmed, "#") ||
// 		   strings.HasPrefix(trimmed, "/*") ||
// 		   strings.HasPrefix(trimmed, "--") ||
// 		   strings.HasPrefix(trimmed, "SET ") ||
// 		   strings.HasPrefix(trimmed, "DELIMITER") ||
// 		   strings.HasPrefix(trimmed, "START TRANSACTION") ||
// 		   strings.HasPrefix(trimmed, "COMMIT") ||
// 		   strings.HasPrefix(trimmed, "ROLLBACK") ||
// 		   strings.HasPrefix(trimmed, "use `") {
// 			continue
// 		}
		
// 		upperLine := strings.ToUpper(trimmed)
// 		for _, keyword := range sqlKeywords {
// 			if strings.Contains(upperLine, keyword) {
// 				log.Printf("✓ SQL ejecutable detectado: %s", trimmed[:min(len(trimmed), 100)])
// 				return true
// 			}
// 		}
// 	}
	
// 	return false
// }

// func min(a, b int) int {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }


