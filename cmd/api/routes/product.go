package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App){
	prod := app.Group("/api/v1/product", middleware.AuthMiddleware(), middleware.AuthPointSaleMiddleware())

	prod.Get("/get_all", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductGetAll(c)
	}))

	prod.Get("/get_by_name", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductGetByName(c)
	}))
	
	prod.Get("/get_by_code", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductGetByCode(c)
	}))

	prod.Get("/get_all", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductGetAll(c)
	}))

	prod.Post("/create", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductCreate(c)
	}))

	prod.Put("/update", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductUpdate(c)
	}))

	prod.Put("/update_price", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductPriceUpdate(c)
	}))

	prod.Delete("/delete/:id", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductDelete(c)
	}))

	prod.Get("/get_by_category/:id", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductGetByCategoryID(c)
	}))

	prod.Get("/:id", GetController("ProductController", func(c *fiber.Ctx, ctrl *controllers.ProductController) error {
		return ctrl.ProductGetByID(c)
	}))

}
