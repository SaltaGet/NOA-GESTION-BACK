package services

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/DanielChachagua/GestionCar/pkg/utils"
)

func (a *AuthService) AuthLogin(username, password string) (string, error) {
	var authResult models.AuthResult
	userName, identifier := utils.ParseUsername(username)
	if identifier != "" {
		tenant, err := a.TenantService.TenantGetByIdentifier(identifier)
		if err != nil {
			return "", err
		}

		connection, err := utils.Decrypt(tenant.Connection)
		if err != nil {
			return "", err
		}

		authResultPtr, err := a.AuthRepository.AuthLogin(username, password, connection)
		if err != nil {
			return "", err
		}
		authResult = *authResultPtr
		authResult.Tenant = tenant
	} else {
			authResultPtr, err := a.AuthRepository.AuthLogin(userName, password, "")
		if err != nil {
			return "", err
		}
		authResult = *authResultPtr
	}

	token, err := utils.GenerateToken(&authResult)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) AuthGetTenant(user *models.AuthenticatedUser, tenantID string) (string, error) {
	tenant, err := a.AuthRepository.AuthGetTenant(user.ID, tenantID)
	if err != nil {
		return "", err
	}

	token, err := utils.GenerateToken(&models.AuthResult{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Username:   user.Username,
		IsAdmin:    user.IsAdminTenant,
		Tenant:     tenant,
		Role:       nil,
		Permissions: nil,
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) CurrentUser(userID string) (*models.User, error) {
	user, err := a.AuthRepository.CurrentUser(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// func (a *AuthService) AuthWorkplace(id string) (string, error) {
// 	workplace, err := repositories.Repo.GetWorkplaceByID(id)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return "", models.ErrorResponse(404, "Lugar de trabajo no encontrado", err)
// 		}
// 		return "", models.ErrorResponse(500, "Error al buscar lugar de trabajo", err)
// 	}

// 	token, err := utils.GenerateWorkplaceToken(workplace)

// 	if err != nil {
// 		return "", models.ErrorResponse(500, "Error al generar token", err)
// 	}

// 	return token, nil
// }

// func (a *AuthService) CurrentUser(userId string) (*models.User, error) {
// 	user, err := a.UserRepository.UserGetByID(userId)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, models.ErrorResponse(404, "Usuario no encontrado", err)
// 		}
// 		return nil, models.ErrorResponse(500, "Error al buscar usuario", err)
// 	}

// 	return user, nil
// }

// func (a *AuthService) CurrentWorkplace(workplaceId string) (*models.Workplace, error) {
// 	workplace, err := repositories.Repo.GetWorkplaceByID(workplaceId)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, models.ErrorResponse(404, "Usuario no encontrado", err)
// 		}
// 		return nil, models.ErrorResponse(500, "Error al buscar usuario", err)
// 	}

// 	return workplace, nil
// }

// func GetWorkplaceByRole(role string) (*models.Workplace, error) {
// 	workplace, err := repositories.Repo.GetWorkplaceByRole(role)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, models.ErrorResponse(404, "Rol no encontrado", err)
// 		}
// 		return nil, models.ErrorResponse(500, "Error al buscar rol", err)
// 	}

// 	return workplace, nil
// }
