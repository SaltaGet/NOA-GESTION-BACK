package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func CashRegisterRoutes(app *fiber.App) {
	cashRegister := app.Group("/api/v1/cash_register", middleware.AuthMiddleware(), middleware.InjectionDependsTenant(), middleware.AuthPointSaleMiddleware())

	cashRegister.Get("/exist_open", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.CashRegisterController.CashRegisterExistOpen(c)
	})

	cashRegister.Post("/open", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.CashRegisterController.CashRegisterOpen(c)
	})
	cashRegister.Get("/inform", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.CashRegisterController.CashRegiterInform(c)
	})

	cashRegister.Post("/close", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.CashRegisterController.CashRegisterClose(c)
	})

	cashRegister.Get("/get/:id", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.CashRegisterController.CashRegisterGetByID(c)
	})

}
