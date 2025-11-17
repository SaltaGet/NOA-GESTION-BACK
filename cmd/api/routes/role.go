package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func RoleRoutes(app *fiber.App){
	role := app.Group("/api/v1/role", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	role.Get("/get_all", GetController("RoleController", func(c *fiber.Ctx, ctrl *controllers.RoleController) error {
		return ctrl.RoleGetAll(c)
	}))

	role.Post("/create", GetController("RoleController", func(c *fiber.Ctx, ctrl *controllers.RoleController) error {
		return ctrl.RoleCreate(c)
	}))

}
