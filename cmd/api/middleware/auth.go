package middleware

import (
	"fmt"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/cache"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/key"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	// "github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("access_token")
		if token == "" {
			return schemas.HandleError(c, schemas.ErrorResponse(401, "No autenticado", fmt.Errorf("no autenticado")))
		}

		claims, err := utils.VerifyToken(token)
		if err != nil {
			return schemas.HandleError(c, schemas.ErrorResponse(401, "Token inválido", err))
		}

		mapClaims, ok := claims.(jwt.MapClaims)
		if !ok {
			return schemas.HandleError(c, schemas.ErrorResponse(401, "Claims inválidos", fmt.Errorf("claims invalidos")))
		}

		tenantID := utils.GetIntClaim(mapClaims, "tenant_id")
		pointSaleID := utils.GetIntClaim(mapClaims, "point_sale_id")
		memberID := utils.GetIntClaim(mapClaims, "member_id")

		deps, ok := c.Locals(key.AppKey).(*dependencies.MainContainer)
		if !ok {
			return schemas.HandleError(c, fmt.Errorf("error al obtener dependencias"))
		}

		// 🔥 Estrategia simple con cache
		authUser, err := getAuthUserWithCache(tenantID, memberID, pointSaleID, deps)
		if err != nil {
			return schemas.HandleError(c, err)
		}

		c.Locals("user", authUser)
		c.Locals("point_sale_id", pointSaleID)
		return c.Next()
	}
}

// getAuthUserWithCache versión simplificada sin versión
func getAuthUserWithCache(tenantID, memberID, pointSaleID int64, deps *dependencies.MainContainer) (*schemas.AuthenticatedUser, error) {
	// 1️⃣ Intentar obtener del cache de Redis (sin versión)
	// if cache.IsAvailable() {
	// 	authUser, err := cache.GetAuthUser(memberID, tenantID)
	// 	if err == nil {
	// 		logging.INFO("Obteniendo usuario autenticado desde cache")
	// 		return authUser, nil
	// 	}
	// }

	// 2️⃣ No está en cache, consultar DB
	authUser, err := deps.AuthController.AuthService.AuthCurrentUser(tenantID, memberID, pointSaleID)
	if err != nil {
		return nil, err
	}

	// 3️⃣ Guardar en cache para futuros requests
	if cache.IsAvailable() {
		_ = cache.SetAuthUser(authUser)
	}

	return authUser, nil
}

// RateLimitMiddleware middleware de rate limiting
func RateLimitMiddleware(maxRequests int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Identificar por IP o por usuario autenticado
		identifier := c.IP()

		if user := c.Locals("user"); user != nil {
			if authUser, ok := user.(*schemas.AuthenticatedUser); ok {
				identifier = fmt.Sprintf("user:%d", authUser.ID)
			}
		}

		allowed, err := cache.CheckRateLimit(identifier, maxRequests, window)
		if err != nil {
			// En caso de error de Redis, permitir
			return c.Next()
		}

		if !allowed {
			return c.Status(429).JSON(fiber.Map{
				"error": "Demasiadas peticiones, intenta más tarde",
			})
		}

		return c.Next()
	}
}
