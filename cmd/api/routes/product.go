package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App){
	prod := app.Group("/api/v1/product", middleware.AuthMiddleware(), middleware.PointSaleMiddleware())

	prod.Get("/get_all", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductGetAll(c)
	}))

	prod.Get("/get_by_name", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductGetByName(c)
	}))

	prod.Get("/get_by_identifier", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductGetByIdentifier(c)
	}))

	prod.Post("/create", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductCreate(c)
	}))

	prod.Put("/update", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductUpdate(c)
	}))

	prod.Put("/update_stock", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductUpdateStock(c)
	}))

	prod.Delete("/delete/:id", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductDelete(c)
	}))

	prod.Get("/:id", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductGetByID(c)
	}))

}
