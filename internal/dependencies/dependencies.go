package dependencies

import (
	"github.com/DanielChachagua/GestionCar/cmd/api/controllers"
	"github.com/DanielChachagua/GestionCar/pkg/repositories"
	"github.com/DanielChachagua/GestionCar/pkg/services"
	"gorm.io/gorm"
)

// import (
// 	"sync"

// 	"github.com/DanielChachagua/GestionCar/cmd/api/controllers"
// 	"github.com/DanielChachagua/GestionCar/pkg/repositories"
// 	"github.com/DanielChachagua/GestionCar/pkg/services"
// 	"gorm.io/gorm"
// )


var App *Application

type Application struct {
	AuthController *controllers.AuthController
	UserController *controllers.UserController
	TenantController *controllers.TenantController
}

func NewApplication(mainDB *gorm.DB) *Application {
	mainRepo := &repositories.MainRepository{DB: mainDB}

	authServ := &services.AuthService{AuthRepository: mainRepo, UserRepository: mainRepo, TenantService: mainRepo}
	userServ := &services.UserService{UserRepository: mainRepo}
	tenantServ := &services.TenantService{TenantRepository: mainRepo}

	return &Application{
		AuthController: &controllers.AuthController{AuthService: authServ},
		UserController: &controllers.UserController{UserService: userServ},
		TenantController: &controllers.TenantController{TenantService: tenantServ},
	}
}

// var TenantApp *TenantApplication

// type TenantApplication struct {
// 	AttendanceController *controllers.AttendanceController
// 	ClientController *controllers.ClientController
// 	EmployeeController *controllers.EmployeeController
// 	ExpenseController *controllers.ExpenseController	
// 	IncomeController *controllers.IncomeController
// 	MemberController *controllers.MemberController
// 	MovementTypeController *controllers.MovementTypeController
// 	PermissionController *controllers.PermissionController
// 	ProductController *controllers.ProductController
// 	PurchaseOrderController *controllers.PurchaseOrderController
// 	PurchaseProductController *controllers.PurchaseProductController
// 	ResumeController *controllers.ResumeController
// 	RoleController *controllers.RoleController
// 	ServiceController *controllers.ServiceController
// 	SupplierController *controllers.SupplierController
// 	VehicleController *controllers.VehicleController
// }

// func TenantDBRepository(db *gorm.DB) *TenantApplication {
// 	tenantRepo := &repositories.TenantRepository{DB: db,}
	
// 	attendanceService := &services.AttendanceService{AttendanceRepository: tenantRepo}
// 	clientService := &services.ClientService{ClientRepository: tenantRepo}
// 	employeeService := &services.EmployeeService{EmployeeRepository: tenantRepo}
// 	expenseService := &services.ExpenseService{ExpenseRepository: tenantRepo}
// 	incomeService := &services.IncomeService{IncomeRepository: tenantRepo}
// 	memberService := &services.MemberService{MemberRepository: tenantRepo}
// 	movementService := &services.MovementTypeService{MovementTypeRepository: tenantRepo}
// 	permissionService := &services.PermissionService{PermissionRepository: tenantRepo}
// 	productService := &services.ProductService{ProductRepository: tenantRepo}
// 	purchaseOrderService := &services.PurchaseOrderService{PurchaseOrderRepository: tenantRepo}
// 	purchaseProductService := &services.PurchaseProductService{PurchaseProductRepository: tenantRepo}
// 	resumeService := &services.ResumeService{ResumeExpenseRepository: tenantRepo, ResumeIncomeRepository: tenantRepo}
// 	roleService := &services.RoleService{RoleRepository: tenantRepo}
// 	serviceService := &services.ServiceService{ServiceRepository: tenantRepo}
// 	supplierService := &services.SupplierService{SupplierRepository: tenantRepo}
// 	vehicleService := &services.VehicleService{VehicleRepository: tenantRepo}

// 	return &TenantApplication{
// 		AttendanceController: &controllers.AttendanceController{AttendanceService: attendanceService},
// 		ClientController: &controllers.ClientController{ClientService: clientService},
// 		EmployeeController: &controllers.EmployeeController{EmployeeService: employeeService},
// 		ExpenseController: &controllers.ExpenseController{ExpenseService: expenseService},
// 		IncomeController: &controllers.IncomeController{IncomeService: incomeService},
// 		MemberController: &controllers.MemberController{MemberService: memberService},
// 		MovementTypeController: &controllers.MovementTypeController{MovementTypeService: movementService},
// 		PermissionController: &controllers.PermissionController{PermissionService: permissionService},
// 		ProductController: &controllers.ProductController{ProductService: productService},
// 		PurchaseOrderController: &controllers.PurchaseOrderController{PurchaseOrderService: purchaseOrderService},
// 		PurchaseProductController: &controllers.PurchaseProductController{PurchaseProductService: purchaseProductService},
// 		ResumeController: &controllers.ResumeController{ResumeExpenseService: resumeService, ResumeIncomeService: resumeService},
// 		RoleController: &controllers.RoleController{RoleService: roleService},
// 		ServiceController: &controllers.ServiceController{ServiceService: serviceService},
// 		SupplierController: &controllers.SupplierController{SupplierService: supplierService},
// 		VehicleController: &controllers.VehicleController{VehicleService: vehicleService},
// 	}
// }

// func NewTenantDBRepository(db *gorm.DB) *TenantApplication {
// 	if db == nil {
// 		return &TenantApplication{}
// 	}
// 	return TenantDBRepository(db)
// }


// func (app *TenantApplication) SetDBTenantRepository(db *gorm.DB) {
// 	tenantRepo := &repositories.TenantRepository{DB: db,}
	
// 	attendanceService := &services.AttendanceService{AttendanceRepository: tenantRepo}
// 	clientService := &services.ClientService{ClientRepository: tenantRepo}
// 	employeeService := &services.EmployeeService{EmployeeRepository: tenantRepo}
// 	expenseService := &services.ExpenseService{ExpenseRepository: tenantRepo}
// 	incomeService := &services.IncomeService{IncomeRepository: tenantRepo}
// 	memberService := &services.MemberService{MemberRepository: tenantRepo}
// 	movementService := &services.MovementTypeService{MovementTypeRepository: tenantRepo}
// 	permissionService := &services.PermissionService{PermissionRepository: tenantRepo}
// 	productService := &services.ProductService{ProductRepository: tenantRepo}
// 	purchaseOrderService := &services.PurchaseOrderService{PurchaseOrderRepository: tenantRepo}
// 	purchaseProductService := &services.PurchaseProductService{PurchaseProductRepository: tenantRepo}
// 	resumeService := &services.ResumeService{ResumeExpenseRepository: tenantRepo, ResumeIncomeRepository: tenantRepo}
// 	roleService := &services.RoleService{RoleRepository: tenantRepo}
// 	serviceService := &services.ServiceService{ServiceRepository: tenantRepo}
// 	supplierService := &services.SupplierService{SupplierRepository: tenantRepo}
// 	vehicleService := &services.VehicleService{VehicleRepository: tenantRepo}

// 	app.AttendanceController.AttendanceService = attendanceService
// 	app.ClientController.ClientService = clientService
// 	app.EmployeeController.EmployeeService = employeeService
// 	app.ExpenseController.ExpenseService = expenseService
// 	app.IncomeController.IncomeService = incomeService
// 	app.MemberController.MemberService = memberService
// 	app.MovementTypeController.MovementTypeService = movementService
// 	app.PermissionController.PermissionService = permissionService
// 	app.ProductController.ProductService = productService
// 	app.PurchaseOrderController.PurchaseOrderService = purchaseOrderService
// 	app.PurchaseProductController.PurchaseProductService = purchaseProductService
// 	app.ResumeController.ResumeExpenseService = resumeService
// 	app.ResumeController.ResumeIncomeService = resumeService
// 	app.RoleController.RoleService = roleService
// 	app.ServiceController.ServiceService = serviceService
// 	app.SupplierController.SupplierService = supplierService
// 	app.VehicleController.VehicleService = vehicleService
// }