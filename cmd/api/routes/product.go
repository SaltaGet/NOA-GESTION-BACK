package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App) {
	prod := app.Group("/api/v1/product", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	prod.Get("/get_all",
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ProductController.ProductGetAll(c)
		})

	prod.Get("/get_by_name",
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ProductController.ProductGetByName(c)
		})

	prod.Get("/get_by_code",
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ProductController.ProductGetByCode(c)
		})

	prod.Post("/create",
		middleware.RolePermissionMiddleware("PR01"),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ProductController.ProductCreate(c)
		})

	prod.Put("/update",
		middleware.RolePermissionMiddleware("PR02"),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ProductController.ProductUpdate(c)
		})

	prod.Put("/list_price",
		middleware.RolePermissionMiddleware("PR02"),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ProductController.ProductPriceUpdate(c)
		})

	prod.Delete("/delete/:id",
		middleware.RolePermissionMiddleware("PR03"),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ProductController.ProductDelete(c)
		})

	prod.Get("/get_by_category/:category_id",
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ProductController.ProductGetByCategoryID(c)
		})

	prod.Get("/get/:id",
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ProductController.ProductGetByID(c)
		})
}
