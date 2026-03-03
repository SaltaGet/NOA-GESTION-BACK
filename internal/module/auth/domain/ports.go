package domain

import (
	modelAdmin "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/admin/domain"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/module/auth/infrastructure/repository"
	modelTenant "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/tenant/domain"
	modelMember "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/member/domain"
	modelPointSale "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/point_sale/domain"
	repositoryPlan "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/plan/infrastructure/repository"
	repositoryTenant "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/tenant/infrastructure/repository"
)

type AuthRepository interface {
	AuthAdminGetByUsername(username string) (*modelAdmin.Admin, error)
	AuthAdminGetByID(id int64) (*modelAdmin.Admin, error)
	AuthTenantGetByID(tenantID int64) (*modelTenant.Tenant, error)
	AuthTenantGetByIdentifier(identifier string) (*modelTenant.Tenant, error)
	AuthMemberGetByUserID(userID int64, connection string, tenantID int64) (*modelMember.Member, error)
	AuthMemberGetByID(id int64, connection string, tenantID int64) (*modelMember.Member, *[]string, error)
	AuthMemberGetByUsername(username string, connection string, tenantID int64) (*modelMember.Member, error)
	AuthPointSale(pointSaleID int64, connection string, tenantID, memberID int64) (*modelPointSale.PointSale, error)
	AuthForgotPassword(forgotPassword *repository.AuthForgotPassword) (*modelMember.Member, *modelTenant.Tenant, error)
	AuthResetPassword(memberID, tenantID int64, newPass string) error
}

type AuthService interface {
	AuthAdminGetByID(id int64) (*modelAdmin.Admin, error)
	AuthLogin(username, password string) (string, error)
	AuthLoginAdmin(username, password string) (string, error)
	AuthCurrentUser(tenantID, memberID, pointSaleID int64) (*repository.AuthenticatedUser, error)
	AuthCurrentPlan(tenantID int64) (*repositoryPlan.PlanResponseDTO, error)
	AuthCurrentTenant(tenantID int64) (*repositoryTenant.TenantResponse, error)
	AuthPointSale(member *repository.AuthenticatedUser, pointSaleID int64) (string, error)
	LogoutPointSale(member *repository.AuthenticatedUser) (string, error)
	AuthForgotPassword(forgotPassword *repository.AuthForgotPassword) error
	AuthResetPassword(resetPassword *repository.AuthResetPassword) error
}
