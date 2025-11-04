package middleware

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

func PointSaleMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*schemas.AuthenticatedUser)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(schemas.Response{
				Status:  false,
				Body:    nil,
				Message: "Teanant No autenticado",
			})
		}

		if user.MemberID == nil {
			return c.Status(401).JSON(schemas.Response{
				Status:  false,
				Body:    nil,
				Message: "login tenant is required",
			})
		}

		return c.Next()
	}
}
