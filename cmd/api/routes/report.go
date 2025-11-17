package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func ReportRoutes(app *fiber.App) {
	report := app.Group("/api/v1/report", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())
	report.Get("/get_excel", GetController("ReportController", func(c *fiber.Ctx, ctrl *controllers.ReportController) error {
		return ctrl.ReportExcelGet(c)
	}))
	
	report.Post("/get_profitable_products", GetController("ReportController", func(c *fiber.Ctx, ctrl *controllers.ReportController) error {
		return ctrl.ReportProfitableProducts(c)
	}))

	report.Post("/get_by_date", GetController("ReportController", func(c *fiber.Ctx, ctrl *controllers.ReportController) error {
		return ctrl.ReportMovementByDate(c)
	}))
}