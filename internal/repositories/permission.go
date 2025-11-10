package repositories

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)


func (t *PermissionRepository) PermissionByRoleID(roleID int64) (*[]string, error) {
	var permission []string
	err := t.DB.Model(&models.Permission{}).
		Select("permissions.code").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Pluck("permissions.code", &permission).Error
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al obtener permisos", err)
	}
	return &permission, nil
}

func (t *PermissionRepository) PermissionGetAll() (*[]schemas.PermissionResponse, error) {
	var permission []schemas.PermissionResponse
	err := t.DB.Model(&models.Permission{}).Select(`permissions.id, permissions.code, permissions.details, permissions."group"`).
		Find(&permission).Error
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al obtener permisos", err)
	}
	return &permission, nil
}

func (t *PermissionRepository) PermissionGetToMe(roleID int64) (*[]schemas.PermissionResponse, error) {
	var permissions []schemas.PermissionResponse
	err := t.DB.
	  Select(`permissions.id, permissions.code, permissions.details, permissions."group"`).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al obtener permisos", err)
	}
	return &permissions, nil
}
