package ports

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/master"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type AuthRepository interface {
	AuthAdminGetByUsername(username string) (*master.Admin, error)
	AuthAdminGetByID(id int64) (*master.Admin, error)
	AuthTenantGetByID(tenantID int64) (*master.Tenant, error)
	AuthTenantGetByIdentifier(identifier string) (*master.Tenant, error)
	AuthMemberGetByUserID(userID int64, connection string, tenantID int64) (*tenant.Member, error)
	AuthMemberGetByID(id int64, connection string, tenantID int64) (*tenant.Member, *tenant.PermissionSlice, error)
	AuthMemberGetByUsername(username string, connection string, tenantID int64) (*tenant.Member, error)
	AuthPointSale(pointSaleID int64, connection string, tenantID, memberID int64) (*tenant.PointSale, error)
	AuthForgotPassword(forgotPassword *schemas.AuthForgotPassword) (*tenant.Member, *master.Tenant, error)
	AuthResetPassword(memberID, tenantID int64, newPass string) error
}

type AuthService interface {
	AuthAdminGetByID(id int64) (*master.Admin, error)
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
