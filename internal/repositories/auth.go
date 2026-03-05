package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/utils"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func mapToModelTenant(c *boilmodels.Tenant) *models.Tenant {
	if c == nil {
		return nil
	}
	t := &models.Tenant{
		ID:            c.ID,
		Name:          c.Name,
		Identifier:    c.Identifier,
		Address:       c.Address,
		Phone:         c.Phone,
		Email:         c.Email,
		CuitPdv:       c.CuitPdv,
		IsActive:      c.IsActive,
		PlanID:        c.PlanID,
		Connection:    c.Connection,
		AcceptedTerms: c.AcceptedTerms,
	}
	if c.Expiration.Valid {
		t.Expiration = &c.Expiration.Time
	}
	if c.Ip.Valid {
		t.IP = &c.Ip.String
	}
	if c.DateAccepted.Valid {
		t.DateAccepted = &c.DateAccepted.Time
	}
	if c.CreatedAt.Valid {
		t.CreatedAt = c.CreatedAt.Time
	}
	if c.UpdatedAt.Valid {
		t.UpdatedAt = c.UpdatedAt.Time
	}
	return t
}

func mapToModelAdmin(c *boilmodels.Admin) *models.Admin {
	if c == nil {
		return nil
	}
	a := &models.Admin{
		ID:           c.ID,
		FirstName:    c.FirstName,
		LastName:     c.LastName,
		Username:     c.Username,
		Email:        c.Email,
		Password:     c.Password,
		IsSuperAdmin: c.IsSuperAdmin,
		IsActive:     c.IsActive,
	}
	if c.CreatedAt.Valid {
		a.CreatedAt = c.CreatedAt.Time
	}
	if c.UpdatedAt.Valid {
		a.UpdatedAt = c.UpdatedAt.Time
	}
	return a
}

func mapToModelMember(c *boilmodels.Member) *models.Member {
	if c == nil {
		return nil
	}
	m := &models.Member{
		ID:        c.ID,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Username:  c.Username,
		Email:     c.Email,
		Password:  c.Password,
		IsAdmin:   c.IsAdmin,
		IsActive:  c.IsActive,
		RoleID:    c.RoleID,
	}
	if c.Address.Valid {
		m.Address = &c.Address.String
	}
	if c.Phone.Valid {
		m.Phone = &c.Phone.String
	}
	if c.CreatedAt.Valid {
		m.CreatedAt = c.CreatedAt.Time
	}
	if c.UpdatedAt.Valid {
		m.UpdatedAt = c.UpdatedAt.Time
	}
	if c.R != nil {
		if c.R.Role != nil {
			m.Role = models.Role{
				ID:   c.R.Role.ID,
				Name: c.R.Role.Name,
			}
		}
	}
	return m
}

func (r *MainRepository) AuthTenantGetByID(tenantID int64) (*models.Tenant, error) {
	ctx := context.Background()
	boilt, err := boilmodels.Tenants(boilmodels.TenantWhere.ID.EQ(tenantID)).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al obtener los tenants", err)
	}

	if !boilt.IsActive {
		return nil, schemas.ErrorResponse(403, "Tenant esta inactivo", fmt.Errorf("credenciales incorrectas"))
	}

	tenant := mapToModelTenant(boilt)

	c, err := boilmodels.Credentials(
		qm.Select(boilmodels.CredentialColumns.ResponsibilityFrontIva),
		boilmodels.CredentialWhere.TenantID.EQ(tenantID),
	).One(ctx, r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msgf("Tenant %d no tiene credenciales configuradas, usando valor por defecto", tenantID)
			tenant.Credentials.ResponsibilityFrontIVA = nil
		} else {
			return nil, schemas.ErrorResponse(500, "Error interno al obtener las credenciales", err)
		}
	} else {
		if c.ResponsibilityFrontIva.Valid {
			tenant.Credentials.ResponsibilityFrontIVA = &c.ResponsibilityFrontIva.String
		}
	}

	return tenant, nil
}

func (r *MainRepository) AuthTenantGetByIdentifier(identifier string) (*models.Tenant, error) {
	ctx := context.Background()
	boilt, err := boilmodels.Tenants(boilmodels.TenantWhere.Identifier.EQ(identifier)).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al obtener los tenants", err)
	}

	if !boilt.IsActive {
		return nil, schemas.ErrorResponse(403, "Tenant esta inactivo", fmt.Errorf("credenciales incorrectas"))
	}

	return mapToModelTenant(boilt), nil
}

func (r *MainRepository) AuthMemberGetByUserID(userID int64, connection string, tenantID int64) (*models.Member, error) {
	ctx := context.Background()
	db, err := database.GetTenantDB(connection, tenantID)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al recibir la conexión", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al recibir la sql db", err)
	}

	boilMember, err := boilmodels.Members(boilmodels.MemberWhere.ID.EQ(userID)).One(ctx, sqlDB)
	if err != nil {
		return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
	}

	if !boilMember.IsActive {
		return nil, schemas.ErrorResponse(403, "Miembro inactivo", fmt.Errorf("miembro inactivo"))
	}

	return mapToModelMember(boilMember), nil
}

func (r *MainRepository) AuthMemberGetByID(id int64, connection string, tenantID int64) (*models.Member, *[]string, error) {
	ctx := context.Background()
	db, err := database.GetTenantDB(connection, tenantID)
	if err != nil {
		return nil, nil, schemas.ErrorResponse(500, "Error al recibir la conexión", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, schemas.ErrorResponse(500, "Error al recibir la sql db", err)
	}

	boilMember, err := boilmodels.Members(
		boilmodels.MemberWhere.ID.EQ(id),
		qm.Load(boilmodels.MemberRels.Role),
	).One(ctx, sqlDB)

	if err != nil {
		return nil, nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
	}

	if !boilMember.IsActive {
		return nil, nil, schemas.ErrorResponse(403, "Miembro inactivo", fmt.Errorf("miembro inactivo"))
	}

	modelsMember := mapToModelMember(boilMember)
	var boilPerms boilmodels.PermissionSlice

	if boilMember.R != nil && boilMember.R.Role != nil {
		if boilMember.R.Role.Name == "admin" {
			boilPerms, err = boilmodels.Permissions().All(ctx, sqlDB)
		} else {
			boilPerms, err = boilmodels.Permissions(
				qm.InnerJoin("role_permissions on permissions.id = role_permissions.permission_id"),
				qm.Where("role_permissions.role_id = ?", boilMember.R.Role.ID),
			).All(ctx, sqlDB)
		}

		if err != nil {
			return nil, nil, schemas.ErrorResponse(500, "Error al obtener los permisos", err)
		}

		perm := make([]string, len(boilPerms))
		for i, p := range boilPerms {
			modelsMember.Role.Permissions = append(modelsMember.Role.Permissions, models.Permission{
				ID:   p.ID,
				Name: p.Name,
				Code: p.Code,
			})
			perm[i] = p.Code
		}
		return modelsMember, &perm, nil
	}

	emptyPerms := make([]string, 0)
	return modelsMember, &emptyPerms, nil
}

func (r *MainRepository) AuthMemberGetByUsername(username string, connection string, tenantID int64) (*models.Member, error) {
	ctx := context.Background()
	db, err := database.GetTenantDB(connection, tenantID)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al recibir la conexión", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al recibir la sql db", err)
	}

	boilMember, err := boilmodels.Members(boilmodels.MemberWhere.Username.EQ(username)).One(ctx, sqlDB)
	if err != nil {
		return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
	}

	if !boilMember.IsActive {
		return nil, schemas.ErrorResponse(403, "Miembro inactivo", fmt.Errorf("miembro inactivo"))
	}

	return mapToModelMember(boilMember), nil
}

func (r *MainRepository) AuthAdminGetByUsername(username string) (*models.Admin, error) {
	ctx := context.Background()
	boilAdmin, err := boilmodels.Admins(boilmodels.AdminWhere.Username.EQ(username)).One(ctx, r.DB)
	if err != nil {
		return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
	}

	return mapToModelAdmin(boilAdmin), nil
}

func (r *MainRepository) AuthAdminGetByID(id int64) (*models.Admin, error) {
	ctx := context.Background()
	boilAdmin, err := boilmodels.Admins(boilmodels.AdminWhere.ID.EQ(id)).One(ctx, r.DB)
	if err != nil {
		return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
	}

	return mapToModelAdmin(boilAdmin), nil
}

func (r *MainRepository) AuthPointSale(pointSaleID int64, connection string, tenantID, memberID int64) (*models.PointSale, error) {
	ctx := context.Background()
	db, err := database.GetTenantDB(connection, tenantID)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al recibir la conexión", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al recibir la sql db", err)
	}

	ps, err := boilmodels.PointSales(
		qm.InnerJoin("member_point_sales mp ON mp.point_sale_id = point_sales.id"),
		qm.Where("mp.member_id = ?", memberID),
		boilmodels.PointSaleWhere.ID.EQ(pointSaleID),
	).One(ctx, sqlDB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return nil, schemas.ErrorResponse(500, "Error al obtener tenant", err)
	}

	pointSale := &models.PointSale{
		ID:        ps.ID,
		Name:      ps.Name,
		IsDeposit: ps.IsDeposit,
	}
	if ps.Address.Valid {
		pointSale.Address = &ps.Address.String
	}
	if ps.CreatedAt.Valid {
		pointSale.CreatedAt = ps.CreatedAt.Time
	}
	if ps.UpdatedAt.Valid {
		pointSale.UpdatedAt = ps.UpdatedAt.Time
	}

	return pointSale, nil
}

func (r *MainRepository) AuthForgotPassword(forgotPassword *schemas.AuthForgotPassword) (*models.Member, *models.Tenant, error) {
	ctx := context.Background()
	t, err := boilmodels.Tenants(
		qm.Select(boilmodels.TenantColumns.ID, boilmodels.TenantColumns.Connection),
		boilmodels.TenantWhere.Identifier.EQ(forgotPassword.TenantIdentifier),
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
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, schemas.ErrorResponse(500, "Error al recibir la sql db", err)
	}

	boilMember, err := boilmodels.Members(
		qm.Select(boilmodels.MemberColumns.ID, boilmodels.MemberColumns.Username, boilmodels.MemberColumns.Email, boilmodels.MemberColumns.IsActive),
		boilmodels.MemberWhere.Username.EQ(forgotPassword.Username),
	).One(ctx, sqlDB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return nil, nil, schemas.ErrorResponse(500, "Error al obtener el miembro", err)
	}

	if !boilMember.IsActive {
		return nil, nil, schemas.ErrorResponse(403, "Miembro inactivo", fmt.Errorf("miembro inactivo"))
	}

	return mapToModelMember(boilMember), mapToModelTenant(t), nil
}

func (r *MainRepository) AuthResetPassword(memberID, tenantID int64, newPass string) error {
	ctx := context.Background()
	t, err := boilmodels.Tenants(
		qm.Select(boilmodels.TenantColumns.ID, boilmodels.TenantColumns.Connection),
		boilmodels.TenantWhere.ID.EQ(tenantID),
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
	sqlDB, err := db.DB()
	if err != nil {
		return schemas.ErrorResponse(500, "Error al recibir la sql db", err)
	}

	boilMember, err := boilmodels.FindMember(ctx, sqlDB, memberID)
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
	if _, err := boilMember.Update(ctx, sqlDB, boil.Whitelist(boilmodels.MemberColumns.Password, boilmodels.MemberColumns.UpdatedAt)); err != nil {
		return schemas.ErrorResponse(500, "Error al actualizar la contraseña", err)
	}

	return nil
}
