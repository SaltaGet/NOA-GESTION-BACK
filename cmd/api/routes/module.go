package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/gofiber/fiber/v2"
)

func ModuleRoutes(app *fiber.App, ctrl *controllers.ModuleController) {
	module := app.Group("/api/v1/module")

	module.Get("/get_all", ctrl.ModuleGetAll)

	module.Post("/create", ctrl.ModuleCreate)

	module.Put("/update", ctrl.ModuleUpdate)

	module.Post("/add_tenant_expiration", ctrl.ModuleAddTenant)
	
	module.Delete("/delete/:id", ctrl.ModuleDelete)

	module.Get("/get/:id", ctrl.ModuleGet)
}
