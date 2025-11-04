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
		Client          *services.ClientService
		Expense         *services.ExpenseService
		Income          *services.IncomeService
		Member          *services.MemberService
		Movement        *services.MovementTypeService
		Permission      *services.PermissionService
		Product         *services.ProductService
		Role            *services.RoleService
		Supplier        *services.SupplierService
	}
	Repositories struct {
		Client          *repositories.ClientRepository
		Employee        *repositories.EmployeeRepository
		Expense         *repositories.ExpenseRepository
		Income          *repositories.IncomeRepository
		Member          *repositories.MemberRepository
		Movement        *repositories.MovementTypeRepository
		Permission      *repositories.PermissionRepository
		Product         *repositories.ProductRepository
		Role            *repositories.RoleRepository
		Supplier        *repositories.SupplierRepository
	}
}

func NewTenantContainer(db *gorm.DB) *TenantContainer {
	c := &TenantContainer{DB: db}

	// Inicializar repositorios
	c.Repositories.Client = &repositories.ClientRepository{DB: db}
	c.Repositories.Employee = &repositories.EmployeeRepository{DB: db}
	c.Repositories.Expense = &repositories.ExpenseRepository{DB: db}
	c.Repositories.Income = &repositories.IncomeRepository{DB: db}
	c.Repositories.Member = &repositories.MemberRepository{DB: db}
	c.Repositories.Movement = &repositories.MovementTypeRepository{DB: db}
	c.Repositories.Permission = &repositories.PermissionRepository{DB: db}
	c.Repositories.Product = &repositories.ProductRepository{DB: db}
	c.Repositories.Role = &repositories.RoleRepository{DB: db}
	c.Repositories.Supplier = &repositories.SupplierRepository{DB: db}

	// Inicializar servicios
	c.Services.Client = &services.ClientService{
		ClientRepository: c.Repositories.Client,
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
	c.Services.Role = &services.RoleService{
		RoleRepository: c.Repositories.Role,
	}
	c.Services.Supplier = &services.SupplierService{
		SupplierRepository: c.Repositories.Supplier,
	}

	return c
}
