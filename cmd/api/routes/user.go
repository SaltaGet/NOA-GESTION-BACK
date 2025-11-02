package routes

import (
	"github.com/DanielChachagua/GestionCar/cmd/api/controllers"
	"github.com/DanielChachagua/GestionCar/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, controllers *controllers.UserController) {
	auth := app.Group("/user")
	auth.Post(
		"/create", 
		middleware.AuthMiddleware(), 
		// middleware.RoleAuthMiddleware([]string{"super_admin","admin"}), 
		controllers.CreateUser,
	)
}