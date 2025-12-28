package utils

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(memberID int64, tenantID, pointSaleID *int64) (string, error) {
	claims := jwt.MapClaims{
		"member_id": memberID,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	if tenantID != nil {
		claims["tenant_id"] = tenantID
	}

	if pointSaleID != nil {
		claims["point_sale_id"] = pointSaleID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func GenerateTokenAdmin(adminID int64) (string, error) {
	claims := jwt.MapClaims{
		"admin_id": adminID,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func VerifyToken(tokenString string) (jwt.Claims, error) {
	cleanToken := CleanToken(tokenString)
	token, err := jwt.Parse(cleanToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, schemas.ErrorResponse(401, "Token inválido", err)
	}

	return token.Claims, nil
}

func CleanToken(bearerToken string) string {
	const prefix = "Bearer "
	if strings.HasPrefix(bearerToken, prefix) {
		return strings.TrimPrefix(bearerToken, prefix)
	}
	return bearerToken
}

func GenerateTokenEmail(memberID, tenantID int64) (string, error) {
	claims := jwt.MapClaims{
		"member_id": memberID,
		"tenant_id": tenantID,
		"exp":       time.Now().Add(30 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_KEY_EMAIL")))
}

func VerifyTokenEmail(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY_EMAIL")), nil
	})

	if err != nil {
		// Si el error es especificamente por expiración
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, schemas.ErrorResponse(401, "Token vencido", err)
		}
		return nil, schemas.ErrorResponse(401, "Token inválido", err)
	}

	// Validar claims y que el token sea válido
	if !token.Valid {
		return nil, schemas.ErrorResponse(401, "Token inválido", nil)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, schemas.ErrorResponse(500, "No se pudieron obtener los claims", nil)
	}

	return claims, nil
}

func GenerateTokenToGrpc(tenantID, productID int64) (string, error) {
	claims := jwt.MapClaims{
		"tenant_identifier": tenantID,
		"product_id":        productID,
		"exp":       time.Now().Add(30 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("KEY_VALIDATOR")))
}

