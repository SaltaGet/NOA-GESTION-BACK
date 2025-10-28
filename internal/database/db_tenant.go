package database

// import (
// 	// "log"
// 	"os"
// 	"strings"

// 	"fmt"

// 	"database/sql"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// func PrepareDB(uri string) error {
// 	env := os.Getenv("ENV")

// 	var driver gorm.Dialector
// 	if err := ensureDatabaseExists(uri); err != nil {
// 		return fmt.Errorf("error al crear la base: %w", err)
// 	}
// 	driver = mysql.Open(uri)

// 	// Conexión GORM
// 	db, err := gorm.Open(driver, &gorm.Config{})
// 	if err != nil {
// 		handleDBCreationError(env, uri)
// 		return fmt.Errorf("error al inicializar DB: %w", err)
// 	}

// 	// Bajo nivel
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		handleDBCreationError(env, uri)
// 		return fmt.Errorf("error al obtener conexión de bajo nivel: %w", err)
// 	}
// 	defer sqlDB.Close()

// 	// Migraciones
// 	if err := db.AutoMigrate(
// 		&models.Attendance{},
// 		&models.Client{},
// 		&models.Employee{},
// 		&models.Expense{},
// 		&models.Income{},
// 		&models.Member{},
// 		&models.MovementType{},
// 		&models.Permission{},
// 		&models.Product{},
// 		&models.PurchaseOrder{},
// 		&models.PurchaseProduct{},
// 		&models.ResumeExpense{},
// 		&models.ResumeIncome{},
// 		&models.Role{},
// 		&models.Service{},
// 		&models.Supplier{},
// 		&models.Vehicle{},
// 	); err != nil {
// 		handleDBCreationError(env, uri)
// 		return fmt.Errorf("error al migrar tablas: %w", err)
// 	}

// 	if err := db.Create(&permissions).Error; err != nil {
// 		handleDBCreationError(env, uri)
// 		return fmt.Errorf("error al migrar permisos: %w", err)
// 	}

// 	return nil
// }

// func handleDBCreationError(env, uri string) {
// 	if env == "prod" {
// 		_ = dropDatabase(uri)
// 	} else {
// 		_ = os.Remove(filePathFromURI(uri))
// 	}
// }

// func dropDatabase(dsn string) error {
// 	dbName, err := extractDBName(dsn)
// 	if err != nil {
// 		return fmt.Errorf("no se pudo extraer el nombre de la base: %w", err)
// 	}

// 	baseDSN := removeDBFromDSN(dsn)
// 	sqlDB, err := sql.Open("mysql", baseDSN)
// 	if err != nil {
// 		return fmt.Errorf("no se pudo conectar al servidor MySQL: %w", err)
// 	}
// 	defer sqlDB.Close()

// 	_, err = sqlDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", dbName))
// 	if err != nil {
// 		return fmt.Errorf("error al ejecutar DROP DATABASE: %w", err)
// 	}
// 	return nil
// }

// func extractDBName(dsn string) (string, error) {
// 	beforeParams := strings.SplitN(dsn, "?", 2)[0]
// 	parts := strings.Split(beforeParams, "/")
// 	if len(parts) < 2 {
// 		return "", fmt.Errorf("formato de DSN inválido (no se encontró la base)")
// 	}
// 	return parts[1], nil
// }

// func removeDBFromDSN(dsn string) string {
// 	i := strings.Index(dsn, "/")
// 	if i == -1 {
// 		return dsn
// 	}

// 	paramStart := strings.Index(dsn[i:], "?")
// 	if paramStart != -1 {
// 		return dsn[:i+1] + dsn[i+paramStart:]
// 	}
// 	return dsn[:i+1]
// }

// func UpdateModels(uri string) error {
// 	env := os.Getenv("ENV")

// 	var driver gorm.Dialector
// 	if err := ensureDatabaseExists(uri); err != nil {
// 		return fmt.Errorf("error al crear la base: %w", err)
// 	}
// 	driver = mysql.Open(uri)

// 	db, err := gorm.Open(driver, &gorm.Config{})
// 	if err != nil {
// 		handleDBCreationError(env, uri)
// 		return fmt.Errorf("error al inicializar DB: %w", err)
// 	}

// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		handleDBCreationError(env, uri)
// 		return fmt.Errorf("error al obtener conexión de bajo nivel: %w", err)
// 	}
// 	defer sqlDB.Close()

// 	if err := db.AutoMigrate(
// 		&models.Attendance{},
// 		&models.Client{},
// 		&models.Employee{},
// 		&models.Expense{},
// 		&models.Income{},
// 		&models.Member{},
// 		&models.MovementType{},
// 		&models.Permission{},
// 		&models.Product{},
// 		&models.PurchaseOrder{},
// 		&models.PurchaseProduct{},
// 		&models.ResumeExpense{},
// 		&models.ResumeIncome{},
// 		&models.Role{},
// 		&models.Service{},
// 		&models.Supplier{},
// 		&models.Vehicle{},
// 	); err != nil {
// 		handleDBCreationError(env, uri)
// 		return fmt.Errorf("error al migrar tablas: %w", err)
// 	}

// 	return nil
// }
