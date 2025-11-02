package middleware

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func TenantMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*models.AuthenticatedUser)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
				Status:  false,
				Body:    nil,
				Message: "Teanant No autenticado",
			})
		}

		if user.TenantID == nil {
			return c.Status(401).JSON(models.Response{
				Status:  false,
				Body:    nil,
				Message: "login tenant is required",
			})
		}

		return c.Next()
	}
}
