package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func PermissionRoutes(app *fiber.App){
	permission := app.Group("/api/v1/permission", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	permission.Get("/get_all", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.PermissionController.PermissionGetAll(c)
	})

	permission.Get("/get_to_me", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.PermissionController.PermissionGetToMe(c)
	})
}
