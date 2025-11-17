package controllers

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/ports"


type AuthController struct {
	AuthService ports.AuhtService
}

type CashRegisterController struct {
	CashRegisterService ports.CashRegisterService
}

type CategoryController struct {
	CategoryService ports.CategoryService
}

type ClientController struct {
	ClientService ports.ClientService
}

type DepositController struct {
	DepositService ports.DepositService
}

type ExpenseBuyController struct {
	ExpenseBuyService ports.ExpenseBuyService
}

type ExpenseOtherController struct {
	ExpenseOtherService ports.ExpenseOtherService
}

type IncomeOtherController struct {
	IncomeOtherService ports.IncomeOtherService
}

type IncomeSaleController struct {
	IncomeSaleService ports.IncomeSaleService
}

type MemberController struct {
	MemberService ports.MemberService
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

type ReportController struct {
	ReportService ports.ReportService
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

type TypeMovementController struct {
	TypeMovementService ports.TypeMovementService
}

type UserController struct {
	UserService ports.UserService
}

