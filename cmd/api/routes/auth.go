package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, controllers *controllers.AuthController) {
	auth := app.Group("/api/v1/auth")

	auth.Post("/login", controllers.AuthLogin)
	auth.Post("/login_admin", controllers.AuthLoginAdmin)
	
	auth.Post("/tenant_logout", controllers.LogoutPointSale)
	
	auth.Post("/logout", controllers.Logout)
	auth.Post("/logout_admin", middleware.AdminAuthMiddleware(), controllers.LogoutAdmin)
	auth.Post("/logout_point_sale", middleware.AuthMiddleware(), middleware.AuthPointSaleMiddleware(), controllers.LogoutPointSale)

	auth.Get("/current_user", middleware.AuthMiddleware(), controllers.CurrentUser)

	auth.Post("/forgot_password", controllers.AuthForgotPassword)
	auth.Post("/reset_password", controllers.AuthResetPassword)
	
	auth.Post("/login_point_sale/:point_sale_id", middleware.AuthMiddleware(), controllers.AuthPointSale)
}
