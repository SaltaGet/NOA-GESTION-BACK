package routes

import (
	// "github.com/DanielChachagua/GestionCar/cmd/api/controllers"
	// "github.com/DanielChachagua/GestionCar/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func PurchaseProductRoutes(app *fiber.App){
	// att := app.Group("/purchase_product", middleware.AuthMiddleware(), middleware.TenantMiddleware())
	// att.Get("/get_purchase/:purchase_id", controllers.PurchaseProductGetAllByPurhcaseID)
	// att.Post("/create", controllers.PurchaseProductCreate)
	// att.Put("/update", controllers.PurchaseProductUpdate)
	// att.Delete("/delete/:id", controllers.PurchaseProductDelete)
	// att.Get("/:id", controllers.PurchaseProductGetByID)
}
// func PurchaseProductRoutes(app *fiber.App, controllers *controllers.PurchaseProductController){
// 	// att := app.Group("/purchase_product", middleware.AuthMiddleware(), middleware.TenantMiddleware())
// 	// att.Get("/get_purchase/:purchase_id", controllers.PurchaseProductGetAllByPurhcaseID)
// 	// att.Post("/create", controllers.PurchaseProductCreate)
// 	// att.Put("/update", controllers.PurchaseProductUpdate)
// 	// att.Delete("/delete/:id", controllers.PurchaseProductDelete)
// 	// att.Get("/:id", controllers.PurchaseProductGetByID)
// }