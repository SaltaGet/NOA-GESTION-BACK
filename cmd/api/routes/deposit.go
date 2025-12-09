package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func DepositRoutes(app *fiber.App) {
	deposit := app.Group("/api/v1/deposit", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	deposit.Get("/get_all",
		middleware.RolePermissionMiddleware("DEP04"),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.DepositController.DepositGetAll(c)
		})

	deposit.Put("/update_stock",
		middleware.RolePermissionMiddleware("DEP02"),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.DepositController.DepositUpdateStock(c)
		})

	deposit.Get("/get_by_name",
		middleware.RolePermissionMiddleware("DEP04"),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.DepositController.DepositGetByName(c)
		})

	deposit.Get("/get_by_code",
		middleware.RolePermissionMiddleware("DEP04"),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.DepositController.DepositGetByCode(c)
		})

	deposit.Get("/get/:id",
		middleware.RolePermissionMiddleware("DEP04"),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.DepositController.DepositGetByID(c)
		})
}
