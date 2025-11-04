package middleware

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func InjectionDepends(deps *dependencies.MainContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("deps", deps)
		return c.Next()
	}
}