package routes

import (
	"github.com/gofiber/fiber/v2"
)

func ExpenseRoutes(app *fiber.App){
	// exp := app.Group("/api/v1/expense", middleware.AuthMiddleware(), middleware.AuthPointSaleMiddleware())

	// exp.Get("/get_all", GetController("ExpenseController", func(c *fiber.Ctx, ctrl *controllers.ExpenseBuyController) error {
	// 	return ctrl.GetAllExpenses(c)
	// }))

	// exp.Get("/get_today", GetController("ExpenseController", func(c *fiber.Ctx, ctrl *controllers.ExpenseBuyController) error {
	// 	return ctrl.GetExpenseToday(c)
	// }))

	// exp.Post("/create", GetController("ExpenseController", func(c *fiber.Ctx, ctrl *controllers.ExpenseBuyController) error {
	// 	return ctrl.CreateExpense(c)
	// }))

	// exp.Put("/update", GetController("ExpenseController", func(c *fiber.Ctx, ctrl *controllers.ExpenseBuyController) error {
	// 	return ctrl.UpdateExpense(c)
	// }))

	// exp.Delete("/delete/:id", GetController("ExpenseController", func(c *fiber.Ctx, ctrl *controllers.ExpenseBuyController) error {
	// 	return ctrl.DeleteExpense(c)
	// }))

	// exp.Get("/:id", GetController("ExpenseController", func(c *fiber.Ctx, ctrl *controllers.ExpenseBuyController) error {
	// 	return ctrl.GetExpenseByID(c)
	// }))

}
