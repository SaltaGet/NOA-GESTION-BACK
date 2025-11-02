package repositories

import (
	"errors"
	"fmt"

	"github.com/DanielChachagua/GestionCar/pkg/database"
	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/DanielChachagua/GestionCar/pkg/utils"
	"gorm.io/gorm"
)

func (r *MainRepository) AuthLogin(username, password string, connection string) (*models.AuthResult, error) {
	if connection != "" {
		db, err := database.GetTenantDB(connection)
		if err != nil {
			return nil, models.ErrorResponse(500, "No se pudo conectar a la base de datos del tenant", err)
		}

		var member models.Member
		if err := db.Where("username = ? AND is_deleted = false", username).Preload("Role").First(&member).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, models.ErrorResponse(401, "Credenciales incorrectas", err)
			}
			return nil, models.ErrorResponse(500, "Error al buscar el miembro", err)
		}

		if !utils.CheckPasswordHash(password, member.Password) {
			return nil, models.ErrorResponse(401, "Credenciales incorrectas", fmt.Errorf("credenciales incorrectas"))
		}

		var permissions []string
		for _, p := range member.Role.Permissions {
			permissions = append(permissions, p.Code)
		}

		return &models.AuthResult{
			ID:          member.ID,
			FirstName:   member.FirstName,
			LastName:    member.LastName,
			Username:    member.Username,
			IsAdmin:     false,
			Tenant:      nil, 
			Role:        &member.Role,
			Permissions: &permissions,
		}, nil

	} else {
		var user models.User
		if err := r.DB.Where("username = ?", username).Preload("UserTenants").First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, models.ErrorResponse(401, "Credenciales incorrectas", err)
			}
			return nil, models.ErrorResponse(500, "Error al buscar el usuario", err)
		}

		if !utils.CheckPasswordHash(password, user.Password) {
			return nil, models.ErrorResponse(401, "Credenciales incorrectas", fmt.Errorf("credenciales incorrectas"))
		}

		return &models.AuthResult{
			ID:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Username:    user.Username,
			IsAdmin:     true,
			Tenant:      nil, 
			Role:        nil,
			Permissions: nil,
		}, nil
	}
}

// func (r *MainRepository) AuthLoginMember(username, password, connection string) (*models.Member, *models.Role, *[]string, error) {
// 	db, err := database.GetTenantDB(connection)
// 	if err != nil {
// 		return nil, nil, nil, models.ErrorResponse(500, "Error al recibir la conexión", err)
// 	}
// 	var member models.Member
// 	if err := db.Where("username = ?", username).First(&member).Error; err != nil {
// 		return nil, nil, nil, models.ErrorResponse(401, "Credenciales incorrectas", err)
// 	}

// 	var role models.Role
// 	if err := db.Where("id = ?", member.RoleID).First(&role).Error; err != nil {
// 		return nil, nil, nil, models.ErrorResponse(500, "Error al obtenr rol", err)
// 	}

// 	var permissions []string
// 	err = db.Model(&models.Permission{}).
// 		Select("permissions.name").
// 		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
// 		Where("role_permissions.role_id = ?", role.ID).
// 		Pluck("permissions.name", &permissions).Error
// 	if err != nil {
// 		return nil, nil, nil, models.ErrorResponse(500, "Error al obtener los permisos", err)
// 	}

// 	if !utils.CheckPasswordHash(password, member.Password) {
// 		return nil, nil, nil, models.ErrorResponse(401, "Credenciales no válidas", nil)
// 	}

// 	return &member, &role, &permissions, nil
// }

func (r *MainRepository) AuthGetTenant(userID string, tenantID string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.DB.
		Preload("UserTenants", "user_id = ? AND tenant_id = ?", userID, tenantID).
		Where("id = ?", tenantID).
		First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return nil, models.ErrorResponse(500, "Error al obtener tenant", err)
	}

	if !tenant.IsActive {
		return nil, models.ErrorResponse(403, "Tenant esta inactivo", fmt.Errorf("credenciales incorrectas"))
	}

	if len(tenant.UserTenants) == 0 {
		return nil, models.ErrorResponse(403, "No tienes permiso para acceder al tenant", fmt.Errorf("sin permiso para acceder al tenant"))
	}

	if !tenant.UserTenants[0].IsActive {
		return nil, models.ErrorResponse(403, "No tienes permiso para acceder al tenant", fmt.Errorf("sin permiso para acceder al tenant"))

	}

	return &tenant, nil
}

func (r *MainRepository) CurrentUser(userID string) (*models.User, error) {
	var user models.User
	if err := r.DB.Preload("UserTenants").Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrorResponse(401, "No autenticado", err)
		}
		return nil, models.ErrorResponse(500, "Error retrieving user", err)
	}
	return &user, nil
}

// func (r *MainRepository) UserGetRolePermissions(connection, userID string) (*models.Member, *models.Role, *[]string, error) {
// 	db, err := database.GetTenantDB(connection)
// 	if err != nil {
// 		return nil, nil, nil, models.ErrorResponse(500, "Error al conectarse a la base de datos", err)
// 	}

// 	var member models.Member
// 	if err := db.Where("user_id = ?", userID).First(&member).Error; err != nil {
// 		return nil, nil, nil, models.ErrorResponse(401, "Credenciales incorrectas", err)
// 	}

// 	var role models.Role
// 	if err := db.Where("id = ?", member.RoleID).First(&role).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil, nil, models.ErrorResponse(404, "Role not found", err)
// 		}
// 		return nil, nil, nil, models.ErrorResponse(500, "Errorinterno al obtener los roles", err)
// 	}

// 	var permissions []string
// 	err = db.Model(&models.Permission{}).
// 		Select("permissions.name").
// 		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
// 		Where("role_permissions.role_id = ?", role.ID).
// 		Pluck("permissions.name", &permissions).Error
// 	if err != nil {
// 		return nil, nil, nil, models.ErrorResponse(500, "Error al obtener los permisos", err)
// 	}

// 	return &member, &role, &permissions, nil
// }
