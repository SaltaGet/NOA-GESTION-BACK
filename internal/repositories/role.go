package repositories

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func (r *RoleRepository) RoleGetByID(id int64) (*schemas.RoleResponse, error) {
	var role models.Role
	if err := r.DB.
		Preload("Permissions"). // ← Agregar esto
		Where("roles.id = ?", id).
		First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Rol no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al obtener rol", err)
	}

	var roleResponse schemas.RoleResponse
	copier.Copy(&roleResponse, &role)

	return &roleResponse, nil
}

func (r *RoleRepository) RoleGetAll() (*[]schemas.RoleResponse, error) {
	var rows []schemas.RolePermissionRow
	if err := r.DB.Table("roles").
		Select("roles.id as role_id, roles.name as role_name, permissions.id as perm_id, permissions.code as perm_code, permissions.`group` as perm_group, permissions.environment as environment, permissions.details as detail").
		Joins("left join role_permissions on roles.id = role_permissions.role_id").
		Joins("left join permissions on permissions.id = role_permissions.permission_id").
		Find(&rows).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al obtener roles", err)
	}

	roleMap := make(map[string]*schemas.RoleResponse)
	for _, row := range rows {
		role, exists := roleMap[strconv.FormatInt(row.RoleID, 10)]
		if !exists {
			idInt := row.RoleID

			role = &schemas.RoleResponse{
				ID:          idInt,
				Name:        row.RoleName,
				Permissions: []schemas.PermissionResponse{},
			}
			roleMap[strconv.FormatInt(row.RoleID, 10)] = role
		}

		role.Permissions = append(role.Permissions, schemas.PermissionResponse{
			ID:          row.PermID,
			Code:        row.PermCode,
			Group:       row.PermGroup,
			Environment: row.Environment,
			Details: row.Detail,
		})
	}
	var allRoles []schemas.RoleResponse
	for _, role := range roleMap {
		allRoles = append(allRoles, *role)
	}
	return &allRoles, nil
}

func (t *RoleRepository) RoleCreate(roleCreate *schemas.RoleCreate) (int64, error) {
	var permissions []models.Permission
	if err := t.DB.Where("id IN ?", roleCreate.PermissionsID).Find(&permissions).Error; err != nil {
		return 0, schemas.ErrorResponse(500, "Error interno al buscar permisos", err)
	}
	if len(permissions) != len(roleCreate.PermissionsID) {
		return 0, schemas.ErrorResponse(400, "Algunos permisos no existen", fmt.Errorf("se esperaban %d permisos, pero se encontraron %d", len(roleCreate.PermissionsID), len(permissions)))
	}

	newRole := &models.Role{Name: roleCreate.Name, Permissions: permissions}

	err := t.DB.Create(&newRole).Error
	if err != nil {
		return 0, schemas.ErrorResponse(500, "Error interno al crear el rol", err)
	}
	return newRole.ID, nil
}

func (t *RoleRepository) RoleUpdate(roleUpdate *schemas.RoleUpdate) error {
	// Verificar que el rol existe
	var existingRole models.Role
	if err := t.DB.First(&existingRole, roleUpdate.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "Rol no encontrado", err)
		}
		return schemas.ErrorResponse(500, "Error interno al buscar rol", err)
	}

	// Verificar que todos los permisos existen
	var permissions []models.Permission
	if err := t.DB.Where("id IN ?", roleUpdate.PermissionsID).Find(&permissions).Error; err != nil {
		return schemas.ErrorResponse(500, "Error interno al buscar permisos", err)
	}
	if len(permissions) != len(roleUpdate.PermissionsID) {
		return schemas.ErrorResponse(400, "Algunos permisos no existen", 
			fmt.Errorf("se esperaban %d permisos, pero se encontraron %d", 
				len(roleUpdate.PermissionsID), len(permissions)))
	}

	// Actualizar el rol dentro de una transacción
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		// Actualizar el nombre del rol
		if err := tx.Model(&existingRole).Update("name", roleUpdate.Name).Error; err != nil {
			return err
		}

		// Reemplazar las asociaciones de permisos
		if err := tx.Model(&existingRole).Association("Permissions").Replace(permissions); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return schemas.ErrorResponse(500, "Error interno al actualizar el rol", err)
	}

	return nil
}