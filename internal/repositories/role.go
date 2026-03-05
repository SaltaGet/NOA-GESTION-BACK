package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func mapToPermissionResponse(p *boilmodels.Permission) schemas.PermissionResponse {
	if p == nil {
		return schemas.PermissionResponse{}
	}
	return schemas.PermissionResponse{
		ID:          p.ID,
		Code:        p.Code,
		Group:       p.Group.String,
		Environment: p.Environment.String,
		Details:     p.Details.String,
	}
}

func mapToRoleResponse(r *boilmodels.Role) schemas.RoleResponse {
	if r == nil {
		return schemas.RoleResponse{}
	}

	res := schemas.RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		Permissions: []schemas.PermissionResponse{},
	}

	if r.R != nil && len(r.R.Permissions) > 0 {
		for _, p := range r.R.Permissions {
			res.Permissions = append(res.Permissions, mapToPermissionResponse(p))
		}
	}

	return res
}

func (r *RoleRepository) RoleGetByID(id int64) (*schemas.RoleResponse, error) {
	ctx := context.Background()

	role, err := boilmodels.Roles(
		qm.Load(boilmodels.RoleRels.Permissions),
		boilmodels.RoleWhere.ID.EQ(id),
	).One(ctx, r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Rol", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Rol", schemas.Read)
	}

	res := mapToRoleResponse(role)
	return &res, nil
}

func (r *RoleRepository) RoleGetAll() (*[]schemas.RoleResponse, error) {
	ctx := context.Background()

	roles, err := boilmodels.Roles(
		qm.Load(boilmodels.RoleRels.Permissions),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Rol", schemas.Read)
	}

	var allRoles []schemas.RoleResponse
	for _, role := range roles {
		allRoles = append(allRoles, mapToRoleResponse(role))
	}

	return &allRoles, nil
}

func expandPermissions(ctx context.Context, exec boil.ContextExecutor, permissionIDs []int64) ([]*boilmodels.Permission, error) {
	if len(permissionIDs) == 0 {
		return []*boilmodels.Permission{}, nil
	}

	// Make ids list for querying
	var interfaces []interface{}
	for _, id := range permissionIDs {
		interfaces = append(interfaces, id)
	}

	requestedPermissions, err := boilmodels.Permissions(
		qm.WhereIn("id IN ?", interfaces...),
	).All(ctx, exec)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Rol", schemas.Read)
	}

	permissionMap := make(map[int64]*boilmodels.Permission)
	for _, perm := range requestedPermissions {
		permissionMap[perm.ID] = perm
	}

	var groupsToExpand []interface{}
	for _, perm := range requestedPermissions {
		if len(perm.Code) >= 2 && perm.Code[len(perm.Code)-2:] == "02" {
			groupsToExpand = append(groupsToExpand, perm.Group.String)
		}
	}

	if len(groupsToExpand) > 0 {
		readPermissions, err := boilmodels.Permissions(
			qm.Where("code LIKE ?", "%04"),
			qm.WhereIn("\"group\" IN ?", groupsToExpand...), // Note: group is a reserved pgx word
		).All(ctx, exec)

		if err != nil {
			return nil, schemas.HandlerErrorDB(err, "Rol", schemas.Read)
		}

		for _, perm := range readPermissions {
			permissionMap[perm.ID] = perm
		}
	}

	var expandedPermissions []*boilmodels.Permission
	for _, perm := range permissionMap {
		expandedPermissions = append(expandedPermissions, perm)
	}

	return expandedPermissions, nil
}

func (t *RoleRepository) RoleCreate(memberID int64, roleCreate *schemas.RoleCreate) (int64, error) {
	ctx := context.Background()
	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return 0, err
	}

	permissions, err := expandPermissions(ctx, tx, roleCreate.PermissionsID)
	if err != nil {
		return 0, err
	}

	// Check if all requested permissions actually exist before expanding
	if len(permissions) < len(roleCreate.PermissionsID) {
		return 0, schemas.ErrorResponse(400, "Algunos permisos no existen",
			fmt.Errorf("se esperaban al menos %d permisos, pero se encontraron %d",
				len(roleCreate.PermissionsID), len(permissions)))
	}

	newRole := &boilmodels.Role{
		Name: roleCreate.Name,
	}

	if err := newRole.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Rol", schemas.Create)
	}

	if len(permissions) > 0 {
		if err := newRole.SetPermissions(ctx, tx, false, permissions...); err != nil {
			return 0, schemas.HandlerErrorDB(err, "Rol", schemas.Update)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newRole.ID, nil
}

func (t *RoleRepository) RoleUpdate(memberID int64, roleUpdate *schemas.RoleUpdate) error {
	ctx := context.Background()
	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	existingRole, err := boilmodels.Roles(
		boilmodels.RoleWhere.ID.EQ(roleUpdate.ID),
	).One(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Rol", schemas.Read)
	}

	permissions, err := expandPermissions(ctx, tx, roleUpdate.PermissionsID)
	if err != nil {
		return err
	}

	if len(permissions) < len(roleUpdate.PermissionsID) {
		return schemas.ErrorResponse(400, "Algunos permisos no existen",
			fmt.Errorf("se esperaban al menos %d permisos, pero se encontraron %d",
				len(roleUpdate.PermissionsID), len(permissions)))
	}

	existingRole.Name = roleUpdate.Name
	if _, err := existingRole.Update(ctx, tx, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Rol", schemas.Update)
	}

	// Reemplazar asociaciones
	if err := existingRole.SetPermissions(ctx, tx, false, permissions...); err != nil {
		return schemas.HandlerErrorDB(err, "Rol", schemas.Update)
	}

	return tx.Commit()
}
