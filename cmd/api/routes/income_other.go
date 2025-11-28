package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func IncomeOtherRoutes(app *fiber.App){
	incomeSale := app.Group("/api/v1/income_other", middleware.AuthMiddleware())

	incomeSale.Get("/get_by_date", GetController("IncomeOtherController", func(c *fiber.Ctx, ctrl *controllers.IncomeOtherController) error {
		return ctrl.IncomeOtherGetByDate(c)
	}))
	
	incomeSale.Get("/get_by_date_point_sale", GetController("IncomeOtherController", func(c *fiber.Ctx, ctrl *controllers.IncomeOtherController) error {
		return ctrl.IncomeOtherGetByDateByPointSale(c)
	}))

	incomeSale.Post("/create", GetController("IncomeOtherController", func(c *fiber.Ctx, ctrl *controllers.IncomeOtherController) error {
		return ctrl.IncomeOtherCreate(c)
	}))
	
	incomeSale.Post("/create_point_sale", middleware.InjectionDependsTenant(), middleware.AuthPointSaleMiddleware(), GetController("IncomeOtherController", func(c *fiber.Ctx, ctrl *controllers.IncomeOtherController) error {
		return ctrl.IncomeOtherCreateByPointSale(c)
	}))

	incomeSale.Put("/update", GetController("IncomeOtherController", func(c *fiber.Ctx, ctrl *controllers.IncomeOtherController) error {
		return ctrl.IncomeOtherUpdate(c)
	}))
	
	incomeSale.Put("/update_point_sale", middleware.InjectionDependsTenant(), middleware.AuthPointSaleMiddleware(), GetController("IncomeOtherController", func(c *fiber.Ctx, ctrl *controllers.IncomeOtherController) error {
		return ctrl.IncomeOtherUpdateByPointSale(c)
	}))

	incomeSale.Delete("/delete/:id", GetController("IncomeOtherController", func(c *fiber.Ctx, ctrl *controllers.IncomeOtherController) error {
		return ctrl.IncomeOtherDelete(c)
	}))
	
	incomeSale.Delete("/delete_point_sale/:id", middleware.InjectionDependsTenant(), middleware.AuthPointSaleMiddleware(), GetController("IncomeOtherController", func(c *fiber.Ctx, ctrl *controllers.IncomeOtherController) error {
		return ctrl.IncomeOtherDeleteByPointSale(c)
	}))

	incomeSale.Get("get/:id", GetController("IncomeOtherController", func(c *fiber.Ctx, ctrl *controllers.IncomeOtherController) error {
		return ctrl.IncomeOtherGetByID(c)
	}))
	
	incomeSale.Get("get_point_sale/:id", middleware.InjectionDependsTenant(), middleware.AuthPointSaleMiddleware(), GetController("IncomeOtherController", func(c *fiber.Ctx, ctrl *controllers.IncomeOtherController) error {
		return ctrl.IncomeOtherGetByIDByPointSale(c)
	}))

}
