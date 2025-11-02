package routes

import (
	"github.com/DanielChachagua/GestionCar/cmd/api/controllers"
	"github.com/DanielChachagua/GestionCar/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func VehicleRoutes(app *fiber.App){
	vehicle := app.Group("/vehicle", middleware.AuthMiddleware(), middleware.TenantMiddleware())

	vehicle.Get("/get_all", GetController("VehicleController", func(c *fiber.Ctx, ctrl *controllers.VehicleController) error {
		return ctrl.VehicleGetAll(c)
	}))

	vehicle.Get("/get_by_domain", GetController("VehicleController", func(c *fiber.Ctx, ctrl *controllers.VehicleController) error {
		return ctrl.VehicleGetByDomain(c)
	}))

	vehicle.Post("/create", GetController("VehicleController", func(c *fiber.Ctx, ctrl *controllers.VehicleController) error {
		return ctrl.VehicleCreate(c)
	}))

	vehicle.Put("/update", GetController("VehicleController", func(c *fiber.Ctx, ctrl *controllers.VehicleController) error {
		return ctrl.VehicleUpdate(c)
	}))

	vehicle.Get("/get_by_client/:client_id", GetController("VehicleController", func(c *fiber.Ctx, ctrl *controllers.VehicleController) error {
		return ctrl.VehicleGetByClientID(c)
	}))

	vehicle.Delete("/delete/:id", GetController("VehicleController", func(c *fiber.Ctx, ctrl *controllers.VehicleController) error {
		return ctrl.VehicleDelete(c)
	}))

	vehicle.Get("/:id", GetController("VehicleController", func(c *fiber.Ctx, ctrl *controllers.VehicleController) error {
		return ctrl.VehicleGetByID(c)
	}))

}
