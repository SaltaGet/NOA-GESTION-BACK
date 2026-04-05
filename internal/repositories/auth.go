package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/master"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/utils"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/rs/zerolog/log"
)

func (r *MainRepository) AuthTenantGetByID(tenantID int64) (*master.Tenant, error) {
	ctx := context.Background()
	boilt, err := master.Tenants(
		master.TenantWhere.ID.EQ(tenantID),
		qm.Load(master.TenantRels.Credential,
			qm.Select(master.CredentialColumns.ResponsibilityFrontIva),
		),
	).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al obtener los tenants", err)
	}

	if !boilt.IsActive {
		return nil, schemas.ErrorResponse(403, "Tenant esta inactivo", fmt.Errorf("credenciales incorrectas"))
	}

	if boilt.R.Credential == nil {
		log.Warn().Msgf("Tenant %d no tiene credenciales configuradas, usando valor por defecto", tenantID)
		boilt.R.Credential = &master.Credential{
			ResponsibilityFrontIva: null.String{},
		}
	}
	
	return boilt, nil
}

func (r *MainRepository) AuthTenantGetByIdentifier(identifier string) (*master.Tenant, error) {
	ctx := context.Background()
	boilt, err := master.Tenants(master.TenantWhere.Identifier.EQ(identifier)).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al obtener los tenants", err)
	}

	if !boilt.IsActive {
		return nil, schemas.ErrorResponse(403, "Tenant esta inactivo", fmt.Errorf("credenciales incorrectas"))
	}

	return boilt, nil
}

func (r *MainRepository) AuthMemberGetByUserID(userID int64, connection string, tenantID int64) (*tenant.Member, error) {
	ctx := context.Background()
	db, err := database.GetTenantDB(connection, tenantID)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al recibir la conexión", err)
	}

	boilMember, err := tenant.Members(tenant.MemberWhere.ID.EQ(userID)).One(ctx, db)
	if err != nil {
		return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
	}

	if !boilMember.IsActive {
		return nil, schemas.ErrorResponse(403, "Miembro inactivo", fmt.Errorf("miembro inactivo"))
	}

	return boilMember, nil
}

func (r *MainRepository) AuthMemberGetByID(id int64, connection string, tenantID int64) (*tenant.Member, *tenant.PermissionSlice, error) {
	ctx := context.Background()
	db, err := database.GetTenantDB(connection, tenantID)
	if err != nil {
		return nil, nil, schemas.ErrorResponse(500, "Error al recibir la conexión", err)
	}

	boilMember, err := tenant.Members(
		tenant.MemberWhere.ID.EQ(id),
		qm.Load(tenant.MemberRels.Role),
	).One(ctx, db)

	if err != nil {
		return nil, nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
	}

	if !boilMember.IsActive {
		return nil, nil, schemas.ErrorResponse(403, "Miembro inactivo", fmt.Errorf("miembro inactivo"))
	}

	var boilPerms tenant.PermissionSlice

	if boilMember.R != nil && boilMember.R.Role != nil {
		if boilMember.R.Role.Name == "admin" {
			boilPerms, err = tenant.Permissions().All(ctx, db)
		} else {
			boilPerms, err = tenant.Permissions(
				qm.InnerJoin("role_permissions on permissions.id = role_permissions.permission_id"),
				qm.Where("role_permissions.role_id = ?", boilMember.R.Role.ID),
			).All(ctx, db)
		}

		if err != nil {
			return nil, nil, schemas.ErrorResponse(500, "Error al obtener los permisos", err)
		}

		return boilMember, &boilPerms, nil
	}

	emptyPerms := make(tenant.PermissionSlice, 0)
	return boilMember, &emptyPerms, nil
}

func (r *MainRepository) AuthMemberGetByUsername(username string, connection string, tenantID int64) (*tenant.Member, error) {
	ctx := context.Background()
	db, err := database.GetTenantDB(connection, tenantID)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al recibir la conexión", err)
	}

	boilMember, err := tenant.Members(tenant.MemberWhere.Username.EQ(username)).One(ctx, db)
	if err != nil {
		return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
	}

	if !boilMember.IsActive {
		return nil, schemas.ErrorResponse(403, "Miembro inactivo", fmt.Errorf("miembro inactivo"))
	}

	return boilMember, nil
}

func (r *MainRepository) AuthAdminGetByUsername(username string) (*master.Admin, error) {
	ctx := context.Background()
	boilAdmin, err := master.Admins(master.AdminWhere.Username.EQ(username)).One(ctx, r.DB)
	if err != nil {
		return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
	}

	return boilAdmin, nil
}

func (r *MainRepository) AuthAdminGetByID(id int64) (*master.Admin, error) {
	ctx := context.Background()
	boilAdmin, err := master.Admins(master.AdminWhere.ID.EQ(id)).One(ctx, r.DB)
	if err != nil {
		return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
	}

	return boilAdmin, nil
}

func (r *MainRepository) AuthPointSale(pointSaleID int64, connection string, tenantID, memberID int64) (*tenant.PointSale, error) {
	ctx := context.Background()
	db, err := database.GetTenantDB(connection, tenantID)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al recibir la conexión", err)
	}

	ps, err := tenant.PointSales(
		qm.InnerJoin("member_point_sales mp ON mp.point_sale_id = point_sales.id"),
		qm.Where("mp.member_id = ?", memberID),
		tenant.PointSaleWhere.ID.EQ(pointSaleID),
	).One(ctx, db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return nil, schemas.ErrorResponse(500, "Error al obtener tenant", err)
	}

	// pointSale := &tenant.PointSale{
	// 	ID:        ps.ID,
	// 	Name:      ps.Name,
	// 	IsDeposit: ps.IsDeposit,
	// }
	// if ps.Address.Valid {
	// 	pointSale.Address = &ps.Address.String
	// }
	// if ps.CreatedAt.Valid {
	// 	pointSale.CreatedAt = ps.CreatedAt.Time
	// }
	// if ps.UpdatedAt.Valid {
	// 	pointSale.UpdatedAt = ps.UpdatedAt.Time
	// }

	return ps, nil
}

func (r *MainRepository) AuthForgotPassword(forgotPassword *schemas.AuthForgotPassword) (*tenant.Member, *master.Tenant, error) {
	ctx := context.Background()
	t, err := master.Tenants(
		qm.Select(master.TenantColumns.ID, master.TenantColumns.Connection),
		master.TenantWhere.Identifier.EQ(forgotPassword.TenantIdentifier),
	).One(ctx, r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return nil, nil, schemas.ErrorResponse(500, "Error interno al obtener los tenants", err)
	}

	db, err := database.GetTenantDB(t.Connection, t.ID)
	if err != nil {
		return nil, nil, schemas.ErrorResponse(500, "Error al recibir la conexión", err)
	}

	boilMember, err := tenant.Members(
		qm.Select(tenant.MemberColumns.ID, tenant.MemberColumns.Username, tenant.MemberColumns.Email, tenant.MemberColumns.IsActive),
		tenant.MemberWhere.Username.EQ(forgotPassword.Username),
	).One(ctx, db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return nil, nil, schemas.ErrorResponse(500, "Error al obtener el miembro", err)
	}

	if !boilMember.IsActive {
		return nil, nil, schemas.ErrorResponse(403, "Miembro inactivo", fmt.Errorf("miembro inactivo"))
	}

	return boilMember, t, nil
}

func (r *MainRepository) AuthResetPassword(memberID, tenantID int64, newPass string) error {
	ctx := context.Background()
	t, err := master.Tenants(
		qm.Select(master.TenantColumns.ID, master.TenantColumns.Connection),
		master.TenantWhere.ID.EQ(tenantID),
	).One(ctx, r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schemas.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return schemas.ErrorResponse(500, "Error interno al obtener los tenants", err)
	}

	db, err := database.GetTenantDB(t.Connection, t.ID)
	if err != nil {
		return schemas.ErrorResponse(500, "Error al recibir la conexión", err)
	}

	boilMember, err := tenant.FindMember(ctx, db, memberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schemas.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return schemas.ErrorResponse(500, "Error al obtener el miembro", err)
	}

	if !boilMember.IsActive {
		return schemas.ErrorResponse(403, "Miembro inactivo", fmt.Errorf("miembro inactivo"))
	}

	passHash, err := utils.HashPassword(newPass)
	if err != nil {
		return schemas.ErrorResponse(500, "Error al generar la contraseña", err)
	}

	boilMember.Password = passHash
	if _, err := boilMember.Update(ctx, db, boil.Whitelist(tenant.MemberColumns.Password, tenant.MemberColumns.UpdatedAt)); err != nil {
		return schemas.ErrorResponse(500, "Error al actualizar la contraseña", err)
	}

	return nil
}
