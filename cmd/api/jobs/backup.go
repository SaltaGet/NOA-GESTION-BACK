package jobs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
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
	Position   int    `json:"position"`
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

func runCommand(name string, args []string, outputPath string) error {
	cmd := exec.Command(name, args...)
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	cmd.Stdout = file
	return cmd.Run()
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
		"--routines", "--events", "--single-transaction"}
	return runCommand("mysqldump", args, path)
}

func getBinlogStatus(cfg *Config) (Checkpoint, error) {
	commands := []string{"mariadb", "mysql"}
	
	var output []byte
	var err error
	var successCmd string
	
	args := []string{
		"-u" + cfg.User,
		"-p" + cfg.Password,
		"-h", cfg.Host,
		"-P", cfg.Port,
		"-e", "SHOW MASTER STATUS\\G",
	}
	
	for _, cmd := range commands {
		output, err = exec.Command(cmd, args...).CombinedOutput()
		if err == nil {
			successCmd = cmd
			break
		}
	}
	
	if err != nil {
		return Checkpoint{}, fmt.Errorf("failed to execute mysql/mariadb: %w, output: %s", err, string(output))
	}
	
	log.Printf("===== SHOW MASTER STATUS OUTPUT =====")
	log.Printf("Comando usado: %s", successCmd)
	log.Printf("Output completo:\n%s", string(output))
	log.Printf("=====================================")
	
	lines := strings.Split(string(output), "\n")
	var cp Checkpoint
	
	for i, l := range lines {
		log.Printf("Línea %d: '%s'", i, l)
		
		if strings.Contains(l, "File:") {
			parts := strings.SplitN(l, ":", 2)
			if len(parts) == 2 {
				cp.BinlogFile = strings.TrimSpace(parts[1])
				log.Printf("✓ File encontrado: '%s'", cp.BinlogFile)
			}
		} else if strings.Contains(l, "Position:") {
			parts := strings.SplitN(l, ":", 2)
			if len(parts) == 2 {
				posStr := strings.TrimSpace(parts[1])
				cp.Position, _ = strconv.Atoi(posStr)
				log.Printf("✓ Position encontrado: %d (string: '%s')", cp.Position, posStr)
			}
		}
	}
	
	log.Printf("===== CHECKPOINT FINAL =====")
	log.Printf("BinlogFile: '%s' (empty=%v)", cp.BinlogFile, cp.BinlogFile == "")
	log.Printf("Position: %d", cp.Position)
	log.Printf("============================")
	
	if cp.BinlogFile == "" || cp.Position == 0 {
		return cp, fmt.Errorf("binlog inválido: File='%s', Position=%d", cp.BinlogFile, cp.Position)
	}
	
	return cp, nil
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

func filterBinlogForDB(input, output, db string) error {
	in, err := os.Open(input)
	if err != nil {
		return err
	}
	defer in.Close()
	
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	scan := bufio.NewScanner(in)
	inSection := false
	lineCount := 0
	matchCount := 0
	
	log.Printf("🔍 Filtrando binlog para DB: '%s'", db)
	
	for scan.Scan() {
		lineCount++
		line := scan.Text()
		lineLower := strings.ToLower(line)
		
		// Detectar inicio de sección para esta DB (más flexible)
		isMatch := false
		
		// 1. Buscar `db`.tabla o `db`
		if strings.Contains(line, fmt.Sprintf("`%s`.", db)) || 
		   strings.Contains(line, fmt.Sprintf("`%s`", db)) {
			isMatch = true
		}
		
		// 2. Buscar "use `db`" (case insensitive)
		if strings.Contains(lineLower, fmt.Sprintf("use `%s`", strings.ToLower(db))) {
			isMatch = true
		}
		
		// 3. Buscar queries que mencionen la DB
		if strings.Contains(lineLower, strings.ToLower(db)) && 
		   (strings.Contains(lineLower, "create") || 
		    strings.Contains(lineLower, "alter") || 
		    strings.Contains(lineLower, "drop") ||
		    strings.Contains(lineLower, "insert") ||
		    strings.Contains(lineLower, "update") ||
		    strings.Contains(lineLower, "delete")) {
			isMatch = true
		}
		
		if isMatch {
			inSection = true
			matchCount++
			if matchCount <= 5 {
				log.Printf("✅ Match %d en línea %d: %s", matchCount, lineCount, line)
			}
		}
		
		if inSection {
			fmt.Fprintln(out, line)
			// Finalizar sección en COMMIT o después de unos comentarios
			if strings.Contains(line, "COMMIT") || 
			   (strings.HasPrefix(line, "# ") && matchCount > 0 && !isMatch) {
				inSection = false
			}
		}
	}
	
	log.Printf("📊 Filtrado completado: %d líneas leídas, %d matches encontrados", lineCount, matchCount)
	
	return scan.Err()
}

func runIncrementalBackup(cfg *Config, db string, cp Checkpoint) (Checkpoint, string, error) {
    ts := time.Now().Format("2006-01-02_15-04-05")
    raw := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_incremental_raw_%s.sql", db, ts))
    final := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_incremental_%s.sql", db, ts))
    
    args := []string{
        "--read-from-remote-server",
        fmt.Sprintf("--host=%s", cfg.Host),
        fmt.Sprintf("--port=%s", cfg.Port),
        fmt.Sprintf("--user=%s", cfg.User),
        fmt.Sprintf("--password=%s", cfg.Password),
        fmt.Sprintf("--start-position=%d", cp.Position),
        cp.BinlogFile,
    }
    
    // Log del comando para debugging
    log.Printf("Ejecutando: mariadb-binlog %s", strings.Join(args, " "))
    
    // NO uses runCommand, hazlo directamente aquí
    cmd := exec.Command("mariadb-binlog", args...)
    
    // Capturar TANTO stdout como stderr
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    err := cmd.Run()
    
    // Imprimir SIEMPRE stderr y stdout para debugging
    if stderr.Len() > 0 {
        log.Printf("STDERR: %s", stderr.String())
    }
    if stdout.Len() > 0 {
        log.Printf("STDOUT length: %d bytes", stdout.Len())
    }
    
    if err != nil {
        return cp, "", fmt.Errorf("mariadb-binlog failed: %w\nSTDERR: %s\nSTDOUT: %s", 
            err, stderr.String(), stdout.String())
    }
    
    // Guardar el output a archivo
    if err := os.WriteFile(raw, stdout.Bytes(), 0644); err != nil {
        return cp, "", fmt.Errorf("failed to write output file: %w", err)
    }
    
    if err := filterBinlogForDB(raw, final, db); err != nil {
        os.Remove(raw)
        return cp, "", err
    }
    
    os.Remove(raw)
    
    newCp, err := getBinlogStatus(cfg)
    return newCp, final, err
}

func hasRealChanges(content string) bool {
	log.Printf("===== ANÁLISIS DE CONTENIDO =====")
	log.Printf("Longitud total: %d bytes", len(content))
	
	if len(content) > 500 {
		log.Printf("Primeros 500 chars:\n%s", content[:500])
	} else {
		log.Printf("Contenido completo:\n%s", content)
	}
	log.Printf("===== FIN ANÁLISIS =====")
	
	// Si el archivo tiene más que solo headers, consideralo válido
	minSize := 200 // Los headers del binlog ocupan ~184 bytes
	if len(content) > minSize {
		log.Printf("✅ Archivo tiene contenido significativo (%d bytes > %d)", len(content), minSize)
		return true
	}
	
	log.Printf("⚠️  Archivo muy pequeño o vacío (%d bytes)", len(content))
	return false
}
func resetCheckpoint(cfg *Config, db string) error {
    cpPath := checkpointPath(db, cfg.BackupDir)
    if err := os.Remove(cpPath); err != nil && !os.IsNotExist(err) {
        return err
    }
    log.Printf("⚠️  Checkpoint eliminado para %s, se hará backup full", db)
    return nil
}

// func RunBackup(cfg *Config) {
//     log.Printf("⏰ [CRON] Iniciando backup de %d bases de datos...", len(cfg.Databases))
    
//     // 1. Verificar que todas las DBs tengan checkpoint
//     needsFullBackup := false
//     for _, db := range cfg.Databases {
//         if !backupExists(db, cfg.BackupDir) {
//             log.Printf("DB '%s' necesita backup full", db)
//             needsFullBackup = true
//         }
//     }
    
//     // 2. Si alguna DB necesita full backup, hacerlo primero
//     if needsFullBackup {
//         for _, db := range cfg.Databases {
//             if !backupExists(db, cfg.BackupDir) {
//                 log.Printf("📦 Ejecutando backup full para: %s", db)
//                 if err := runFullBackup(cfg, db); err != nil {
//                     log.Printf("❌ Error en backup full de %s: %v", db, err)
//                     continue
//                 }
                
//                 cp, err := getBinlogStatus(cfg)
//                 if err != nil {
//                     log.Printf("❌ Error obteniendo binlog status: %v", err)
//                     continue
//                 }
                
//                 log.Printf("✅ Checkpoint inicial para %s: File=%s, Position=%d", db, cp.BinlogFile, cp.Position)
//                 saveCheckpoint(cfg, db, cp)
//             }
//         }
//         return // Esperar al siguiente ciclo para incrementales
//     }
    
//     // 3. Cargar todos los checkpoints
//     checkpoints := make(map[string]Checkpoint)
//     oldestCheckpoint := Checkpoint{Position: int(^uint(0) >> 1)} // Max int
    
//     for _, db := range cfg.Databases {
//         cp, err := loadCheckpoint(cfg, db)
//         if err != nil {
//             log.Printf("❌ Error cargando checkpoint de %s: %v", db, err)
//             continue
//         }
        
//         if cp.BinlogFile == "" {
//             log.Printf("⚠️  Checkpoint inválido para %s, regenerando...", db)
//             resetCheckpoint(cfg, db)
//             continue
//         }
        
//         checkpoints[db] = cp
        
//         // Encontrar el checkpoint más antiguo (menor posición)
//         if cp.Position < oldestCheckpoint.Position {
//             oldestCheckpoint = cp
//         }
//     }
    
//     if len(checkpoints) == 0 {
//         log.Printf("⚠️  No hay checkpoints válidos")
//         return
//     }
    
//     log.Printf("📋 Usando checkpoint más antiguo: File=%s, Position=%d", 
//         oldestCheckpoint.BinlogFile, oldestCheckpoint.Position)
    
//     // 4. Extraer binlog UNA SOLA VEZ desde la posición más antigua
//     ts := time.Now().Format("2006-01-02_15-04-05")
//     rawFile := filepath.Join(cfg.BackupDir, fmt.Sprintf("binlog_raw_%s.sql", ts))
    
//     if err := extractBinlog(cfg, oldestCheckpoint, rawFile); err != nil {
//         log.Printf("❌ Error extrayendo binlog: %v", err)
//         return
//     }
    
//     // 5. Distribuir cambios a cada base de datos
//     dbChanges := distributeBinlogChanges(rawFile, cfg.Databases)
    
//     // 6. Obtener nuevo checkpoint
//     newCp, err := getBinlogStatus(cfg)
//     if err != nil {
//         log.Printf("❌ Error obteniendo nuevo checkpoint: %v", err)
//         return
//     }
    
//     // 7. Guardar archivos de backup para cada DB que tenga cambios
//     for db, changes := range dbChanges {
//         if len(changes) == 0 {
//             log.Printf("ℹ️  %s: Sin cambios", db)
//             continue
//         }
        
//         finalFile := filepath.Join(cfg.BackupDir, 
//             fmt.Sprintf("%s_incremental_%s.sql", db, ts))
        
//         if err := os.WriteFile(finalFile, []byte(strings.Join(changes, "\n")), 0644); err != nil {
//             log.Printf("❌ Error guardando backup de %s: %v", db, err)
//             continue
//         }
        
//         log.Printf("✅ %s: Backup incremental generado (%d líneas, %d bytes)", 
//             db, len(changes), len(strings.Join(changes, "\n")))
        
//         // Actualizar checkpoint
//         saveCheckpoint(cfg, db, newCp)
//     }
    
//     // 8. Limpiar archivo raw
//     os.Remove(rawFile)
    
//     log.Printf("✅ [CRON] Backup completado")
// }
func RunBackup(cfg *Config) {
    log.Printf("⏰ [CRON] Iniciando backup de %d bases de datos...", len(cfg.Databases))
    
    // 1. Verificar que todas las DBs tengan checkpoint
    needsFullBackup := false
    for _, db := range cfg.Databases {
        if !backupExists(db, cfg.BackupDir) {
            log.Printf("DB '%s' necesita backup full", db)
            needsFullBackup = true
        }
    }
    
    // 2. Si alguna DB necesita full backup, hacerlo primero
    if needsFullBackup {
        for _, db := range cfg.Databases {
            if !backupExists(db, cfg.BackupDir) {
                log.Printf("📦 Ejecutando backup full para: %s", db)
                if err := runFullBackup(cfg, db); err != nil {
                    log.Printf("❌ Error en backup full de %s: %v", db, err)
                    continue
                }
                
                cp, err := getBinlogStatus(cfg)
                if err != nil {
                    log.Printf("❌ Error obteniendo binlog status: %v", err)
                    continue
                }
                
                log.Printf("✅ Checkpoint inicial para %s: File=%s, Position=%d", db, cp.BinlogFile, cp.Position)
                saveCheckpoint(cfg, db, cp)
            }
        }
        return
    }
    
    // 3. Obtener checkpoint ACTUAL (antes de extraer binlog)
    currentBinlogStatus, err := getBinlogStatus(cfg)
    if err != nil {
        log.Printf("❌ Error obteniendo binlog status actual: %v", err)
        return
    }
    
    log.Printf("📊 Estado actual del binlog: File=%s, Position=%d", 
        currentBinlogStatus.BinlogFile, currentBinlogStatus.Position)
    
    // 4. Procesar cada DB individualmente con su propio checkpoint
    for _, db := range cfg.Databases {
        cp, err := loadCheckpoint(cfg, db)
        if err != nil {
            log.Printf("❌ Error cargando checkpoint de %s: %v", db, err)
            continue
        }
        
        if cp.BinlogFile == "" {
            log.Printf("⚠️  Checkpoint inválido para %s, regenerando...", db)
            resetCheckpoint(cfg, db)
            continue
        }
        
        // Si no hay cambios nuevos, skip
        if cp.Position >= currentBinlogStatus.Position {
            log.Printf("ℹ️  %s: Sin cambios nuevos (checkpoint: %d >= actual: %d)", 
                db, cp.Position, currentBinlogStatus.Position)
            continue
        }
        
        log.Printf("📋 %s: Procesando desde position %d hasta %d (%d bytes)", 
            db, cp.Position, currentBinlogStatus.Position, 
            currentBinlogStatus.Position - cp.Position)
        
        // Extraer binlog para esta DB específica
        ts := time.Now().Format("2006-01-02_15-04-05")
        rawFile := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_raw_%s.sql", db, ts))
        
        if err := extractBinlog(cfg, cp, rawFile); err != nil {
            log.Printf("❌ Error extrayendo binlog para %s: %v", db, err)
            continue
        }
        
        // Filtrar solo los cambios de esta DB
        changes := filterBinlogForDatabase(rawFile, db)
        
        // Limpiar archivo raw
        os.Remove(rawFile)
        
        if len(changes) == 0 {
            log.Printf("ℹ️  %s: Sin cambios relevantes", db)
            // Actualizar checkpoint aunque no haya cambios (para avanzar la posición)
            saveCheckpoint(cfg, db, currentBinlogStatus)
            continue
        }
        
        // Guardar backup incremental
        finalFile := filepath.Join(cfg.BackupDir, 
            fmt.Sprintf("%s_incremental_%s.sql", db, ts))
        
        content := strings.Join(changes, "\n")
        if err := os.WriteFile(finalFile, []byte(content), 0644); err != nil {
            log.Printf("❌ Error guardando backup de %s: %v", db, err)
            continue
        }
        
        log.Printf("✅ %s: Backup incremental generado (%d líneas, %d bytes)", 
            db, len(changes), len(content))
        
        // Actualizar checkpoint con la posición ACTUAL
        saveCheckpoint(cfg, db, currentBinlogStatus)
    }
    
    log.Printf("✅ [CRON] Backup completado")
}

// Filtrar cambios solo para una base de datos específica
func filterBinlogForDatabase(binlogFile string, targetDB string) []string {
    file, err := os.Open(binlogFile)
    if err != nil {
        log.Printf("❌ Error abriendo binlog: %v", err)
        return nil
    }
    defer file.Close()
    
    var changes []string
    var currentSection []string
    inTargetDB := false
    
    scanner := bufio.NewScanner(file)
    targetDBLower := strings.ToLower(targetDB)
    
    for scanner.Scan() {
        line := scanner.Text()
        lineLower := strings.ToLower(line)
        
        // Detectar cuando cambiamos a nuestra DB
        if strings.Contains(lineLower, fmt.Sprintf("use `%s`", targetDBLower)) {
            // Guardar sección anterior si no era nuestra DB
            if !inTargetDB && len(currentSection) > 0 {
                currentSection = []string{}
            }
            inTargetDB = true
            currentSection = append(currentSection, line)
            continue
        }
        
        // Detectar cuando cambiamos a otra DB
        if strings.HasPrefix(lineLower, "use `") && !strings.Contains(lineLower, targetDBLower) {
            // Guardar sección de nuestra DB
            if inTargetDB && len(currentSection) > 0 {
                changes = append(changes, currentSection...)
                currentSection = []string{}
            }
            inTargetDB = false
            continue
        }
        
        // Si estamos en nuestra DB, agregar línea
        if inTargetDB {
            currentSection = append(currentSection, line)
            
            // Si encontramos COMMIT, cerrar y guardar sección
            if strings.Contains(line, "COMMIT") {
                changes = append(changes, currentSection...)
                currentSection = []string{}
            }
        }
    }
    
    // Guardar última sección si existe
    if inTargetDB && len(currentSection) > 0 {
        changes = append(changes, currentSection...)
    }
    
    log.Printf("📋 %s: %d líneas filtradas", targetDB, len(changes))
    
    return changes
}

// Extraer binlog desde una posición
func extractBinlog(cfg *Config, cp Checkpoint, outputFile string) error {
    args := []string{
        "--read-from-remote-server",
        fmt.Sprintf("--host=%s", cfg.Host),
        fmt.Sprintf("--port=%s", cfg.Port),
        fmt.Sprintf("--user=%s", cfg.User),
        fmt.Sprintf("--password=%s", cfg.Password),
        fmt.Sprintf("--start-position=%d", cp.Position),
        cp.BinlogFile,
    }
    
    log.Printf("🔄 Extrayendo binlog desde posición %d...", cp.Position)
    
    cmd := exec.Command("mariadb-binlog", args...)
    
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("mariadb-binlog failed: %w\nSTDERR: %s", err, stderr.String())
    }
    
    log.Printf("📊 Binlog extraído: %d bytes", stdout.Len())
    
    return os.WriteFile(outputFile, stdout.Bytes(), 0644)
}

// Distribuir cambios del binlog por base de datos
func distributeBinlogChanges(binlogFile string, databases []string) map[string][]string {
    file, err := os.Open(binlogFile)
    if err != nil {
        log.Printf("❌ Error abriendo binlog: %v", err)
        return nil
    }
    defer file.Close()
    
    // Mapa para guardar líneas por DB
    dbChanges := make(map[string][]string)
    for _, db := range databases {
        dbChanges[db] = []string{}
    }
    
    scanner := bufio.NewScanner(file)
    var currentDB string
    var currentSection []string
    
    for scanner.Scan() {
        line := scanner.Text()
        lineLower := strings.ToLower(line)
        
        // Detectar cambio de base de datos
        for _, db := range databases {
            dbLower := strings.ToLower(db)
            
            // Patrón: use `database`
            if strings.Contains(lineLower, fmt.Sprintf("use `%s`", dbLower)) {
                // Guardar sección anterior si existe
                if currentDB != "" && len(currentSection) > 0 {
                    dbChanges[currentDB] = append(dbChanges[currentDB], currentSection...)
                }
                currentDB = db
                currentSection = []string{line}
                break
            }
        }
        
        // Agregar línea a la sección actual
        if currentDB != "" {
            currentSection = append(currentSection, line)
            
            // Si encontramos COMMIT, cerrar sección
            if strings.Contains(line, "COMMIT") {
                dbChanges[currentDB] = append(dbChanges[currentDB], currentSection...)
                currentSection = []string{}
                currentDB = ""
            }
        }
    }
    
    // Guardar última sección si existe
    if currentDB != "" && len(currentSection) > 0 {
        dbChanges[currentDB] = append(dbChanges[currentDB], currentSection...)
    }
    
    // Log de resumen
    for db, changes := range dbChanges {
        if len(changes) > 0 {
            log.Printf("📋 %s: %d líneas de cambios detectadas", db, len(changes))
        }
    }
    
    return dbChanges
}
