package middleware

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/key"
	"github.com/gofiber/fiber/v2"
)

func InjectionDepends(deps *dependencies.MainContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(key.AppKey, deps)
		return c.Next()
	}
}