package routes

import (
	"github.com/DanielChachagua/GestionCar/cmd/api/controllers"
	"github.com/DanielChachagua/GestionCar/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func IncomeRoutes(app *fiber.App){
	inc := app.Group("/income", middleware.AuthMiddleware(), middleware.TenantMiddleware())

	inc.Get("/get_all", GetController("IncomeController", func(c *fiber.Ctx, ctrl *controllers.IncomeController) error {
		return ctrl.GetAllIncomes(c)
	}))

	inc.Get("/get_today", GetController("IncomeController", func(c *fiber.Ctx, ctrl *controllers.IncomeController) error {
		return ctrl.GetIncomeToday(c)
	}))

	inc.Post("/create", GetController("IncomeController", func(c *fiber.Ctx, ctrl *controllers.IncomeController) error {
		return ctrl.CreateIncome(c)
	}))

	inc.Put("/update", GetController("IncomeController", func(c *fiber.Ctx, ctrl *controllers.IncomeController) error {
		return ctrl.UpdateIncome(c)
	}))

	inc.Delete("/delete/:id", GetController("IncomeController", func(c *fiber.Ctx, ctrl *controllers.IncomeController) error {
		return ctrl.DeleteIncome(c)
	}))

	inc.Get("/:id", GetController("IncomeController", func(c *fiber.Ctx, ctrl *controllers.IncomeController) error {
		return ctrl.GetIncomeByID(c)
	}))

}
