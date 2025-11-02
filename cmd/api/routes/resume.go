package routes

import (
	"github.com/DanielChachagua/GestionCar/cmd/api/controllers"
	"github.com/DanielChachagua/GestionCar/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func ResumeRoutes(app *fiber.App) {
	resume := app.Group("/resume", middleware.AuthMiddleware(), middleware.TenantMiddleware())

	resume.Put("/update/expense", GetController("ResumeController", func(c *fiber.Ctx, ctrl *controllers.ResumeController) error {
		return ctrl.ExpenseResumeUpdate(c)
	}))

	resume.Put("/update/income", GetController("ResumeController", func(c *fiber.Ctx, ctrl *controllers.ResumeController) error {
		return ctrl.IncomeResumeUpdate(c)
	}))

	resume.Get("/get_by_date/expense", GetController("ResumeController", func(c *fiber.Ctx, ctrl *controllers.ResumeController) error {
		return ctrl.ExpenseResumeGetByDateBetween(c)
	}))

	resume.Get("/get_by_date/income", GetController("ResumeController", func(c *fiber.Ctx, ctrl *controllers.ResumeController) error {
		return ctrl.IncomeResumeGetByDateBetween(c)
	}))

	resume.Get("/get/expense/:id", GetController("ResumeController", func(c *fiber.Ctx, ctrl *controllers.ResumeController) error {
		return ctrl.ExpenseResumeGetByID(c)
	}))

	resume.Get("/get/income/:id", GetController("ResumeController", func(c *fiber.Ctx, ctrl *controllers.ResumeController) error {
		return ctrl.IncomeResumeGetByID(c)
	}))

}
