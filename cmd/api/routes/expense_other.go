package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func ExpenseOtherRoutes(app *fiber.App) {
	exp := app.Group("/api/v1/expense_other", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	exp.Get("/get_by_date", 
	middleware.RolePermissionMiddleware("EO04"),
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ExpenseOtherController.ExpenseOtherGetByDate(c)
	})

	exp.Get("/get_by_date_point_sale",
middleware.RolePermissionMiddleware("EOPS04"),
		middleware.AuthPointSaleMiddleware(),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ExpenseOtherController.ExpenseOtherGetByDatePointSale(c)
		},
	)

	exp.Post("/create", 
	middleware.RolePermissionMiddleware("EO01"),
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ExpenseOtherController.ExpenseOtherCreate(c)
	})

	exp.Post("/create_point_sale",
		middleware.RolePermissionMiddleware("EOPS01"),
		middleware.AuthPointSaleMiddleware(),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ExpenseOtherController.ExpenseOtherCreatePointSale(c)
		},
	)

	exp.Put("/update", 
	middleware.RolePermissionMiddleware("EO02"),
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ExpenseOtherController.ExpenseOtherUpdate(c)
	})

	exp.Put("/update_point_sale",
		middleware.RolePermissionMiddleware("EOPS02"),
		middleware.AuthPointSaleMiddleware(),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ExpenseOtherController.ExpenseOtherUpdatePointSale(c)
		},
	)

	exp.Delete("/delete/:id", 
	middleware.RolePermissionMiddleware("EO03"),
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ExpenseOtherController.ExpenseOtherDelete(c)
	})

	exp.Delete("/delete_point_sale/:id",
		middleware.RolePermissionMiddleware("EOPS03"),
		middleware.AuthPointSaleMiddleware(),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ExpenseOtherController.ExpenseOtherDeletePointSale(c)
		},
	)

	exp.Get("/get/:id", 
	middleware.RolePermissionMiddleware("EO04"),
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ExpenseOtherController.ExpenseOtherGetByID(c)
	})

	exp.Get("/get_point_sale/:id",
		middleware.RolePermissionMiddleware("EOPS04"),
		middleware.AuthPointSaleMiddleware(),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ExpenseOtherController.ExpenseOtherGetByIDPointSale(c)
		},
	)
}
