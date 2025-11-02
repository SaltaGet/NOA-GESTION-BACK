package jobs

import (
	"os"

	"github.com/DanielChachagua/GestionCar/pkg/database"
	"github.com/DanielChachagua/GestionCar/pkg/dependencies"
	"github.com/DanielChachagua/GestionCar/pkg/migrations"
)

func Migrations(deps *dependencies.Application) error {
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
