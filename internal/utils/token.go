package utils

import (
	"os"
	"strings"
	"time"

	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(result *models.AuthResult) (string, error) {
	claims := jwt.MapClaims{
		"id":         result.ID,
		"first_name": result.FirstName,
		"last_name":  result.LastName,
		"username":   result.Username,
		"is_admin_tenant":   result.IsAdmin,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	}

	if result.Tenant != nil {
		claims["tenant_id"] = result.Tenant.ID
		claims["tenant_name"] = result.Tenant.Name
		claims["identifier"] = result.Tenant.Identifier
	}
	if result.Role != nil {
		claims["role_id"] = result.Role.ID
		claims["role_name"] = result.Role.Name
	}
	if result.Permissions != nil {
		claims["permissions"] = *result.Permissions
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}


// func GenerateToken(user *models.User, tenant *models.Tenant, member *models.Member, role *models.Role, permissions *[]string) (string, error) {
// 	claims := jwt.MapClaims{}
// 	if user == nil {
// 		claims["id"] = member.ID
// 		claims["first_name"] = member.FirstName
// 		claims["last_name"] = member.LastName
// 		claims["username"] = member.Username
// 		claims["is_admin"] = false
// 		claims["exp"] = time.Now().Add(24 * time.Hour).Unix()

// 	} else {
// 		claims["id"] = user.ID
// 		claims["first_name"] = user.FirstName
// 		claims["last_name"] = user.LastName
// 		claims["username"] = user.Username
// 		claims["is_admin"] = true
// 		claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
// 	}

// 	if tenant != nil && role != nil && permissions != nil {
// 		claims["tenant_id"] = tenant.ID
// 		claims["tenant_name"] = tenant.Name
// 		claims["role_id"] = role.ID
// 		claims["role_name"] = role.Name
// 		claims["permissions"] = permissions
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	t, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
// 	if err != nil {
// 		return "", err
// 	}

// 	return t, nil
// }

func VerifyToken(tokenString string) (jwt.Claims, error) {
	cleanToken := CleanToken(tokenString)
	token, err := jwt.Parse(cleanToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, models.ErrorResponse(401, "Token inválido", err)
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
