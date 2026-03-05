package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib" // Driver para sql.Open
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PrepareDB(uri string, memberAdmin models.Member) error {
	var driver gorm.Dialector
	if err := EnsureDatabaseExists(uri); err != nil {
		return fmt.Errorf("error al crear la base: %w", err)
	}
	driver = postgres.Open(uri)

	// Conexión GORM
	db, err := gorm.Open(driver, &gorm.Config{})
	if err != nil {
		handleDBCreationError(uri)
		return fmt.Errorf("error al inicializar DB: %w", err)
	}

	// Bajo nivel
	sqlDB, err := db.DB()
	if err != nil {
		handleDBCreationError(uri)
		return fmt.Errorf("error al obtener conexión de bajo nivel: %w", err)
	}
	defer sqlDB.Close()

	// Migraciones
	if err := db.AutoMigrate(
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
		&models.IncomeEcommerce{},
		&models.IncomeEcommerceItem{},
		&models.Invoice{},
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
		handleDBCreationError(uri)
		return fmt.Errorf("error al migrar tablas: %w", err)
	}

	if err := db.Create(&models.Role{ID: 1, Name: "admin"}).Error; err != nil {
		handleDBCreationError(uri)
		return fmt.Errorf("error al crear rol admin: %w", err)
	}

	if err := db.Create(&memberAdmin).Error; err != nil {
		handleDBCreationError(uri)
		return fmt.Errorf("error al crear member admin: %w", err)
	}

	if err := db.Create(&Permissions).Error; err != nil {
		handleDBCreationError(uri)
		return fmt.Errorf("error al migrar permisos: %w", err)
	}

	responsabilityFrontIVA := "consumidor_final"
	if err := db.Create(&models.Client{
		ID:                     1,
		FirstName:              "Consumidor",
		LastName:               "Final",
		CompanyName:            nil,
		Identifier:             nil,
		Email:                  nil,
		Phone:                  nil,
		Address:                nil,
		ResponsabilityFrontIVA: &responsabilityFrontIVA,
		MemberCreateID:         memberAdmin.ID,
	}).Error; err != nil {
		handleDBCreationError(uri)
		return fmt.Errorf("error al crear cliente consumidor final: %w", err)
	}

	if err := db.Create(&models.Category{
		Name: "Sin categoría",
	}).Error; err != nil {
		handleDBCreationError(uri)
		return fmt.Errorf("error al crear categoría sin categoría: %w", err)
	}

	if err := db.Create(&models.Supplier{
		ID:          1,
		Name:        "Sin proveedor",
		CompanyName: "Sin nombre",
	}).Error; err != nil {
		handleDBCreationError(uri)
		return fmt.Errorf("error al crear proveedor sin proveedor: %w", err)
	}

	if err := db.Create(&models.TypeExpense{
		ID:   1,
		Name: "Otros",
	}).Error; err != nil {
		handleDBCreationError(uri)
		return fmt.Errorf("error al crear tipo de gasto sin tipo de gasto: %w", err)
	}

	if err := db.Create(&models.TypeIncome{
		ID:   1,
		Name: "Otros",
	}).Error; err != nil {
		handleDBCreationError(uri)
		return fmt.Errorf("error al crear tipo de gasto sin tipo de gasto: %w", err)
	}

	description := "Mi primer punto de venta de NOA Gestión"
	pointSales := []models.PointSale{
		{
			ID:          1,
			Name:        "Mi punto de venta",
			Description: &description,
			Number:      1,
			IsDeposit:   true,
			IsMain:      true,
		},
	}
	if err := db.Create(&pointSales).Error; err != nil {
		handleDBCreationError(uri)
		return fmt.Errorf("error al crear punto de venta inicial: %w", err)
	}

	if err := db.Model(&memberAdmin).Association("PointSales").Append(&pointSales); err != nil {
		return fmt.Errorf("error al crear relacion entre miembro y punto de venta: %w", err)
	}

	tables := []string{"roles", "members", "clients", "suppliers", "type_expenses", "type_incomes", "point_sales"}
	for _, table := range tables {
		query := fmt.Sprintf("SELECT setval(pg_get_serial_sequence('%s', 'id'), COALESCE(MAX(id), 1)) FROM %s", table, table)
		if err := db.Exec(query).Error; err != nil {
			return fmt.Errorf("error sincronizando secuencia de %s: %w", table, err)
		}
	}

	if err := ApplyAuditTriggers(db); err != nil {
		return fmt.Errorf("error al aplicar triggers de auditoria: %w", err)
	}

	return nil
}

func handleDBCreationError(uri string) {
	_ = dropDatabase(uri)
}

func dropDatabase(dsn string) error {
	dbName, err := extractDBName(dsn)
	if err != nil {
		return fmt.Errorf("no se pudo extraer el nombre de la base: %w", err)
	}

	baseDSN := removeDBFromDSN(dsn)
	sqlDB, err := sql.Open("pgx", baseDSN)
	if err != nil {
		return fmt.Errorf("no se pudo conectar al servidor Postgres: %w", err)
	}
	defer sqlDB.Close()

	_, err = sqlDB.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, dbName))
	if err != nil {
		return fmt.Errorf("error al ejecutar DROP DATABASE: %w", err)
	}
	return nil
}

func extractDBName(dsn string) (string, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return "", fmt.Errorf("error al parsear DSN: %w", err)
	}
	dbName := strings.TrimPrefix(u.Path, "/")
	if dbName == "" {
		return "", fmt.Errorf("formato de DSN inválido (no se encontró la base)")
	}
	return dbName, nil
}

func removeDBFromDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		return dsn
	}
	u.Path = "/postgres"
	return u.String()
}

func UpdateModels(uri string) error {
	var driver gorm.Dialector
	if err := EnsureDatabaseExists(uri); err != nil {
		return fmt.Errorf("error al crear la base: %w", err)
	}
	driver = postgres.Open(uri)

	db, err := gorm.Open(driver, &gorm.Config{})
	if err != nil {
		handleDBCreationError(uri)
		return fmt.Errorf("error al inicializar DB: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		handleDBCreationError(uri)
		return fmt.Errorf("error al obtener conexión de bajo nivel: %w", err)
	}
	defer sqlDB.Close()

	if err := db.AutoMigrate(
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
		&models.IncomeEcommerce{},
		&models.IncomeEcommerceItem{},
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
		handleDBCreationError(uri)
		return fmt.Errorf("error al migrar tablas: %w", err)
	}

	tables := []string{"roles", "members", "clients", "suppliers", "type_expenses", "type_incomes", "point_sales"}
	for _, table := range tables {
		query := fmt.Sprintf("SELECT setval(pg_get_serial_sequence('%s', 'id'), COALESCE(MAX(id), 1)) FROM %s", table, table)
		if err := db.Exec(query).Error; err != nil {
			return fmt.Errorf("error sincronizando secuencia de %s: %w", table, err)
		}
	}

	return nil
}

var Permissions []models.Permission = []models.Permission{
	//CASH REGISTER
	{Code: "CR01", Name: "apertura y cierre de caja", Details: "Apertura y cierre de caja del punto de venta", Group: "caja", Environment: "point_sale"},
	{Code: "CR04", Name: "informes", Details: "Obtener infomes de caja", Group: "caja", Environment: "point_sale"},

	//CATEGORY
	{Code: "CAT01", Name: "crear", Details: "Crear nueva categoría", Group: "categoría", Environment: "dashboard"},
	{Code: "CAT02", Name: "actualizar", Details: "Actualizar categoria existente", Group: "categoría", Environment: "dashboard"},
	{Code: "CAT03", Name: "eliminar", Details: "Eliminar cagotería", Group: "categoría", Environment: "dashboard"},

	//CLIENT
	{Code: "CL01", Name: "crear", Details: "Crear nuevo cliente", Group: "cliente", Environment: "dashboard"},
	{Code: "CL02", Name: "actualizar", Details: "Actualizar cliente existente", Group: "cliente", Environment: "dashboard"},
	{Code: "CL03", Name: "eliminar", Details: "Eliminar clientes", Group: "cliente", Environment: "dashboard"},
	// {Code: "CL04", Name: "obtener clientes", Details: "Obtener clientes", Group: "cliente", Environment: "point_sale"},

	//DEPOSIT
	{Code: "DEP02", Name: "actualizar stock de productos", Details: "Actualizar el stock de productos en el depósito", Group: "deposito", Environment: "dashboard"},
	{Code: "DEP04", Name: "obtener depópsito", Details: "Obtener información de los productos del depósito", Group: "deposito", Environment: "dashboard"},

	//EXPENSE BUY
	{Code: "EB01", Name: "crear", Details: "Crear  nuevo gasto de compra", Group: "gastos de compra", Environment: "dashboard"},
	{Code: "EB02", Name: "actualizar", Details: "Actualizar  gasto de compra existente", Group: "gastos de compra", Environment: "dashboard"},
	{Code: "EB03", Name: "eliminar", Details: "Eliminar  gasto de compra", Group: "gastos de compra", Environment: "dashboard"},
	{Code: "EB04", Name: "obtener gatos de compra", Details: "Obtener gastos de compra", Group: "gastos de compra", Environment: "dashboard"},

	//EXPENSE OTHER
	{Code: "EO01", Name: "crear", Details: "Crear un nuevo otros gastos", Group: "otros gastos", Environment: "dashboard"},
	{Code: "EO02", Name: "actualizar", Details: "Actualizar un otros gastos existente", Group: "otros gastos", Environment: "dashboard"},
	{Code: "EO03", Name: "eliminar", Details: "Eliminar otros gastos", Group: "otros gastos", Environment: "dashboard"},
	{Code: "EO04", Name: "obtener otros gastos", Details: "Obtener un otros gastos", Group: "otros gastos", Environment: "dashboard"},

	//EXPENSE OTHER POINT SALE
	{Code: "EOPS01", Name: "crear", Details: "Crear otro gasto para un punto de venta", Group: "otros grastos - punto de venta", Environment: "point_sale"},
	{Code: "EOPS02", Name: "actualizar", Details: "Actualizar gasto de un punto de venta", Group: "otros grastos - punto de venta", Environment: "point_sale"},
	{Code: "EOPS03", Name: "eliminar", Details: "Eliminar gasto de un punto de venta", Group: "otros grastos - punto de venta", Environment: "point_sale"},
	{Code: "EOPS04", Name: "obtener otros gastos", Details: "Obtener otros gastos de un punto de venta", Group: "otros grastos - punto de venta", Environment: "point_sale"},

	//INCOME OTHER
	{Code: "INO01", Name: "crear", Details: "Crear otro ingreso", Group: "otros ingresos", Environment: "dashboard"},
	{Code: "INO02", Name: "actualizar", Details: "Actualizar ingreso existente", Group: "otros ingresos", Environment: "dashboard"},
	{Code: "INO03", Name: "eliminar", Details: "Eliminar ingreso", Group: "otros ingresos", Environment: "dashboard"},
	{Code: "INO04", Name: "obtener otros ingresos", Details: "Obtener otros ingresos", Group: "otros ingresos", Environment: "dashboard"},

	//INCOME OTHER POINT SALE
	{Code: "INOPS01", Name: "crear", Details: "Crear otro ingreso para un punto de venta", Group: "otros ingresos - punto de venta", Environment: "point_sale"},
	{Code: "INOPS02", Name: "actualizar", Details: "Actualizar ingreso de un punto de venta", Group: "otros ingresos - punto de venta", Environment: "point_sale"},
	{Code: "INOPS03", Name: "eliminar", Details: "Eliminar ingreso de un punto de venta", Group: "otros ingresos - punto de venta", Environment: "point_sale"},
	{Code: "INOPS04", Name: "obtener otros ingresos", Details: "Obtener otros ingresos de un punto de venta", Group: "otros ingresos - punto de venta", Environment: "point_sale"},

	//INCOME SALE
	{Code: "INS01", Name: "crear", Details: "Crear nuevo ingreso de venta", Group: "ingreso de venta", Environment: "point_sale"},
	{Code: "INS02", Name: "actualizar", Details: "Actualizar ingreso de venta existente", Group: "ingreso de venta", Environment: "point_sale"},
	{Code: "INS03", Name: "eliminar", Details: "Eliminar ingreso de venta", Group: "ingreso de venta", Environment: "point_sale"},
	{Code: "INS04", Name: "obtener ingresos de ventas", Details: "Obtener ingresos de ventas", Group: "ingreso de venta", Environment: "point_sale"},

	//MEMBER
	{Code: "MB01", Name: "crear", Details: "Crear un nuevo miembro", Group: "miembro", Environment: "dashboard"},
	{Code: "MB02", Name: "actualizar", Details: "Actualizar un miembro existente", Group: "miembro", Environment: "dashboard"},
	{Code: "MB03", Name: "eliminar", Details: "Eliminar miembro existente", Group: "miembro", Environment: "dashboard"},
	{Code: "MB04", Name: "obtener miembros", Details: "Obtener miembros", Group: "miembro", Environment: "dashboard"},

	//MOVEMENT STOCK
	{Code: "MS02", Name: "mover", Details: "Mover stock de productos", Group: "movimiento de stock", Environment: "dashboard"},
	{Code: "MS04", Name: "obtener movimientos", Details: "Obtener movimientos de stock", Group: "movimiento de stock", Environment: "dashboard"},

	//PERMISSION
	{Code: "PER04", Name: "obtener permisos", Details: "Obtener permisos", Group: "permisos", Environment: "dashboard"},

	//POINT SALE
	{Code: "PS01", Name: "crear", Details: "Crear un nuevo punto de venta", Group: "punto de venta", Environment: "dashboard"},
	{Code: "PS02", Name: "actualizar", Details: "Actualizar un punto de venta existente", Group: "punto de venta", Environment: "dashboard"},
	{Code: "PS04", Name: "obtener puntos de ventas", Details: "Obtener puntos de ventas", Group: "punto de venta", Environment: "dashboard"},

	//PRODUCT
	{Code: "PR01", Name: "crear", Details: "Crear un nuevo producto", Group: "producto", Environment: "dashboard"},
	{Code: "PR02", Name: "actualizar", Details: "Actualizar un producto existente", Group: "producto", Environment: "dashboard"},
	{Code: "PR03", Name: "eliminar", Details: "Eliminar un producto", Group: "producto", Environment: "dashboard"},
	// {Code: "PR04", Name: "obtener uno en específico", Details: "Obtener un producto en específico", Group: "producto", Environment: "dashboard"},

	//REPORT
	{Code: "RP01", Name: "obtener reportes", Details: "Obtener reporte de productos y balances", Group: "report", Environment: "dashboard"},
	// {Code: "RP02", Name: "obtener rentabilidad de productos", Details: "Obtener rentabilidad de productos", Group: "report", Environment: "dashboard"},
	// {Code: "RP03", Name: "obtener por fecha", Details: "Obtener reporte por rango de fechas", Group: "report", Environment: "dashboard"},

	//ROLE
	{Code: "RL01", Name: "crear", Details: "Crear un nuevo rol", Group: "rol", Environment: "dashboard"},
	{Code: "RL02", Name: "editar", Details: "Editar rol existente", Group: "rol", Environment: "dashboard"},
	{Code: "RL04", Name: "obtener todos", Details: "Obtener todos los roles", Group: "rol", Environment: "dashboard"},

	// //STOCK POINT SALE
	// {Code: "ST01", Name: "crear", Details: "Crear nuevo stock de producto", Group: "stock", Environment: "point_sale"},
	// {Code: "ST04", Name: "obtener stock", Details: "Obtener stock de un producto en punto de venta", Group: "stock", Environment: "point_sale"},

	//SUPPLIER
	{Code: "SP01", Name: "crear", Details: "Crear un nuevo proveedor", Group: "proveedor", Environment: "dashboard"},
	{Code: "SP02", Name: "actualizar", Details: "Actualizar un proveedor existente", Group: "proveedor", Environment: "dashboard"},
	{Code: "SP03", Name: "eliminar", Details: "Eliminar un proveedor", Group: "proveedor", Environment: "dashboard"},
	{Code: "SP04", Name: "obtener proveedores", Details: "Obtener proveedores", Group: "proveedor", Environment: "dashboard"},

	//TYPE MOVEMENT
	{Code: "TM01", Name: "crear", Details: "Crear un nuevo tipo de movimiento", Group: "tipo de movimiento", Environment: "dashboard"},
	{Code: "TM02", Name: "actualizar", Details: "Actualizar un tipo de movimiento existente", Group: "tipo de movimiento", Environment: "dashboard"},
}

var TriggerAudit = `
CREATE OR REPLACE FUNCTION audit_trigger_function()
RETURNS TRIGGER AS $$
DECLARE
    current_member TEXT;
		current_tx_id BIGINT;
BEGIN
    current_member := current_setting('app.current_member_id', true);
		current_tx_id := txid_current();

    IF current_member IS NULL OR current_member = '' OR current_member = '0' THEN
        RETURN NEW;
    END IF;

    INSERT INTO audit_logs (
				transaction_id,
        member_id,
        method,
        path,
        old_value,
        new_value,
        created_at
    )
    VALUES (
				current_tx_id,
        current_member::BIGINT,
        LOWER(TG_OP),
        TG_TABLE_NAME,
        CASE WHEN TG_OP = 'INSERT' THEN NULL ELSE to_jsonb(OLD) END,
        CASE WHEN TG_OP = 'DELETE' THEN NULL ELSE to_jsonb(NEW) END,
        NOW()
    );

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
`

func ApplyAuditTriggers(db *gorm.DB) error {
	// 1. Crear la función del trigger
	if err := db.Exec(TriggerAudit).Error; err != nil {
		return err
	}

	// 2. Obtener TODAS las tablas del esquema público
	var tables []string
	queryTables := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_type = 'BASE TABLE'
		AND table_name != 'audit_logs' 
		AND table_name != 'migrations';
	`
	if err := db.Raw(queryTables).Scan(&tables).Error; err != nil {
		return fmt.Errorf("error al obtener lista de tablas: %w", err)
	}

	// 3. Aplicar el trigger a cada una
	for _, table := range tables {
		db.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS tr_audit_%s ON %s", table, table))

		createQuery := fmt.Sprintf(`
      CREATE TRIGGER tr_audit_%s
      AFTER INSERT OR UPDATE OR DELETE ON "%s"
      FOR EACH ROW EXECUTE FUNCTION audit_trigger_function();`,
			table, table)

		if err := db.Exec(createQuery).Error; err != nil {
			return fmt.Errorf("error al crear trigger en tabla %s: %w", table, err)
		}
	}
	return nil
}

// package database

// import (
// 	"strings"

// 	"fmt"

// 	"database/sql"

// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// func PrepareDB(uri string, memberAdmin models.Member) error {
// 	var driver gorm.Dialector
// 	if err := EnsureDatabaseExists(uri); err != nil {
// 		return fmt.Errorf("error al crear la base: %w", err)
// 	}
// 	driver = mysql.Open(uri)

// 	// Conexión GORM
// 	db, err := gorm.Open(driver, &gorm.Config{})
// 	if err != nil {
// 		handleDBCreationError(uri)
// 		return fmt.Errorf("error al inicializar DB: %w", err)
// 	}

// 	// Bajo nivel
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		handleDBCreationError(uri)
// 		return fmt.Errorf("error al obtener conexión de bajo nivel: %w", err)
// 	}
// 	defer sqlDB.Close()

// 	// Migraciones
// 	if err := db.AutoMigrate(
// 		&models.AuditLog{},
// 		&models.CashRegister{},
// 		&models.Category{},
// 		&models.Client{},
// 		&models.Deposit{},

// 		&models.ExpenseBuy{},
// 		&models.ExpenseBuyItem{},
// 		&models.ExpenseOther{},
// 		&models.PayExpenseBuy{},
// 		&models.PayExpenseOther{},
// 		&models.TypeExpense{},

// 		&models.IncomeSale{},
// 		&models.IncomeSaleItem{},
// 		&models.IncomeEcommerce{},
// 		&models.IncomeEcommerceItem{},
// 		&models.PayIncome{},
// 		&models.IncomeOther{},
// 		&models.TypeIncome{},

// 		&models.Member{},
// 		&models.MovementStock{},
// 		&models.Permission{},
// 		&models.PointSale{},
// 		&models.Product{},
// 		&models.Role{},
// 		&models.StockPointSale{},
// 		&models.Supplier{},
// 	); err != nil {
// 		handleDBCreationError(uri)
// 		return fmt.Errorf("error al migrar tablas: %w", err)
// 	}

// 	if err := db.Create(&models.Role{ID: 1, Name: "admin"}).Error; err != nil {
// 		handleDBCreationError(uri)
// 		return fmt.Errorf("error al crear rol admin: %w", err)
// 	}

// 	if err := db.Create(&memberAdmin).Error; err != nil {
// 		handleDBCreationError(uri)
// 		return fmt.Errorf("error al crear member admin: %w", err)
// 	}

// 	if err := db.Create(&Permissions).Error; err != nil {
// 		handleDBCreationError(uri)
// 		return fmt.Errorf("error al migrar permisos: %w", err)
// 	}

// 	if err := db.Create(&models.Client{
// 		ID:             1,
// 		FirstName:      "Consumidor",
// 		LastName:       "Final",
// 		CompanyName:    nil,
// 		Identifier:     nil,
// 		Email:          nil,
// 		Phone:          nil,
// 		Address:        nil,
// 		MemberCreateID: memberAdmin.ID,
// 	}).Error; err != nil {
// 		handleDBCreationError(uri)
// 		return fmt.Errorf("error al crear cliente consumidor final: %w", err)
// 	}

// 	if err := db.Create(&models.Category{
// 		Name: "Sin categoría",
// 	}).Error; err != nil {
// 		handleDBCreationError(uri)
// 		return fmt.Errorf("error al crear categoría sin categoría: %w", err)
// 	}

// 	if err := db.Create(&models.Supplier{
// 		ID:          1,
// 		Name:        "Sin proveedor",
// 		CompanyName: "Sin nombre",
// 	}).Error; err != nil {
// 		handleDBCreationError(uri)
// 		return fmt.Errorf("error al crear proveedor sin proveedor: %w", err)
// 	}

// 	if err := db.Create(&models.TypeExpense{
// 		ID:   1,
// 		Name: "Otros",
// 	}).Error; err != nil {
// 		handleDBCreationError(uri)
// 		return fmt.Errorf("error al crear tipo de gasto sin tipo de gasto: %w", err)
// 	}

// 	if err := db.Create(&models.TypeIncome{
// 		ID:   1,
// 		Name: "Otros",
// 	}).Error; err != nil {
// 		handleDBCreationError(uri)
// 		return fmt.Errorf("error al crear tipo de gasto sin tipo de gasto: %w", err)
// 	}

// 	description := "Mi primer punto de venta de NOA Gestión"
// 	pointSales := []models.PointSale{
// 		{
// 			ID:          1,
// 			Name:        "Mi punto de venta",
// 			Description: &description,
// 			IsDeposit:   true,
// 			IsMain:      true,
// 		},
// 	}
// 	if err := db.Create(&pointSales).Error; err != nil {
// 		handleDBCreationError(uri)
// 		return fmt.Errorf("error al crear punto de venta inicial: %w", err)
// 	}

// 	if err := db.Model(&memberAdmin).Association("PointSales").Append(&pointSales); err != nil {
// 		return fmt.Errorf("error al crear relacion entre miembro y punto de venta: %w", err)
// 	}

// 	return nil
// }

// func handleDBCreationError(uri string) {
// 	_ = dropDatabase(uri)
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
// 	var driver gorm.Dialector
// 	if err := EnsureDatabaseExists(uri); err != nil {
// 		return fmt.Errorf("error al crear la base: %w", err)
// 	}
// 	driver = mysql.Open(uri)

// 	db, err := gorm.Open(driver, &gorm.Config{})
// 	if err != nil {
// 		handleDBCreationError(uri)
// 		return fmt.Errorf("error al inicializar DB: %w", err)
// 	}

// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		handleDBCreationError(uri)
// 		return fmt.Errorf("error al obtener conexión de bajo nivel: %w", err)
// 	}
// 	defer sqlDB.Close()

// 	if err := db.AutoMigrate(
// 		&models.AuditLog{},
// 		&models.CashRegister{},
// 		&models.Category{},
// 		&models.Client{},
// 		&models.Deposit{},

// 		&models.ExpenseBuy{},
// 		&models.ExpenseBuyItem{},
// 		&models.ExpenseOther{},
// 		&models.PayExpenseBuy{},
// 		&models.PayExpenseOther{},
// 		&models.TypeExpense{},

// 		&models.IncomeSale{},
// 		&models.IncomeSaleItem{},
// 		&models.IncomeEcommerce{},
// 		&models.IncomeEcommerceItem{},
// 		&models.PayIncome{},
// 		&models.IncomeOther{},
// 		&models.TypeIncome{},

// 		&models.Member{},
// 		&models.MovementStock{},
// 		&models.Permission{},
// 		&models.PointSale{},
// 		&models.Product{},
// 		&models.Role{},
// 		&models.StockPointSale{},
// 		&models.Supplier{},
// 	); err != nil {
// 		handleDBCreationError(uri)
// 		return fmt.Errorf("error al migrar tablas: %w", err)
// 	}

// 	return nil
// }

// var Permissions []models.Permission = []models.Permission{
// 	//CASH REGISTER
// 	{Code: "CR01", Name: "apertura y cierre de caja", Details: "Apertura y cierre de caja del punto de venta", Group: "caja", Environment: "point_sale"},
// 	{Code: "CR04", Name: "informes", Details: "Obtener infomes de caja", Group: "caja", Environment: "point_sale"},

// 	//CATEGORY
// 	{Code: "CAT01", Name: "crear", Details: "Crear nueva categoría", Group: "categoría", Environment: "dashboard"},
// 	{Code: "CAT02", Name: "actualizar", Details: "Actualizar categoria existente", Group: "categoría", Environment: "dashboard"},
// 	{Code: "CAT03", Name: "eliminar", Details: "Eliminar cagotería", Group: "categoría", Environment: "dashboard"},

// 	//CLIENT
// 	{Code: "CL01", Name: "crear", Details: "Crear nuevo cliente", Group: "cliente", Environment: "dashboard"},
// 	{Code: "CL02", Name: "actualizar", Details: "Actualizar cliente existente", Group: "cliente", Environment: "dashboard"},
// 	{Code: "CL03", Name: "eliminar", Details: "Eliminar clientes", Group: "cliente", Environment: "dashboard"},
// 	// {Code: "CL04", Name: "obtener clientes", Details: "Obtener clientes", Group: "cliente", Environment: "point_sale"},

// 	//DEPOSIT
// 	{Code: "DEP02", Name: "actualizar stock de productos", Details: "Actualizar el stock de productos en el depósito", Group: "deposito", Environment: "dashboard"},
// 	{Code: "DEP04", Name: "obtener depópsito", Details: "Obtener información de los productos del depósito", Group: "deposito", Environment: "dashboard"},

// 	//EXPENSE BUY
// 	{Code: "EB01", Name: "crear", Details: "Crear  nuevo gasto de compra", Group: "gastos de compra", Environment: "dashboard"},
// 	{Code: "EB02", Name: "actualizar", Details: "Actualizar  gasto de compra existente", Group: "gastos de compra", Environment: "dashboard"},
// 	{Code: "EB03", Name: "eliminar", Details: "Eliminar  gasto de compra", Group: "gastos de compra", Environment: "dashboard"},
// 	{Code: "EB04", Name: "obtener gatos de compra", Details: "Obtener gastos de compra", Group: "gastos de compra", Environment: "dashboard"},

// 	//EXPENSE OTHER
// 	{Code: "EO01", Name: "crear", Details: "Crear un nuevo otros gastos", Group: "otros gastos", Environment: "dashboard"},
// 	{Code: "EO02", Name: "actualizar", Details: "Actualizar un otros gastos existente", Group: "otros gastos", Environment: "dashboard"},
// 	{Code: "EO03", Name: "eliminar", Details: "Eliminar otros gastos", Group: "otros gastos", Environment: "dashboard"},
// 	{Code: "EO04", Name: "obtener otros gastos", Details: "Obtener un otros gastos", Group: "otros gastos", Environment: "dashboard"},

// 	//EXPENSE OTHER POINT SALE
// 	{Code: "EOPS01", Name: "crear", Details: "Crear otro gasto para un punto de venta", Group: "otros grastos - punto de venta", Environment: "point_sale"},
// 	{Code: "EOPS02", Name: "actualizar", Details: "Actualizar gasto de un punto de venta", Group: "otros grastos - punto de venta", Environment: "point_sale"},
// 	{Code: "EOPS03", Name: "eliminar", Details: "Eliminar gasto de un punto de venta", Group: "otros grastos - punto de venta", Environment: "point_sale"},
// 	{Code: "EOPS04", Name: "obtener otros gastos", Details: "Obtener otros gastos de un punto de venta", Group: "otros grastos - punto de venta", Environment: "point_sale"},

// 	//INCOME OTHER
// 	{Code: "INO01", Name: "crear", Details: "Crear otro ingreso", Group: "otros ingresos", Environment: "dashboard"},
// 	{Code: "INO02", Name: "actualizar", Details: "Actualizar ingreso existente", Group: "otros ingresos", Environment: "dashboard"},
// 	{Code: "INO03", Name: "eliminar", Details: "Eliminar ingreso", Group: "otros ingresos", Environment: "dashboard"},
// 	{Code: "INO04", Name: "obtener otros ingresos", Details: "Obtener otros ingresos", Group: "otros ingresos", Environment: "dashboard"},

// 	//INCOME OTHER POINT SALE
// 	{Code: "INOPS01", Name: "crear", Details: "Crear otro ingreso para un punto de venta", Group: "otros ingresos - punto de venta", Environment: "point_sale"},
// 	{Code: "INOPS02", Name: "actualizar", Details: "Actualizar ingreso de un punto de venta", Group: "otros ingresos - punto de venta", Environment: "point_sale"},
// 	{Code: "INOPS03", Name: "eliminar", Details: "Eliminar ingreso de un punto de venta", Group: "otros ingresos - punto de venta", Environment: "point_sale"},
// 	{Code: "INOPS04", Name: "obtener otros ingresos", Details: "Obtener otros ingresos de un punto de venta", Group: "otros ingresos - punto de venta", Environment: "point_sale"},

// 	//INCOME SALE
// 	{Code: "INS01", Name: "crear", Details: "Crear nuevo ingreso de venta", Group: "ingreso de venta", Environment: "point_sale"},
// 	{Code: "INS02", Name: "actualizar", Details: "Actualizar ingreso de venta existente", Group: "ingreso de venta", Environment: "point_sale"},
// 	{Code: "INS03", Name: "eliminar", Details: "Eliminar ingreso de venta", Group: "ingreso de venta", Environment: "point_sale"},
// 	{Code: "INS04", Name: "obtener ingresos de ventas", Details: "Obtener ingresos de ventas", Group: "ingreso de venta", Environment: "point_sale"},

// 	//MEMBER
// 	{Code: "MB01", Name: "crear", Details: "Crear un nuevo miembro", Group: "miembro", Environment: "dashboard"},
// 	{Code: "MB02", Name: "actualizar", Details: "Actualizar un miembro existente", Group: "miembro", Environment: "dashboard"},
// 	{Code: "MB03", Name: "eliminar", Details: "Eliminar miembro existente", Group: "miembro", Environment: "dashboard"},
// 	{Code: "MB04", Name: "obtener miembros", Details: "Obtener miembros", Group: "miembro", Environment: "dashboard"},

// 	//MOVEMENT STOCK
// 	{Code: "MS02", Name: "mover", Details: "Mover stock de productos", Group: "movimiento de stock", Environment: "dashboard"},
// 	{Code: "MS04", Name: "obtener movimientos", Details: "Obtener movimientos de stock", Group: "movimiento de stock", Environment: "dashboard"},

// 	//PERMISSION
// 	{Code: "PER04", Name: "obtener permisos", Details: "Obtener permisos", Group: "permisos", Environment: "dashboard"},

// 	//POINT SALE
// 	{Code: "PS01", Name: "crear", Details: "Crear un nuevo punto de venta", Group: "punto de venta", Environment: "dashboard"},
// 	{Code: "PS02", Name: "actualizar", Details: "Actualizar un punto de venta existente", Group: "punto de venta", Environment: "dashboard"},
// 	{Code: "PS04", Name: "obtener puntos de ventas", Details: "Obtener puntos de ventas", Group: "punto de venta", Environment: "dashboard"},

// 	//PRODUCT
// 	{Code: "PR01", Name: "crear", Details: "Crear un nuevo producto", Group: "producto", Environment: "dashboard"},
// 	{Code: "PR02", Name: "actualizar", Details: "Actualizar un producto existente", Group: "producto", Environment: "dashboard"},
// 	{Code: "PR03", Name: "eliminar", Details: "Eliminar un producto", Group: "producto", Environment: "dashboard"},
// 	// {Code: "PR04", Name: "obtener uno en específico", Details: "Obtener un producto en específico", Group: "producto", Environment: "dashboard"},

// 	//REPORT
// 	{Code: "RP01", Name: "obtener reportes", Details: "Obtener reporte de productos y balances", Group: "report", Environment: "dashboard"},
// 	// {Code: "RP02", Name: "obtener rentabilidad de productos", Details: "Obtener rentabilidad de productos", Group: "report", Environment: "dashboard"},
// 	// {Code: "RP03", Name: "obtener por fecha", Details: "Obtener reporte por rango de fechas", Group: "report", Environment: "dashboard"},

// 	//ROLE
// 	{Code: "RL01", Name: "crear", Details: "Crear un nuevo rol", Group: "rol", Environment: "dashboard"},
// 	{Code: "RL02", Name: "editar", Details: "Editar rol existente", Group: "rol", Environment: "dashboard"},
// 	{Code: "RL04", Name: "obtener todos", Details: "Obtener todos los roles", Group: "rol", Environment: "dashboard"},

// 	// //STOCK POINT SALE
// 	// {Code: "ST01", Name: "crear", Details: "Crear nuevo stock de producto", Group: "stock", Environment: "point_sale"},
// 	// {Code: "ST04", Name: "obtener stock", Details: "Obtener stock de un producto en punto de venta", Group: "stock", Environment: "point_sale"},

// 	//SUPPLIER
// 	{Code: "SP01", Name: "crear", Details: "Crear un nuevo proveedor", Group: "proveedor", Environment: "dashboard"},
// 	{Code: "SP02", Name: "actualizar", Details: "Actualizar un proveedor existente", Group: "proveedor", Environment: "dashboard"},
// 	{Code: "SP03", Name: "eliminar", Details: "Eliminar un proveedor", Group: "proveedor", Environment: "dashboard"},
// 	{Code: "SP04", Name: "obtener proveedores", Details: "Obtener proveedores", Group: "proveedor", Environment: "dashboard"},

// 	//TYPE MOVEMENT
// 	{Code: "TM01", Name: "crear", Details: "Crear un nuevo tipo de movimiento", Group: "tipo de movimiento", Environment: "dashboard"},
// 	{Code: "TM02", Name: "actualizar", Details: "Actualizar un tipo de movimiento existente", Group: "tipo de movimiento", Environment: "dashboard"},
// }
