package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func CategoryRoutes(app *fiber.App) {
	category := app.Group("/api/v1/category", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	category.Post("/create", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.CategoryController.CategoryCreate(c)
	})

	category.Get("/get_all", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.CategoryController.CategoryGetAll(c)
	})

	category.Put("/update", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.CategoryController.CategoryUpdate(c)
	})

	category.Get("/get/:id", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.CategoryController.CategoryGet(c)
	})

	category.Delete("/delete/:id", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.CategoryController.CategoryDelete(c)
	})
}
