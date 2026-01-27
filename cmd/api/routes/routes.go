package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, appDependencies *dependencies.MainContainer) {
	AuthRoutes(app, appDependencies.AuthController)
	CashRegisterRoutes(app)
	CategoryRoutes(app)
	ClientRoutes(app)
	CredentialsRoutes(app, appDependencies.CredentialController)
	DepositRoutes(app)
	EcommerceRoutes(app)
	ExpenseBuyRoutes(app)
	ExpenseOtherRoutes(app)
	FeedbackRoutes(app, appDependencies.FeedbackController)
	IncomeOtherRoutes(app)
	IncomeSaleRoutes(app)
	MemberRoutes(app)
	ModuleRoutes(app, appDependencies.ModuleController)
	MovementStockRoutes(app)
	NewsRoutes(app, appDependencies.NewsController)
	NotificationRoutes(app)
	PermissionRoutes(app)
	PlanRoutes(app, appDependencies.PlanController)
	PointSaleRoutes(app)
	ProductRoutes(app)
	ReportRoutes(app)
	RoleRoutes(app)
	StockRoutes(app)
	SupplierRoutes(app)
	UserRoutes(app, appDependencies.UserController)
	TenantRoutes(app, appDependencies.TenantController)
	TypeMovementRoutes(app)
}

// func GetController[T any](key string, handler func(c *fiber.Ctx, ctrl *T) error) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		ctrlInterface := c.Locals(key)
// 		if ctrlInterface == nil {
// 			return c.Status(fiber.StatusInternalServerError).SendString(key + " no inicializado correctamente")
// 		}

// 		ctrl, ok := ctrlInterface.(*T)
// 		if !ok || ctrl == nil {
// 			return c.Status(fiber.StatusInternalServerError).SendString("Error de tipo para controlador " + key)
// 		}

// 		return handler(c, ctrl)
// 	}
// }