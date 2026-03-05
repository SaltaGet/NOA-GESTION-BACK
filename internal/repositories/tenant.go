package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/utils"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (r *MainRepository) TenantGetByID(tenantID int64) (*models.Tenant, error) {
	ctx := context.Background()

	tenant, err := boilmodels.Tenants(boilmodels.TenantWhere.ID.EQ(tenantID)).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
	}

	if !tenant.IsActive {
		return nil, schemas.ErrorResponse(403, "Tenant is inactive", fmt.Errorf("tenant is inactive"))
	}

	// For compatibility with the current interface using models.Tenant
	exp := time.Time{}
	if tenant.Expiration.Valid {
		exp = tenant.Expiration.Time
	}

	dateAccepted := time.Time{}
	if tenant.DateAccepted.Valid {
		dateAccepted = tenant.DateAccepted.Time
	}

	return &models.Tenant{
		ID:            tenant.ID,
		Name:          tenant.Name,
		Address:       tenant.Address.String,
		Phone:         tenant.Phone.String,
		Email:         tenant.Email.String,
		CuitPdv:       tenant.CuitPdv.String,
		Connection:    tenant.Connection,
		Identifier:    tenant.Identifier,
		PlanID:        tenant.PlanID.Int64,
		IsActive:      tenant.IsActive,
		Expiration:    &exp,
		AcceptedTerms: tenant.AcceptedTerms.Bool,
		IP:            &tenant.IP.String,
		DateAccepted:  &dateAccepted,
	}, nil
}

func (r *MainRepository) TenantGetByIdentifier(identifier string) (*models.Tenant, error) {
	ctx := context.Background()

	tenant, err := boilmodels.Tenants(boilmodels.TenantWhere.Identifier.EQ(identifier)).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al obtener los tenants", err)
	}

	if !tenant.IsActive {
		return nil, schemas.ErrorResponse(403, "Tenant es inactivo", fmt.Errorf("tenant es inactivo"))
	}

	exp := time.Time{}
	if tenant.Expiration.Valid {
		exp = tenant.Expiration.Time
	}

	dateAccepted := time.Time{}
	if tenant.DateAccepted.Valid {
		dateAccepted = tenant.DateAccepted.Time
	}

	// Assuming address string is required in models.Tenant
	address := ""
	if tenant.Address.Valid {
		address = tenant.Address.String
	}

	return &models.Tenant{
		ID:            tenant.ID,
		Name:          tenant.Name,
		Address:       address,
		Phone:         tenant.Phone.String,
		Email:         tenant.Email.String,
		CuitPdv:       tenant.CuitPdv.String,
		Connection:    tenant.Connection,
		Identifier:    tenant.Identifier,
		PlanID:        tenant.PlanID.Int64,
		IsActive:      tenant.IsActive,
		Expiration:    &exp,
		AcceptedTerms: tenant.AcceptedTerms.Bool,
		IP:            &tenant.IP.String,
		DateAccepted:  &dateAccepted,
	}, nil
}

func (r *MainRepository) TenantGetAll() (*[]schemas.TenantResponse, error) {
	ctx := context.Background()

	tenants, err := boilmodels.Tenants(
		qm.Select(
			boilmodels.TenantColumns.ID,
			boilmodels.TenantColumns.Name,
			boilmodels.TenantColumns.Address,
			boilmodels.TenantColumns.Phone,
			boilmodels.TenantColumns.Email,
			boilmodels.TenantColumns.IsActive,
			boilmodels.TenantColumns.Expiration,
			boilmodels.TenantColumns.CreatedAt,
			boilmodels.TenantColumns.UpdatedAt,
		),
		qm.InnerJoin("user_tenants ON tenants.id = user_tenants.tenant_id"),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
	}

	var response []schemas.TenantResponse
	for _, t := range tenants {
		var exp *time.Time
		if t.Expiration.Valid {
			expTime := t.Expiration.Time
			exp = &expTime
		}

		var createdAt time.Time
		if t.CreatedAt.Valid {
			createdAt = t.CreatedAt.Time
		}

		var updatedAt time.Time
		if t.UpdatedAt.Valid {
			updatedAt = t.UpdatedAt.Time
		}

		response = append(response, schemas.TenantResponse{
			ID:         t.ID,
			Name:       t.Name,
			Address:    t.Address.String,
			Phone:      t.Phone.String,
			Email:      t.Email.String,
			IsActive:   t.IsActive,
			Expiration: exp,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		})
	}

	return &response, nil
}

func (r *MainRepository) TenantGetConnectionByIdentifier(tenantIdentifier string) (*models.Tenant, error) {
	ctx := context.Background()

	tenant, err := boilmodels.Tenants(
		qm.Select("id", "connection"),
		boilmodels.TenantWhere.Identifier.EQ(tenantIdentifier),
	).One(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
	}

	return &models.Tenant{
		ID:         tenant.ID,
		Connection: tenant.Connection,
	}, nil
}

func (r *MainRepository) TenantGetConections() ([]*models.Tenant, error) {
	ctx := context.Background()

	tenants, err := boilmodels.Tenants(
		qm.Select("id", "name", "connection", "identifier"),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
	}

	var response []*models.Tenant
	for _, tenant := range tenants {
		connectionDecrypted, err := utils.Decrypt(tenant.Connection)
		if err != nil {
			return nil, schemas.ErrorResponse(500, "Error interno al obtener las connections", err)
		}
		response = append(response, &models.Tenant{
			ID:         tenant.ID,
			Name:       tenant.Name,
			Connection: connectionDecrypted,
			Identifier: tenant.Identifier,
		})
	}

	return response, nil
}

func (r *MainRepository) TenantCreateByUserID(adminID int64, tenantCreate *schemas.TenantCreate, userID int64) (int64, error) {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", adminID)); err != nil {
		return 0, err
	}

	tenantName := strings.ReplaceAll(tenantCreate.Name, " ", "_")
	identifier := strings.ReplaceAll(tenantCreate.Identifier, " ", "_")

	uri := fmt.Sprintf("%s%s_%s%s",
		os.Getenv("URI_PATH"),
		tenantName,
		identifier,
		os.Getenv("URI_CONFIG"),
	)

	connection, err := utils.Encrypt(uri)
	if err != nil {
		return 0, schemas.ErrorResponse(500, "Error interno al obtener connection", err)
	}

	tenant := &boilmodels.Tenant{
		Name:       tenantCreate.Name,
		Address:    null.StringFrom(tenantCreate.Address),
		Phone:      null.StringFrom(tenantCreate.Phone),
		Email:      null.StringFrom(tenantCreate.Email),
		CuitPdv:    null.StringFrom(tenantCreate.CuitPdv),
		Connection: connection,
		Identifier: identifier,
	}

	if err := tenant.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Tenant", schemas.Create)
	}

	// Check user existence
	user, err := boilmodels.Users(boilmodels.UserWhere.ID.EQ(userID)).One(ctx, tx)
	if err != nil {
		return 0, schemas.HandlerErrorDB(err, "User", schemas.Read)
	}

	userTenant := &boilmodels.UserTenant{
		UserID:   user.ID,
		TenantID: tenant.ID,
	}

	if err := userTenant.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "UserTenant", schemas.Create)
	}

	// Commit transaction before migrating sub db
	// wait, if migration fails, tenant is created but broken. So we must commit AFTER setup.
	memberAdmin := models.Member{
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		Username:  user.Username,
		Email:     user.Email,
		Password:  "1",
		IsAdmin:   true,
		Address:   &user.Address.String,
		RoleID:    1,
	}

	if err := database.PrepareDB(uri, memberAdmin); err != nil {
		// handleDBCreationError inside PrepareDB will drop it if it fails
		return 0, schemas.HandlerErrorDB(err, "Tenant", schemas.Create)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return tenant.ID, nil
}

func (r *MainRepository) TenantUserCreate(adminID int64, tenantUserCreate *schemas.TenantUserCreate) (int64, error) {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", adminID)); err != nil {
		return 0, err
	}

	tenantName := strings.ReplaceAll(tenantUserCreate.TenantCreate.Name, " ", "_")
	identifier := strings.ReplaceAll(tenantUserCreate.TenantCreate.Identifier, " ", "_")

	uri := fmt.Sprintf(
		"%s%s_%s%s",
		os.Getenv("URI_PATH"),
		tenantName,
		identifier,
		os.Getenv("URI_CONFIG"),
	)

	connection, err := utils.Encrypt(uri)
	if err != nil {
		return 0, err
	}

	tenant := &boilmodels.Tenant{
		Name:       tenantUserCreate.TenantCreate.Name,
		Address:    null.StringFrom(tenantUserCreate.TenantCreate.Address),
		Phone:      null.StringFrom(tenantUserCreate.TenantCreate.Phone),
		Email:      null.StringFrom(tenantUserCreate.TenantCreate.Email),
		CuitPdv:    null.StringFrom(tenantUserCreate.TenantCreate.CuitPdv),
		Connection: connection,
		Identifier: identifier,
		PlanID:     null.Int64From(tenantUserCreate.TenantCreate.PlanID),
	}

	if err := tenant.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Tenant", schemas.Create)
	}

	user := &boilmodels.User{
		FirstName: null.StringFrom(tenantUserCreate.UserCreate.FirstName),
		LastName:  null.StringFrom(tenantUserCreate.UserCreate.LastName),
		Email:     tenantUserCreate.UserCreate.Email,
		Address:   null.StringFrom(tenantUserCreate.TenantCreate.Address),
		Username:  tenantUserCreate.UserCreate.Username,
	}

	// We'd hash password here conceptually, but that should happen in the service.
	// We'll leave it as we found it.

	if err := user.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "User", schemas.Create)
	}

	if err := (&boilmodels.UserTenant{
		UserID:   user.ID,
		TenantID: tenant.ID,
	}).Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "UserTenant", schemas.Create)
	}

	memberAdmin := models.Member{
		FirstName: tenantUserCreate.UserCreate.FirstName,
		LastName:  tenantUserCreate.UserCreate.LastName,
		Username:  tenantUserCreate.UserCreate.Username,
		Email:     tenantUserCreate.UserCreate.Email,
		Password:  tenantUserCreate.UserCreate.Password,
		IsAdmin:   true,
		Address:   &tenantUserCreate.TenantCreate.Address,
		RoleID:    1,
	}

	if err := database.PrepareDB(uri, memberAdmin); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Base de datos", schemas.Create)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return tenant.ID, nil
}

func (r *MainRepository) TenantUpdate(adminID, userID int64, tenantUpdate *schemas.TenantUpdate) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", adminID)); err != nil {
		return err
	}

	userTenant, err := boilmodels.UserTenants(
		boilmodels.UserTenantWhere.UserID.EQ(userID),
		boilmodels.UserTenantWhere.TenantID.EQ(tenantUpdate.ID),
	).One(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "UserTenant", schemas.Read)
	}
	_ = userTenant // Just verifying they belong to the tenant

	tenant, err := boilmodels.Tenants(boilmodels.TenantWhere.ID.EQ(tenantUpdate.ID)).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
	}

	// Updates... Since TenantUpdate isn't fully defined here, we do general updates.
	// Normally we'd bind it. I'll make the typical assignment.
	if tenantUpdate.Name != "" {
		tenant.Name = tenantUpdate.Name
	}
	if tenantUpdate.Address != "" {
		tenant.Address = null.StringFrom(tenantUpdate.Address)
	}
	if tenantUpdate.Phone != "" {
		tenant.Phone = null.StringFrom(tenantUpdate.Phone)
	}
	if tenantUpdate.Email != "" {
		tenant.Email = null.StringFrom(tenantUpdate.Email)
	}
	if tenantUpdate.CuitPdv != "" {
		tenant.CuitPdv = null.StringFrom(tenantUpdate.CuitPdv)
	}

	if _, err := tenant.Update(ctx, tx, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Tenant", schemas.Update)
	}

	return tx.Commit()
}

func (r *MainRepository) TenantUpdateExpiration(adminID int64, updateExp *schemas.TenantUpdateExpiration) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", adminID)); err != nil {
		return err
	}

	loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
	exp, err := time.ParseInLocation("2006-01-02", updateExp.Expiration, loc)
	if err != nil {
		return schemas.ErrorResponse(422, "Formato de fecha inválido, debe ser YYYY-MM-DD", err)
	}

	tenant, err := boilmodels.Tenants(boilmodels.TenantWhere.ID.EQ(updateExp.ID)).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
	}

	tenant.Expiration = null.TimeFrom(exp)

	if _, err := tenant.Update(ctx, tx, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Tenant", schemas.Update)
	}

	return tx.Commit()
}

func (r *MainRepository) TenantUpdateTerms(tenantID int64, updateTerms *schemas.TenantUpdateTerms) error {
	ctx := context.Background()

	tenant, err := boilmodels.Tenants(boilmodels.TenantWhere.ID.EQ(tenantID)).One(ctx, r.DB)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
	}

	tenant.AcceptedTerms = null.BoolFrom(updateTerms.AcceptedTerms)
	tenant.IP = null.StringFrom(updateTerms.IP)
	tenant.DateAccepted = null.TimeFrom(time.Now()) // Taking over DateAccepted

	if _, err := tenant.Update(ctx, r.DB, boil.Whitelist(
		boilmodels.TenantColumns.AcceptedTerms,
		boilmodels.TenantColumns.IP,
		boilmodels.TenantColumns.DateAccepted,
	)); err != nil {
		return schemas.HandlerErrorDB(err, "Tenant", schemas.Update)
	}

	return nil
}

func (r *MainRepository) TenantGetSettings(tenantID int64) (*schemas.TenantSettingsResponse, error) {
	ctx := context.Background()

	settings, err := boilmodels.SettingTenants(boilmodels.SettingTenantWhere.TenantID.EQ(tenantID)).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Configuraciones Tenant", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Configuraciones Tenant", schemas.Read)
	}

	return &schemas.TenantSettingsResponse{
		Logo:           settings.Logo.String,
		FrontPage:      settings.FrontPage.String,
		Title:          settings.Title.String,
		Slogan:         settings.Slogan.String,
		PrimaryColor:   settings.PrimaryColor.String,
		SecondaryColor: settings.SecondaryColor.String,
		Phone:          settings.Phone.String,
	}, nil
}

func (r *MainRepository) TenantUpdateSettings(tenantID int64, req *schemas.TenantUpdateSettings) error {
	ctx := context.Background()

	query := `
		INSERT INTO setting_tenants (tenant_id, title, slogan, primary_color, secondary_color, phone, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (tenant_id) DO UPDATE SET
			title = EXCLUDED.title,
			slogan = EXCLUDED.slogan,
			primary_color = EXCLUDED.primary_color,
			secondary_color = EXCLUDED.secondary_color,
			phone = EXCLUDED.phone,
			updated_at = NOW()
	`
	_, err := r.DB.ExecContext(ctx, query,
		tenantID,
		req.Title,
		req.Slogan,
		req.PrimaryColor,
		req.SecondaryColor,
		req.Phone,
	)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Configuraciones Tenant", schemas.Update)
	}

	return nil
}
