package http


import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func ArcaRoutes(app *fiber.App) {
	arca := app.Group("/api/v1/arca", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	arca.Post("/emit_invoice/:income_sale_id",
		middleware.AuthPointSaleMiddleware(),
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ArcaController.ArcaEmitInvoice(c)
		})

	arca.Post("/generate_key",
		func(c *fiber.Ctx) error {
			tenant := c.Locals("tenant").(*dependencies.TenantContainer)
			return tenant.Controllers.ArcaController.ArcaGenerateKey(c)
		})

	// arca.Post("/open",
	// 	middleware.RolePermissionMiddleware("CR01"),
	// 	func(c *fiber.Ctx) error {
	// 		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
	// 		return tenant.Controllers.ArcaController.ArcaOpen(c)
	// 	})

	// arca.Get("/inform",
	// 	middleware.RolePermissionMiddleware("CR04"),
	// 	func(c *fiber.Ctx) error {
	// 		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
	// 		return tenant.Controllers.ArcaController.ArcaInform(c)
	// 	})

	// arca.Post("/close",
	// 	middleware.RolePermissionMiddleware("CR01"),
	// 	func(c *fiber.Ctx) error {
	// 		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
	// 		return tenant.Controllers.ArcaController.ArcaClose(c)
	// 	})

	// arca.Get("/get/:id",
	// 	middleware.RolePermissionMiddleware("CR04"),
	// 	func(c *fiber.Ctx) error {
	// 		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
	// 		return tenant.Controllers.ArcaController.ArcaGetByID(c)
	// 	})

}
