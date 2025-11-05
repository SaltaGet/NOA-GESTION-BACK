package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func PermissionRoutes(app *fiber.App){
	permission := app.Group("/api/v1/permission", middleware.AuthMiddleware(), middleware.PointSaleMiddleware())

	permission.Get("/get_all", GetController("PermissionController", func(c *fiber.Ctx, ctrl *controllers.PermissionController) error {
		return ctrl.PermissionGetAll(c)
	}))

	permission.Get("/get_to_me", GetController("PermissionController", func(c *fiber.Ctx, ctrl *controllers.PermissionController) error {
		return ctrl.PermissionGetToMe(c)
	}))

}
