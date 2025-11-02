package routes

import (
	"github.com/DanielChachagua/GestionCar/cmd/api/controllers"
	"github.com/DanielChachagua/GestionCar/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func PurchaseOrderRoutes(app *fiber.App){
	order := app.Group("/purchase_order", middleware.AuthMiddleware(), middleware.TenantMiddleware())

	order.Get("/get_all", GetController("PurchaseOrderController", func(c *fiber.Ctx, ctrl *controllers.PurchaseOrderController) error {
		return ctrl.PurchaseOrderGetAll(c)
	}))

	order.Post("/create", GetController("PurchaseOrderController", func(c *fiber.Ctx, ctrl *controllers.PurchaseOrderController) error {
		return ctrl.PurchaseOrderCreate(c)
	}))

	// order.Put("/update", GetController("PurchaseOrderController", func(c *fiber.Ctx, ctrl *controllers.PurchaseOrderController) error {
	// 	return ctrl.PurchaseOrderUpdate(c)
	// }))

	order.Delete("/delete/:id", GetController("PurchaseOrderController", func(c *fiber.Ctx, ctrl *controllers.PurchaseOrderController) error {
		return ctrl.PurchaseOrderDelete(c)
	}))

	order.Get("/:id", GetController("PurchaseOrderController", func(c *fiber.Ctx, ctrl *controllers.PurchaseOrderController) error {
		return ctrl.PurchaseOrderGetByID(c)
	}))

}
