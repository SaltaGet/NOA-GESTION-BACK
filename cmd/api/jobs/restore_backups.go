package jobs

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

// BackupFile representa un archivo de backup con su tipo y fecha.
type BackupFile struct {
	Path string
	Time time.Time
	Type string // "full" o "incremental"
}

// GetLastBackup devuelve el backup más reciente (full o incremental) de una DB.
func GetLastBackup(backupDir string, dbName string) (*BackupFile, error) {
	files, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el directorio: %w", err)
	}

	// Regex para detectar archivos del tipo nombreDB_full_YYYY-MM-DD_HH-MM-SS.sql
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

	// Ordenamos por fecha descendente
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Time.After(backups[j].Time)
	})

	return &backups[0], nil
}
