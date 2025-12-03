package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func PointSaleRoutes(app *fiber.App){
	pointSale := app.Group("/api/v1/point_sale", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	pointSale.Get("/get_all", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.PointSaleController.PointSaleGetAll(c)
	})

	pointSale.Get("/get_all_by_member", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.PointSaleController.PointSaleGetAllByMember(c)
	})

	pointSale.Post("/create", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.PointSaleController.PointSaleCreate(c)
	})
	
pointSale.Put("/update", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.PointSaleController.PointSaleUpdate(c)
	})
}
