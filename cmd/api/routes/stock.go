package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func StockRoutes(app *fiber.App) {
	stock := app.Group("/api/v1/stock", middleware.AuthMiddleware(), middleware.InjectionDependsTenant(), middleware.AuthPointSaleMiddleware())

	stock.Get("/get_all", 
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.StockController.StockGetAll(c)
	})

	stock.Get("/get_by_name", 
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.StockController.StockGetByName(c)
	})

	stock.Get("/get_by_code", 
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.StockController.StockGetByCode(c)
	})

	stock.Get("/get_by_category/:category_id", 
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.StockController.StockGetByCategoryID(c)
	})

	stock.Get("/get/:id", 
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.StockController.StockGetByID(c)
	})
}

