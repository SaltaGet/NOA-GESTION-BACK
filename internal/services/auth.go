package services

import (
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
)

func (a *AuthService) AuthLogin(username, password string) (string, error) {
	userName, identifier := utils.ParseUsername(username)

	user, err := a.AuthRepository.AuthUserGetByUsername(userName)
	if err != nil {
		return "", err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", schemas.ErrorResponse(401, "Credenciales incorrectas", fmt.Errorf("credenciales incorrectas"))
	}

	if user.IsAdmin {
		token, err := utils.GenerateToken(user.ID, nil, nil, nil)
		if err != nil {
			return "", err
		}
		return token, nil
	}

	tenant, err := a.AuthRepository.AuthTenantGetByIdentifier(identifier)
	if err != nil {
		return "", err
	}

	connection, err := utils.Decrypt(tenant.Connection)
	if err != nil {
		return "", err
	}

	member, _, err := a.AuthRepository.AuthMemberGetByID(user.ID, connection, tenant.ID)
	if err != nil {
		return "", err
	}

	token, err := utils.GenerateToken(user.ID, &tenant.ID, &member.ID, nil)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) AuthPointSale(user *schemas.AuthenticatedUser, pointSaleID int64) (string, error) {
	tenant, err := a.AuthRepository.AuthTenantGetByIdentifier(*user.TenantIdentifier)
	if err != nil {
		return "", err
	}

	connection, err := utils.Decrypt(tenant.Connection)
	if err != nil {
		return "", err
	}

	pointSale, err := a.AuthRepository.AuthPointSale(user.ID, connection, tenant.ID)
	if err != nil {
		return "", err
	}

	token, err := utils.GenerateToken(user.ID, &tenant.ID, user.MemberID, &pointSale.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) AuthCurrentUser(userID, tenantID, memberID, pointSaleID int64) (*schemas.AuthenticatedUser, error) {
	user, err := a.AuthRepository.AuthUserGetByID(userID)
	if err != nil {
		return nil, err
	}

	authUser := schemas.AuthenticatedUser{
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Username:      user.Username,
		IsAdmin:       user.IsAdmin,
		IsAdminTenant: user.IsAdminTenant,
	}

	if user.IsAdmin {
		return &authUser, nil
	}

	// if tenantID == -1 || memberID == -1 || pointSaleID == -1 {
	// 	return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", fmt.Errorf("credenciales incorrectas"))
	// }

	var tenant *models.Tenant
	var member *models.Member
	var permissions *[]string

	if tenantID != -1 {
		tenant, err = a.AuthRepository.AuthTenantGetByID(tenantID)
		if err != nil {
			return nil, err
		}
		// connection, err := utils.Decrypt(tenant.Connection)
		// if err != nil {
		// 	return nil, err
		// }

		if memberID == -1 {
			member, permissions, err = a.AuthRepository.AuthMemberGetByID(userID, tenant.Connection, tenantID)
			if err != nil {
				return nil, err
			}
		}
	}

	authUser.MemberID = &member.ID
	authUser.RoleID = &member.Role.ID
	authUser.RoleName = &member.Role.Name
	authUser.Permissions = permissions
	authUser.TenantID = &tenant.ID
	authUser.TenantName = &tenant.Name
	authUser.TenantIdentifier = &tenant.Identifier

	return &authUser, nil
}
