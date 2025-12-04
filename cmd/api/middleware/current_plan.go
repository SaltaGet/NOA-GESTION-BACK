package middleware

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/key"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

func CurrentPlan() fiber.Handler {
	return func(c *fiber.Ctx) error {
		member := c.Locals("user").(*schemas.AuthenticatedUser)

		deps := c.Locals(key.AppKey).(*dependencies.MainContainer)
		plan, err := deps.AuthController.AuthService.AuthCurrentPlan(member.TenantID)
		if err != nil {
			return schemas.HandleError(c, err)
		}

		c.Locals("current_plan", plan)

		return c.Next()
	}
}