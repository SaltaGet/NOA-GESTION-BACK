package routes

import (
	"github.com/DanielChachagua/GestionCar/cmd/api/controllers"
	"github.com/DanielChachagua/GestionCar/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func PermissionRoutes(app *fiber.App){
	permission := app.Group("/permission", middleware.AuthMiddleware(), middleware.TenantMiddleware())

	permission.Get("/get_all", GetController("PermissionController", func(c *fiber.Ctx, ctrl *controllers.PermissionController) error {
		return ctrl.PermissionGetAll(c)
	}))

	permission.Get("/get_to_me", GetController("PermissionController", func(c *fiber.Ctx, ctrl *controllers.PermissionController) error {
		return ctrl.PermissionGetToMe(c)
	}))

}
