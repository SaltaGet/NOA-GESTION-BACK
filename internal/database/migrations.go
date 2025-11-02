package database

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func ApplyMigrations(dbURI string, migration embed.FS, dirName string) error {
	var databaseURI string
	env := os.Getenv("ENV")

	if env == "prod" {
		if err := ensureDatabaseExists(dbURI); err != nil {
			return fmt.Errorf("error al asegurar que la base de datos exista: %w", err)
		}

		if !strings.HasPrefix(dbURI, "mysql") {
			databaseURI = "mysql://" + dbURI
		} else {
			databaseURI = dbURI
		}
	} else {
		// Para SQLite, migrate necesita la ruta del archivo con el prefijo "sqlite3://"
		databaseURI = "sqlite3://" + filePathFromURI(dbURI)
	}

	source, err := iofs.New(migration, dirName)
	if err != nil {
		return fmt.Errorf("error al crear la fuente de migración desde embed: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, databaseURI)
	if err != nil {
		return fmt.Errorf("error al crear instancia de migrate: %w", err)
	}

	// Aplicar todas las migraciones pendientes
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error al aplicar migraciones: %w", err)
	}

	fmt.Println("Migraciones aplicadas exitosamente.")
	return nil
}
