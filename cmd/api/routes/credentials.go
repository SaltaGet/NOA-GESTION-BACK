package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func CredentialsRoutes(app *fiber.App, ctrl *controllers.CredentialController) {
	credential := app.Group("/api/v1/credential", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	credential.Get("/get_mercado_pago_token",
		middleware.AuthMiddleware(),
		middleware.AdminTenantMiddleware(),
		middleware.CurrentPlan(),
		middleware.AuthModule("ecommerce"),
		ctrl.CredentialGetMPToken,
	)

	credential.Put("/set_mercado_pago_token",
		middleware.AuthMiddleware(),
		middleware.AdminTenantMiddleware(),
		middleware.CurrentPlan(),
		middleware.AuthModule("ecommerce"),
		ctrl.CredentialSetMPToken,
	)
	
	credential.Get("/get_arca_token",
		middleware.AuthMiddleware(),
		middleware.AdminTenantMiddleware(),
		ctrl.CredentialGetArca,
	)

	credential.Put("/set_arca_token",
		middleware.AuthMiddleware(),
		middleware.AdminTenantMiddleware(),
		ctrl.CredentialSetArca,
	)

}
