package repositories

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
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
			Details:     row.Detail,
		})
	}
	var allRoles []schemas.RoleResponse
	for _, role := range roleMap {
		allRoles = append(allRoles, *role)
	}
	return &allRoles, nil
}

// func (t *RoleRepository) RoleCreate(roleCreate *schemas.RoleCreate) (int64, error) {
// 	var permissions []models.Permission
// 	if err := t.DB.Where("id IN ?", roleCreate.PermissionsID).Find(&permissions).Error; err != nil {
// 		return 0, schemas.ErrorResponse(500, "Error interno al buscar permisos", err)
// 	}
// 	if len(permissions) != len(roleCreate.PermissionsID) {
// 		return 0, schemas.ErrorResponse(400, "Algunos permisos no existen", fmt.Errorf("se esperaban %d permisos, pero se encontraron %d", len(roleCreate.PermissionsID), len(permissions)))
// 	}

// 	newRole := &models.Role{Name: roleCreate.Name, Permissions: permissions}

// 	err := t.DB.Create(&newRole).Error
// 	if err != nil {
// 		return 0, schemas.ErrorResponse(500, "Error interno al crear el rol", err)
// 	}
// 	return newRole.ID, nil
// }

// func (t *RoleRepository) RoleUpdate(roleUpdate *schemas.RoleUpdate) error {
// 	// Verificar que el rol existe
// 	var existingRole models.Role
// 	if err := t.DB.First(&existingRole, roleUpdate.ID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return schemas.ErrorResponse(404, "Rol no encontrado", err)
// 		}
// 		return schemas.ErrorResponse(500, "Error interno al buscar rol", err)
// 	}

// 	// Verificar que todos los permisos existen
// 	var permissions []models.Permission
// 	if err := t.DB.Where("id IN ?", roleUpdate.PermissionsID).Find(&permissions).Error; err != nil {
// 		return schemas.ErrorResponse(500, "Error interno al buscar permisos", err)
// 	}
// 	if len(permissions) != len(roleUpdate.PermissionsID) {
// 		return schemas.ErrorResponse(400, "Algunos permisos no existen",
// 			fmt.Errorf("se esperaban %d permisos, pero se encontraron %d",
// 				len(roleUpdate.PermissionsID), len(permissions)))
// 	}

// 	// Actualizar el rol dentro de una transacción
// 	err := t.DB.Transaction(func(tx *gorm.DB) error {
// 		// Actualizar el nombre del rol
// 		if err := tx.Model(&existingRole).Update("name", roleUpdate.Name).Error; err != nil {
// 			return err
// 		}

// 		// Reemplazar las asociaciones de permisos
// 		if err := tx.Model(&existingRole).Association("Permissions").Replace(permissions); err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return schemas.ErrorResponse(500, "Error interno al actualizar el rol", err)
// 	}

// 	return nil
// }



// expandPermissions expande los permisos para incluir automáticamente los permisos de lectura (04)
// cuando se otorgan permisos de actualización (02) del mismo grupo
func expandPermissions(db *gorm.DB, permissionIDs []int64) ([]models.Permission, error) {
	// Obtener los permisos solicitados
	var requestedPermissions []models.Permission
	if err := db.Where("id IN ?", permissionIDs).Find(&requestedPermissions).Error; err != nil {
		return nil, err
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
			return nil, err
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
		// Expandir permisos automáticamente
		permissions, err := expandPermissions(tx, roleCreate.PermissionsID)
		if err != nil {
			return schemas.ErrorResponse(500, "Error interno al buscar permisos", err)
		}
		
		// Validar que al menos los permisos solicitados existan
		if len(permissions) < len(roleCreate.PermissionsID) {
			return schemas.ErrorResponse(400, "Algunos permisos no existen", 
				fmt.Errorf("se esperaban al menos %d permisos, pero se encontraron %d", 
					len(roleCreate.PermissionsID), len(permissions)))
		}

		newRole := &models.Role{Name: roleCreate.Name, Permissions: permissions}

		if err := tx.Create(&newRole).Error; err != nil {
			return schemas.ErrorResponse(500, "Error interno al crear el rol", err)
		}

		newRoleSave = newRole
		return nil
	})

	if err != nil {
		return 0, err
	}

	go database.SaveAuditAsync(t.DB, models.AuditLog{
		MemberID: memberID,
		Method:   "create",
		Path:     "role",
	}, nil, newRoleSave)

	return newRoleSave.ID, nil
}

// RoleUpdate actualizado
func (t *RoleRepository) RoleUpdate(memberID int64, roleUpdate *schemas.RoleUpdate) error {
	var oldRole, newRole models.Role
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el rol existe
		var existingRole models.Role
		if err := tx.Preload("Permissions").First(&existingRole, roleUpdate.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Rol no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error interno al buscar rol", err)
		}

		oldRole = existingRole

		// Expandir permisos automáticamente
		permissions, err := expandPermissions(tx, roleUpdate.PermissionsID)
		if err != nil {
			return schemas.ErrorResponse(500, "Error interno al buscar permisos", err)
		}
		
		// Validar que al menos los permisos solicitados existan
		if len(permissions) < len(roleUpdate.PermissionsID) {
			return schemas.ErrorResponse(400, "Algunos permisos no existen",
				fmt.Errorf("se esperaban al menos %d permisos, pero se encontraron %d",
					len(roleUpdate.PermissionsID), len(permissions)))
		}

		// Actualizar el nombre del rol
		if err := tx.Model(&existingRole).Update("name", roleUpdate.Name).Error; err != nil {
			return schemas.ErrorResponse(500, "Error interno al actualizar el nombre del rol", err)
		}

		// Reemplazar las asociaciones de permisos
		if err := tx.Model(&existingRole).Association("Permissions").Replace(permissions); err != nil {
			return schemas.ErrorResponse(500, "Error interno al actualizar los permisos del rol", err)
		}

		// Recargar el rol con los permisos actualizados
		tx.Preload("Permissions").First(&newRole, roleUpdate.ID)
		return nil
	})
	
	if err == nil {
		// Guardar auditoría
		go database.SaveAuditAsync(t.DB, models.AuditLog{
			MemberID: memberID,
			Method:   "update",
			Path:     "role",
		}, oldRole, newRole)
	}

	return err
}

// RoleCreate crea un nuevo rol con auditoría
// func (t *RoleRepository) RoleCreate(memberID int64, roleCreate *schemas.RoleCreate) (int64, error) {
// 	var newRoleSave *models.Role
// 	err := t.DB.Transaction(func(tx *gorm.DB) error {
// 		var permissions []models.Permission
// 		if err := tx.Where("id IN ?", roleCreate.PermissionsID).Find(&permissions).Error; err != nil {
// 			return schemas.ErrorResponse(500, "Error interno al buscar permisos", err)
// 		}
// 		if len(permissions) != len(roleCreate.PermissionsID) {
// 			return schemas.ErrorResponse(400, "Algunos permisos no existen", fmt.Errorf("se esperaban %d permisos, pero se encontraron %d", len(roleCreate.PermissionsID), len(permissions)))
// 		}

// 		newRole := &models.Role{Name: roleCreate.Name, Permissions: permissions}

// 		if err := tx.Create(&newRole).Error; err != nil {
// 			return schemas.ErrorResponse(500, "Error interno al crear el rol", err)
// 		}

// 		newRoleSave = newRole
// 		return nil
// 	})

// 	if err != nil {
// 		return 0, err
// 	}

// 	go database.SaveAuditAsync(t.DB, models.AuditLog{
// 		MemberID: memberID,
// 		Method:   "create",
// 		Path:     "role",
// 	}, nil, newRoleSave)

// 	return newRoleSave.ID, nil
// }

// // RoleUpdate actualiza un rol con auditoría
// func (t *RoleRepository) RoleUpdate(memberID int64, roleUpdate *schemas.RoleUpdate) error {
// 	var oldRole, newRole models.Role
// 	err := t.DB.Transaction(func(tx *gorm.DB) error {
// 		// Verificar que el rol existe
// 		var existingRole models.Role
// 		if err := tx.Preload("Permissions").First(&existingRole, roleUpdate.ID).Error; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return schemas.ErrorResponse(404, "Rol no encontrado", err)
// 			}
// 			return schemas.ErrorResponse(500, "Error interno al buscar rol", err)
// 		}

// 		oldRole = existingRole

// 		// Verificar que todos los permisos existen
// 		var permissions []models.Permission
// 		if err := tx.Where("id IN ?", roleUpdate.PermissionsID).Find(&permissions).Error; err != nil {
// 			return schemas.ErrorResponse(500, "Error interno al buscar permisos", err)
// 		}
// 		if len(permissions) != len(roleUpdate.PermissionsID) {
// 			return schemas.ErrorResponse(400, "Algunos permisos no existen",
// 				fmt.Errorf("se esperaban %d permisos, pero se encontraron %d",
// 					len(roleUpdate.PermissionsID), len(permissions)))
// 		}

// 		// Actualizar el nombre del rol
// 		if err := tx.Model(&existingRole).Update("name", roleUpdate.Name).Error; err != nil {
// 			return schemas.ErrorResponse(500, "Error interno al actualizar el nombre del rol", err)
// 		}

// 		// Reemplazar las asociaciones de permisos
// 		if err := tx.Model(&existingRole).Association("Permissions").Replace(permissions); err != nil {
// 			return schemas.ErrorResponse(500, "Error interno al actualizar los permisos del rol", err)
// 		}

// 		// Recargar el rol con los permisos actualizados
// 		tx.Preload("Permissions").First(&newRole, roleUpdate.ID)
// 		return nil
// 	})
	
// 	if err == nil {
// 		// Guardar auditoría
// 		go database.SaveAuditAsync(t.DB, models.AuditLog{
// 			MemberID: memberID,
// 			Method:   "update",
// 			Path:     "role",
// 		}, oldRole, newRole)
// 	}

// 	return err
// }
