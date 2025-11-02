package repositories

import "github.com/DanielChachagua/GestionCar/pkg/models"

func (t *PermissionRepository) PermissionByRoleID(roleID string) (*[]string, error) {
	var permission []string
	err := t.DB.Model(&models.Permission{}).
		Select("permissions.code").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Pluck("permissions.code", &permission).Error
	if err != nil {
		return nil, models.ErrorResponse(500, "Error interno al obtener permisos", err)
	}
	return &permission, nil
}

func (t *PermissionRepository) PermissionGetAll() (*[]models.PermissionResponse, error) {
	var permission []models.PermissionResponse
	err := t.DB.Model(&models.Permission{}).Select(`permissions.id, permissions.code, permissions.details, permissions."group"`).
		Find(&permission).Error
	if err != nil {
		return nil, models.ErrorResponse(500, "Error interno al obtener permisos", err)
	}
	return &permission, nil
}

func (t *PermissionRepository) PermissionGetToMe(roleID string) (*[]models.PermissionResponse, error) {
	var permissions []models.PermissionResponse
	err := t.DB.
	  Select(`permissions.id, permissions.code, permissions.details, permissions."group"`).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	if err != nil {
		return nil, models.ErrorResponse(500, "Error interno al obtener permisos", err)
	}
	return &permissions, nil
}
