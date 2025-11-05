package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func TenantRoutes(app *fiber.App, controllers *controllers.TenantController){
	tenant := app.Group("/api/v1/tenant")
	tenant.Get("/get_all", middleware.AdminAuthMiddleware(), controllers.GetTenants)
	tenant.Post("/create", middleware.AdminAuthMiddleware(), controllers.TenantCreateByUserID)
	tenant.Post("/create_tenant_user", middleware.AdminAuthMiddleware(), controllers.TenantUserCreate)
}