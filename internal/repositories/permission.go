package repositories

import (
	"context"

	boiltenant "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

func (t *PermissionRepository) PermissionByRoleID(roleID int64) (*[]string, error) {
	ctx := context.Background()

	permissions, err := boiltenant.Permissions(
		qm.Select(boiltenant.PermissionColumns.Code),
		qm.InnerJoin("role_permissions rp ON rp.permission_id = permissions.id"),
		qm.Where("rp.role_id = ?", roleID),
	).All(ctx, t.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Permiso", schemas.Read)
	}

	var permCodes []string
	for _, p := range permissions {
		permCodes = append(permCodes, p.Code)
	}

	return &permCodes, nil
}

func (t *PermissionRepository) PermissionGetAll() (*[]schemas.PermissionResponse, error) {
	ctx := context.Background()

	permissions, err := boiltenant.Permissions(
		qm.Select(
			boiltenant.PermissionColumns.ID,
			boiltenant.PermissionColumns.Code,
			boiltenant.PermissionColumns.Details,
			boiltenant.PermissionColumns.Group,
			boiltenant.PermissionColumns.Environment,
		),
	).All(ctx, t.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Permiso", schemas.Read)
	}

	var response []schemas.PermissionResponse
	for _, p := range permissions {
		response = append(response, schemas.PermissionResponse{
			ID:          p.ID,
			Code:        p.Code,
			Details:     p.Details,
			Group:       p.Group,
			Environment: p.Environment,
		})
	}

	return &response, nil
}

func (t *PermissionRepository) PermissionGetToMe(roleID int64) (*[]schemas.PermissionResponse, error) {
	ctx := context.Background()

	permissions, err := boiltenant.Permissions(
		qm.Select(
			"permissions."+boiltenant.PermissionColumns.ID,
			"permissions."+boiltenant.PermissionColumns.Code,
			"permissions."+boiltenant.PermissionColumns.Details,
			"permissions."+"\"group\"", // Escaped group because it's a reserved word in PG
		),
		qm.InnerJoin("role_permissions rp ON rp.permission_id = permissions.id"),
		qm.Where("rp.role_id = ?", roleID),
	).All(ctx, t.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Permiso", schemas.Read)
	}

	var response []schemas.PermissionResponse
	for _, p := range permissions {
		response = append(response, schemas.PermissionResponse{
			ID:      p.ID,
			Code:    p.Code,
			Details: p.Details,
			Group:   p.Group,
		})
	}

	return &response, nil
}

func (t *PermissionRepository) PermissionUpdateAll() error {
	ctx := context.Background()
	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Obtener todos los permisos existentes en la BD
	permissionsAll, err := boiltenant.Permissions().All(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Permiso", schemas.Read)
	}

	// Crear mapas para búsqueda eficiente
	type PermissionWrapper struct {
		Name        string
		Code        string
		Details     string
		Group       string
		Environment string
	}

	definedMap := make(map[string]PermissionWrapper)
	// Como database.Permissions es un slice de models.Permission (GORM),
	// necesitamos mapearlo para la comparación con SQLBoiler
	for _, p := range database.Permissions {
		definedMap[p.Code] = PermissionWrapper{
			Name:        p.Name,
			Code:        p.Code,
			Details:     p.Details,
			Group:       p.Group,
			Environment: p.Environment,
		}
	}

	existingMap := make(map[string]*boiltenant.Permission)
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

				existingPerm.Name = definedPerm.Name
				existingPerm.Details = definedPerm.Details
				existingPerm.Group = definedPerm.Group
				existingPerm.Environment = definedPerm.Environment

				if _, err := existingPerm.Update(ctx, tx, boil.Infer()); err != nil {
					return schemas.HandlerErrorDB(err, "Permiso", schemas.Update)
				}
			}
		} else {
			// No existe: crear
			newPerm := &boiltenant.Permission{
				Code: definedPerm.Code,
			}
			newPerm.Name = definedPerm.Name
			newPerm.Details = definedPerm.Details
			newPerm.Group = definedPerm.Group
			newPerm.Environment = definedPerm.Environment

			if err := newPerm.Insert(ctx, tx, boil.Infer()); err != nil {
				return schemas.HandlerErrorDB(err, "Permiso", schemas.Create)
			}
		}
	}

	// 2. ELIMINAR permisos que NO están en la lista definida
	for _, existingPerm := range permissionsAll {
		if _, stillDefined := definedMap[existingPerm.Code]; !stillDefined {
			// No está en la lista definida: eliminar
			if _, err := existingPerm.Delete(ctx, tx); err != nil {
				return schemas.HandlerErrorDB(err, "Permiso", schemas.Delete)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
