package repositories

import (
	"fmt"
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (r *RoleRepository) RoleGetAll() (*[]schemas.RoleResponse, error) {
	var rows []schemas.RolePermissionRow
	if err := r.DB.Table("roles").
		Select(`roles.id as role_id, roles.name as role_name, permissions.id as perm_id, permissions.code as perm_code, permissions."group" as perm_group`).
		Joins("left join role_permissions on roles.id = role_permissions.role_id").
		Joins("left join permissions on permissions.id = role_permissions.permission_id").
		Find(&rows).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al obtener roles", err)
	}

	roleMap := make(map[string]*schemas.RoleResponse)
	for _, row := range rows {
		role, exists := roleMap[row.RoleID]
		if !exists {
			idInt, err := strconv.ParseInt(row.RoleID, 10, 64)
			if err != nil {
				return nil, schemas.ErrorResponse(500, "Error interno al obtener roles", err)
			}

			role = &schemas.RoleResponse{
				ID:          idInt,
				Name:        row.RoleName,
				Permissions: []schemas.PermissionResponseDTO{},
			}
			roleMap[row.RoleID] = role
		}

		idPerm, err := strconv.ParseInt(row.PermID, 10, 64)
		if err != nil {
			return nil, schemas.ErrorResponse(500, "Error interno al obtener roles", err)
		}
			role.Permissions = append(role.Permissions, schemas.PermissionResponseDTO{
				ID:      idPerm,
				Code:    row.PermCode,
				Group:   row.PermGroup,
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