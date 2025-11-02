package routes

import (
	"github.com/DanielChachagua/GestionCar/cmd/api/controllers"
	"github.com/DanielChachagua/GestionCar/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func TenantRoutes(app *fiber.App, controllers *controllers.TenantController){
	tenant := app.Group("/tenant")
	tenant.Get("/get_all", middleware.AuthMiddleware(), controllers.GetTenants)
	tenant.Post("/create", middleware.AuthMiddleware(), controllers.TenantCreateByUserID)
	tenant.Post("/create_tenant_user", middleware.AuthMiddleware(), controllers.TenantUserCreate)
}