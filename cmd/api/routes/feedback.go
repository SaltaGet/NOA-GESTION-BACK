package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func FeedbackRoutes(app *fiber.App, controllers *controllers.FeedbackController) {
	feedback := app.Group("/api/v1/feedback")

	feedback.Get("/get_all", middleware.AdminAuthMiddleware(), controllers.FeedbackGetAll)

	feedback.Post("/create", middleware.AuthMiddleware(), controllers.FeedbackCreate)

	feedback.Get("/get/:id", middleware.AdminAuthMiddleware(), controllers.FeedbackGetByID)
	
}
