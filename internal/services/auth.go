package services

import (
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
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

	member, permsissions, err := a.AuthRepository.AuthMemberGetByID(memberID, tenant.Connection, tenantID)
	if err != nil {
		return nil, err
	}

	authUser := schemas.AuthenticatedUser{
		ID:               member.ID,
		FirstName:        member.FirstName,
		LastName:         member.LastName,
		Username:         member.Username,
		IsAdmin:          member.IsAdmin,
		Permissions:      *permsissions,
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