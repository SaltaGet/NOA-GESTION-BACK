package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, appDependencies *dependencies.Application) {
	AuthRoutes(app, appDependencies.AuthController)
	ClientRoutes(app)
	EmployeeRoutes(app)
	ExpenseRoutes(app)
	IncomeRoutes(app)
	MemberRoutes(app)
	MovementRoutes(app)
	PermissionRoutes(app)
	ProductRoutes(app)
	PurchaseOrderRoutes(app)
	// PurchaseProductRoutes(app, tenantDependencies.PurchaseProductController)
	RoleRoutes(app)
	ServiceRoutes(app)
	SupplierRoutes(app)
	UserRoutes(app, appDependencies.UserController)
	VehicleRoutes(app)
	TenantRoutes(app, appDependencies.TenantController)
	// AttendanceRoutes(app, tenantDependencies.AttendanceController)
	// AuthRoutes(app, appDependencies.AuthController)
	// ClientRoutes(app, tenantDependencies.ClientController)
	// EmployeeRoutes(app, tenantDependencies.EmployeeController)
	// ExpenseRoutes(app, tenantDependencies.ExpenseController)
	// IncomeRoutes(app, tenantDependencies.IncomeController)
	// MemberRoutes(app, tenantDependencies.MemberController)
	// MovementRoutes(app, tenantDependencies.MovementTypeController)
	// ProductRoutes(app, tenantDependencies.ProductController)
	// PurchaseOrderRoutes(app, tenantDependencies.PurchaseOrderController)
	// // PurchaseProductRoutes(app, tenantDependencies.PurchaseProductController)
	// RoleRoutes(app, tenantDependencies.RoleController)
	// ServiceRoutes(app, tenantDependencies.ServiceController)
	// SupplierRoutes(app, tenantDependencies.SupplierController)
	// UserRoutes(app, appDependencies.UserController)
	// VehicleRoutes(app, tenantDependencies.VehicleController)
	// TenantRoutes(app, appDependencies.TenantController)
}

func GetController[T any](key string, handler func(c *fiber.Ctx, ctrl *T) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctrlInterface := c.Locals(key)
		if ctrlInterface == nil {
			return c.Status(fiber.StatusInternalServerError).SendString(key + " no inicializado correctamente")
		}

		ctrl, ok := ctrlInterface.(*T)
		if !ok || ctrl == nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error de tipo para controlador " + key)
		}

		return handler(c, ctrl)
	}
}