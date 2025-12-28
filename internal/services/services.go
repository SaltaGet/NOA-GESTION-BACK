package services

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/ports"
	"gopkg.in/gomail.v2"
)

type AuthService struct {
	AuthRepository   ports.AuthRepository
	UserRepository   ports.UserRepository
	TenantService    ports.TenantService
	EmailService     ports.EmailService
	PlanRepository   ports.PlanRepository
	ModuleRepository ports.ModuleRepository
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

type EmailService struct {
	Dialer *gomail.Dialer
}

type ExpenseBuyService struct {
	ExpenseBuyRepository ports.ExpenseBuyRepository
}

type ExpenseOtherService struct {
	ExpenseOtherRepository ports.ExpenseOtherRepository
}

type FeedbackService struct {
	FeedbackRepository ports.FeedbackRepository
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

type ModuleService struct {
	ModuleRepository ports.ModuleRepository
}

type MovementStockService struct {
	MovementStockRepository ports.MovementStockRepository
	NotifyService           ports.NotificationService
}

type NewsService struct {
	NewsRepository ports.NewsRepository
}

type NotificationService struct {
	NotificationRepository ports.NotificationRepository
}

type PermissionService struct {
	PermissionRepository ports.PermissionRepository
}

type PlanService struct {
	PlanRepository ports.PlanRepository
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
	EmailService     ports.EmailService
}

type TypeMovementService struct {
	TypeMovementRepository ports.TypeMovementRepository
}

type UserService struct {
	UserRepository ports.UserRepository
}
