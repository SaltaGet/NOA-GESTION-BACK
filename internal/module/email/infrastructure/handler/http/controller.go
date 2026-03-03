package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/module/email/domain"
)

type EmailController struct {
	service domain.EmailService
}

func NewEmailController(service domain.EmailService) *EmailController {
	return &EmailController{
		service: service,
	}
}

func (c *EmailController) Create(ctx *fiber.Ctx) error {
	// TODO: implement logic
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create Email",
	})
}
