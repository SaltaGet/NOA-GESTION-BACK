package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func EcommerceRoutes(app *fiber.App) {
	ecommerce := app.Group("/api/v1/ecommerce", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	ecommerce.Get("/get_all",
		// middleware.RolePermissionMiddleware("ECO04"),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.EcommerceController.EcommerceGetAll(c)
		})

	ecommerce.Put("/update_status",
		// middleware.RolePermissionMiddleware("ECO04"),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.EcommerceController.EcommerceUpdateStatus(c)
		})

	ecommerce.Get("/get/:id",
		// middleware.RolePermissionMiddleware("ECO02"),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.EcommerceController.EcommerceGetByID(c)
		})

	ecommerce.Get("/get_by_reference/:reference",
		// middleware.RolePermissionMiddleware("ECO04"),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.EcommerceController.EcommerceGetByReference(c)
		})
}
