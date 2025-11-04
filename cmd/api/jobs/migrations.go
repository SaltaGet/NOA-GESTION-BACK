package jobs

import (
	"os"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/migrations"
)

func Migrations(deps *dependencies.MainContainer) error {
	connections, err := deps.TenantController.TenantService.TenantGetConections()
	if err != nil {
		return err
	}

	err = database.ApplyMigrations(os.Getenv("URI_DB"), migrations.MainMigrationsFS, "main")
	if err != nil {
		return err
	}

	for _, connection := range *connections {
		err = database.ApplyMigrations(connection, migrations.TenantMigrations, "tenant")
		if err != nil {
			return err
		}

	}
	return nil
}

// func Migrations(deps *dependencies.Application, path string) error {
// 	connections, err := deps.TenantController.TenantService.TenantGetConections()
// 	if err != nil {
// 		return err
// 	}

// 	tenant := filepath.Join(path, "tenant")
// 	main := filepath.Join(path, "main")

// 	err = database.ApplyMigrations(os.Getenv("URI_DB"), main)
// 	if err != nil {
// 		return err
// 	}

// 	for _, connection := range *connections {
// 		err = database.ApplyMigrations(connection, tenant)
// 		if err != nil {
// 			return err
// 		}

// 	}
// 	return nil
// }
