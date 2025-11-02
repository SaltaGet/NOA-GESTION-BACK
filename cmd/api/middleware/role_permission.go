package middleware

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func RolePermissionMiddleware(code string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*models.AuthenticatedUser)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
				Status:  false,
				Body:    nil,
				Message: "No autenticado",
			})
		}

		if user.IsAdminTenant {
			return c.Next()
		}

		for _, permission := range user.Permissions {
			if permission == code {
				return c.Next()
			}
		} 

		return c.Status(fiber.StatusForbidden).JSON(models.Response{
			Status:  false,
			Body:    nil,
			Message: "No autorizado",
		})
	}
}
