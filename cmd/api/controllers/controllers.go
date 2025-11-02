package controllers

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/ports"


type AuthController struct {
	AuthService ports.AuhtService
}

type ClientController struct {
	ClientService ports.ClientService
}

type EmployeeController struct {
	EmployeeService ports.EmployeeService
}

type ExpenseController struct {
	ExpenseService ports.ExpenseService
}

type IncomeController struct {
	IncomeService ports.IncomeService
}

type MemberController struct {
	MemberService ports.MemberService
}

type MovementTypeController struct {
	MovementTypeService ports.MovementTypeService
}

type PermissionController struct {
	PermissionService ports.PermissionService
}

type ProductController struct {
	ProductService ports.ProductService
}

type PurchaseOrderController struct {
	PurchaseOrderService ports.PurchaseOrderService
}

type PurchaseProductController struct {
	PurchaseProductService ports.PurchaseProductService
}

type ResumeController struct {
	ResumeExpenseService ports.ResumeExpenseService
	ResumeIncomeService ports.ResumeIncomeService
}

type RoleController struct {
	RoleService ports.RoleService
}

type ServiceController struct {
	ServiceService ports.ServiceService
}

type SupplierController struct {
	SupplierService ports.SupplierService
}

type TenantController struct {
	TenantService ports.TenantService
}

type UserController struct {
	UserService ports.UserService
}

type VehicleController struct {
	VehicleService ports.VehicleService
}
