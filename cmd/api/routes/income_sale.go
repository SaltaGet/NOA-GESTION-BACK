package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func IncomeSaleRoutes(app *fiber.App){
	incomeSale := app.Group("/api/v1/income_sale", middleware.AuthMiddleware(), middleware.InjectionDependsTenant(), middleware.AuthPointSaleMiddleware())

	incomeSale.Get("/get_by_date", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.IncomeSaleController.IncomeSaleGetByDate(c)
	})

	incomeSale.Post("/create", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.IncomeSaleController.IncomeSaleCreate(c)
	})

	incomeSale.Put("/update", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.IncomeSaleController.IncomeSaleUpdate(c)
	})

	incomeSale.Delete("/delete/:id", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.IncomeSaleController.IncomeSaleDelete(c)
	})

	incomeSale.Get("/:id", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.IncomeSaleController.IncomeSaleGetByID(c)
	})
}
