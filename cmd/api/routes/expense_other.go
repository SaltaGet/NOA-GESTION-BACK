package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func ExpenseOtherRoutes(app *fiber.App){
	incomeSale := app.Group("/api/v1/expense_other", middleware.AuthMiddleware(), middleware.InjectionDependsTenant(), middleware.AuthPointSaleMiddleware())

	incomeSale.Get("/get_by_date", GetController("ExpenseOtherController", func(c *fiber.Ctx, ctrl *controllers.ExpenseOtherController) error {
		return ctrl.ExpenseOtherGetByDate(c)
	}))

	incomeSale.Post("/create", GetController("ExpenseOtherController", func(c *fiber.Ctx, ctrl *controllers.ExpenseOtherController) error {
		return ctrl.ExpenseOtherCreate(c)
	}))

	incomeSale.Put("/update", GetController("ExpenseOtherController", func(c *fiber.Ctx, ctrl *controllers.ExpenseOtherController) error {
		return ctrl.ExpenseOtherUpdate(c)
	}))

	incomeSale.Delete("/delete/:id", GetController("ExpenseOtherController", func(c *fiber.Ctx, ctrl *controllers.ExpenseOtherController) error {
		return ctrl.ExpenseOtherDelete(c)
	}))

	incomeSale.Get("/:id", GetController("ExpenseOtherController", func(c *fiber.Ctx, ctrl *controllers.ExpenseOtherController) error {
		return ctrl.ExpenseOtherGetByID(c)
	}))

}
