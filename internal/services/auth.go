package services

import (
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
)

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
		return "", schemas.ErrorResponse(401, "Credenciales incorrectas", fmt.Errorf("credenciales incorrectas"))
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
		return "", schemas.ErrorResponse(401, "Credenciales incorrectas", fmt.Errorf("credenciales incorrectas"))
	}

	token, err := utils.GenerateTokenAdmin(admin.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) AuthPointSale(member *schemas.AuthenticatedUser, pointSaleID int64) (string, error) {
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

func (a *AuthService) AuthCurrentUser(tenantID, memberID, pointSaleID int64) (*schemas.AuthenticatedUser, error) {
	tenant, err := a.AuthRepository.AuthTenantGetByID(tenantID)
	if err != nil {
		return nil, err
	}

	member, permissions, err := a.AuthRepository.AuthMemberGetByID(memberID, tenant.Connection, tenantID)
	if err != nil {
		return nil, err
	}

	authUser := schemas.AuthenticatedUser{
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
	}

	// if user.IsAdmin {
	// 	return &authUser, nil
	// }

	// if tenantID == -1 || memberID == -1 || pointSaleID == -1 {
	// 	return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", fmt.Errorf("credenciales incorrectas"))
	// }

	return &authUser, nil
}

func (a *AuthService) AuthCurrentPlan(tenantID int64) (*schemas.PlanResponseDTO, error) {
	tenant, err := a.AuthRepository.AuthTenantGetByID(tenantID)
	if err != nil {
		return nil, err
	}

	plan, err := a.PlanRepository.PlanGetByID(tenant.PlanID)
	if err != nil {
		return nil, err
	}

	var planResponse schemas.PlanResponseDTO
	copier.Copy(&planResponse, &plan)
	

	return &planResponse, nil
}

func (a *AuthService) AuthCurrentTenant(tenantID int64) (*schemas.TenantResponse, error) {
	tenant, err := a.AuthRepository.AuthTenantGetByID(tenantID)
	if err != nil {
		return nil, err
	}

	var tenantResponse schemas.TenantResponse
	copier.Copy(&tenantResponse, &tenant)

	return &tenantResponse, nil
}

func (a *AuthService) AuthAdminGetByID(userID int64) (*models.Admin, error) {
	admin, err := a.AuthRepository.AuthAdminGetByID(userID)
	if err != nil {
		return nil, err
	}

	admin.Password = ""
	return admin, nil
}

func (a *AuthService) LogoutPointSale(member *schemas.AuthenticatedUser) (string, error) {
	return utils.GenerateToken(member.ID, &member.TenantID, nil)
}

func BuildUserPermissions(perms []models.Permission) []schemas.EnvironmentPermissions {
	envMap := make(map[string]map[string][]string)

	for _, p := range perms {
		if _, ok := envMap[p.Environment]; !ok {
			envMap[p.Environment] = make(map[string][]string)
		}
		envMap[p.Environment][p.Group] = append(envMap[p.Environment][p.Group], p.Code)
	}

	// convertir map en estructura final
	result := make([]schemas.EnvironmentPermissions, 0)

	for env, groups := range envMap {
		grpList := make([]schemas.GroupPermissions, 0)

		for group, permCodes := range groups {
			grpList = append(grpList, schemas.GroupPermissions{
				Group:       group,
				Permissions: permCodes,
			})
		}

		result = append(result, schemas.EnvironmentPermissions{
			Environment: env,
			Groups:      grpList,
		})
	}

	return result
}

func (a *AuthService) AuthForgotPassword(forgotPassword *schemas.AuthForgotPassword) error {
	member, tenant, err := a.AuthRepository.AuthForgotPassword(forgotPassword)
	if err != nil {
		return err
	}

	token, err := utils.GenerateTokenEmail(member.ID, tenant.ID)
	if err != nil {
		return err
	}

	body := utils.ForgotPassword(member.Username, member.Email, token)

	err = a.EmailService.SendEmail(member, "Restablecimiento de contraseña", body)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) AuthResetPassword(resetPassword *schemas.AuthResetPassword) error {
	claims, err := utils.VerifyTokenEmail(resetPassword.Token)
	if err != nil {
		return err
	}

	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return schemas.ErrorResponse(401, "Claims inválidos", fmt.Errorf("claims invalidos"))
	}

	tenantID := utils.GetIntClaim(mapClaims, "tenant_id")
	memberID := utils.GetIntClaim(mapClaims, "member_id")
	if tenantID == -1 || memberID == -1 {
		return schemas.ErrorResponse(401, "Claims inválidos", fmt.Errorf("claims invalidos"))
	}

	err = a.AuthRepository.AuthResetPassword(memberID, tenantID, resetPassword.NewPassword)
	if err != nil {
		return err
	}

	return nil
}
