package routes

import (
	"github.com/DanielChachagua/GestionCar/cmd/api/controllers"
	"github.com/DanielChachagua/GestionCar/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func SupplierRoutes(app *fiber.App){
	supplier := app.Group("/supplier", middleware.AuthMiddleware(), middleware.TenantMiddleware())

	supplier.Get("/get_all", GetController("SupplierController", func(c *fiber.Ctx, ctrl *controllers.SupplierController) error {
		return ctrl.SupplierGetAll(c)
	}))

	supplier.Get("/get_by_name", GetController("SupplierController", func(c *fiber.Ctx, ctrl *controllers.SupplierController) error {
		return ctrl.SupplierGetByName(c)
	}))

	supplier.Post("/create", GetController("SupplierController", func(c *fiber.Ctx, ctrl *controllers.SupplierController) error {
		return ctrl.SupplierCreate(c)
	}))

	supplier.Put("/update", GetController("SupplierController", func(c *fiber.Ctx, ctrl *controllers.SupplierController) error {
		return ctrl.SupplierUpdate(c)
	}))

	supplier.Delete("/delete/:id", GetController("SupplierController", func(c *fiber.Ctx, ctrl *controllers.SupplierController) error {
		return ctrl.SupplierDeleteByID(c)
	}))

	supplier.Get("/:id", GetController("SupplierController", func(c *fiber.Ctx, ctrl *controllers.SupplierController) error {
		return ctrl.SupplierGetByID(c)
	}))

}
