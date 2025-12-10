package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func TypeMovementRoutes(app *fiber.App) {
	typeMovement := app.Group("/api/v1/type_movement", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	typeMovement.Get("/get_all", 
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.TypeMovementController.TypeMovementGetAll(c)
	})

	typeMovement.Post("/create", 
	middleware.RolePermissionMiddleware("TM01"),
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.TypeMovementController.TypeMovementCreate(c)
	})

	typeMovement.Put("/update", 
	middleware.RolePermissionMiddleware("TM02"),
	func(c *fiber.Ctx) error {
		tenant := c.Locals("tenant").(*dependencies.TenantContainer)
		return tenant.Controllers.TypeMovementController.TypeMovementUpdate(c)
	})
}

