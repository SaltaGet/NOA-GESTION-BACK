package jobs

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Optimized backup using pg_dump in Directory format with parallelism
// This matches "Opción 2" from the recommendations.
//
// PARA RESTAURAR (RESTORE):
// Usar el comando 'pg_restore' con la opción '-d' para la base de datos destino y '-j' para paralelismo.
// El formato Directorio (-Fd) requiere usar pg_restore, NO psql.
//
// Comando ejemplo:
// pg_restore -h localhost -p 5432 -U postgres -d mi_base_datos -j 4 -v --clean "/ruta/al/backup_dir/"
//
// Flags útiles:
// -j 4      : Usa 4 hilos para restaurar (más rápido).
// --clean   : Borra los objetos existentes en la DB antes de crear los nuevos.
// --if-exists: Se usa junto con --clean para no fallar si los objetos no existen.
// -v        : Verbose (muestra progreso).
func RunOptimizedBackupPostgres(cfg *ConfigPostgres) {
	log.Info().Msgf("⏰ [CRON PG OPTIMIZED] Iniciando backup de %d bases de datos...", len(cfg.Databases))

	// Ensure main backup directory exists
	// We assume cfg.BackupDir is already set
	if err := os.MkdirAll(cfg.BackupDir, 0755); err != nil {
		log.Error().Err(err).Msgf("❌ Error creando directorio base de backups: %s", cfg.BackupDir)
		return
	}

	// We'll process databases with some concurrency (e.g. 2 at a time)
	// inside each database worker, pg_dump will use -j 4
	maxConcurrentDBs := 2
	tasks := make(chan string, len(cfg.Databases))
	results := make(chan string, len(cfg.Databases))
	var wg sync.WaitGroup

	for i := 0; i < maxConcurrentDBs; i++ {
		wg.Add(1)
		go runOptimizedBackupWorker(cfg, tasks, &wg, results)
	}

	for _, db := range cfg.Databases {
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
	log.Info().Msgf("✅ Backups OPTIMIZED PG completados: %d/%d exitosos", successCount, len(cfg.Databases))

	// Cleanup old backups immediately after backup
	RunCleanupBackups(cfg, 7)
}

func runOptimizedBackupWorker(cfg *ConfigPostgres, tasks <-chan string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	for db := range tasks {
		ts := time.Now().Format("2006-01-02_15-04-05")
		// The backup "file" for directory format is actually a directory.
		backupPath := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s_opt_pg_%s", db, ts))

		// pg_dump arguments for Option 2:
		// -F d : Directory format
		// -j 4 : Parallel dump jobs
		// -Z 5 : Compression level
		args := []string{
			"-h", cfg.Host,
			"-p", cfg.Port,
			"-U", cfg.User,
			"-d", db,
			"-F", "d",
			"-j", "4",
			"-Z", "5",
			"-f", backupPath,
		}

		cmd := exec.Command("pg_dump", args...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", cfg.Password))

		// Capture standard error to see pg_dump logs
		cmd.Stderr = os.Stderr

		startTime := time.Now()
		if err := cmd.Run(); err != nil {
			log.Err(err).Msgf("❌ [Worker PG] Error en backup optimizado de %s", db)
			// Cleanup on failure
			os.RemoveAll(backupPath)
			results <- fmt.Sprintf("ERROR:%s", db)
		} else {
			duration := time.Since(startTime)
			log.Info().Msgf("✅ [Worker PG] Backup optimizado completado para %s en %s (%v)", db, backupPath, duration)
			results <- fmt.Sprintf("SUCCESS:%s", db)
		}
	}
}

// Helper to restore an optimized backup
// WARNING: This will overwrite data if --clean is used.
func RestoreOptimizedBackupPostgres(cfg *ConfigPostgres, dbName string, backupDir string) error {
	log.Info().Msgf("🔄 Iniciando restauración de %s desde %s", dbName, backupDir)

	// Check if backup dir exists
	info, err := os.Stat(backupDir)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("directorio de backup no válido: %s", backupDir)
	}

	args := []string{
		"-h", cfg.Host,
		"-p", cfg.Port,
		"-U", cfg.User,
		"-d", dbName,
		"-j", "4", // Parallel restore
		"-v",          // Verbose
		"--clean",     // Clean existing objects
		"--if-exists", // Don't fail if objects don't exist
		backupDir,     // Source directory
	}

	cmd := exec.Command("pg_restore", args...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", cfg.Password))

	// Pipe stdout and stderr to see progress
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	startTime := time.Now()
	if err := cmd.Run(); err != nil {
		log.Err(err).Msgf("❌ [Restore PG] Error restaurando %s", dbName)
		return err
	}

	log.Info().Msgf("✅ [Restore PG] Restauración completada para %s en %v", dbName, time.Since(startTime))
	return nil
}

// Cleanup backups keeping only the last N per database
func RunCleanupBackups(cfg *ConfigPostgres, retentionCount int) {
	log.Info().Msgf("🧹 [CLEANUP PG] Iniciando limpieza de backups antiguos (Retención: %d)...", retentionCount)

	entries, err := os.ReadDir(cfg.BackupDir)
	if err != nil {
		log.Error().Err(err).Msg("❌ [CLEANUP PG] Error leyendo directorio de backups")
		return
	}

	// Group backups by database name
	// Patterns: dbname_opt_pg_TIMESTAMP, dbname_full_pg_TIMESTAMP.sql, dbname_incremental_pg_TIMESTAMP.sql
	backups := make(map[string][]fs.DirEntry)

	// Separators to reliably identify database name
	separators := []string{"_opt_pg_", "_full_pg_", "_incremental_pg_"}

	for _, entry := range entries {
		name := entry.Name()

		// Skip non-backup files (e.g. checkpoints or unrelated files)
		isBackup := false
		var dbName string

		for _, sep := range separators {
			if strings.Contains(name, sep) {
				parts := strings.Split(name, sep)
				if len(parts) >= 2 {
					dbName = parts[0]
					isBackup = true
					break
				}
			}
		}

		if !isBackup {
			continue
		}

		backups[dbName] = append(backups[dbName], entry)
	}

	deletedCount := 0
	for db, entries := range backups {
		if len(entries) <= retentionCount {
			continue
		}

		// Sort by name (which acts as date sort due to YYYY-MM-DD format) descending
		// Newest first
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name() > entries[j].Name()
		})

		// Delete from index retentionCount onwards (older files)
		toDelete := entries[retentionCount:]
		log.Info().Msgf("🧹 [CLEANUP PG] Encontrados %d backups para '%s'. Eliminando %d antiguos...", len(entries), db, len(toDelete))

		for _, entry := range toDelete {
			fullPath := filepath.Join(cfg.BackupDir, entry.Name())
			// log.Info().Msgf("🗑️ [CLEANUP PG] Eliminando backup antiguo: %s", entry.Name())
			if err := os.RemoveAll(fullPath); err != nil {
				log.Error().Err(err).Msgf("❌ [CLEANUP PG] Error eliminando %s", entry.Name())
			} else {
				deletedCount++
			}
		}
	}

	log.Info().Msgf("✅ [CLEANUP PG] Limpieza completada. %d backups eliminados.", deletedCount)
}
