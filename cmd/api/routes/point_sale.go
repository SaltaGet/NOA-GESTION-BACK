package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func PointSaleRoutes(app *fiber.App){
	pointSale := app.Group("/api/v1/point_sale", middleware.AuthMiddleware(), middleware.InjectionDependsTenant())

	pointSale.Get("/get_all", GetController("PointSaleController", func(c *fiber.Ctx, ctrl *controllers.PointSaleController) error {
		return ctrl.PointSaleGetAll(c)
	}))

	pointSale.Get("/get_all_by_member", GetController("PointSaleController", func(c *fiber.Ctx, ctrl *controllers.PointSaleController) error {
		return ctrl.PointSaleGetAllByMember(c)
	}))
	
	pointSale.Post("/create", GetController("PointSaleController", func(c *fiber.Ctx, ctrl *controllers.PointSaleController) error {
		return ctrl.PointSaleCreate(c)
	}))

}
