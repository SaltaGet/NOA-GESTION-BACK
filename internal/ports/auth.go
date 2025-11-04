package ports

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type AuthRepository interface {
	AuthUserGetByID(userID int64) (*models.User, error)
	AuthUserGetByUsername(username string) (*models.User, error)
	AuthTenantGetByID(tenantID int64) (*models.Tenant, error)
	AuthTenantGetByIdentifier(identifier string) (*models.Tenant, error)
	AuthMemberGetByUserID(userID int64, connection string, tenantID int64) (*models.Member, error)
	AuthMemberGetByID(id int64, connection string, tenantID int64) (*models.Member, *[]string, error)
	AuthPointSale(pointSaleID int64, connection string, tenantID int64) (*models.PointSale, error)
}

type AuhtService interface {
	AuthLogin(username, password string) (string, error)
	AuthCurrentUser(userID, tenantID, memberID, pointSaleID int64) (*schemas.AuthenticatedUser, error)
	AuthPointSale(user *schemas.AuthenticatedUser, pointSaleID int64) (string, error)
}
