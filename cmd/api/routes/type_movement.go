package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func TypeMovementRoutes(app *fiber.App){
	typeMovement := app.Group("/api/v1/type_movement", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	typeMovement.Get("/get_all", GetController("TypeMovementController", func(c *fiber.Ctx, ctrl *controllers.TypeMovementController) error {
		return ctrl.TypeMovementGetAll(c)
	}))

	typeMovement.Post("/create", GetController("TypeMovementController", func(c *fiber.Ctx, ctrl *controllers.TypeMovementController) error {
		return ctrl.TypeMovementCreate(c)
	}))

	typeMovement.Put("/update", GetController("TypeMovementController", func(c *fiber.Ctx, ctrl *controllers.TypeMovementController) error {
		return ctrl.TypeMovementUpdate(c)
	}))
}
