package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, appDependencies *dependencies.MainContainer) {
	AuthRoutes(app, appDependencies.AuthController)
	ClientRoutes(app)
	ExpenseRoutes(app)
	IncomeRoutes(app)
	MemberRoutes(app)
	MovementRoutes(app)
	PermissionRoutes(app)
	PointSaleRoutes(app)
	ProductRoutes(app)
	RoleRoutes(app)
	SupplierRoutes(app)
	UserRoutes(app, appDependencies.UserController)
	TenantRoutes(app, appDependencies.TenantController)
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