package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func ExpenseBuyRoutes(app *fiber.App){
	exp := app.Group("/api/v1/expense_buy", middleware.AuthMiddleware(), middleware.InjectionDependsTenant(), middleware.AuthPointSaleMiddleware())

	exp.Get("/get_by_date", GetController("ExpenseBuyController", func(c *fiber.Ctx, ctrl *controllers.ExpenseBuyController) error {
		return ctrl.ExpenseBuyGetByDate(c)
	}))

	exp.Post("/create", GetController("ExpenseBuyController", func(c *fiber.Ctx, ctrl *controllers.ExpenseBuyController) error {
		return ctrl.ExpenseBuyCreate(c)
	}))

	exp.Put("/update", GetController("ExpenseBuyController", func(c *fiber.Ctx, ctrl *controllers.ExpenseBuyController) error {
		return ctrl.ExpenseBuyUpdate(c)
	}))

	exp.Delete("/delete/:id", GetController("ExpenseBuyController", func(c *fiber.Ctx, ctrl *controllers.ExpenseBuyController) error {
		return ctrl.ExpenseBuyDelete(c)
	}))

	exp.Get("/:id", GetController("ExpenseBuyController", func(c *fiber.Ctx, ctrl *controllers.ExpenseBuyController) error {
		return ctrl.ExpenseBuyGetByID(c)
	}))

}
