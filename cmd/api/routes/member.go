package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func MemberRoutes(app *fiber.App){
	member := app.Group("/api/v1/member", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	member.Get("/get_all", 
	middleware.RolePermissionMiddleware("MB04"),
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.MemberController.MemberGetAll(c)
	})

	member.Post("/create", 
	middleware.RolePermissionMiddleware("MB01"),	
	middleware.CurrentPlan(), 
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.MemberController.MemberCreate(c)
	})

	member.Put("/update", 
	middleware.RolePermissionMiddleware("MB02"),	
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.MemberController.MemberUpdate(c)
	})
	
	member.Delete("/delete/:id", 
	middleware.RolePermissionMiddleware("MB03"),	
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.MemberController.MemberDelete(c)
	})

	member.Get("/get/:id", 
	middleware.RolePermissionMiddleware("MB04"),	
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.MemberController.MemberGetByID(c)
	})
}
