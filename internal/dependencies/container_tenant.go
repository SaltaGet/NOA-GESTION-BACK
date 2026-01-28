// pkg/dependencies/container.go
package dependencies

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/repositories"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/services"
	"gorm.io/gorm"
)

type TenantContainer struct {
	DB          *gorm.DB
	Controllers struct {
		ArcaController          *controllers.ArcaController
		CashRegisterController  *controllers.CashRegisterController
		CategoryController      *controllers.CategoryController
		ClientController        *controllers.ClientController
		DepositController       *controllers.DepositController
		EcommerceController     *controllers.EcommerceController
		ExpenseBuyController    *controllers.ExpenseBuyController
		ExpenseOtherController  *controllers.ExpenseOtherController
		IncomeSaleController    *controllers.IncomeSaleController
		IncomeOtherController   *controllers.IncomeOtherController
		MemberController        *controllers.MemberController
		MovementStockController *controllers.MovementStockController
		PermissionController    *controllers.PermissionController
		PointSaleController     *controllers.PointSaleController
		ProductController       *controllers.ProductController
		ReportController        *controllers.ReportController
		RoleController          *controllers.RoleController
		StockController         *controllers.StockController
		SupplierController      *controllers.SupplierController
		TypeMovementController  *controllers.TypeMovementController
	}
	Services struct {
		Arca          *services.ArcaService
		CashRegister  *services.CashRegisterService
		Category      *services.CategoryService
		Client        *services.ClientService
		Deposit       *services.DepositService
		Ecommerce     *services.EcommerceService
		ExpenseBuy    *services.ExpenseBuyService
		ExpenseOther  *services.ExpenseOtherService
		IncomeSale    *services.IncomeSaleService
		IncomeOther   *services.IncomeOtherService
		Member        *services.MemberService
		MovementStock *services.MovementStockService
		Permission    *services.PermissionService
		PointSale     *services.PointSaleService
		Product       *services.ProductService
		Report        *services.ReportService
		Role          *services.RoleService
		Stock         *services.StockService
		Supplier      *services.SupplierService
		TypeMovement  *services.TypeMovementService
	}
	Repositories struct {
		Arca          *repositories.ArcaRepository
		CashRegister  *repositories.CashRegisterRepository
		Category      *repositories.CategoryRepository
		Client        *repositories.ClientRepository
		Deposit       *repositories.DepositRepository
		Ecommerce     *repositories.EcommerceRepository
		Employee      *repositories.EmployeeRepository
		ExpenseBuy    *repositories.ExpenseBuyRepository
		ExpenseOther  *repositories.ExpenseOtherRepository
		IncomeSale    *repositories.IncomeSaleRepository
		IncomeOther   *repositories.IncomeOtherRepository
		Member        *repositories.MemberRepository
		Movement      *repositories.MovementTypeRepository
		MovementStock *repositories.MovementStockRepository
		Permission    *repositories.PermissionRepository
		PointSale     *repositories.PointSaleRepository
		Product       *repositories.ProductRepository
		Report        *repositories.ReportRepository
		Role          *repositories.RoleRepository
		Stock         *repositories.StockRepository
		Supplier      *repositories.SupplierRepository
		TypeMovement  *repositories.TypeMovementRepository
	}
}

func NewTenantContainer(db *gorm.DB) *TenantContainer {
	c := &TenantContainer{DB: db}

	// Inicializar repositorios
	c.Repositories.Arca = &repositories.ArcaRepository{DB: db}
	c.Repositories.CashRegister = &repositories.CashRegisterRepository{DB: db}
	c.Repositories.Category = &repositories.CategoryRepository{DB: db}
	c.Repositories.Client = &repositories.ClientRepository{DB: db}
	c.Repositories.Deposit = &repositories.DepositRepository{DB: db}
	c.Repositories.Ecommerce = &repositories.EcommerceRepository{DB: db}
	c.Repositories.Employee = &repositories.EmployeeRepository{DB: db}
	c.Repositories.ExpenseBuy = &repositories.ExpenseBuyRepository{DB: db}
	c.Repositories.ExpenseOther = &repositories.ExpenseOtherRepository{DB: db}
	c.Repositories.IncomeSale = &repositories.IncomeSaleRepository{DB: db}
	c.Repositories.IncomeOther = &repositories.IncomeOtherRepository{DB: db}
	c.Repositories.Member = &repositories.MemberRepository{DB: db}
	c.Repositories.Movement = &repositories.MovementTypeRepository{DB: db}
	c.Repositories.MovementStock = &repositories.MovementStockRepository{DB: db}
	c.Repositories.Permission = &repositories.PermissionRepository{DB: db}
	c.Repositories.PointSale = &repositories.PointSaleRepository{DB: db}
	c.Repositories.Product = &repositories.ProductRepository{DB: db}
	c.Repositories.Report = &repositories.ReportRepository{DB: db}
	c.Repositories.Role = &repositories.RoleRepository{DB: db}
	c.Repositories.Stock = &repositories.StockRepository{DB: db}
	c.Repositories.Supplier = &repositories.SupplierRepository{DB: db}
	c.Repositories.TypeMovement = &repositories.TypeMovementRepository{DB: db}

	// Inicializar servicios
	c.Services.Arca = &services.ArcaService{
		ArcaRepository: c.Repositories.Arca,
	}
	c.Services.CashRegister = &services.CashRegisterService{
		CashRegisterRepository: c.Repositories.CashRegister,
	}
	c.Services.Category = &services.CategoryService{
		CategoryRepository: c.Repositories.Category,
	}
	c.Services.Client = &services.ClientService{
		ClientRepository: c.Repositories.Client,
	}
	c.Services.Deposit = &services.DepositService{
		DepositRepository: c.Repositories.Deposit,
	}
	c.Services.Ecommerce = &services.EcommerceService{
		EcommerceRepository: c.Repositories.Ecommerce,
	}
	c.Services.ExpenseBuy = &services.ExpenseBuyService{
		ExpenseBuyRepository: c.Repositories.ExpenseBuy,
	}
	c.Services.ExpenseOther = &services.ExpenseOtherService{
		ExpenseOtherRepository: c.Repositories.ExpenseOther,
	}
	c.Services.IncomeSale = &services.IncomeSaleService{
		IncomeSaleRepository: c.Repositories.IncomeSale,
	}
	c.Services.IncomeOther = &services.IncomeOtherService{
		IncomeOtherRepository: c.Repositories.IncomeOther,
	}
	c.Services.Member = &services.MemberService{
		MemberRepository: c.Repositories.Member,
	}
	c.Services.MovementStock = &services.MovementStockService{
		MovementStockRepository: c.Repositories.MovementStock,
		NotifyService:           nil,
	}
	c.Services.Permission = &services.PermissionService{
		PermissionRepository: c.Repositories.Permission,
	}
	c.Services.PointSale = &services.PointSaleService{
		PointSaleRepository: c.Repositories.PointSale,
	}
	c.Services.Product = &services.ProductService{
		ProductRepository: c.Repositories.Product,
	}
	c.Services.Report = &services.ReportService{
		ReportRepository: c.Repositories.Report,
	}
	c.Services.Role = &services.RoleService{
		RoleRepository: c.Repositories.Role,
	}
	c.Services.Stock = &services.StockService{
		StockRepository: c.Repositories.Stock,
	}
	c.Services.Supplier = &services.SupplierService{
		SupplierRepository: c.Repositories.Supplier,
	}
	c.Services.TypeMovement = &services.TypeMovementService{
		TypeMovementRepository: c.Repositories.TypeMovement,
	}

	// Inicializar controladores
	c.Controllers.ArcaController = &controllers.ArcaController{
		ArcaService: c.Services.Arca,
	}
	c.Controllers.CashRegisterController = &controllers.CashRegisterController{
		CashRegisterService: c.Services.CashRegister,
	}
	c.Controllers.CategoryController = &controllers.CategoryController{
		CategoryService: c.Services.Category,
	}
	c.Controllers.ClientController = &controllers.ClientController{
		ClientService: c.Services.Client,
	}
	c.Controllers.DepositController = &controllers.DepositController{
		DepositService: c.Services.Deposit,
	}
	c.Controllers.EcommerceController = &controllers.EcommerceController{
		EcommerceService: c.Services.Ecommerce,
	}
	c.Controllers.ExpenseBuyController = &controllers.ExpenseBuyController{
		ExpenseBuyService: c.Services.ExpenseBuy,
	}
	c.Controllers.ExpenseOtherController = &controllers.ExpenseOtherController{
		ExpenseOtherService: c.Services.ExpenseOther,
	}
	c.Controllers.IncomeSaleController = &controllers.IncomeSaleController{
		IncomeSaleService: c.Services.IncomeSale,
	}
	c.Controllers.IncomeOtherController = &controllers.IncomeOtherController{
		IncomeOtherService: c.Services.IncomeOther,
	}
	c.Controllers.MemberController = &controllers.MemberController{
		MemberService: c.Services.Member,
	}
	c.Controllers.MovementStockController = &controllers.MovementStockController{
		MovementStockService: c.Services.MovementStock,
	}
	c.Controllers.PermissionController = &controllers.PermissionController{
		PermissionService: c.Services.Permission,
	}
	c.Controllers.PointSaleController = &controllers.PointSaleController{
		PointSaleService: c.Services.PointSale,
	}
	c.Controllers.ProductController = &controllers.ProductController{
		ProductService: c.Services.Product,
	}
	c.Controllers.ReportController = &controllers.ReportController{
		ReportService: c.Services.Report,
	}
	c.Controllers.RoleController = &controllers.RoleController{
		RoleService: c.Services.Role,
	}
	c.Controllers.StockController = &controllers.StockController{
		StockService: c.Services.Stock,
	}
	c.Controllers.SupplierController = &controllers.SupplierController{
		SupplierService: c.Services.Supplier,
	}
	c.Controllers.TypeMovementController = &controllers.TypeMovementController{
		TypeMovementService: c.Services.TypeMovement,
	}

	return c
}
