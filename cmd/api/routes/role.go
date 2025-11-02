package routes

import (
	"github.com/DanielChachagua/GestionCar/cmd/api/controllers"
	"github.com/DanielChachagua/GestionCar/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func RoleRoutes(app *fiber.App){
	role := app.Group("/role", middleware.AuthMiddleware(), middleware.TenantMiddleware())

	role.Get("/get_all", GetController("RoleController", func(c *fiber.Ctx, ctrl *controllers.RoleController) error {
		return ctrl.RoleGetAll(c)
	}))

	role.Post("/create", GetController("RoleController", func(c *fiber.Ctx, ctrl *controllers.RoleController) error {
		return ctrl.RoleCreate(c)
	}))

}
