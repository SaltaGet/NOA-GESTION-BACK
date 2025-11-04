package middleware

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

func AdminTenantMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*schemas.AuthenticatedUser)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(schemas.Response{
				Status:  false,
				Body:    nil,
				Message: "Teanant No autenticado",
			})
		}
		
		if user.IsAdminTenant {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "No autorizado",
		})
	}
}