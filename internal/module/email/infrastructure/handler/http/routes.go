package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/module/email/infrastructure/handler/http"
)

func SetupEmailRoutes(router fiber.Router, controller *http.EmailController) {
	group := router.Group("/email")
	
	group.Post("/", controller.Create)
	// TODO: Add more routes
}
