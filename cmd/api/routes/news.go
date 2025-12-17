package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func NewsRoutes(app *fiber.App, controllers *controllers.NewsController) {
	news := app.Group("/api/v1/news")

	news.Get("/get_all", controllers.NewsGetAll)

	news.Post("/create", middleware.AdminAuthMiddleware(), controllers.NewsCreate)

	news.Put("/update", middleware.AdminAuthMiddleware(), controllers.NewsUpdate)

	news.Get("/get/:id", controllers.NewsGetByID)
	
	news.Delete("/delete/:id", middleware.AdminAuthMiddleware(), controllers.NewsDelete)
}
