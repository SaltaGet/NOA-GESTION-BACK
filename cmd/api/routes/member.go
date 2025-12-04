package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func MemberRoutes(app *fiber.App){
	member := app.Group("/api/v1/member", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	member.Get("/get_all", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.MemberController.MemberGetAll(c)
	})

	member.Post("/create", middleware.CurrentPlan(), func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.MemberController.MemberCreate(c)
	})

	member.Put("/update", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.MemberController.MemberUpdate(c)
	})

	member.Get("/get/:id", func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.MemberController.MemberGetByID(c)
	})
}
