package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func ReportRoutes(app *fiber.App) {
	report := app.Group("/api/v1/report", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	report.Get("/get_excel", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ReportController.ReportExcelGet(c)
	})

	report.Post("/get_profitable_products", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ReportController.ReportProfitableProducts(c)
	})

	report.Post("/get_by_date", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ReportController.ReportMovementByDate(c)
	})
}
