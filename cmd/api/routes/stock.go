package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func StockRoutes(app *fiber.App){
	prod := app.Group("/api/v1/stock", middleware.AuthMiddleware(), middleware.InjectionDependsTenant(), middleware.AuthPointSaleMiddleware())

	prod.Get("/get_all", GetController("StockController", func(c *fiber.Ctx, ctrl *controllers.StockController) error {
		return ctrl.StockGetAll(c)
	}))

	prod.Get("/get_by_name", GetController("StockController", func(c *fiber.Ctx, ctrl *controllers.StockController) error {
		return ctrl.StockGetByName(c)
	}))
	
	prod.Get("/get_by_code", GetController("StockController", func(c *fiber.Ctx, ctrl *controllers.StockController) error {
		return ctrl.StockGetByCode(c)
	}))

	prod.Get("/get_by_category/:category_id", GetController("StockController", func(c *fiber.Ctx, ctrl *controllers.StockController) error {
		return ctrl.StockGetByCategoryID(c)
	}))

	prod.Get("/get/:id", GetController("StockController", func(c *fiber.Ctx, ctrl *controllers.StockController) error {
		return ctrl.StockGetByID(c)
	}))
}
