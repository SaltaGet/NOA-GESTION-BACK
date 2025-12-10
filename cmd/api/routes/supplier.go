package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func SupplierRoutes(app *fiber.App) {
	supplier := app.Group("/api/v1/supplier", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	supplier.Get("/get_all",
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.SupplierController.SupplierGetAll(c)
		})

	supplier.Post("/create",
		middleware.RolePermissionMiddleware("SP01"),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.SupplierController.SupplierCreate(c)
		})

	supplier.Put("/update",
		middleware.RolePermissionMiddleware("SP02"),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.SupplierController.SupplierUpdate(c)
		})

	supplier.Delete("/delete/:id",
		middleware.RolePermissionMiddleware("SP03"),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.SupplierController.SupplierDeleteByID(c)
		})

	supplier.Get("/:id",
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.SupplierController.SupplierGetByID(c)
		})
}
