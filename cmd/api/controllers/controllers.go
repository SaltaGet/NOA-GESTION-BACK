package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/ports"
	"github.com/alexandrevicenzi/go-sse"
)

type AuthController struct {
	AuthService  ports.AuthService
	EmailService ports.EmailService
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

type FeedbackController struct {
	FeedbackService ports.FeedbackServices
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

type ModuleController struct {
	ModuleService ports.ModuleService
}

type MovementStockController struct {
	MovementStockService   ports.MovementStockService
	NotificationController *NotificationController
}

type NewsController struct {
	NewsService ports.NewsServices
}

type NotificationController struct {
	NotificationService ports.NotificationService
	SSEServer           *sse.Server
}

type PermissionController struct {
	PermissionService ports.PermissionService
}

type PlanController struct {
	PlanService ports.PlanService
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

type StockController struct {
	StockService ports.StockService
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
