package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func NewsRoutes(app *fiber.App, controllers *controllers.NewsController) {
	news := app.Group("/api/v1/news")

	news.Get("/get_all", middleware.AuthMiddleware(), controllers.NewsGetAll)

	news.Post("/create", middleware.AdminAuthMiddleware(), controllers.NewsGetByID)

	news.Put("/update", middleware.AdminAuthMiddleware(), controllers.NewsUpdate)

	news.Get("/get/:id", middleware.AuthMiddleware(), controllers.NewsGetByID)
	
	news.Put("/delete/:id", middleware.AdminAuthMiddleware(), controllers.NewsDelete)
}
