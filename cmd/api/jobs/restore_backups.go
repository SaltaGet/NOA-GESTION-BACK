package jobs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"time"
	"github.com/rs/zerolog/log"
)

// BackupFile representa un archivo de backup con su tipo y fecha.
type BackupFile struct {
	Path string
	Time time.Time
	Type string // "full" o "incremental"
}

// RestoreChain representa la cadena completa de backups para restaurar
type RestoreChain struct {
	FullBackup          *BackupFile
	IncrementalBackups  []BackupFile
}

// GetRestoreChain devuelve el último backup FULL + todos los incrementales posteriores
func GetRestoreChain(backupDir string, dbName string) (*RestoreChain, error) {
	files, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el directorio: %w", err)
	}

	// Regex para detectar archivos del tipo nombreDB_full_YYYY-MM-DD_HH-MM-SS.sql
	pattern := fmt.Sprintf(`^%s_(full|incremental)_(\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2})\.sql$`, regexp.QuoteMeta(dbName))
	re := regexp.MustCompile(pattern)

	var fullBackups []BackupFile
	var incrementalBackups []BackupFile

	for _, entry := range files {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		matches := re.FindStringSubmatch(name)
		if len(matches) != 3 {
			continue
		}

		backupType := matches[1]
		timestamp := matches[2]
		parsedTime, err := time.Parse("2006-01-02_15-04-05", timestamp)
		if err != nil {
			continue
		}

		backup := BackupFile{
			Path: filepath.Join(backupDir, name),
			Time: parsedTime,
			Type: backupType,
		}

		if backupType == "full" {
			fullBackups = append(fullBackups, backup)
		} else {
			incrementalBackups = append(incrementalBackups, backup)
		}
	}

	if len(fullBackups) == 0 {
		return nil, fmt.Errorf("no se encontró ningún backup FULL para %s", dbName)
	}

	// Ordenar backups full por fecha (más reciente primero)
	sort.Slice(fullBackups, func(i, j int) bool {
		return fullBackups[i].Time.After(fullBackups[j].Time)
	})

	lastFull := fullBackups[0]

	// Filtrar incrementales posteriores al último full
	var relevantIncrementals []BackupFile
	for _, inc := range incrementalBackups {
		if inc.Time.After(lastFull.Time) {
			relevantIncrementals = append(relevantIncrementals, inc)
		}
	}

	// Ordenar incrementales por fecha (más antiguo primero)
	sort.Slice(relevantIncrementals, func(i, j int) bool {
		return relevantIncrementals[i].Time.Before(relevantIncrementals[j].Time)
	})

	return &RestoreChain{
		FullBackup:         &lastFull,
		IncrementalBackups: relevantIncrementals,
	}, nil
}

// RestoreDatabase restaura una base de datos usando la cadena de backups
func RestoreDatabase(cfg *Config, dbName string, chain *RestoreChain) error {
	log.Info().Msgf("📦 Restaurando %s desde backup FULL: %s\n", dbName, chain.FullBackup.Path)
	
	// 1. Restaurar backup FULL (sin especificar DB, el SQL contiene CREATE DATABASE)
	if err := executeSQLFile(cfg, dbName, chain.FullBackup.Path, true); err != nil {
		return fmt.Errorf("error restaurando backup full: %w", err)
	}
	
	log.Info().Msgf("✅ Backup FULL restaurado\n")
	
	// 2. Aplicar backups incrementales en orden
	if len(chain.IncrementalBackups) > 0 {
		log.Info().Msgf("🔄 Aplicando %d backups incrementales...\n", len(chain.IncrementalBackups))
		
		for i, inc := range chain.IncrementalBackups {
			log.Info().Msgf("   %d/%d: %s\n", i+1, len(chain.IncrementalBackups), filepath.Base(inc.Path))
			
			// Para incrementales sí especificamos la DB
			if err := executeSQLFile(cfg, dbName, inc.Path, false); err != nil {
				return fmt.Errorf("error aplicando incremental %s: %w", inc.Path, err)
			}
		}
		
		log.Info().Msgf("✅ Backups incrementales aplicados\n")
	} else {
		log.Info().Msgf("ℹ️  No hay backups incrementales posteriores al FULL\n")
	}
	
	return nil
}

// executeSQLFile ejecuta un archivo SQL en la base de datos
func executeSQLFile(cfg *Config, dbName string, sqlFile string, isFullBackup bool) error {
	// Verificar que el archivo existe
	if _, err := os.Stat(sqlFile); os.IsNotExist(err) {
		return fmt.Errorf("archivo no existe: %s", sqlFile)
	}
	
	var args []string
	
	if isFullBackup {
		// Para backup FULL: NO especificar DB (el SQL contiene CREATE DATABASE y USE)
		// mariadb -u USER -pPASS -h HOST -P PORT < file.sql
		args = []string{
			"-u" + cfg.User,
			"-p" + cfg.Password,
			"-h", cfg.Host,
			"-P", cfg.Port,
		}
	} else {
		// Para backup INCREMENTAL: especificar DB (el SQL solo tiene DML)
		// mariadb -u USER -pPASS -h HOST -P PORT DB < file.sql
		args = []string{
			"-u" + cfg.User,
			"-p" + cfg.Password,
			"-h", cfg.Host,
			"-P", cfg.Port,
			dbName,
		}
	}
	
	cmd := exec.Command("mariadb", args...)
	
	// Abrir archivo SQL
	file, err := os.Open(sqlFile)
	if err != nil {
		return fmt.Errorf("no se pudo abrir %s: %w", sqlFile, err)
	}
	defer file.Close()
	
	// Redirigir archivo SQL a stdin
	cmd.Stdin = file
	
	// Capturar salida para debug
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error ejecutando SQL: %w\nOutput: %s", err, string(output))
	}
	
	return nil
}

// GetLastBackup (DEPRECATED) - Mantener por compatibilidad pero usar GetRestoreChain
// Esta función solo devuelve el último backup individual, NO LA CADENA COMPLETA
func GetLastBackup(backupDir string, dbName string) (*BackupFile, error) {
	files, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el directorio: %w", err)
	}

	pattern := fmt.Sprintf(`^%s_(full|incremental)_(\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2})\.sql$`, regexp.QuoteMeta(dbName))
	re := regexp.MustCompile(pattern)

	var backups []BackupFile

	for _, entry := range files {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		matches := re.FindStringSubmatch(name)
		if len(matches) != 3 {
			continue
		}

		backupType := matches[1]
		timestamp := matches[2]
		parsedTime, err := time.Parse("2006-01-02_15-04-05", timestamp)
		if err != nil {
			continue
		}

		backups = append(backups, BackupFile{
			Path: filepath.Join(backupDir, name),
			Time: parsedTime,
			Type: backupType,
		})
	}

	if len(backups) == 0 {
		return nil, fmt.Errorf("no se encontraron backups para la base %s", dbName)
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Time.After(backups[j].Time)
	})

	return &backups[0], nil
}

// Ejemplo de uso
func ExampleRestore(cfg *Config, dbName string) error {
	// 1. Obtener la cadena de backups
	chain, err := GetRestoreChain(cfg.BackupDir, dbName)
	if err != nil {
		return fmt.Errorf("error obteniendo cadena de backups: %w", err)
	}
	
	// 2. Mostrar información
	log.Info().Msgf("═══════════════════════════════════════\n")
	log.Info().Msgf("📋 Base de datos: %s\n", dbName)
	log.Info().Msgf("📦 Backup FULL: %s (%s)\n", 
		filepath.Base(chain.FullBackup.Path), 
		chain.FullBackup.Time.Format("2006-01-02 15:04:05"))
	
	if len(chain.IncrementalBackups) > 0 {
		log.Info().Msgf("🔄 Incrementales: %d archivos\n", len(chain.IncrementalBackups))
		for _, inc := range chain.IncrementalBackups {
			log.Info().Msgf("   - %s (%s)\n", 
				filepath.Base(inc.Path), 
				inc.Time.Format("2006-01-02 15:04:05"))
		}
	}
	log.Info().Msgf("═══════════════════════════════════════\n")
	
	// 4. Ejecutar restauración
	log.Info().Msgf("🚀 Iniciando restauración...")
	
	if err := RestoreDatabase(cfg, dbName, chain); err != nil {
		return fmt.Errorf("error durante la restauración: %w", err)
	}
	
	log.Info().Msgf("✅ Base de datos %s restaurada correctamente\n", dbName)
	
	return nil
}

// package jobs

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"regexp"
// 	"sort"
// 	"time"
// )

// // BackupFile representa un archivo de backup con su tipo y fecha.
// type BackupFile struct {
// 	Path string
// 	Time time.Time
// 	Type string // "full" o "incremental"
// }

// // GetLastBackup devuelve el backup más reciente (full o incremental) de una DB.
// func GetLastBackup(backupDir string, dbName string) (*BackupFile, error) {
// 	files, err := os.ReadDir(backupDir)
// 	if err != nil {
// 		return nil, fmt.Errorf("no se pudo leer el directorio: %w", err)
// 	}

// 	// Regex para detectar archivos del tipo nombreDB_full_YYYY-MM-DD_HH-MM-SS.sql
// 	pattern := fmt.Sprintf(`^%s_(full|incremental)_(\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2})\.sql$`, regexp.QuoteMeta(dbName))
// 	re := regexp.MustCompile(pattern)

// 	var backups []BackupFile

// 	for _, entry := range files {
// 		if entry.IsDir() {
// 			continue
// 		}

// 		name := entry.Name()
// 		matches := re.FindStringSubmatch(name)
// 		if len(matches) != 3 {
// 			continue
// 		}

// 		backupType := matches[1]
// 		timestamp := matches[2]
// 		parsedTime, err := time.Parse("2006-01-02_15-04-05", timestamp)
// 		if err != nil {
// 			continue
// 		}

// 		backups = append(backups, BackupFile{
// 			Path: filepath.Join(backupDir, name),
// 			Time: parsedTime,
// 			Type: backupType,
// 		})
// 	}

// 	if len(backups) == 0 {
// 		return nil, fmt.Errorf("no se encontraron backups para la base %s", dbName)
// 	}

// 	// Ordenamos por fecha descendente
// 	sort.Slice(backups, func(i, j int) bool {
// 		return backups[i].Time.After(backups[j].Time)
// 	})

// 	return &backups[0], nil
// }
