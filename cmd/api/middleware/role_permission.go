package middleware

import (
	"slices"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

func RolePermissionMiddleware(code string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		member, ok := c.Locals("user").(*schemas.AuthenticatedUser)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(schemas.Response{
				Status:  false,
				Body:    nil,
				Message: "No autenticado",
			})
		}

		if member.IsAdmin {
			return c.Next()
		}

		if slices.Contains(member.ListPermissions, code) {
			return c.Next()
		}

		return c.Status(403).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "No autorizado",
		})
	}
}
