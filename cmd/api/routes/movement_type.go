package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func MovementRoutes(app *fiber.App){
	mov := app.Group("/api/v1/movement", middleware.AuthMiddleware(), middleware.AuthPointSaleMiddleware())

	mov.Get("/get_all", GetController("MovementTypeController", func(c *fiber.Ctx, ctrl *controllers.MovementTypeController) error {
		return ctrl.GetAllMovementTypes(c)
	}))

	mov.Post("/create", GetController("MovementTypeController", func(c *fiber.Ctx, ctrl *controllers.MovementTypeController) error {
		return ctrl.MovementTypeCreate(c)
	}))

	mov.Put("/update", GetController("MovementTypeController", func(c *fiber.Ctx, ctrl *controllers.MovementTypeController) error {
		return ctrl.MovementTypeUpdate(c)
	}))

	mov.Delete("/delete/:id", GetController("MovementTypeController", func(c *fiber.Ctx, ctrl *controllers.MovementTypeController) error {
		return ctrl.MovementTypeDelete(c)
	}))

	mov.Get("/:id", GetController("MovementTypeController", func(c *fiber.Ctx, ctrl *controllers.MovementTypeController) error {
		return ctrl.GetMovementTypeByID(c)
	}))

}
