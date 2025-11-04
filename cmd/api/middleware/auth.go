// package middleware

// import (
// 	"fmt"

// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/key"
// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/golang-jwt/jwt/v5"
// )

// func AuthMiddleware() fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		token := c.Cookies("access_token")
// 		if token == "" {
// 			return schemas.HandleError(c, schemas.ErrorResponse(401, "No  autenticado", fmt.Errorf("No autenticado")))
// 		}

// 		claims, err := utils.VerifyToken(token)
// 		if err != nil {
// 			return schemas.HandleError(c, schemas.ErrorResponse(401, "Token inválido", err))
// 		}

// 		mapClaims, ok := claims.(jwt.MapClaims)
// 		if !ok {
// 			return schemas.HandleError(c, schemas.ErrorResponse(401, "Claims inválidos", err))
// 		}

// 		userID := getIntClaim(mapClaims, "user_id")
// 		tenantID := getIntClaim(mapClaims, "tenant_id")
// 		pointSaleID := getIntClaim(mapClaims, "point_sale_id")
// 		memberID := getIntClaim(mapClaims, "member_id")

// 		deps, ok := c.Locals(key.AppKey).(*dependencies.MainContainer)
// 		if !ok {
// 			return schemas.HandleError(c, fmt.Errorf("error al obtener dependencias"))
// 		}

// 		authUser, err := deps.AuthController.AuthService.AuthCurrentUser(userID, tenantID, memberID, pointSaleID)
// 		if err != nil {
// 			return schemas.HandleError(c, err)
// 		}

// 		c.Locals("user", &authUser)
// 		return c.Next()
// 	}
// }

// func getIntClaim(claims jwt.MapClaims, key string) int64 {
// 	val, ok := claims[key].(float64)
// 	if ok {
// 		return int64(val)
// 	}
// 	return -1
// }

package middleware

import (
	"fmt"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/cache"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/key"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
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
			return schemas.HandleError(c, schemas.ErrorResponse(401, "Claims inválidos", nil))
		}

		userID := getIntClaim(mapClaims, "user_id")
		tenantID := getIntClaim(mapClaims, "tenant_id")
		pointSaleID := getIntClaim(mapClaims, "point_sale_id")
		memberID := getIntClaim(mapClaims, "member_id")
		version := getIntClaim(mapClaims, "version")

		deps, ok := c.Locals(key.AppKey).(*dependencies.MainContainer)
		if !ok {
			return schemas.HandleError(c, fmt.Errorf("error al obtener dependencias"))
		}

		// 🔥 Estrategia con cache
		authUser, err := getAuthUserWithCache(userID, tenantID, memberID, pointSaleID, version, deps)
		if err != nil {
			return schemas.HandleError(c, err)
		}

		c.Locals("user", authUser)
		return c.Next()
	}
}

// getAuthUserWithCache intenta obtener del cache, fallback a DB
func getAuthUserWithCache(userID, tenantID, memberID, pointSaleID, version int64, deps *dependencies.MainContainer) (*schemas.AuthenticatedUser, error) {
	// 1️⃣ Intentar obtener del cache de Redis
	if cache.IsAvailable() {
		authUser, err := cache.GetAuthUser(userID, version)
		if err == nil {
			return authUser, nil
		}
		// Si hay error, continuar a DB
	}

	// 2️⃣ No está en cache, consultar DB
	authUser, err := deps.AuthController.AuthService.AuthCurrentUser(userID, tenantID, memberID, pointSaleID)
	if err != nil {
		return nil, err
	}

	// 3️⃣ Guardar en cache para futuros requests
	if cache.IsAvailable() {
		_ = cache.SetAuthUser(authUser, version)
	}

	return authUser, nil
}

func getIntClaim(claims jwt.MapClaims, key string) int64 {
	val, ok := claims[key].(float64)
	if ok {
		return int64(val)
	}
	return -1
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