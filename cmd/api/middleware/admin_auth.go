package middleware

import (
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/key"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AdminAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("access_token_admin")
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

		userID := getIntClaim(mapClaims, "admin_id")

		deps, ok := c.Locals(key.AppKey).(*dependencies.MainContainer)
		if !ok {
			return schemas.HandleError(c, fmt.Errorf("error al obtener dependencias"))
		}

		user, err := deps.AuthController.AuthService.AuthAdminGetByID(userID)
		if err != nil {
			return schemas.HandleError(c, err)
		}

		c.Locals("user_admin", user)
		return c.Next()
	}
}
