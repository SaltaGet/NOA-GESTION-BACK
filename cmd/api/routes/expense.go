package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/gofiber/fiber/v2"
)

func ExpenseRoutes(app *fiber.App){
	exp := app.Group("/expense", middleware.AuthMiddleware(), middleware.PointSaleMiddleware())

	exp.Get("/get_all", GetController("ExpenseController", func(c *fiber.Ctx, ctrl *controllers.ExpenseController) error {
		return ctrl.GetAllExpenses(c)
	}))

	exp.Get("/get_today", GetController("ExpenseController", func(c *fiber.Ctx, ctrl *controllers.ExpenseController) error {
		return ctrl.GetExpenseToday(c)
	}))

	exp.Post("/create", GetController("ExpenseController", func(c *fiber.Ctx, ctrl *controllers.ExpenseController) error {
		return ctrl.CreateExpense(c)
	}))

	exp.Put("/update", GetController("ExpenseController", func(c *fiber.Ctx, ctrl *controllers.ExpenseController) error {
		return ctrl.UpdateExpense(c)
	}))

	exp.Delete("/delete/:id", GetController("ExpenseController", func(c *fiber.Ctx, ctrl *controllers.ExpenseController) error {
		return ctrl.DeleteExpense(c)
	}))

	exp.Get("/:id", GetController("ExpenseController", func(c *fiber.Ctx, ctrl *controllers.ExpenseController) error {
		return ctrl.GetExpenseByID(c)
	}))

}
