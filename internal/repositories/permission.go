package repositories

import (
	"context"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (t *PermissionRepository) PermissionByRoleID(roleID int64) (*[]string, error) {
	ctx := context.Background()

	permissions, err := boilmodels.Permissions(
		qm.Select(boilmodels.PermissionColumns.Code),
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

	permissions, err := boilmodels.Permissions(
		qm.Select(
			boilmodels.PermissionColumns.ID,
			boilmodels.PermissionColumns.Code,
			boilmodels.PermissionColumns.Details,
			boilmodels.PermissionColumns.Group,
			boilmodels.PermissionColumns.Environment,
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
			Details:     p.Details.String,
			Group:       p.Group.String,
			Environment: p.Environment.String,
		})
	}

	return &response, nil
}

func (t *PermissionRepository) PermissionGetToMe(roleID int64) (*[]schemas.PermissionResponse, error) {
	ctx := context.Background()

	permissions, err := boilmodels.Permissions(
		qm.Select(
			"permissions."+boilmodels.PermissionColumns.ID,
			"permissions."+boilmodels.PermissionColumns.Code,
			"permissions."+boilmodels.PermissionColumns.Details,
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
			Details: p.Details.String,
			Group:   p.Group.String,
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
	permissionsAll, err := boilmodels.Permissions().All(ctx, tx)
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

	existingMap := make(map[string]*boilmodels.Permission)
	for _, p := range permissionsAll {
		existingMap[p.Code] = p
	}

	// 1. CREAR o ACTUALIZAR permisos que están en la lista definida
	for _, definedPerm := range database.Permissions {
		if existingPerm, exists := existingMap[definedPerm.Code]; exists {
			// Existe: actualizar si hay cambios
			if existingPerm.Name.String != definedPerm.Name ||
				existingPerm.Details.String != definedPerm.Details ||
				existingPerm.Group.String != definedPerm.Group ||
				existingPerm.Environment.String != definedPerm.Environment {

				existingPerm.Name.String = definedPerm.Name
				existingPerm.Name.Valid = true
				existingPerm.Details.String = definedPerm.Details
				existingPerm.Details.Valid = true
				existingPerm.Group.String = definedPerm.Group
				existingPerm.Group.Valid = true
				existingPerm.Environment.String = definedPerm.Environment
				existingPerm.Environment.Valid = true

				if _, err := existingPerm.Update(ctx, tx, boil.Infer()); err != nil {
					return schemas.HandlerErrorDB(err, "Permiso", schemas.Update)
				}
			}
		} else {
			// No existe: crear
			newPerm := &boilmodels.Permission{
				Code: definedPerm.Code,
			}
			newPerm.Name.String = definedPerm.Name
			newPerm.Name.Valid = true
			newPerm.Details.String = definedPerm.Details
			newPerm.Details.Valid = true
			newPerm.Group.String = definedPerm.Group
			newPerm.Group.Valid = true
			newPerm.Environment.String = definedPerm.Environment
			newPerm.Environment.Valid = true

			if err := newPerm.Insert(ctx, tx, boil.Infer()); err != nil {
				return schemas.HandlerErrorDB(err, "Permiso", schemas.Create)
			}
		}
	}

	// 2. ELIMINAR permisos que NO están en la lista definida
	for _, existingPerm := range permissionsAll {
		if _, stillDefined := definedMap[existingPerm.Code]; !stillDefined {
			// No está en la lista definida: eliminar
			if _, err := existingPerm.Delete(ctx, tx, false); err != nil {
				return schemas.HandlerErrorDB(err, "Permiso", schemas.Delete)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
