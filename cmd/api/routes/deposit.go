package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func DepositRoutes(app *fiber.App) {
	deposit := app.Group("/api/v1/deposit", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	deposit.Get("/get_all", GetController("DepositController", func(c *fiber.Ctx, ctrl *controllers.DepositController) error {
		return ctrl.DepositGetAll(c)
	}))

	deposit.Put("/update_stock", GetController("DepositController", func(c *fiber.Ctx, ctrl *controllers.DepositController) error {
		return ctrl.DepositUpdateStock(c)
	}))

	deposit.Get("/get_by_name", GetController("DepositController", func(c *fiber.Ctx, ctrl *controllers.DepositController) error {
		return ctrl.DepositGetByName(c)
	}))

	deposit.Get("/get_by_code", GetController("DepositController", func(c *fiber.Ctx, ctrl *controllers.DepositController) error {
		return ctrl.DepositGetByCode(c)
	}))

	deposit.Get("/get/:id", GetController("DepositController", func(c *fiber.Ctx, ctrl *controllers.DepositController) error {
		return ctrl.DepositGetByID(c)
	}))
}