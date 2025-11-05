package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, controllers *controllers.UserController) {
	auth := app.Group("/api/v1/user")
	auth.Post(
		"/create", 
		middleware.AdminAuthMiddleware(), 
		// middleware.RoleAuthMiddleware([]string{"super_admin","admin"}), 
		controllers.CreateUser,
	)
}