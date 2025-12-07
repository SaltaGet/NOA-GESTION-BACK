package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func MovementStockRoutes(app *fiber.App) {
	movementStock := app.Group("/api/v1/movement_stock", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	movementStock.Post("/move_list", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.MovementStockController.MoveStockList(c)
	})

	movementStock.Get("/get_by_date", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.MovementStockController.MovementStockGetByDate(c)
	})

	movementStock.Get("/get/:id", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.MovementStockController.MovementStockGet(c)
	})
}
