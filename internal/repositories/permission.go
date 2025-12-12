package repositories

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"gorm.io/gorm"
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
	var permission *[]schemas.PermissionResponse
	err := t.DB.Model(&models.Permission{}).Select("id", "code", "details", "group", "environment").Scan(&permission).Error
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al obtener permisos", err)
	}
	return permission, nil
}

func (t *PermissionRepository) PermissionGetToMe(roleID int64) (*[]schemas.PermissionResponse, error) {
	var permissions []schemas.PermissionResponse
	err := t.DB.
	  Select(`permissions.id, permissions.code, permissions.details, permissions."group"`).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Scan(&permissions).Error
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al obtener permisos", err)
	}
	return &permissions, nil
}

func (t *PermissionRepository) PermissionUpdateAll() error {
	return t.DB.Transaction(func(tx *gorm.DB) error {
		// Obtener todos los permisos existentes en la BD
		var permissionsAll []models.Permission
		if err := tx.Find(&permissionsAll).Error; err != nil {
			return schemas.ErrorResponse(500, "Error interno al obtener permisos", err)
		}

		// Crear mapas para búsqueda eficiente
		definedMap := make(map[string]models.Permission)
		for _, p := range database.Permissions {
			definedMap[p.Code] = p
		}

		existingMap := make(map[string]models.Permission)
		for _, p := range permissionsAll {
			existingMap[p.Code] = p
		}

		// 1. CREAR o ACTUALIZAR permisos que están en la lista definida
		for _, definedPerm := range database.Permissions {
			if existingPerm, exists := existingMap[definedPerm.Code]; exists {
				// Existe: actualizar si hay cambios
				if existingPerm.Name != definedPerm.Name ||
					existingPerm.Details != definedPerm.Details ||
					existingPerm.Group != definedPerm.Group ||
					existingPerm.Environment != definedPerm.Environment {
					
					// Actualizar manteniendo el ID
					definedPerm.ID = existingPerm.ID
					if err := tx.Save(&definedPerm).Error; err != nil {
						return schemas.ErrorResponse(500, "Error al actualizar permiso: "+definedPerm.Code, err)
					}
				}
			} else {
				// No existe: crear
				if err := tx.Create(&definedPerm).Error; err != nil {
					return schemas.ErrorResponse(500, "Error al crear permiso: "+definedPerm.Code, err)
				}
			}
		}

		// 2. ELIMINAR permisos que NO están en la lista definida
		for _, existingPerm := range permissionsAll {
			if _, stillDefined := definedMap[existingPerm.Code]; !stillDefined {
				// No está en la lista definida: eliminar
				if err := tx.Delete(&existingPerm).Error; err != nil {
					return schemas.ErrorResponse(500, "Error al eliminar permiso: "+existingPerm.Code, err)
				}
			}
		}

		return nil
	})
}


