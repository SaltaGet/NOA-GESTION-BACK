package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func ExpenseBuyRoutes(app *fiber.App) {
	exp := app.Group("/api/v1/expense_buy", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	exp.Get("/get_by_date", 
	middleware.RolePermissionMiddleware("EB04"),
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ExpenseBuyController.ExpenseBuyGetByDate(c)
	})

	exp.Post("/create", 
	middleware.RolePermissionMiddleware("EB01"),
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ExpenseBuyController.ExpenseBuyCreate(c)
	})

	exp.Put("/update", 
	middleware.RolePermissionMiddleware("EB02"),
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ExpenseBuyController.ExpenseBuyUpdate(c)
	})

	exp.Delete("/delete/:id", 
	middleware.RolePermissionMiddleware("EB03"),
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ExpenseBuyController.ExpenseBuyDelete(c)
	})

	exp.Get("/:id", 
	middleware.RolePermissionMiddleware("EB04"),
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ExpenseBuyController.ExpenseBuyGetByID(c)
	})
}
