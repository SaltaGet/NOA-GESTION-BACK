package services

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/ports"

type AuthService struct {
	AuthRepository ports.AuthRepository
	UserRepository ports.UserRepository
	TenantService  ports.TenantService
}

type CashRegisterService struct {
	CashRegisterRepository ports.CashRegisterRepository
}

type CategoryService struct {
	CategoryRepository ports.CategoryRepository
}

type ClientService struct {
	ClientRepository ports.ClientRepository
}

type DepositService struct {
	DepositRepository ports.DepositRepository
}

type ExpenseBuyService struct {
	ExpenseBuyRepository ports.ExpenseBuyRepository
}

type ExpenseOtherService struct {
	ExpenseOtherRepository ports.ExpenseOtherRepository
}

type IncomeSaleService struct {
	IncomeSaleRepository ports.IncomeSaleRepository
}

type IncomeOtherService struct {
	IncomeOtherRepository ports.IncomeOtherRepository
}

type MemberService struct {
	MemberRepository ports.MemberRepository
	// UserRepository ports.UserRepository
}

type PermissionService struct {
	PermissionRepository ports.PermissionRepository
}

type PointSaleService struct {
	PointSaleRepository ports.PointSaleRepository
}

type ProductService struct {
	ProductRepository ports.ProductRepository
}

type ReportService struct {
	ReportRepository ports.ReportRepository
}

type RoleService struct {
	RoleRepository ports.RoleRepository
}

type StockService struct {
	StockRepository ports.StockRepository
}

type SupplierService struct {
	SupplierRepository ports.SupplierRepository
}

type TenantService struct {
	UserRepository   ports.UserRepository
	TenantRepository ports.TenantRepository
}

type TypeMovementService struct {
	TypeMovementRepository ports.TypeMovementRepository
}

type UserService struct {
	UserRepository ports.UserRepository
}
