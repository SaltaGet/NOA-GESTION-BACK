package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func ClientRoutes(app *fiber.App) {
	cli := app.Group("/api/v1/client", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	cli.Get("/get_all", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ClientController.ClientGetAll(c)
	})

	cli.Get("/", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ClientController.ClientGetByFilter(c)
	})

	cli.Post("/create", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ClientController.ClientCreate(c)
	})

	cli.Put("/update", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ClientController.ClientUpdate(c)
	})

	cli.Delete("/delete/:id", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ClientController.ClientDelete(c)
	})

	cli.Get("/:id", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ClientController.ClientGetByID(c)
	})
}
