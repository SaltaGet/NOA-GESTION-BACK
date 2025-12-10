package middleware

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/jobs"
	"github.com/gofiber/fiber/v2"
)

func ReadOnlyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if jobs.IsReadOnly() && c.Method() != fiber.MethodGet && c.Method() != fiber.MethodHead && c.Method() != fiber.MethodOptions {
			return c.Status(503).JSON(fiber.Map{
				"error": "Server in read-only mode",
				"reason": "Backup in progress",
			})
		}
		return c.Next()
	}
}