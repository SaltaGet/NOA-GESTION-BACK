package services

import (
	"github.com/DanielChachagua/GestionCar/pkg/ports"
)

type AttendanceService struct {
	AttendanceRepository ports.AttendanceRespository
}

type AuthService struct {
	AuthRepository ports.AuthRepository
	UserRepository ports.UserRepository
	TenantService ports.TenantService
}

type ClientService struct {
	ClientRepository ports.ClientRepository
}

type EmployeeService struct {
	EmployeeRepository ports.EmployeeRepository
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

type ProductService struct {
	ProductRepository ports.ProductRepository
}

type PurchaseOrderService struct{
	PurchaseOrderRepository ports.PurchaseOrderRepository
}

type PurchaseProductService struct{
	PurchaseProductRepository ports.PurchaseProductRepository
}

type ResumeService struct {
	ResumeExpenseRepository ports.ResumeExpenseRepository
	ResumeIncomeRepository ports.ResumeIncomeRepository
}

type RoleService struct {
	RoleRepository ports.RoleRepository
}

type ServiceService struct {
	ServiceRepository ports.ServiceRepository
}

type SupplierService struct {
	SupplierRepository ports.SupplierRepository
}

type TenantService struct {
	UserRepository ports.UserRepository
	TenantRepository ports.TenantRepository
}

type UserService struct {
	UserRepository ports.UserRepository
}

type VehicleService struct {
	VehicleRepository ports.VehicleRepository
}
