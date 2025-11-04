package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func ClientRoutes(app *fiber.App){
	cli := app.Group("/client", middleware.AuthMiddleware(), middleware.PointSaleMiddleware())

	cli.Get("/get_all", GetController("ClientController", func(c *fiber.Ctx, ctrl *controllers.ClientController) error {
		return ctrl.ClientGetAll(c)
	}))

	cli.Get("/get_by_name", GetController("ClientController", func(c *fiber.Ctx, ctrl *controllers.ClientController) error {
		return ctrl.ClientGetByName(c)
	}))

	cli.Post("/create", GetController("ClientController", func(c *fiber.Ctx, ctrl *controllers.ClientController) error {
		return ctrl.CreateClient(c)
	}))

	cli.Put("/update", GetController("ClientController", func(c *fiber.Ctx, ctrl *controllers.ClientController) error {
		return ctrl.ClientUpdate(c)
	}))

	cli.Delete("/delete/:id", GetController("ClientController", func(c *fiber.Ctx, ctrl *controllers.ClientController) error {
		return ctrl.ClientDelete(c)
	}))

	cli.Get("/:id", GetController("ClientController", func(c *fiber.Ctx, ctrl *controllers.ClientController) error {
		return ctrl.ClientGetByID(c)
	}))

}
