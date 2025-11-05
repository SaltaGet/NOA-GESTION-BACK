package controllers

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/ports"


type AuthController struct {
	AuthService ports.AuhtService
}

type ClientController struct {
	ClientService ports.ClientService
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

type PointSaleController struct {
	PointSaleService ports.PointSaleService
}

type ProductController struct {
	ProductService ports.ProductService
}

type RoleController struct {
	RoleService ports.RoleService
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

