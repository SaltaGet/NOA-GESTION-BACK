package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func IncomeSaleRoutes(app *fiber.App){
	incomeSale := app.Group("/api/v1/income_sale", middleware.AuthMiddleware(), middleware.InjectionDependsTenant(), middleware.AuthPointSaleMiddleware())

	incomeSale.Get("/get_by_date", GetController("IncomeSaleController", func(c *fiber.Ctx, ctrl *controllers.IncomeSaleController) error {
		return ctrl.IncomeSaleGetByDate(c)
	}))

	incomeSale.Post("/create", GetController("IncomeSaleController", func(c *fiber.Ctx, ctrl *controllers.IncomeSaleController) error {
		return ctrl.IncomeSaleCreate(c)
	}))

	incomeSale.Put("/update", GetController("IncomeSaleController", func(c *fiber.Ctx, ctrl *controllers.IncomeSaleController) error {
		return ctrl.IncomeSaleUpdate(c)
	}))

	incomeSale.Delete("/delete/:id", GetController("IncomeSaleController", func(c *fiber.Ctx, ctrl *controllers.IncomeSaleController) error {
		return ctrl.IncomeSaleDelete(c)
	}))

	incomeSale.Get("/:id", GetController("IncomeSaleController", func(c *fiber.Ctx, ctrl *controllers.IncomeSaleController) error {
		return ctrl.IncomeSaleGetByID(c)
	}))

}
