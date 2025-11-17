// pkg/dependencies/container.go
package dependencies

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/repositories"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/services"
	"gorm.io/gorm"
)

type TenantContainer struct {
	DB       *gorm.DB
	Services struct {
		CashRegister *services.CashRegisterService
		Category     *services.CategoryService
		Client       *services.ClientService
		Deposit      *services.DepositService
		ExpenseBuy   *services.ExpenseBuyService
		ExpenseOther *services.ExpenseOtherService
		IncomeSale   *services.IncomeSaleService
		IncomeOther  *services.IncomeOtherService
		Member       *services.MemberService
		Permission   *services.PermissionService
		PointSale    *services.PointSaleService
		Product      *services.ProductService
		Report       *services.ReportService
		Role         *services.RoleService
		Supplier     *services.SupplierService
		TypeMovement *services.TypeMovementService
	}
	Repositories struct {
		CashRegister *repositories.CashRegisterRepository
		Category     *repositories.CategoryRepository
		Client       *repositories.ClientRepository
		Deposit      *repositories.DepositRepository
		Employee     *repositories.EmployeeRepository
		ExpenseBuy   *repositories.ExpenseBuyRepository
		ExpenseOther *repositories.ExpenseOtherRepository
		IncomeSale   *repositories.IncomeSaleRepository
		IncomeOther  *repositories.IncomeOtherRepository
		Member       *repositories.MemberRepository
		Movement     *repositories.MovementTypeRepository
		Permission   *repositories.PermissionRepository
		PointSale    *repositories.PointSaleRepository
		Product      *repositories.ProductRepository
		Report       *repositories.ReportRepository
		Role         *repositories.RoleRepository
		Supplier     *repositories.SupplierRepository
		TypeMovement *repositories.TypeMovementRepository
	}
}

func NewTenantContainer(db *gorm.DB) *TenantContainer {
	c := &TenantContainer{DB: db}

	// Inicializar repositorios
	c.Repositories.CashRegister = &repositories.CashRegisterRepository{DB: db}
	c.Repositories.Category = &repositories.CategoryRepository{DB: db}
	c.Repositories.Client = &repositories.ClientRepository{DB: db}
	c.Repositories.Deposit = &repositories.DepositRepository{DB: db}
	c.Repositories.Employee = &repositories.EmployeeRepository{DB: db}
	c.Repositories.ExpenseBuy = &repositories.ExpenseBuyRepository{DB: db}
	c.Repositories.ExpenseOther = &repositories.ExpenseOtherRepository{DB: db}
	c.Repositories.IncomeSale = &repositories.IncomeSaleRepository{DB: db}
	c.Repositories.IncomeOther = &repositories.IncomeOtherRepository{DB: db}
	c.Repositories.Member = &repositories.MemberRepository{DB: db}
	c.Repositories.Movement = &repositories.MovementTypeRepository{DB: db}
	c.Repositories.Permission = &repositories.PermissionRepository{DB: db}
	c.Repositories.PointSale = &repositories.PointSaleRepository{DB: db}
	c.Repositories.Product = &repositories.ProductRepository{DB: db}
	c.Repositories.Report = &repositories.ReportRepository{DB: db}
	c.Repositories.Role = &repositories.RoleRepository{DB: db}
	c.Repositories.Supplier = &repositories.SupplierRepository{DB: db}
	c.Repositories.TypeMovement = &repositories.TypeMovementRepository{DB: db}

	// Inicializar servicios
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
	c.Services.Supplier = &services.SupplierService{
		SupplierRepository: c.Repositories.Supplier,
	}
	c.Services.TypeMovement = &services.TypeMovementService{
		TypeMovementRepository: c.Repositories.TypeMovement,
	}

	return c
}
