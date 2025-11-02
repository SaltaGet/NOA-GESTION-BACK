package repositories

import (
	"fmt"

	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/google/uuid"
)

func (r *RoleRepository) RoleGetAll() (*[]models.RoleResponse, error) {
	var rows []models.RolePermissionRow
	if err := r.DB.Table("roles").
		Select(`roles.id as role_id, roles.name as role_name, permissions.id as perm_id, permissions.code as perm_code, permissions.details as perm_details, permissions."group" as perm_group`).
		Joins("left join role_permissions on roles.id = role_permissions.role_id").
		Joins("left join permissions on permissions.id = role_permissions.permission_id").
		Find(&rows).Error; err != nil {
		return nil, models.ErrorResponse(500, "Error interno al obtener roles", err)
	}
	
	roleMap := make(map[string]*models.RoleResponse)
	for _, row := range rows {
		role, exists := roleMap[row.RoleID]
		if !exists {
			role = &models.RoleResponse{
				ID:          row.RoleID,
				Name:        row.RoleName,
				Permissions: []models.PermissionResponse{},
			}
			roleMap[row.RoleID] = role
		}
		if row.PermID != "" {
			role.Permissions = append(role.Permissions, models.PermissionResponse{
				ID:      row.PermID,
				Code:    row.PermCode,
				Details: row.PermDetails,
				Group:   row.PermGroup,
			})
		}
	}
	var allRoles []models.RoleResponse
	for _, role := range roleMap {
		allRoles = append(allRoles, *role)
	}
	return &allRoles, nil
}

func (t *RoleRepository) RoleCreate(roleCreate *models.RoleCreate) (string, error) {
	var permissions []models.Permission
	if err := t.DB.Where("id IN ?", roleCreate.PermissionsID).Find(&permissions).Error; err != nil {
		return "", models.ErrorResponse(500, "Error interno al buscar permisos", err)
	}
	if len(permissions) != len(roleCreate.PermissionsID) {
		return "", models.ErrorResponse(400, "Algunos permisos no existen", fmt.Errorf("se esperaban %d permisos, pero se encontraron %d", len(roleCreate.PermissionsID), len(permissions)))
	}

	newID := uuid.NewString()
	err := t.DB.Create(&models.Role{ID: newID, Name: roleCreate.Name, Permissions: permissions}).Error
	if err != nil {
		return "", models.ErrorResponse(500, "Error interno al crear el rol", err)
	}
	return newID, nil
}
