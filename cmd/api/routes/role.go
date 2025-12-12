package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func RoleRoutes(app *fiber.App) {
	role := app.Group("/api/v1/role", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	role.Get("/get_all", 
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.RoleController.RoleGetAll(c)
	})

	role.Post("/create", 
	middleware.RolePermissionMiddleware("RL01"),
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.RoleController.RoleCreate(c)
	})
	
	role.Put("/update", 
	middleware.RolePermissionMiddleware("RL02"),
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.RoleController.RoleUpdate(c)
	})
	
	role.Get("/get/:id", 
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.RoleController.RoleGetByID(c)
	})
}

