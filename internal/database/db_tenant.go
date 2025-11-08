package database

import (
	// "log"
	"os"
	"strings"

	"fmt"

	"database/sql"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func PrepareDB(uri string, memberAdmin models.Member) error {
	env := os.Getenv("ENV")

	var driver gorm.Dialector
	if err := ensureDatabaseExists(uri); err != nil {
		return fmt.Errorf("error al crear la base: %w", err)
	}
	driver = mysql.Open(uri)

	// Conexión GORM
	db, err := gorm.Open(driver, &gorm.Config{})
	if err != nil {
		handleDBCreationError(env, uri)
		return fmt.Errorf("error al inicializar DB: %w", err)
	}

	// Bajo nivel
	sqlDB, err := db.DB()
	if err != nil {
		handleDBCreationError(env, uri)
		return fmt.Errorf("error al obtener conexión de bajo nivel: %w", err)
	}
	defer sqlDB.Close()

	// Migraciones
	if err := db.AutoMigrate(
		&models.Client{},
		&models.ExpenseBuy{},
		&models.Member{},
		&models.Permission{},
		&models.Product{},
		&models.Role{},
		&models.Supplier{},
	); err != nil {
		handleDBCreationError(env, uri)
		return fmt.Errorf("error al migrar tablas: %w", err)
	}

	if err := db.Create(&models.Role{ID: 1, Name: "admin"}).Error; err != nil {
		handleDBCreationError(env, uri)
		return fmt.Errorf("error al crear rol admin: %w", err)
	}

	if err := db.Create(&memberAdmin).Error; err != nil {
		handleDBCreationError(env, uri)
		return fmt.Errorf("error al crear member admin: %w", err)
	}

	if err := db.Create(&permissions).Error; err != nil {
		handleDBCreationError(env, uri)
		return fmt.Errorf("error al migrar permisos: %w", err)
	}

	return nil
}

func handleDBCreationError(env, uri string) {
	if env == "prod" {
		_ = dropDatabase(uri)
	} else {
		_ = os.Remove(filePathFromURI(uri))
	}
}

func dropDatabase(dsn string) error {
	dbName, err := extractDBName(dsn)
	if err != nil {
		return fmt.Errorf("no se pudo extraer el nombre de la base: %w", err)
	}

	baseDSN := removeDBFromDSN(dsn)
	sqlDB, err := sql.Open("mysql", baseDSN)
	if err != nil {
		return fmt.Errorf("no se pudo conectar al servidor MySQL: %w", err)
	}
	defer sqlDB.Close()

	_, err = sqlDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", dbName))
	if err != nil {
		return fmt.Errorf("error al ejecutar DROP DATABASE: %w", err)
	}
	return nil
}

func extractDBName(dsn string) (string, error) {
	beforeParams := strings.SplitN(dsn, "?", 2)[0]
	parts := strings.Split(beforeParams, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("formato de DSN inválido (no se encontró la base)")
	}
	return parts[1], nil
}

func removeDBFromDSN(dsn string) string {
	i := strings.Index(dsn, "/")
	if i == -1 {
		return dsn
	}

	paramStart := strings.Index(dsn[i:], "?")
	if paramStart != -1 {
		return dsn[:i+1] + dsn[i+paramStart:]
	}
	return dsn[:i+1]
}

func UpdateModels(uri string) error {
	env := os.Getenv("ENV")

	var driver gorm.Dialector
	if err := ensureDatabaseExists(uri); err != nil {
		return fmt.Errorf("error al crear la base: %w", err)
	}
	driver = mysql.Open(uri)

	db, err := gorm.Open(driver, &gorm.Config{})
	if err != nil {
		handleDBCreationError(env, uri)
		return fmt.Errorf("error al inicializar DB: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		handleDBCreationError(env, uri)
		return fmt.Errorf("error al obtener conexión de bajo nivel: %w", err)
	}
	defer sqlDB.Close()

	if err := db.AutoMigrate(
		&models.Client{},
		&models.ExpenseBuy{},
		&models.Member{},
		&models.Permission{},
		&models.Product{},
		&models.Role{},
		&models.Supplier{},
	); err != nil {
		handleDBCreationError(env, uri)
		return fmt.Errorf("error al migrar tablas: %w", err)
	}

	return nil
}

var permissions []models.Permission = []models.Permission{
	{Code: "create_client", Details: "Crear clientes", Group: "clients"},
	{Code: "update_client", Details: "Actualizar clientes", Group: "clients"},
	{Code: "delete_client", Details: "Eliminar clientes", Group: "clients"},
	{Code: "create_expense", Details: "Crear gastos", Group: "expenses"},
	{Code: "update_expense", Details: "Actualizar gastos", Group: "expenses"},
}