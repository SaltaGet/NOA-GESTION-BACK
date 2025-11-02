package routes

import (
	"github.com/DanielChachagua/GestionCar/cmd/api/controllers"
	"github.com/DanielChachagua/GestionCar/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func ServiceRoutes(app *fiber.App){
	service := app.Group("/service", middleware.AuthMiddleware(), middleware.TenantMiddleware())

	service.Get("/get_all", GetController("ServiceController", func(c *fiber.Ctx, ctrl *controllers.ServiceController) error {
		return ctrl.ServiceGetAll(c)
	}))

	service.Post("/create", GetController("ServiceController", func(c *fiber.Ctx, ctrl *controllers.ServiceController) error {
		return ctrl.ServiceCreate(c)
	}))

	service.Put("/update", GetController("ServiceController", func(c *fiber.Ctx, ctrl *controllers.ServiceController) error {
		return ctrl.ServiceUpdate(c)
	}))

	service.Delete("/delete/:id", GetController("ServiceController", func(c *fiber.Ctx, ctrl *controllers.ServiceController) error {
		return ctrl.ServiceDeleteByID(c)
	}))

	service.Get("/:id", GetController("ServiceController", func(c *fiber.Ctx, ctrl *controllers.ServiceController) error {
		return ctrl.ServiceGetByID(c)
	}))

}
