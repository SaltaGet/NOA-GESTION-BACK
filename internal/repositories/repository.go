package repositories

import (
	"database/sql"
)

// type Repository struct {
// 	DB *sql.DB
// }

type MainRepository struct {
	DB *sql.DB
}

// type TenantRepository struct {
// 	DB *sql.DB
// }

//	func NewTenantRepository(db *sql.DB) *TenantRepository {
//	    return &TenantRepository{DB: db}
//	}

type ArcaRepository struct {
	DB *sql.DB
}

type CashRegisterRepository struct {
	DB *sql.DB
}

type CategoryRepository struct {
	DB *sql.DB
}

type ClientRepository struct {
	DB *sql.DB
}

type DepositRepository struct {
	DB *sql.DB
}

type EcommerceRepository struct {
	DB *sql.DB
}

type EmployeeRepository struct {
	DB *sql.DB
}

type ExpenseBuyRepository struct {
	DB *sql.DB
}

type ExpenseOtherRepository struct {
	DB *sql.DB
}

type IncomeSaleRepository struct {
	DB *sql.DB
}

type IncomeOtherRepository struct {
	DB *sql.DB
}

type MemberRepository struct {
	DB *sql.DB
}

type MovementStockRepository struct {
	DB *sql.DB
}

type MovementTypeRepository struct {
	DB *sql.DB
}

// type NotificationRepository struct {
// 	DB *sql.DB
// }

type PermissionRepository struct {
	DB *sql.DB
}

type PointSaleRepository struct {
	DB *sql.DB
}

type ProductRepository struct {
	DB *sql.DB
}

type PurchaseOrderRepository struct {
	DB *sql.DB
}

type PurchaseProductRepository struct {
	DB *sql.DB
}

type ReportRepository struct {
	DB *sql.DB
}

type RoleRepository struct {
	DB *sql.DB
}

type ServiceRepository struct {
	DB *sql.DB
}

type StockRepository struct {
	DB *sql.DB
}

type SupplierRepository struct {
	DB *sql.DB
}

type TypeMovementRepository struct {
	DB *sql.DB
}
