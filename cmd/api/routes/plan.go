package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func PlanRoutes(app *fiber.App, controllers *controllers.PlanController) {
	plan := app.Group("/api/v1/plan")

	plan.Get("/get_all", controllers.PlanGetAll)

	plan.Post("/create", middleware.AdminAuthMiddleware(), controllers.PlanCreate)

	plan.Put("/update", middleware.AdminAuthMiddleware(), controllers.PlanUpdate)

	plan.Get("/get/:id", middleware.AdminAuthMiddleware(), controllers.PlanGetByID)
}
