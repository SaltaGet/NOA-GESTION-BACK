package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func ModuleRoutes(app *fiber.App, ctrl *controllers.ModuleController) {
	module := app.Group("/api/v1/module")

	module.Get("/get_all", ctrl.ModuleGetAll)

	module.Post("/create", middleware.AdminAuthMiddleware(), ctrl.ModuleCreate)

	module.Put("/update", middleware.AdminAuthMiddleware(), ctrl.ModuleUpdate)

	module.Put("/add_tenant_expiration", middleware.AdminAuthMiddleware(), ctrl.ModuleAddTenant)
	
	module.Delete("/delete/:id", middleware.AdminAuthMiddleware(), ctrl.ModuleDelete)

	module.Get("/get/:id", middleware.AdminAuthMiddleware(), ctrl.ModuleGet)
}
