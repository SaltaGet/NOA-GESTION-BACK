package database

import (
	"embed"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rs/zerolog/log"
)

func ApplyMigrations(dbURI string, migration embed.FS, dirName string) error {
	var databaseURI string

	if err := EnsureDatabaseExists(dbURI); err != nil {
		return fmt.Errorf("error al asegurar que la base de datos exista: %w", err)
	}

	if !strings.HasPrefix(dbURI, "mysql") {
		databaseURI = "mysql://" + dbURI
	} else {
		databaseURI = dbURI
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

	log.Info().Msg("Migraciones aplicadas exitosamente.")
	return nil
}
