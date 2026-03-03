package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
	Auth "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/auth/infrastructure/handler/http"
	Arca "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/arca/infrastructure/handler/http"
	CashRegister "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/cash_register/infrastructure/handler/http"
	Category "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/category/infrastructure/handler/http"
	Client "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/client/infrastructure/handler/http"
	Credentials "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/credential/infrastructure/handler/http"
	Deposit "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/deposit/infrastructure/handler/http"
	Ecommerce "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/ecommerce/infrastructure/handler/http"
	ExpenseBuy "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/expense_buy/infrastructure/handler/http"
	ExpenseOther "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/expense_other/infrastructure/handler/http"
	Feedback "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/feedback/infrastructure/handler/http"
	IncomeOther "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/income_other/infrastructure/handler/http"
	IncomeSale "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/income_sale/infrastructure/handler/http"
	Member "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/member/infrastructure/handler/http"
	Module "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/module/infrastructure/handler/http"
	MovementStock "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/movement_stock/infrastructure/handler/http"
	News "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/news/infrastructure/handler/http"
	Notification "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/notification/infrastructure/handler/http"
	Permission "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/permission/infrastructure/handler/http"
	Plan "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/plan/infrastructure/handler/http"
	PointSale "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/point_sale/infrastructure/handler/http"
	Product "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/product/infrastructure/handler/http"
	Report "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/report/infrastructure/handler/http"
	Role "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/role/infrastructure/handler/http"
	Stock "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/stock/infrastructure/handler/http"
	Supplier "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/supplier/infrastructure/handler/http"
	Tenant "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/tenant/infrastructure/handler/http"
	TypeMovement "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/type_movement/infrastructure/handler/http"
)

func SetupRoutes(app *fiber.App, appDependencies *dependencies.MainContainer) {
	Auth.AuthRoutes(app, appDependencies.AuthController)
	Arca.ArcaRoutes(app)
	CashRegister.CashRegisterRoutes(app)
	Category.CategoryRoutes(app)
	Client.ClientRoutes(app)
	Credentials.CredentialsRoutes(app, appDependencies.CredentialController)
	Deposit.DepositRoutes(app)
	Ecommerce.EcommerceRoutes(app)
	ExpenseBuy.ExpenseBuyRoutes(app)
	ExpenseOther.ExpenseOtherRoutes(app)
	Feedback.FeedbackRoutes(app, appDependencies.FeedbackController)
	IncomeOther.IncomeOtherRoutes(app)
	IncomeSale.IncomeSaleRoutes(app)
	Member.MemberRoutes(app)
	Module.ModuleRoutes(app, appDependencies.ModuleController)
	MovementStock.MovementStockRoutes(app)
	News.NewsRoutes(app, appDependencies.NewsController)
	Notification.NotificationRoutes(app)
	Permission.PermissionRoutes(app)
	Plan.PlanRoutes(app, appDependencies.PlanController)
	PointSale.PointSaleRoutes(app)
	Product.ProductRoutes(app)
	Report.ReportRoutes(app)
	Role.RoleRoutes(app)
	Stock.StockRoutes(app)
	Supplier.SupplierRoutes(app)
	Tenant.TenantRoutes(app, appDependencies.TenantController)
	TypeMovement.TypeMovementRoutes(app)
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