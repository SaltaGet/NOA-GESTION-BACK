package middleware

import (
	"errors"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

func AuthModule(moduleName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentPlan := c.Locals("current_plan").(*schemas.PlanResponseDTO)

		for _, module := range currentPlan.Modules {
			if module.Name == moduleName {
				return c.Next()
			}
		}
		return schemas.HandleError(c, schemas.ErrorResponse(403, "No tienes permiso para acceder a este modulo", errors.New("no tienes permisos para acceder al modulo")))
	}
}
