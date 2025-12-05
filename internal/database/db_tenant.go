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
	if err := EnsureDatabaseExists(uri); err != nil {
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
		&models.Admin{},
		&models.AuditLogAdmin{},
		&models.AuditLog{},
		&models.CashRegister{},
		&models.Category{},
		&models.Client{},
		&models.Deposit{},

		&models.ExpenseBuy{},
		&models.ExpenseBuyItem{},
		&models.ExpenseOther{},
		&models.PayExpenseBuy{},
		&models.PayExpenseOther{},
		&models.TypeExpense{},

		&models.IncomeSale{},
		&models.IncomeSaleItem{},
		&models.PayIncome{},
		&models.IncomeOther{},
		&models.TypeIncome{},

		&models.Member{},
		&models.MovementStock{},
		&models.Permission{},
		&models.PointSale{},
		&models.Product{},
		&models.Role{},
		&models.StockPointSale{},
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

	if err := db.Create(&models.Client{
		ID:             1,
		FirstName:      "Consumidor",
		LastName:       "Final",
		CompanyName:    nil,
		Identifier:     nil,
		Email:          nil,
		Phone:          nil,
		Address:        nil,
		MemberCreateID: memberAdmin.ID,
	}).Error; err != nil {
		handleDBCreationError(env, uri)
		return fmt.Errorf("error al crear cliente consumidor final: %w", err)
	}

	if err := db.Create(&models.Category{
		Name: "Sin categoría",
	}).Error; err != nil {
		handleDBCreationError(env, uri)
		return fmt.Errorf("error al crear categoría sin categoría: %w", err)
	}

	if err := db.Create(&models.Supplier{
		ID:          1,
		Name:        "Sin proveedor",
		CompanyName: "Sin nombre",
	}).Error; err != nil {
		handleDBCreationError(env, uri)
		return fmt.Errorf("error al crear proveedor sin proveedor: %w", err)
	}

	if err := db.Create(&models.TypeExpense{
		ID:   1,
		Name: "Otros",
	}).Error; err != nil {
		handleDBCreationError(env, uri)
		return fmt.Errorf("error al crear tipo de gasto sin tipo de gasto: %w", err)
	}

	if err := db.Create(&models.TypeIncome{
		ID:   1,
		Name: "Otros",
	}).Error; err != nil {
		handleDBCreationError(env, uri)
		return fmt.Errorf("error al crear tipo de gasto sin tipo de gasto: %w", err)
	}

	return nil
}

func handleDBCreationError(env, uri string) {
	if env == "prod" {
		_ = dropDatabase(uri)
	} else {
		_ = os.Remove(FilePathFromURI(uri))
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
	if err := EnsureDatabaseExists(uri); err != nil {
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
		&models.Admin{},
		&models.AuditLogAdmin{},
		&models.AuditLog{},
		&models.CashRegister{},
		&models.Category{},
		&models.Client{},
		&models.Deposit{},

		&models.ExpenseBuy{},
		&models.ExpenseBuyItem{},
		&models.ExpenseOther{},
		&models.PayExpenseBuy{},
		&models.PayExpenseOther{},
		&models.TypeExpense{},

		&models.IncomeSale{},
		&models.IncomeSaleItem{},
		&models.PayIncome{},
		&models.IncomeOther{},
		&models.TypeIncome{},

		&models.Member{},
		&models.MovementStock{},
		&models.Permission{},
		&models.PointSale{},
		&models.Product{},
		&models.Role{},
		&models.StockPointSale{},
		&models.Supplier{},
	); err != nil {
		handleDBCreationError(env, uri)
		return fmt.Errorf("error al migrar tablas: %w", err)
	}

	return nil
}

var permissions []models.Permission = []models.Permission{
	//CASH REGISTER
	{Code: "CR01", Name: "verificar apertura", Details: "Verifica si existe alguna apertura de caja del punto de venta", Group: "caja", Environment: "point_sale"},
	{Code: "CR02", Name: "apertura", Details: "Crea una nueva apertura de caja", Group: "caja", Environment: "point_sale"},
	{Code: "CR03", Name: "informes", Details: "Obtener infomes de caja", Group: "caja", Environment: "point_sale"},
	{Code: "CR04", Name: "cierre", Details: "Cerrar una caja abierta", Group: "caja", Environment: "point_sale"},
	{Code: "CR05", Name: "informe específico", Details: "Obtener informe de una caja específica", Group: "caja", Environment: "point_sale"},
	//CATEGORY
	{Code: "CAT01", Name: "crear", Details: "Crear una nueva categoría", Group: "categoría", Environment: "point_sale"},
	{Code: "CAT02", Name: "actualizar", Details: "Actualizar una categoria existente", Group: "categoría", Environment: "point_sale"},
	{Code: "CAT03", Name: "eliminar", Details: "Eliminar una cagotería", Group: "categoría", Environment: "point_sale"},
	{Code: "CAT04", Name: "obtner uno en específico", Details: "Obtener una categoría en específico", Group: "categoría", Environment: "point_sale"},
	{Code: "CAT05", Name: "obtner todos", Details: "Obtener todas las categorías", Group: "categoría", Environment: "point_sale"},
	//CLIENT
	{Code: "CL01", Name: "crear", Details: "Crear un nuevo cliente", Group: "cliente", Environment: "point_sale"},
	{Code: "CL02", Name: "actualizar", Details: "Actualizar un cliente existente", Group: "cliente", Environment: "point_sale"},
	{Code: "CL03", Name: "eliminar", Details: "Eliminar un cliente", Group: "cliente", Environment: "point_sale"},
	{Code: "CL04", Name: "obtener uno en específico", Details: "Obtener un cliente en específico", Group: "cliente", Environment: "point_sale"},
	{Code: "CL05", Name: "obtener por filtro", Details: "Obtener clientes por filtros", Group: "cliente", Environment: "point_sale"},
	{Code: "CL06", Name: "obtener todos", Details: "Obtener todos lo clientes", Group: "cliente", Environment: "point_sale"},
	//DEPOSIT
	{Code: "DEP01", Name: "obtener uno en específico", Details: "Obtener información de un productos en específico del depósito", Group: "deposito", Environment: "point_sale"},
	{Code: "DEP02", Name: "obtener por nombre de producto", Details: "Obtener información de un producto en especifico del depósito por nombre", Group: "deposito", Environment: "point_sale"},
	{Code: "DEP03", Name: "obtener por código de producto", Details: "Obtener información de un producto en especifico del depósito por codigo", Group: "deposito", Environment: "point_sale"},
	{Code: "DEP04", Name: "obtener todos", Details: "Obtener todos los productos del depósito", Group: "deposito", Environment: "point_sale"},
	{Code: "DEP05", Name: "actualizar stock de producto", Details: "Actualizar el stock de productos del depósito", Group: "deposito", Environment: "point_sale"},
	//EXPENSE BUY
	{Code: "EB01", Name: "crear", Details: "Crear un nuevo gasto de compra", Group: "gastos de compra", Environment: "dashboard"},
	{Code: "EB02", Name: "actualizar", Details: "Actualizar un gasto de compra existente", Group: "gastos de compra", Environment: "dashboard"},
	{Code: "EB03", Name: "eliminar", Details: "Eliminar un gasto de compra", Group: "gastos de compra", Environment: "dashboard"},
	{Code: "EB04", Name: "obtener uno en específico", Details: "Obtener un gasto de compra en específico", Group: "gastos de compra", Environment: "dashboard"},
	{Code: "EB05", Name: "obtener por fecha", Details: "Obtener gastos de compra por rango de fechas", Group: "gastos de compra", Environment: "dashboard"},
	//EXPENSE OTHER
	{Code: "EO01", Name: "crear", Details: "Crear un nuevo otros gastos", Group: "otros gastos", Environment: "dashboard"},
	{Code: "EO02", Name: "actualizar", Details: "Actualizar un otros gastos existente", Group: "otros gastos", Environment: "dashboard"},
	{Code: "EO03", Name: "eliminar", Details: "Eliminar un otros gastos", Group: "otros gastos", Environment: "dashboard"},
	{Code: "EO04", Name: "obtener uno en específico", Details: "Obtener un otros gastos en específico", Group: "otros gastos", Environment: "dashboard"},
	{Code: "EO05", Name: "obtener por fecha", Details: "Obtener otros gastos por rango de fechas", Group: "otros gastos", Environment: "dashboard"},
	//EXPENSE OTHER POINT SALE
	{Code: "EOPS01", Name: "crear", Details: "Crear otro gasto para un punto de venta", Group: "otros grastos - punto de venta", Environment: "point_sale"},
	{Code: "EOPS02", Name: "actualizar", Details: "Actualizar gasto de un punto de venta", Group: "otros grastos - punto de venta", Environment: "point_sale"},
	{Code: "EOPS03", Name: "eliminar", Details: "Eliminar un gasto de un punto de venta", Group: "otros grastos - punto de venta", Environment: "point_sale"},
	{Code: "EOPS04", Name: "obtener uno en específico", Details: "Obtener un gasto en específico de un punto de venta", Group: "otros grastos - punto de venta", Environment: "point_sale"},
	{Code: "EOPS05", Name: "obtener por fecha", Details: "Obtener gastos de compra por rango de fechas de un punto de venta", Group: "otros grastos - punto de venta", Environment: "point_sale"},
	//INCOME OTHER
	{Code: "INO01", Name: "crear", Details: "Crear otro ingreso", Group: "otros ingresos", Environment: "dashboard"},
	{Code: "INO02", Name: "actualizar", Details: "Actualizar un ingreso existente", Group: "otros ingresos", Environment: "dashboard"},
	{Code: "INO03", Name: "eliminar", Details: "Eliminar un ingreso", Group: "otros ingresos", Environment: "dashboard"},
	{Code: "INO04", Name: "obtener uno en específico", Details: "Obtener un ingreso en específico", Group: "otros ingresos", Environment: "dashboard"},
	{Code: "INO05", Name: "obtener por fecha", Details: "Obtener ingresos por rango de fechas", Group: "otros ingresos", Environment: "dashboard"},
	//INCOME OTHER POINT SALE
	{Code: "INOPS01", Name: "crear", Details: "Crear otro ingreso para un punto de venta", Group: "otros ingresos - punto de venta", Environment: "point_sale"},
	{Code: "INOPS02", Name: "actualizar", Details: "Actualizar un ingreso de un punto de venta", Group: "otros ingresos - punto de venta", Environment: "point_sale"},
	{Code: "INOPS03", Name: "eliminar", Details: "Eliminar un ingreso de un punto de venta", Group: "otros ingresos - punto de venta", Environment: "point_sale"},
	{Code: "INOPS04", Name: "obtener uno en específico", Details: "Obtener un ingreso otros en específico", Group: "otros ingresos - punto de venta", Environment: "point_sale"},
	{Code: "INOPS05", Name: "obtener por fecha", Details: "Obtener ingresos otros por rango de fechas", Group: "otros ingresos - punto de venta", Environment: "point_sale"},
	//INCOME SALE
	{Code: "INS01", Name: "crear", Details: "Crear un nuevo ingreso de venta", Group: "ingreso de venta", Environment: "dashboard"},
	{Code: "INS02", Name: "actualizar", Details: "Actualizar un ingreso de venta existente", Group: "ingreso de venta", Environment: "dashboard"},
	{Code: "INS03", Name: "eliminar", Details: "Eliminar un ingreso de venta", Group: "ingreso de venta", Environment: "dashboard"},
	{Code: "INS04", Name: "obtener uno en específico", Details: "Obtener un ingreso de venta en específico", Group: "ingreso de venta", Environment: "dashboard"},
	{Code: "INS05", Name: "obtener por fecha", Details: "Obtener ingresos de venta por rango de fechas", Group: "ingreso de venta", Environment: "dashboard"},
	//MEMBER
	{Code: "MB01", Name: "crear", Details: "Crear un nuevo miembro", Group: "miembro", Environment: "dashboard"},
	{Code: "MB02", Name: "actualizar", Details: "Actualizar un miembro existente", Group: "miembro", Environment: "dashboard"},
	{Code: "MB04", Name: "obtener uno en específico", Details: "Obtener un miembro en específico", Group: "miembro", Environment: "dashboard"},
	{Code: "MB05", Name: "obtener todos", Details: "Obtener todos los miembros", Group: "miembro", Environment: "dashboard"},
	//MOVEMENT STOCK
	{Code: "MS01", Name: "mover", Details: "Mover stock de productos", Group: "movimiento de stock", Environment: "dashboard"},
	{Code: "MS02", Name: "obtener uno en específico", Details: "Obtener un movimiento en específico", Group: "movimiento de stock", Environment: "dashboard"},
	{Code: "MS03", Name: "obtener por fecha", Details: "Obtener movimientos por rango de fechas", Group: "movimiento de stock", Environment: "dashboard"},
	//PERMISSION
	{Code: "PER01", Name: "obtener todos", Details: "Obtener todos los permisos", Group: "permisos", Environment: "dashboard"},
	{Code: "PER02", Name: "obtener propios", Details: "Obtener los permisos propios del usuario", Group: "permisos", Environment: "dashboard"},
	//POINT SALE
	{Code: "PS01", Name: "crear", Details: "Crear un nuevo punto de venta", Group: "punto de venta", Environment: "dashboard"},
	{Code: "PS02", Name: "actualizar", Details: "Actualizar un punto de venta existente", Group: "punto de venta", Environment: "dashboard"},
	{Code: "PS03", Name: "obtener miembros", Details: "Obtener los miembros de un punto de venta", Group: "punto de venta", Environment: "dashboard"},
	{Code: "PS04", Name: "obtener todos", Details: "Obtener todos los puntos de venta", Group: "punto de venta", Environment: "dashboard"},
	//PRODUCT
	{Code: "PR01", Name: "crear", Details: "Crear un nuevo producto", Group: "producto", Environment: "dashboard"},
	{Code: "PR02", Name: "actualizar", Details: "Actualizar un producto existente", Group: "producto", Environment: "dashboard"},
	{Code: "PR03", Name: "actualizar lista de precios", Details: "Actualizar la lista de precios de productos", Group: "producto", Environment: "dashboard"},
	{Code: "PR04", Name: "eliminar", Details: "Eliminar un producto", Group: "producto", Environment: "dashboard"},
	{Code: "PR05", Name: "obtener uno en específico", Details: "Obtener un producto en específico", Group: "producto", Environment: "dashboard"},
	{Code: "PR06", Name: "obtener por nombre", Details: "Obtener productos por nombre", Group: "producto", Environment: "dashboard"},
	{Code: "PR07", Name: "obtener por código", Details: "Obtener producto por codigo", Group: "producto", Environment: "dashboard"},
	{Code: "PR08", Name: "obtener por categoría", Details: "Obtener productos por categoría", Group: "producto", Environment: "dashboard"},
	{Code: "PR09", Name: "obtener todos", Details: "Obtener todos los productos", Group: "producto", Environment: "dashboard"},
	//REPORT
	{Code: "RP01", Name: "obtener excel", Details: "Obtener reporte en excel", Group: "report", Environment: "dashboard"},
	{Code: "RP02", Name: "obtener rentabilidad de productos", Details: "Obtener rentabilidad de productos", Group: "report", Environment: "dashboard"},
	{Code: "RP03", Name: "obtener por fecha", Details: "Obtener reporte por rango de fechas", Group: "report", Environment: "dashboard"},
	//ROLE
	{Code: "RL01", Name: "crear", Details: "Crear un nuevo rol", Group: "rol", Environment: "dashboard"},
	{Code: "RL02", Name: "obtener todos", Details: "Obtener todos los roles", Group: "rol", Environment: "dashboard"},
	//STOCK POINT SALE
	{Code: "ST01", Name: "crear", Details: "Crear nuevo stock de producto", Group: "stock", Environment: "point_sale"},
	{Code: "ST02", Name: "obtener uno en específico", Details: "Obtener stock de un producto en específico", Group: "stock", Environment: "point_sale"},
	{Code: "ST03", Name: "obtener por código", Details: "Obtener stock por codigo de producto", Group: "stock", Environment: "point_sale"},
	{Code: "ST04", Name: "obtener por categoría", Details: "Obtener stock por categoría de producto", Group: "stock", Environment: "point_sale"},
	{Code: "ST05", Name: "obtener todos", Details: "Obtener todos los stocks de productos", Group: "stock", Environment: "point_sale"},
	//SUPPLIER
	{Code: "SP01", Name: "crear", Details: "Crear un nuevo proveedor", Group: "proveedor", Environment: "dashboard"},
	{Code: "SP02", Name: "actualizar", Details: "Actualizar un proveedor existente", Group: "proveedor", Environment: "dashboard"},
	{Code: "SP03", Name: "eliminar", Details: "Eliminar un proveedor", Group: "proveedor", Environment: "dashboard"},
	{Code: "SP04", Name: "obtener uno en específico", Details: "Obtener un proveedor en específico", Group: "proveedor", Environment: "dashboard"},
	{Code: "SP05", Name: "obtener todos", Details: "Obtener todos los proveedores", Group: "proveedor", Environment: "dashboard"},
	//TYPE MOVEMENT
	{Code: "TM01", Name: "crear", Details: "Crear un nuevo tipo de movimiento", Group: "tipo de movimiento", Environment: "dashboard"},
	{Code: "TM02", Name: "actualizar", Details: "Actualizar un tipo de movimiento existente", Group: "tipo de movimiento", Environment: "dashboard"},
	{Code: "TM03", Name: "obtener todos", Details: "Obtener todos los tipos de movimiento", Group: "tipo de movimiento", Environment: "dashboard"},
}
