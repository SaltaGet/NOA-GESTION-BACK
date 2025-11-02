package repositories

import (
	"gorm.io/gorm"
)

// type Repository struct {
// 	DB *gorm.DB
// }

type MainRepository struct {
	DB *gorm.DB
}

// type TenantRepository struct {
// 	DB *gorm.DB
// }

//	func NewTenantRepository(db *gorm.DB) *TenantRepository {
//	    return &TenantRepository{DB: db}
//	}

type AttendanceRepository struct {
	DB *gorm.DB
}

type ClientRepository struct {
	DB *gorm.DB
}

type EmployeeRepository struct {
	DB *gorm.DB
}

type ExpenseRepository struct {
	DB *gorm.DB
}

type IncomeRepository struct {
	DB *gorm.DB
}

type MemberRepository struct {
	DB *gorm.DB
}

type MovementTypeRepository struct {
	DB *gorm.DB
}

type PermissionRepository struct {
	DB *gorm.DB
}

type ProductRepository struct {
	DB *gorm.DB
}

type PurchaseOrderRepository struct {
	DB *gorm.DB
}

type PurchaseProductRepository struct {
	DB *gorm.DB
}

type ResumeRepository struct {
	DB *gorm.DB
}

type RoleRepository struct {
	DB *gorm.DB
}

type ServiceRepository struct {
	DB *gorm.DB
}

type SupplierRepository struct {
	DB *gorm.DB
}

type VehicleRepository struct {
	DB *gorm.DB
}

