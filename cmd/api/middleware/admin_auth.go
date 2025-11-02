package middleware

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func AdminAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*models.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
				Status:  false,
				Body:    nil,
				Message: "Autenticado",
			})
		}
		
		if user.IsAdmin {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(models.Response{
			Status:  false,
			Body:    nil,
			Message: "No autorizado",
		})
	}
}