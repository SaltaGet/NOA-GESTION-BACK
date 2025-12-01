package jobs

import (
	"fmt"
	"log"
	"os"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/migrations"
)

func Migrations(deps *dependencies.MainContainer) error {
	tenants, err := deps.TenantController.TenantService.TenantGetConections()
	if err != nil {
		return err
	}

	log.Println("Aplicando migraciones a la base de datos principal...")
	err = database.ApplyMigrations(os.Getenv("URI_DB"), migrations.MainMigrationsFS, "main")
	if err != nil {
		// La migración de la DB principal es crítica. Si falla, detenemos todo.
		return fmt.Errorf("error crítico al aplicar migraciones a la DB principal: %w", err)
	}

	var migrationErrors []error
	log.Printf("Aplicando migraciones a %d base(s) de datos de tenants...", len(tenants))
	for _, tenantInfo := range tenants {
		err = database.ApplyMigrations(tenantInfo.Connection, migrations.TenantMigrations, "tenant")
		if err != nil {
			// En lugar de detener todo, registramos el error y continuamos.
			log.Printf("⚠️ Error al aplicar migración para Tenant ID: %d, Nombre: '%s', Identificador: '%s': %v", tenantInfo.ID, tenantInfo.Name, tenantInfo.Identifier, err)
			migrationErrors = append(migrationErrors, fmt.Errorf("tenant ID %d (%s): %w", tenantInfo.ID, tenantInfo.Name, err))
		}
	}
	if len(migrationErrors) > 0 {
		return fmt.Errorf("ocurrieron %d errores durante las migraciones de los tenants", len(migrationErrors))
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
