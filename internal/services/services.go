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

type ClientService struct {
	ClientRepository ports.ClientRepository
}

type ExpenseService struct {
	ExpenseRepository ports.ExpenseRepository
}

type IncomeService struct {
	IncomeRepository ports.IncomeRepository
}

type MemberService struct {
	MemberRepository ports.MemberRepository
	// UserRepository ports.UserRepository
}

type MovementTypeService struct {
	MovementTypeRepository ports.MovementTypeRepository
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

type RoleService struct {
	RoleRepository ports.RoleRepository
}

type SupplierService struct {
	SupplierRepository ports.SupplierRepository
}

type TenantService struct {
	UserRepository   ports.UserRepository
	TenantRepository ports.TenantRepository
}

type UserService struct {
	UserRepository ports.UserRepository
}
