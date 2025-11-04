package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, controllers *controllers.AuthController) {
	auth := app.Group("api/v1/auth")

	auth.Post("/login", controllers.AuthLogin)
	
	auth.Post("/tenant_logout", controllers.LogoutPointSale)
	
	auth.Post("/logout", controllers.Logout)

	auth.Get("/current_user", controllers.CurrentUser)
	
	auth.Post("/tenant_login/:tenant_id", controllers.AuthPointSale)
}
