package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func IncomeOtherRoutes(app *fiber.App){
	incomeSale := app.Group("/api/v1/income_other", middleware.AuthMiddleware(), middleware.InjectionDependsTenant(), middleware.AuthPointSaleMiddleware())

	incomeSale.Get("/get_by_date", GetController("IncomeOtherController", func(c *fiber.Ctx, ctrl *controllers.IncomeOtherController) error {
		return ctrl.IncomeOtherGetByDate(c)
	}))

	incomeSale.Post("/create", GetController("IncomeOtherController", func(c *fiber.Ctx, ctrl *controllers.IncomeOtherController) error {
		return ctrl.IncomeOtherCreate(c)
	}))

	incomeSale.Put("/update", GetController("IncomeOtherController", func(c *fiber.Ctx, ctrl *controllers.IncomeOtherController) error {
		return ctrl.IncomeOtherUpdate(c)
	}))

	incomeSale.Delete("/delete/:id", GetController("IncomeOtherController", func(c *fiber.Ctx, ctrl *controllers.IncomeOtherController) error {
		return ctrl.IncomeOtherDelete(c)
	}))

	incomeSale.Get("/:id", GetController("IncomeOtherController", func(c *fiber.Ctx, ctrl *controllers.IncomeOtherController) error {
		return ctrl.IncomeOtherGetByID(c)
	}))

}
