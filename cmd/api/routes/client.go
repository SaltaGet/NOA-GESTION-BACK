package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func ClientRoutes(app *fiber.App) {
	cli := app.Group("/api/v1/client", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	cli.Get("/get_all", 
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ClientController.ClientGetAll(c)
	})

	cli.Get("/get_by_filter", 
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ClientController.ClientGetByFilter(c)
	})

	cli.Post("/create", 
	middleware.RolePermissionMiddleware("CL01"), 
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ClientController.ClientCreate(c)
	})

	cli.Put("/update", 
	middleware.RolePermissionMiddleware("CL02"), 
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ClientController.ClientUpdate(c)
	})
	
	cli.Put("/update_credit", 
	middleware.RolePermissionMiddleware("CL02"), 
	middleware.AuthPointSaleMiddleware(), 
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ClientController.ClientUpdateCredit(c)
	})

	cli.Delete("/delete/:id", 
	middleware.RolePermissionMiddleware("CL03"), 
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ClientController.ClientDelete(c)
	})

	cli.Get("/get/:id", 
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ClientController.ClientGetByID(c)
	})
}
