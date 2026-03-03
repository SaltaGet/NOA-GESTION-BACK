package application

import (
	"fmt"

	domainAdmin "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/admin/domain"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/module/auth/domain"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/module/auth/infrastructure/repository"
	domainE "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/email/domain"
	domainM "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/module/domain"
	domainP "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/plan/domain"
	domainT "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/tenant/domain"
	domainPermission "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/permission/domain"
	repositoryPlan "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/plan/infrastructure/repository"
	repositoryTenant "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/tenant/infrastructure/repository"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/cache"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
)

type AuthService struct {
	AuthRepository   domain.AuthRepository
	TenantService    domainT.TenantService
	EmailService     domainE.EmailService
	PlanRepository   domainP.PlanRepository
	ModuleRepository domainM.ModuleRepository
}

func (a *AuthService) AuthLogin(username, password string) (string, error) {
	userName, identifier := utils.ParseUsername(username)

	tenant, err := a.AuthRepository.AuthTenantGetByIdentifier(identifier)
	if err != nil {
		return "", err
	}

	member, err := a.AuthRepository.AuthMemberGetByUsername(userName, tenant.Connection, tenant.ID)
	if err != nil {
		return "", err
	}

	if !utils.CheckPasswordHash(password, member.Password) {
		return "", utils.ErrorResponse(401, "Credenciales incorrectas", fmt.Errorf("credenciales incorrectas"))
	}

	token, err := utils.GenerateToken(member.ID, &tenant.ID, nil)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) AuthLoginAdmin(username, password string) (string, error) {
	admin, err := a.AuthRepository.AuthAdminGetByUsername(username)
	if err != nil {
		return "", err
	}

	if !utils.CheckPasswordHash(password, admin.Password) {
		return "", utils.ErrorResponse(401, "Credenciales incorrectas", fmt.Errorf("credenciales incorrectas"))
	}

	token, err := utils.GenerateTokenAdmin(admin.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) AuthPointSale(member *repository.AuthenticatedUser, pointSaleID int64) (string, error) {
	tenant, err := a.AuthRepository.AuthTenantGetByIdentifier(member.TenantIdentifier)
	if err != nil {
		return "", err
	}

	pointSale, err := a.AuthRepository.AuthPointSale(pointSaleID, tenant.Connection, tenant.ID, member.ID)
	if err != nil {
		return "", err
	}

	token, err := utils.GenerateToken(member.ID, &tenant.ID, &pointSale.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) AuthCurrentUser(tenantID, memberID, pointSaleID int64) (*repository.AuthenticatedUser, error) {
	tenant, err := a.AuthRepository.AuthTenantGetByID(tenantID)
	if err != nil {
		return nil, err
	}

	member, permissions, err := a.AuthRepository.AuthMemberGetByID(memberID, tenant.Connection, tenantID)
	if err != nil {
		return nil, err
	}

	authUser := repository.AuthenticatedUser{
		ID:               member.ID,
		FirstName:        member.FirstName,
		LastName:         member.LastName,
		Username:         member.Username,
		IsAdmin:          member.IsAdmin,
		Permissions:      BuildUserPermissions(member.Role.Permissions),
		ListPermissions:  *permissions,
		TenantID:         tenant.ID,
		TenantName:       tenant.Name,
		TenantIdentifier: tenant.Identifier,
		RoleID:           member.Role.ID,
		RoleName:         member.Role.Name,
		AcceptedTerms:    tenant.AcceptedTerms,
	}

	// if user.IsAdmin {
	// 	return &authUser, nil
	// }

	// if tenantID == -1 || memberID == -1 || pointSaleID == -1 {
	// 	return nil, utils.ErrorResponse(401, "Credenciales incorrectas", fmt.Errorf("credenciales incorrectas"))
	// }

	return &authUser, nil
}

func (a *AuthService) AuthCurrentPlan(tenantID int64) (*repositoryPlan.PlanResponseDTO, error) {
	tenant, err := a.AuthRepository.AuthTenantGetByID(tenantID)
	if err != nil {
		return nil, err
	}

	plan, err := a.PlanRepository.PlanGetByID(tenant.PlanID)
	if err != nil {
		return nil, err
	}

	modules, err := a.ModuleRepository.ModuleGetByTenantID(tenantID)
	if err != nil {
		return nil, err
	}

	var planResponse repositoryPlan.PlanResponseDTO
	copier.Copy(&planResponse, &plan)
	planResponse.Modules = modules

	return &planResponse, nil
}

func (a *AuthService) AuthCurrentTenant(tenantID int64) (*repositoryTenant.TenantResponse, error) {
	if cache.IsAvailable() {
		if cachedInfo, err := cache.GetTenantInfo(tenantID); err == nil && cachedInfo != nil {
			return cachedInfo, nil
		}
	}

	tenant, err := a.AuthRepository.AuthTenantGetByID(tenantID)
	if err != nil {
		return nil, err
	}

	var tenantResponse repositoryTenant.TenantResponse
	copier.Copy(&tenantResponse, &tenant)

	tenantResponse.ResponsabilityFrontIVA = tenant.Credentials.ResponsibilityFrontIVA

	if cache.IsAvailable() {
		_ = cache.SetTenantInfo(tenantID, &tenantResponse)
	}

	return &tenantResponse, nil
}

func (a *AuthService) AuthAdminGetByID(userID int64) (*domainAdmin.Admin, error) {
	admin, err := a.AuthRepository.AuthAdminGetByID(userID)
	if err != nil {
		return nil, err
	}

	admin.Password = ""
	return admin, nil
}

func (a *AuthService) LogoutPointSale(member *repository.AuthenticatedUser) (string, error) {
	return utils.GenerateToken(member.ID, &member.TenantID, nil)
}

func BuildUserPermissions(perms []domainPermission.Permission) []repository.EnvironmentPermissions {
	envMap := make(map[string]map[string][]string)

	for _, p := range perms {
		if _, ok := envMap[p.Environment]; !ok {
			envMap[p.Environment] = make(map[string][]string)
		}
		envMap[p.Environment][p.Group] = append(envMap[p.Environment][p.Group], p.Code)
	}

	// convertir map en estructura final
	result := make([]repository.EnvironmentPermissions, 0)

	for env, groups := range envMap {
		grpList := make([]repository.GroupPermissions, 0)

		for group, permCodes := range groups {
			grpList = append(grpList, repository.GroupPermissions{
				Group:       group,
				Permissions: permCodes,
			})
		}

		result = append(result, repository.EnvironmentPermissions{
			Environment: env,
			Groups:      grpList,
		})
	}

	return result
}

func (a *AuthService) AuthForgotPassword(forgotPassword *repository.AuthForgotPassword) error {
	member, tenant, err := a.AuthRepository.AuthForgotPassword(forgotPassword)
	if err != nil {
		return err
	}

	token, err := utils.GenerateTokenEmail(member.ID, tenant.ID)
	if err != nil {
		return err
	}

	body := utils.ForgotPassword(member.Username, member.Email, token)

	err = a.EmailService.SendEmail(member.Email, "Restablecimiento de contraseña", body)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) AuthResetPassword(resetPassword *repository.AuthResetPassword) error {
	claims, err := utils.VerifyTokenEmail(resetPassword.Token)
	if err != nil {
		return err
	}

	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return utils.ErrorResponse(401, "Claims inválidos", fmt.Errorf("claims invalidos"))
	}

	tenantID := utils.GetIntClaim(mapClaims, "tenant_id")
	memberID := utils.GetIntClaim(mapClaims, "member_id")
	if tenantID == -1 || memberID == -1 {
		return utils.ErrorResponse(401, "Claims inválidos", fmt.Errorf("claims invalidos"))
	}

	err = a.AuthRepository.AuthResetPassword(memberID, tenantID, resetPassword.NewPassword)
	if err != nil {
		return err
	}

	if cache.IsAvailable() {
		_ = cache.InvalidateAuthUser(memberID, tenantID)
	}

	return nil
}
