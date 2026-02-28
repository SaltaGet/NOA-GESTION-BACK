package repositories

import (
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
		return nil, schemas.HandlerErrorGorm(err, "Rol", schemas.Read)
	}

	var roleResponse schemas.RoleResponse
	copier.Copy(&roleResponse, &role)

	return &roleResponse, nil
}

func (r *RoleRepository) RoleGetAll() (*[]schemas.RoleResponse, error) {
	var rows []schemas.RolePermissionRow
	if err := r.DB.Table("roles").
		Select(`
    roles.id as role_id, 
    roles.name as role_name, 
    permissions.id as perm_id, 
    permissions.code as perm_code, 
    permissions."group" as perm_group,
    permissions.environment as environment, 
    permissions.details as detail
`).
		Joins("left join role_permissions on roles.id = role_permissions.role_id").
		Joins("left join permissions on permissions.id = role_permissions.permission_id").
		Find(&rows).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Rol", schemas.Read)
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
			Details:     row.Detail,
		})
	}
	var allRoles []schemas.RoleResponse
	for _, role := range roleMap {
		allRoles = append(allRoles, *role)
	}
	return &allRoles, nil
}

func expandPermissions(db *gorm.DB, permissionIDs []int64) ([]models.Permission, error) {
	// Obtener los permisos solicitados
	var requestedPermissions []models.Permission
	if err := db.Where("id IN ?", permissionIDs).Find(&requestedPermissions).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Rol", schemas.Read)
	}

	// Mapa para evitar duplicados
	permissionMap := make(map[int64]models.Permission)
	for _, perm := range requestedPermissions {
		permissionMap[perm.ID] = perm
	}

	// Para cada permiso de actualización (02), buscar el correspondiente de lectura (04)
	var groupsToExpand []string
	for _, perm := range requestedPermissions {
		// Si el código termina en 02 (update)
		if len(perm.Code) >= 2 && perm.Code[len(perm.Code)-2:] == "02" {
			groupsToExpand = append(groupsToExpand, perm.Group)
		}
	}

	// Si hay grupos para expandir, buscar los permisos 04 correspondientes
	if len(groupsToExpand) > 0 {
		var readPermissions []models.Permission
		// Buscar permisos que terminen en 04 y pertenezcan a los grupos relevantes
		if err := db.Where("code LIKE ? AND `group` IN ?", "%04", groupsToExpand).Find(&readPermissions).Error; err != nil {
			return nil, schemas.HandlerErrorGorm(err, "Rol", schemas.Read)
		}

		// Agregar los permisos de lectura al mapa (evita duplicados automáticamente)
		for _, perm := range readPermissions {
			permissionMap[perm.ID] = perm
		}
	}

	// Convertir el mapa a slice
	expandedPermissions := make([]models.Permission, 0, len(permissionMap))
	for _, perm := range permissionMap {
		expandedPermissions = append(expandedPermissions, perm)
	}

	return expandedPermissions, nil
}

// RoleCreate actualizado
func (t *RoleRepository) RoleCreate(memberID int64, roleCreate *schemas.RoleCreate) (int64, error) {
	var newRoleSave *models.Role
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}
		// Expandir permisos automáticamente
		permissions, err := expandPermissions(tx, roleCreate.PermissionsID)
		if err != nil {
			return err
		}

		// Validar que al menos los permisos solicitados existan
		if len(permissions) < len(roleCreate.PermissionsID) {
			return schemas.ErrorResponse(400, "Algunos permisos no existen",
				fmt.Errorf("se esperaban al menos %d permisos, pero se encontraron %d",
					len(roleCreate.PermissionsID), len(permissions)))
		}

		newRole := &models.Role{Name: roleCreate.Name, Permissions: permissions}

		if err := tx.Create(&newRole).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Rol", schemas.Create)
		}

		newRoleSave = newRole
		return nil
	})

	if err != nil {
		return 0, err
	}

	return newRoleSave.ID, nil
}

// RoleUpdate actualizado
func (t *RoleRepository) RoleUpdate(memberID int64, roleUpdate *schemas.RoleUpdate) error {
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}
		// Verificar que el rol existe
		var existingRole models.Role
		if err := tx.Preload("Permissions").First(&existingRole, roleUpdate.ID).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Rol", schemas.Read)
		}

		// Expandir permisos automáticamente
		permissions, err := expandPermissions(tx, roleUpdate.PermissionsID)
		if err != nil {
			return err
		}

		// Validar que al menos los permisos solicitados existan
		if len(permissions) < len(roleUpdate.PermissionsID) {
			return schemas.ErrorResponse(400, "Algunos permisos no existen",
				fmt.Errorf("se esperaban al menos %d permisos, pero se encontraron %d",
					len(roleUpdate.PermissionsID), len(permissions)))
		}

		// Actualizar el nombre del rol
		if err := tx.Model(&existingRole).Update("name", roleUpdate.Name).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Rol", schemas.Update)
		}

		// Reemplazar las asociaciones de permisos
		if err := tx.Model(&existingRole).Association("Permissions").Replace(permissions); err != nil {
			return schemas.HandlerErrorGorm(err, "Rol", schemas.Update)
		}

		return nil
	})

	return err
}
