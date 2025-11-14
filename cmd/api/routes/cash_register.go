package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func CashRegisterRoutes(app *fiber.App) {
	cashRegister := app.Group("/api/v1/cash_register", middleware.AuthMiddleware(), middleware.InjectionDependsTenant(), middleware.AuthPointSaleMiddleware())

	cashRegister.Get("/exist_open", GetController("CashRegisterController", func(c *fiber.Ctx, ctrl *controllers.CashRegisterController) error {
		return ctrl.CashRegisterExistOpen(c)
	}))

	cashRegister.Post("/open", GetController("CashRegisterController", func(c *fiber.Ctx, ctrl *controllers.CashRegisterController) error {
		return ctrl.CashRegisterOpen(c)
	}))

	cashRegister.Post("/inform", GetController("CashRegisterController", func(c *fiber.Ctx, ctrl *controllers.CashRegisterController) error {
		return ctrl.CashRegiterInform(c)
	}))

	cashRegister.Post("/close", GetController("CashRegisterController", func(c *fiber.Ctx, ctrl *controllers.CashRegisterController) error {
		return ctrl.CashRegisterClose(c)
	}))

	cashRegister.Get("/get/:id", GetController("CashRegisterController", func(c *fiber.Ctx, ctrl *controllers.CashRegisterController) error {
		return ctrl.CashRegisterGetByID(c)
	}))
}
