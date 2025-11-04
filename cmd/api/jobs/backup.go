package jobs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

	connections, err := deps.TenantController.TenantService.TenantGetConections()
	if err != nil {
		return nil, err
	}

	for _, conn := range *connections {
		dbName, err := extractDBName(conn)
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
	args := []string{"-u", cfg.User, "-p" + cfg.Password, "-h", cfg.Host, "-P", cfg.Port, "-e", "SHOW MASTER STATUS\\G"}
	output, err := exec.Command("mysql", args...).Output()
	if err != nil {
		return Checkpoint{}, err
	}

	lines := strings.Split(string(output), "\n")
	var cp Checkpoint
	for _, l := range lines {
		if strings.Contains(l, "File:") {
			cp.BinlogFile = strings.TrimSpace(strings.SplitN(l, ":", 2)[1])
		} else if strings.Contains(l, "Position:") {
			cp.Position, _ = strconv.Atoi(strings.TrimSpace(strings.SplitN(l, ":", 2)[1]))
		}
	}
	if cp.BinlogFile == "" || cp.Position == 0 {
		return cp, fmt.Errorf("binlog inválido")
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
	in, _ := os.Open(input)
	defer in.Close()
	out, _ := os.Create(output)
	defer out.Close()

	scan := bufio.NewScanner(in)
	inSection := false
	dbRegex := regexp.MustCompile(fmt.Sprintf("`%s`\\.", db))
	for scan.Scan() {
		line := scan.Text()
		if dbRegex.MatchString(line) || strings.Contains(line, fmt.Sprintf("USE `%s`", db)) {
			inSection = true
		}
		if inSection {
			fmt.Fprintln(out, line)
			if strings.Contains(line, "COMMIT") {
				inSection = false
			}
		}
	}
	return scan.Err()
}

func runIncrementalBackup(cfg *Config, db string, cp Checkpoint) (Checkpoint, string, error) {
	ts := time.Now().Format("2006-01-02_15-04-05")
	raw := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_incremental_raw_%s.sql", db, ts))
	final := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_incremental_%s.sql", db, ts))

	args := []string{
		"--read-from-remote-server", 
		"--host=" + cfg.Host, 
		"--port=" + cfg.Port,
		"--user=" + cfg.User, 
		"--password=" + cfg.Password,
		"--database=" + db,
		"--start-position", fmt.Sprintf("%d", cp.Position), 
		cp.BinlogFile}

	if err := runCommand("mysqlbinlog", args, raw); err != nil {
		os.Remove(raw)
		return cp, "", err
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
	log.Println(content)
	keywords := []string{
		"Write_rows",   
		"Update_rows",
		"Delete_rows", 
		"Create_table",
		"ALTER TABLE",
		"DROP TABLE",
		"CREATE TABLE",
		"Query: CREATE", 
		"Query: DROP",
		"Query: ALTER",
	}
	for _, k := range keywords {
		if strings.Contains(content, k) {
			return true
		}
	}
	return false
}

func RunBackup(cfg *Config) {
	for _, db := range cfg.Databases {
		fmt.Println("Procesando DB:", db)
		if !backupExists(db, cfg.BackupDir) {
			runFullBackup(cfg, db)
			cp, _ := getBinlogStatus(cfg)
			saveCheckpoint(cfg, db, cp)
			continue
		}
		cp, err := loadCheckpoint(cfg, db)
		if err != nil {
			fmt.Println("Error cargando checkpoint:", err)
			continue
		}
		newCp, file, err := runIncrementalBackup(cfg, db, cp)
		if err != nil {
			fmt.Println("Error backup incremental:", err)
			continue
		}
		content, _ := os.ReadFile(file)
		if hasRealChanges(string(content)) {
			saveCheckpoint(cfg, db, newCp)
			fmt.Println("Backup incremental exitoso")
		} else {
			os.Remove(file)
			fmt.Println("No hay cambios relevantes, backup descartado")
		}
	}
}

// import (
// 	"bufio"
// 	// "bytes"
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"regexp"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/DanielChachagua/GestionCar/pkg/dependencies"
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

// func LoadConfig(deps *dependencies.Application) (*Config, error) {
// 	mainDB := os.Getenv("URI_DB")
// 	conf, err := parseDSN(mainDB)
// 	if err != nil {
// 		return nil, err
// 	}

// 	connections, err := deps.TenantController.TenantService.TenantGetConections()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, connection := range *connections {
// 		conn, err := extractDBName(connection)
// 		if err != nil {
// 			return nil, err
// 		}
// 		conf.Databases = append(conf.Databases, conn)
// 	}

// 	return &conf, nil
// }

// func parseDSN(dsn string) (Config, error) {
// 	var cfg Config

// 	// Separar la parte user:password y el resto
// 	parts := strings.SplitN(dsn, "@tcp(", 2)
// 	if len(parts) != 2 {
// 		return cfg, fmt.Errorf("DSN inválido: falta '@tcp('")
// 	}

// 	userPass := parts[0]    // root:Qwer1234*
// 	hostAndRest := parts[1] // 127.0.0.1:3306)/gestion_car?charset=...

// 	// Parsear user y password
// 	up := strings.SplitN(userPass, ":", 2)
// 	if len(up) != 2 {
// 		return cfg, fmt.Errorf("DSN inválido: falta ':' entre usuario y password")
// 	}
// 	cfg.User = up[0]
// 	cfg.Password = up[1]

// 	// Separar host y base de datos
// 	// hostAndRest tiene la forma: 127.0.0.1:3306)/gestion_car?...
// 	idx := strings.Index(hostAndRest, ")/")
// 	if idx == -1 {
// 		return cfg, fmt.Errorf("DSN inválido: falta ')/'")
// 	}

// 	hostPort := hostAndRest[:idx] // Esto será "127.0.0.1:3306"

// 	// Separar host y port
// 	hostParts := strings.SplitN(hostPort, ":", 2)
// 	if len(hostParts) != 2 {
// 		return cfg, fmt.Errorf("DSN inválido: falta ':' en host:port")
// 	}
// 	cfg.Host = hostParts[0] // Ahora solo "127.0.0.1"
// 	cfg.Port = hostParts[1] // Y "3306"

// 	dbAndParams := hostAndRest[idx+2:]

// 	// Separar base de datos y parámetros (opcional)
// 	dbName := dbAndParams
// 	if i := strings.Index(dbAndParams, "?"); i != -1 {
// 		dbName = dbAndParams[:i]
// 	}
// 	cfg.Databases = []string{dbName}

// 	return cfg, nil
// }

// func extractDBName(dsn string) (string, error) {
// 	// Separar la parte antes del '?'
// 	beforeParams := strings.SplitN(dsn, "?", 2)[0]

// 	// Separar la parte antes del '/' (lo que está después del host y puerto)
// 	parts := strings.Split(beforeParams, "/")
// 	if len(parts) < 2 {
// 		return "", fmt.Errorf("DSN malformado: no se pudo encontrar la base de datos")
// 	}
// 	dbName := parts[len(parts)-1]
// 	return dbName, nil
// }

// func checkpointPath(db string, dir string) string {
// 	return filepath.Join(dir, fmt.Sprintf("%s_checkpoint.json", db))
// }

// func backupExists(db string, dir string) bool {
// 	_, err := os.Stat(checkpointPath(db, dir))
// 	return err == nil
// }

// func runFullBackup(cfg *Config, db string) error {
// 	timestamp := time.Now().Format("2006-01-02_15-04-05")
// 	backupFile := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_full_%s.sql", db, timestamp))

// 	// cmd := exec.Command("mysqldump",
// 	// 	"-u", cfg.User,
// 	// 	"-p"+cfg.Password,
// 	// 	"-h", cfg.Host,
// 	// 	"--databases", db,
// 	// 	"--routines", "--events", "--single-transaction")

// 	cmdArgs := []string{
// 		"-u", cfg.User,
// 		"-p" + cfg.Password, // Atención: la contraseña se imprimirá.
// 		"-h", cfg.Host,
// 		"-P", cfg.Port,
// 		"--databases", db,
// 		"--routines", "--events", "--single-transaction",
// 	}

// 	cmd := exec.Command("mysqldump", cmdArgs...)

// 	fmt.Println("Comando mysqldump a ejecutar:", "mysqldump", strings.Join(cmdArgs, " "))

// 	outputFile, err := os.Create(backupFile)
// 	if err != nil {
// 		return err
// 	}
// 	defer outputFile.Close()
// 	cmd.Stdout = outputFile

// 	return cmd.Run()
// }

// func getBinlogStatus(cfg *Config) (Checkpoint, error) {
// 	cmd := exec.Command("mysql", "-u", cfg.User, "-p"+cfg.Password, "-h", cfg.Host, "-P", cfg.Port, "-e", "SHOW MASTER STATUS\\G")
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return Checkpoint{}, err
// 	}

// 	var cp Checkpoint
// 	lines := strings.Split(string(output), "\n")
// 	for _, line := range lines {
// 		if strings.Contains(line, "File:") { // Usar Contains es más flexible que HasPrefix por los espacios iniciales
// 			// Eliminar cualquier espacio al inicio y luego buscar "File:"
// 			trimmedLine := strings.TrimSpace(line)
// 			if strings.HasPrefix(trimmedLine, "File:") {
// 				parts := strings.SplitN(trimmedLine, ":", 2)
// 				if len(parts) == 2 {
// 					cp.BinlogFile = strings.TrimSpace(parts[1])
// 				}
// 			}
// 		} else if strings.Contains(line, "Position:") { // Usar Contains
// 			// Eliminar cualquier espacio al inicio y luego buscar "Position:"
// 			trimmedLine := strings.TrimSpace(line)
// 			if strings.HasPrefix(trimmedLine, "Position:") {
// 				parts := strings.SplitN(trimmedLine, ":", 2)
// 				if len(parts) == 2 {
// 					posStr := strings.TrimSpace(parts[1])
// 					parsedPos, err := strconv.Atoi(posStr)
// 					if err != nil {
// 						fmt.Printf("Error convirtiendo Position a int: %v en '%s'\n", err, posStr)
// 						continue // Salta al siguiente iteración del bucle
// 					}
// 					cp.Position = parsedPos
// 				}
// 			}
// 		}
// 	}

// 	if cp.BinlogFile == "" || cp.Position == 0 {
// 		return cp, fmt.Errorf("estado de binlog inválido")
// 	}
// 	return cp, nil
// }

// func saveCheckpoint(cfg *Config, db string, cp Checkpoint) error {
// 	data, err := json.MarshalIndent(cp, "", "  ")
// 	if err != nil {
// 		return err
// 	}
// 	return os.WriteFile(checkpointPath(db, cfg.BackupDir), data, 0644)
// }

// func loadCheckpoint(cfg *Config, db string) (Checkpoint, error) {
// 	data, err := os.ReadFile(checkpointPath(db, cfg.BackupDir))
// 	if err != nil {
// 		return Checkpoint{}, err
// 	}
// 	var cp Checkpoint
// 	err = json.Unmarshal(data, &cp)
// 	return cp, err
// }

// func runIncrementalBackup(cfg *Config, db string, cp Checkpoint) (Checkpoint, string, error) {
// 	timestamp := time.Now().Format("2006-01-02_15-04-05")
// 	rawFile := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_incremental_raw_%s.sql", db, timestamp))
// 	finalFile := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_incremental_%s.sql", db, timestamp))

// 	// 1. Obtener binlog completo sin filtrar por DB
// 	cmd := exec.Command("mysqlbinlog",
// 		"--read-from-remote-server",
// 		"--host="+cfg.Host,
// 		"--port="+cfg.Port,
// 		"--user="+cfg.User,
// 		"--password="+cfg.Password,
// 		"--skip-empty-transactions",
// 		"--start-position", fmt.Sprintf("%d", cp.Position),
// 		cp.BinlogFile)

// 	rawOut, err := os.Create(rawFile)
// 	if err != nil {
// 		return cp, "", err
// 	}
// 	defer rawOut.Close()
// 	cmd.Stdout = rawOut

// 	if err := cmd.Run(); err != nil {
// 		os.Remove(rawFile)
// 		return cp, "", fmt.Errorf("error ejecutando mysqlbinlog: %v", err)
// 	}

// 	// 2. Filtrar solo eventos relevantes para esta DB
// 	if err := filterBinlogForDB(rawFile, finalFile, db); err != nil {
// 		os.Remove(rawFile)
// 		return cp, "", fmt.Errorf("error filtrando binlog: %v", err)
// 	}
// 	os.Remove(rawFile)

// 	newCp, err := getBinlogStatus(cfg)
//     if err != nil {
//         return Checkpoint{}, finalFile, err
//     }
//     return newCp, finalFile, nil
// }

// func filterBinlogForDB(inputFile, outputFile, db string) error {
// 	inFile, err := os.Open(inputFile)
// 	if err != nil {
// 		return err
// 	}
// 	defer inFile.Close()

// 	outFile, err := os.Create(outputFile)
// 	if err != nil {
// 		return err
// 	}
// 	defer outFile.Close()

// 	scanner := bufio.NewScanner(inFile)
// 	inRelevantSection := false
// 	dbPattern := regexp.MustCompile(fmt.Sprintf("`%s`\\.", db))

// 	for scanner.Scan() {
// 		line := scanner.Text()

// 		// Buscar inicio de evento relevante
// 		if strings.Contains(line, fmt.Sprintf("### INSERT INTO `%s`", db)) ||
// 			strings.Contains(line, fmt.Sprintf("### UPDATE `%s`", db)) ||
// 			strings.Contains(line, fmt.Sprintf("### DELETE FROM `%s`", db)) ||
// 			strings.Contains(line, fmt.Sprintf("CREATE TABLE `%s`", db)) ||
// 			strings.Contains(line, fmt.Sprintf("ALTER TABLE `%s`", db)) ||
// 			strings.Contains(line, fmt.Sprintf("USE `%s`", db)) ||
// 			dbPattern.MatchString(line) {

// 			inRelevantSection = true
// 		}

// 		// Mantener encabezados importantes
// 		if strings.HasPrefix(line, "/*!") ||
// 			strings.HasPrefix(line, "SET ") ||
// 			strings.HasPrefix(line, "BINLOG '") {
// 			outFile.WriteString(line + "\n")
// 		}

// 		// Escribir secciones relevantes
// 		if inRelevantSection {
// 			outFile.WriteString(line + "\n")

// 			// Fin de sección
// 			if strings.Contains(line, "COMMIT") ||
// 				strings.Contains(line, "/*!*/;") {
// 				inRelevantSection = false
// 			}
// 		}
// 	}

// 	return scanner.Err()
// }

// func hasRealChanges(content string) bool {
// 	patterns := []string{
// 		"INSERT INTO",
// 		"UPDATE ",
// 		"DELETE FROM",
// 		"CREATE TABLE",
// 		"ALTER TABLE",
// 		"DROP TABLE",
// 	}
// 	for _, pattern := range patterns {
// 		if strings.Contains(content, pattern) {
// 			return true
// 		}
// 	}
// 	return false
// }

// func getBinlogFiles(cfg *Config) ([]string, error) {
// 	cmd := exec.Command("mysql", "-u", cfg.User, "-p"+cfg.Password, "-h", cfg.Host, "-P", cfg.Port, "-e", "SHOW BINARY LOGS")
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var files []string
// 	lines := strings.Split(string(output), "\n")
// 	for i, line := range lines {
// 		if i == 0 { // Saltar cabecera
// 			continue
// 		}
// 		if parts := strings.Fields(line); len(parts) > 0 {
// 			files = append(files, parts[0])
// 		}
// 	}
// 	return files, nil
// }

// func binlogExists(cfg *Config, binlogFile string) bool {
// 	cmd := exec.Command("mysql", "-u", cfg.User, "-p"+cfg.Password, "-h", cfg.Host, "-P", cfg.Port, "-e", "SHOW BINARY LOGS")
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return false
// 	}
// 	return strings.Contains(string(output), binlogFile)
// }

// func RunBackup(cfg *Config) {
// 	for _, db := range cfg.Databases {
// 		fmt.Println("Procesando DB:", db)

// 		if !backupExists(db, cfg.BackupDir) {
// 			fmt.Println("  No existe backup completo. Haciendo uno...")
// 			if err := runFullBackup(cfg, db); err != nil {
// 				fmt.Println("  ❌ Error backup full:", err)
// 				continue
// 			}
// 			cp, err := getBinlogStatus(cfg)
// 			if err != nil {
// 				fmt.Println("  ❌ Error obteniendo binlog:", err)
// 				continue
// 			}
// 			if err := saveCheckpoint(cfg, db, cp); err != nil {
// 				fmt.Println("  ❌ Error guardando checkpoint:", err)
// 			}
// 			fmt.Println("  ✅ Backup full realizado.")
// 		} else {
// 			cp, err := loadCheckpoint(cfg, db)
// 			if err != nil {
// 				fmt.Printf("  ❌ Error leyendo checkpoint para %s: %v\n", db, err)
// 				continue
// 			}

// 			if !binlogExists(cfg, cp.BinlogFile) {
// 				fmt.Println("⚠️ Binlog purgado, haciendo backup completo")
// 				if err := runFullBackup(cfg, db); err != nil {
// 					fmt.Println("❌ Error full backup:", err)
// 				}
// 				continue
// 			}

// 			currentCp, err := getBinlogStatus(cfg)
// 			if err != nil {
// 				fmt.Printf("  ❌ Error obteniendo estado binlog: %v\n", err)
// 				continue
// 			}

// 			// Verifica si hay cambios
// 			if cp.BinlogFile == currentCp.BinlogFile && cp.Position == currentCp.Position {
// 				fmt.Printf("  ✅ No hay cambios para %s (binlog=%s, pos=%d)\n",
// 					db, cp.BinlogFile, cp.Position)
// 				continue
// 			}

// 			// Verifica si el binlog original todavía existe
// 			availableFiles, err := getBinlogFiles(cfg)
// 			if err != nil {
// 				fmt.Printf("  ❌ Error obteniendo binlogs: %v\n", err)
// 				continue
// 			}

// 			found := false
// 			for _, file := range availableFiles {
// 				if file == cp.BinlogFile {
// 					found = true
// 					break
// 				}
// 			}

// 			if !found {
// 				fmt.Println("  🔄 Binlog original purgado, haciendo backup completo")
// 				if err := runFullBackup(cfg, db); err != nil {
// 					fmt.Println("  ❌ Error backup full:", err)
// 					continue
// 				}
// 			} else {
// 				newCp, finalFilePath, err := runIncrementalBackup(cfg, db, cp) // MODIFICADO: ahora devuelve filePath
// 				if err != nil {
// 					fmt.Printf("  ❌ Error en backup incremental para %s: %v\n", db, err)
// 					continue
// 				}

// 				// 1. VERIFICAR SI EL INCREMENTAL TIENE CAMBIOS REALES
// 				content, err := os.ReadFile(finalFilePath)
// 				if err != nil {
// 					fmt.Printf("  ❌ Error leyendo incremental para %s: %v\n", db, err)
// 					continue
// 				}

// 				if hasRealChanges(string(content)) {
// 					// 2. GUARDAR CHECKPOINT SOLO SI HAY CAMBIOS REALES
// 					if err := saveCheckpoint(cfg, db, newCp); err != nil {
// 						fmt.Printf("  ❌ Error guardando checkpoint para %s: %v\n", db, err)
// 					} else {
// 						fmt.Printf("  ✅ Backup incremental válido para %s (nueva posición: %d)\n",
// 							db, newCp.Position)
// 					}
// 				} else {
// 					// 3. ELIMINAR INCREMENTAL VACÍO Y NO ACTUALIZAR CHECKPOINT
// 					os.Remove(finalFilePath)
// 					fmt.Printf("  ⚠️ Cambios irrelevantes para %s, omitiendo incremental\n", db)
// 					// Mantenemos el checkpoint anterior para la próxima ejecución
// 				}
// 			}
// 		}
// 	}
// }



















////////// FUNCIONAL A MEDIAS
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

// func LoadConfig(deps *dependencies.Application) (*Config, error) {
// 	mainDB := os.Getenv("URI_DB")
// 	conf, err := parseDSN(mainDB)
// 	if err != nil {
// 		return nil, err
// 	}

// 	connections, err := deps.TenantController.TenantService.TenantGetConections()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, connection := range *connections {
// 		conn, err := extractDBName(connection)
// 		if err != nil {
// 			return nil, err
// 		}
// 		conf.Databases = append(conf.Databases, conn)
// 	}

// 	return &conf, nil
// }

// func parseDSN(dsn string) (Config, error) {
// 	var cfg Config

// 	// Separar la parte user:password y el resto
// 	parts := strings.SplitN(dsn, "@tcp(", 2)
// 	if len(parts) != 2 {
// 		return cfg, fmt.Errorf("DSN inválido: falta '@tcp('")
// 	}

// 	userPass := parts[0]    // root:Qwer1234*
// 	hostAndRest := parts[1] // 127.0.0.1:3306)/gestion_car?charset=...

// 	// Parsear user y password
// 	up := strings.SplitN(userPass, ":", 2)
// 	if len(up) != 2 {
// 		return cfg, fmt.Errorf("DSN inválido: falta ':' entre usuario y password")
// 	}
// 	cfg.User = up[0]
// 	cfg.Password = up[1]

// 	// Separar host y base de datos
// 	// hostAndRest tiene la forma: 127.0.0.1:3306)/gestion_car?...
// 	idx := strings.Index(hostAndRest, ")/")
// 	if idx == -1 {
// 		return cfg, fmt.Errorf("DSN inválido: falta ')/'")
// 	}

// 	hostPort := hostAndRest[:idx] // Esto será "127.0.0.1:3306"

// 	// Separar host y port
// 	hostParts := strings.SplitN(hostPort, ":", 2)
// 	if len(hostParts) != 2 {
// 		return cfg, fmt.Errorf("DSN inválido: falta ':' en host:port")
// 	}
// 	cfg.Host = hostParts[0] // Ahora solo "127.0.0.1"
// 	cfg.Port = hostParts[1] // Y "3306"

// 	dbAndParams := hostAndRest[idx+2:]

// 	// Separar base de datos y parámetros (opcional)
// 	dbName := dbAndParams
// 	if i := strings.Index(dbAndParams, "?"); i != -1 {
// 		dbName = dbAndParams[:i]
// 	}
// 	cfg.Databases = []string{dbName}

// 	return cfg, nil
// }

// func extractDBName(dsn string) (string, error) {
// 	// Separar la parte antes del '?'
// 	beforeParams := strings.SplitN(dsn, "?", 2)[0]

// 	// Separar la parte antes del '/' (lo que está después del host y puerto)
// 	parts := strings.Split(beforeParams, "/")
// 	if len(parts) < 2 {
// 		return "", fmt.Errorf("DSN malformado: no se pudo encontrar la base de datos")
// 	}
// 	dbName := parts[len(parts)-1]
// 	return dbName, nil
// }

// func checkpointPath(db string, dir string) string {
// 	return filepath.Join(dir, fmt.Sprintf("%s_checkpoint.json", db))
// }

// func backupExists(db string, dir string) bool {
// 	_, err := os.Stat(checkpointPath(db, dir))
// 	return err == nil
// }

// func runFullBackup(cfg *Config, db string) error {
// 	timestamp := time.Now().Format("2006-01-02_15-04-05")
// 	backupFile := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_full_%s.sql", db, timestamp))

// 	// cmd := exec.Command("mysqldump",
// 	// 	"-u", cfg.User,
// 	// 	"-p"+cfg.Password,
// 	// 	"-h", cfg.Host,
// 	// 	"--databases", db,
// 	// 	"--routines", "--events", "--single-transaction")

// 	cmdArgs := []string{
// 		"-u", cfg.User,
// 		"-p" + cfg.Password, // Atención: la contraseña se imprimirá.
// 		"-h", cfg.Host,
// 		"-P", cfg.Port,
// 		"--databases", db,
// 		"--routines", "--events", "--single-transaction",
// 	}

// 	cmd := exec.Command("mysqldump", cmdArgs...)

// 	fmt.Println("Comando mysqldump a ejecutar:", "mysqldump", strings.Join(cmdArgs, " "))

// 	outputFile, err := os.Create(backupFile)
// 	if err != nil {
// 		return err
// 	}
// 	defer outputFile.Close()
// 	cmd.Stdout = outputFile

// 	return cmd.Run()
// }

// func getBinlogStatus(cfg *Config) (Checkpoint, error) {
// 	cmd := exec.Command("mysql", "-u", cfg.User, "-p"+cfg.Password, "-h", cfg.Host, "-P", cfg.Port, "-e", "SHOW MASTER STATUS\\G")
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return Checkpoint{}, err
// 	}

// 	var cp Checkpoint
// 	lines := strings.Split(string(output), "\n")
// 	for _, line := range lines {
// 		if strings.Contains(line, "File:") { // Usar Contains es más flexible que HasPrefix por los espacios iniciales
// 			// Eliminar cualquier espacio al inicio y luego buscar "File:"
// 			trimmedLine := strings.TrimSpace(line)
// 			if strings.HasPrefix(trimmedLine, "File:") {
// 				parts := strings.SplitN(trimmedLine, ":", 2)
// 				if len(parts) == 2 {
// 					cp.BinlogFile = strings.TrimSpace(parts[1])
// 				}
// 			}
// 		} else if strings.Contains(line, "Position:") { // Usar Contains
// 			// Eliminar cualquier espacio al inicio y luego buscar "Position:"
// 			trimmedLine := strings.TrimSpace(line)
// 			if strings.HasPrefix(trimmedLine, "Position:") {
// 				parts := strings.SplitN(trimmedLine, ":", 2)
// 				if len(parts) == 2 {
// 					posStr := strings.TrimSpace(parts[1])
// 					parsedPos, err := strconv.Atoi(posStr)
// 					if err != nil {
// 						fmt.Printf("Error convirtiendo Position a int: %v en '%s'\n", err, posStr)
// 						continue // Salta al siguiente iteración del bucle
// 					}
// 					cp.Position = parsedPos
// 				}
// 			}
// 		}
// 	}

// 	if cp.BinlogFile == "" || cp.Position == 0 {
// 		return cp, fmt.Errorf("estado de binlog inválido")
// 	}
// 	return cp, nil
// }

// func saveCheckpoint(cfg *Config, db string, cp Checkpoint) error {
// 	data, err := json.MarshalIndent(cp, "", "  ")
// 	if err != nil {
// 		return err
// 	}
// 	return os.WriteFile(checkpointPath(db, cfg.BackupDir), data, 0644)
// }

// func loadCheckpoint(cfg *Config, db string) (Checkpoint, error) {
// 	data, err := os.ReadFile(checkpointPath(db, cfg.BackupDir))
// 	if err != nil {
// 		return Checkpoint{}, err
// 	}
// 	var cp Checkpoint
// 	err = json.Unmarshal(data, &cp)
// 	return cp, err
// }

// func runIncrementalBackup(cfg *Config, db string, cp Checkpoint) (Checkpoint, error) {
// 	timestamp := time.Now().Format("2006-01-02_15-04-05")
// 	outFile := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_incremental_%s.sql", db, timestamp))

// 	cmd := exec.Command("mysqlbinlog",
// 		"--read-from-remote-server",
// 		"--host="+cfg.Host,
// 		"--user="+cfg.User,
// 		"--password="+cfg.Password,
// 		"--database="+db,
// 		"--skip-empty-transactions",
// 		"--start-position", fmt.Sprint(cp.Position),
// 		cp.BinlogFile)

// 	out, err := os.Create(outFile)
// 	if err != nil {
// 		return cp, fmt.Errorf("error creando archivo de salida: %v", err)
// 	}
// 	defer out.Close()

// 	cmd.Stdout = out
// 	if err := cmd.Run(); err != nil {
// 		return cp, fmt.Errorf("error ejecutando mysqlbinlog: %v", err)
// 	}

// 	return getBinlogStatus(cfg)
// }


// func hasChangesForDB(cfg *Config, db string, cp Checkpoint) (bool, error) {
// 	currentCp, err := getBinlogStatus(cfg)
// 	if err != nil {
// 		return false, err
// 	}

// 	// Si el binlog/posición no ha cambiado, no hay cambios
// 	if currentCp.BinlogFile == cp.BinlogFile && currentCp.Position == cp.Position {
// 		return false, nil
// 	}

// 	cmd := exec.Command("mysqlbinlog",
// 		"--read-from-remote-server",
// 		"--host="+cfg.Host,
// 		"--user="+cfg.User,
// 		"--password="+cfg.Password,
// 		"--database="+db,
// 		"--start-position", fmt.Sprint(cp.Position),
// 		"--stop-position", fmt.Sprint(currentCp.Position),
// 		"--base64-output=DECODE-ROWS",
// 		"--verbose", // Muestra queries en formato SQL legible
// 		cp.BinlogFile)

// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return false, fmt.Errorf("error ejecutando mysqlbinlog: %v", err)
// 	}

// 	// Busca cualquier operación DML/DDL para la BD específica
// 	patterns := []string{
// 		fmt.Sprintf("`%s`.", db),    // Ej: `gestion_car`.clientes
// 		fmt.Sprintf("USE `%s`", db), // Ej: USE `gestion_car`
// 		fmt.Sprintf("ALTER TABLE `%s`", db),
// 	}

// 	for _, pattern := range patterns {
// 		if strings.Contains(string(output), pattern) {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }

// func RunBackup(cfg *Config) {
// 	for _, db := range cfg.Databases {
// 		fmt.Println("Procesando DB:", db)

// 		if !backupExists(db, cfg.BackupDir) {
// 			fmt.Println("  No existe backup completo. Haciendo uno...")
// 			if err := runFullBackup(cfg, db); err != nil {
// 				fmt.Println("  ❌ Error backup full:", err)
// 				continue
// 			}
// 			cp, err := getBinlogStatus(cfg)
// 			if err != nil {
// 				fmt.Println("  ❌ Error obteniendo binlog:", err)
// 				continue
// 			}
// 			if err := saveCheckpoint(cfg, db, cp); err != nil {
// 				fmt.Println("  ❌ Error guardando checkpoint:", err)
// 			}
// 			fmt.Println("  ✅ Backup full realizado.")
// 		} else {
// 			cp, err := loadCheckpoint(cfg, db)
// 			if err != nil {
// 				fmt.Printf("  ❌ Error leyendo checkpoint para %s: %v\n", db, err)
// 				continue
// 			}

// 			hasChanges, err := hasChangesForDB(cfg, db, cp)
// 			if err != nil {
// 				fmt.Printf("  ❌ Error verificando cambios para %s: %v\n", db, err)
// 				continue
// 			}

// 			if !hasChanges {
// 				fmt.Printf("  ✅ No hay cambios para %s (binlog=%s, pos=%d)\n",
// 					db, cp.BinlogFile, cp.Position)
// 				continue
// 			}

// 			newCp, err := runIncrementalBackup(cfg, db, cp)
// 			if err != nil {
// 				fmt.Printf("  ❌ Error en backup incremental para %s: %v\n", db, err)
// 				continue
// 			}

// 			if err := saveCheckpoint(cfg, db, newCp); err != nil {
// 				fmt.Printf("  ❌ Error guardando checkpoint para %s: %v\n", db, err)
// 			}
// 			fmt.Printf("  ✅ Backup incremental realizado para %s (nueva posición: %d)\n",
// 				db, newCp.Position)
// 		}
// 	}
// }