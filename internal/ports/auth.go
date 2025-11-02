package ports

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
)

type AuthRepository interface {
	AuthLogin(username, password string, connection string) (result *models.AuthResult, err error)
	// AuthLoginMember(username, password, connection string) (member *models.Member, role *models.Role, permissions *[]string, err error)
	AuthGetTenant(userID string, tenantID string) (tenant *models.Tenant, err error)
	CurrentUser(userID string) (user *models.User, err error)
	// UserGetRolePermissions(connection, userID string) (member *models.Member, role *models.Role, permissions *[]string, err error)
}

type AuhtService interface {
	AuthLogin(username string , password string) (token string, err error)
	AuthGetTenant(user *models.AuthenticatedUser, tenantID string) (token string, err error)
	CurrentUser(userID string) (user *models.User, err error)
}
