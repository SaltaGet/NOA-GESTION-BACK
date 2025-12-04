package ports

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type AuthRepository interface {
	AuthAdminGetByUsername(username string) (*models.Admin, error)
	AuthAdminGetByID(id int64) (*models.Admin, error)
	AuthTenantGetByID(tenantID int64) (*models.Tenant, error)
	AuthTenantGetByIdentifier(identifier string) (*models.Tenant, error)
	AuthMemberGetByUserID(userID int64, connection string, tenantID int64) (*models.Member, error)
	AuthMemberGetByID(id int64, connection string, tenantID int64) (*models.Member, *[]string, error)
	AuthMemberGetByUsername(username string, connection string, tenantID int64) (*models.Member, error)
	AuthPointSale(pointSaleID int64, connection string, tenantID, memberID int64) (*models.PointSale, error)
	AuthForgotPassword(forgotPassword *schemas.AuthForgotPassword) (*models.Member, *models.Tenant, error)
	AuthResetPassword(memberID, tenantID int64, newPass string) error
}

type AuhtService interface {
	AuthAdminGetByID(id int64) (*models.Admin, error)
	AuthLogin(username, password string) (string, error)
	AuthLoginAdmin(username, password string) (string, error)
	AuthCurrentUser(tenantID, memberID, pointSaleID int64) (*schemas.AuthenticatedUser, error)
	AuthCurrentPlan(tenantID int64) (*schemas.PlanResponseDTO, error)
	AuthCurrentTenant(tenantID int64) (*schemas.TenantResponse, error)
	AuthPointSale(member *schemas.AuthenticatedUser, pointSaleID int64) (string, error)
	LogoutPointSale(member *schemas.AuthenticatedUser) (string, error)
	AuthForgotPassword(forgotPassword *schemas.AuthForgotPassword) error
	AuthResetPassword(resetPassword *schemas.AuthResetPassword) error
}
