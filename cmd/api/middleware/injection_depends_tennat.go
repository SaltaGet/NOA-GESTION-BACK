package middleware

import (
	tenant_cache "github.com/SaltaGet/NOA-GESTION-BACK/internal/cache/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

func InjectionDependsTenant() fiber.Handler {
	return func(c *fiber.Ctx) error {
		member := c.Locals("user").(*schemas.AuthenticatedUser)

		db, err := database.GetTenantDB("", member.TenantID)
		if err != nil {
			return schemas.ErrorResponse(401, "No autenticado", err)
		}

		// ⚡ Usar el container cacheado
		container := tenant_cache.GetTenantContainer(db, member.TenantID)

		// Guardarlo en el contexto
		c.Locals("tenant", container)

		return c.Next()
	}
}
