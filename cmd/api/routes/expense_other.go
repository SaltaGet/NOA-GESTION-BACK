package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func ExpenseOtherRoutes(app *fiber.App) {
	exp := app.Group("/api/v1/expense_other", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	exp.Get("/get_by_date", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ExpenseOtherController.ExpenseOtherGetByDate(c)
	})

	exp.Get("/get_by_date_point_sale",
		middleware.AuthPointSaleMiddleware(),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ExpenseOtherController.ExpenseOtherGetByDatePointSale(c)
		},
	)

	exp.Post("/create", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ExpenseOtherController.ExpenseOtherCreate(c)
	})

	exp.Post("/create_point_sale",
		middleware.AuthPointSaleMiddleware(),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ExpenseOtherController.ExpenseOtherCreatePointSale(c)
		},
	)

	exp.Put("/update", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ExpenseOtherController.ExpenseOtherUpdate(c)
	})

	exp.Put("/update_point_sale",
		middleware.AuthPointSaleMiddleware(),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ExpenseOtherController.ExpenseOtherUpdatePointSale(c)
		},
	)

	exp.Delete("/delete/:id", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ExpenseOtherController.ExpenseOtherDelete(c)
	})

	exp.Delete("/delete_point_sale/:id",
		middleware.AuthPointSaleMiddleware(),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ExpenseOtherController.ExpenseOtherDeletePointSale(c)
		},
	)

	exp.Get("/get/:id", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.ExpenseOtherController.ExpenseOtherGetByID(c)
	})

	exp.Get("/get_point_sale/:id",
		middleware.AuthPointSaleMiddleware(),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ExpenseOtherController.ExpenseOtherGetByIDPointSale(c)
		},
	)
}
