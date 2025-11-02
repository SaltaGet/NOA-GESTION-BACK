// pkg/dependencies/container.go
package dependencies

import (
	"github.com/DanielChachagua/GestionCar/pkg/repositories"
	"github.com/DanielChachagua/GestionCar/pkg/services"
	"gorm.io/gorm"
)

type TenantContainer struct {
	DB       *gorm.DB
	Services struct {
		Attendance      *services.AttendanceService
		Client          *services.ClientService
		Employee        *services.EmployeeService
		Expense         *services.ExpenseService
		Income          *services.IncomeService
		Member          *services.MemberService
		Movement        *services.MovementTypeService
		Permission      *services.PermissionService
		Product         *services.ProductService
		Purchase        *services.PurchaseOrderService
		PurchaseProduct *services.PurchaseProductService
		Resume          *services.ResumeService
		Role            *services.RoleService
		Service         *services.ServiceService
		Supplier        *services.SupplierService
		Vehicle         *services.VehicleService
	}
	Repositories struct {
		Attendance      *repositories.AttendanceRepository
		Client          *repositories.ClientRepository
		Employee        *repositories.EmployeeRepository
		Expense         *repositories.ExpenseRepository
		Income          *repositories.IncomeRepository
		Member          *repositories.MemberRepository
		Movement        *repositories.MovementTypeRepository
		Permission      *repositories.PermissionRepository
		Product         *repositories.ProductRepository
		Purchase        *repositories.PurchaseOrderRepository
		PurchaseProduct *repositories.PurchaseProductRepository
		Resume          *repositories.ResumeRepository
		Role            *repositories.RoleRepository
		Service         *repositories.ServiceRepository
		Supplier        *repositories.SupplierRepository
		Vehicle         *repositories.VehicleRepository
	}
}

func NewTenantContainer(db *gorm.DB) *TenantContainer {
	c := &TenantContainer{DB: db}

	// Inicializar repositorios
	c.Repositories.Attendance = &repositories.AttendanceRepository{DB: db}
	c.Repositories.Client = &repositories.ClientRepository{DB: db}
	c.Repositories.Employee = &repositories.EmployeeRepository{DB: db}
	c.Repositories.Expense = &repositories.ExpenseRepository{DB: db}
	c.Repositories.Income = &repositories.IncomeRepository{DB: db}
	c.Repositories.Member = &repositories.MemberRepository{DB: db}
	c.Repositories.Movement = &repositories.MovementTypeRepository{DB: db}
	c.Repositories.Permission = &repositories.PermissionRepository{DB: db}
	c.Repositories.Product = &repositories.ProductRepository{DB: db}
	c.Repositories.Purchase = &repositories.PurchaseOrderRepository{DB: db}
	c.Repositories.PurchaseProduct = &repositories.PurchaseProductRepository{DB: db}
	c.Repositories.Resume = &repositories.ResumeRepository{DB: db}
	c.Repositories.Role = &repositories.RoleRepository{DB: db}
	c.Repositories.Service = &repositories.ServiceRepository{DB: db}
	c.Repositories.Supplier = &repositories.SupplierRepository{DB: db}
	c.Repositories.Vehicle = &repositories.VehicleRepository{DB: db}

	// Inicializar servicios
	c.Services.Attendance = &services.AttendanceService{
		AttendanceRepository: c.Repositories.Attendance,
	}
	c.Services.Client = &services.ClientService{
		ClientRepository: c.Repositories.Client,
	}
	c.Services.Employee = &services.EmployeeService{
		EmployeeRepository: c.Repositories.Employee,
	}
	c.Services.Expense = &services.ExpenseService{
		ExpenseRepository: c.Repositories.Expense,
	}
	c.Services.Income = &services.IncomeService{
		IncomeRepository: c.Repositories.Income,
	}
	c.Services.Member = &services.MemberService{
		MemberRepository: c.Repositories.Member,
	}
	c.Services.Movement = &services.MovementTypeService{
		MovementTypeRepository: c.Repositories.Movement,
	}
	c.Services.Permission = &services.PermissionService{
		PermissionRepository: c.Repositories.Permission,
	}
	c.Services.Product = &services.ProductService{
		ProductRepository: c.Repositories.Product,
	}
	c.Services.Purchase = &services.PurchaseOrderService{
		PurchaseOrderRepository: c.Repositories.Purchase,
	}
	c.Services.PurchaseProduct = &services.PurchaseProductService{
		PurchaseProductRepository: c.Repositories.PurchaseProduct,
	}
	c.Services.Resume = &services.ResumeService{
		ResumeExpenseRepository: c.Repositories.Resume,
		ResumeIncomeRepository:  c.Repositories.Resume,
	}
	c.Services.Role = &services.RoleService{
		RoleRepository: c.Repositories.Role,
	}
	c.Services.Service = &services.ServiceService{
		ServiceRepository: c.Repositories.Service,
	}
	c.Services.Supplier = &services.SupplierService{
		SupplierRepository: c.Repositories.Supplier,
	}
	c.Services.Vehicle = &services.VehicleService{
		VehicleRepository: c.Repositories.Vehicle,
	}

	return c
}
