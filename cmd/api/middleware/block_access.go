package middleware

import (
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

func BlockAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		if strings.Contains(path, "/.") {
			return c.Status(fiber.StatusForbidden).JSON(schemas.Response{
				Status:  false,
				Body:    nil,
				Message: "Acceso denegado",
			})
		}

		return c.Next()
	}
}
