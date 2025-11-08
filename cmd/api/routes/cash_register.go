package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, controllers *controllers.CashRegisterController) {
	register := app.Group("/api/v1/register", middleware.AuthMiddleware(), middleware.AuthPointSaleMiddleware())

	register.Get("/exist_open", controllers.CashRegisterExistOpen)
	register.Post("/open", controllers.CashRegisterOpen)
	register.Post("/inform", controllers.CashRegiterInform)
	register.Post("/close", controllers.CashRegisterClose)
	register.Get("/get/:id", controllers.CashRegisterGetByID)
}