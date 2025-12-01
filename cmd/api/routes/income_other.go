package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func IncomeOtherRoutes(app *fiber.App){
	incomeOther := app.Group("/api/v1/income_other", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	incomeOther.Get("/get_by_date", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.IncomeOtherController.IncomeOtherGetByDate(c)
	})

	incomeOther.Get("/get_by_date_point_sale", middleware.AuthPointSaleMiddleware(), func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.IncomeOtherController.IncomeOtherGetByDateByPointSale(c)
	})

	incomeOther.Post("/create", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.IncomeOtherController.IncomeOtherCreate(c)
	})

	incomeOther.Post("/create_point_sale", middleware.AuthPointSaleMiddleware(), func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.IncomeOtherController.IncomeOtherCreateByPointSale(c)
	})

	incomeOther.Put("/update", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.IncomeOtherController.IncomeOtherUpdate(c)
	})

	incomeOther.Put("/update_point_sale", middleware.AuthPointSaleMiddleware(), func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.IncomeOtherController.IncomeOtherUpdateByPointSale(c)
	})

	incomeOther.Delete("/delete/:id", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.IncomeOtherController.IncomeOtherDelete(c)
	})

	incomeOther.Delete("/delete_point_sale/:id", middleware.AuthPointSaleMiddleware(), func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.IncomeOtherController.IncomeOtherDeleteByPointSale(c)
	})

	incomeOther.Get("get/:id", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.IncomeOtherController.IncomeOtherGetByID(c)
	})

	incomeOther.Get("get_point_sale/:id", middleware.AuthPointSaleMiddleware(), func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.IncomeOtherController.IncomeOtherGetByIDByPointSale(c)
	})
}
