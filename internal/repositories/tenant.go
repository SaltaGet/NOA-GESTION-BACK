package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"strconv"

	mastermodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/master"
	tenantmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/utils"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

func (r *MainRepository) TenantGetByID(tenantID int64) (*mastermodels.Tenant, error) {
	ctx := context.Background()

	tenant, err := mastermodels.Tenants(mastermodels.TenantWhere.ID.EQ(tenantID)).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
	}

	if !tenant.IsActive {
		return nil, schemas.ErrorResponse(403, "Tenant is inactive", fmt.Errorf("tenant is inactive"))
	}

	return tenant, nil
}

func (r *MainRepository) TenantGetByIdentifier(identifier string) (*mastermodels.Tenant, error) {
	ctx := context.Background()

	tenant, err := mastermodels.Tenants(mastermodels.TenantWhere.Identifier.EQ(identifier)).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al obtener los tenants", err)
	}

	if !tenant.IsActive {
		return nil, schemas.ErrorResponse(403, "Tenant es inactivo", fmt.Errorf("tenant es inactivo"))
	}

	return tenant, nil
}

func (r *MainRepository) TenantGetAll() (*[]schemas.TenantResponse, error) {
	ctx := context.Background()

	tenants, err := mastermodels.Tenants(
		qm.Select(
			mastermodels.TenantColumns.ID,
			mastermodels.TenantColumns.Name,
			mastermodels.TenantColumns.Address,
			mastermodels.TenantColumns.Phone,
			mastermodels.TenantColumns.Email,
			mastermodels.TenantColumns.IsActive,
			mastermodels.TenantColumns.Expiration,
			mastermodels.TenantColumns.CreatedAt,
			mastermodels.TenantColumns.UpdatedAt,
		),
		qm.InnerJoin("user_tenants ON tenants.id = user_tenants.tenant_id"),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
	}

	var response []schemas.TenantResponse
	for _, t := range tenants {
		response = append(response, schemas.TenantResponse{
			ID:         t.ID,
			Name:       t.Name,
			Address:    t.Address,
			Phone:      t.Phone,
			Email:      t.Email,
			IsActive:   t.IsActive,
			Expiration: t.Expiration.Time,
			CreatedAt:  t.CreatedAt,
			UpdatedAt:  t.UpdatedAt,
		})
	}

	return &response, nil
}

func (r *MainRepository) TenantGetConnectionByIdentifier(tenantIdentifier string) (*mastermodels.Tenant, error) {
	ctx := context.Background()

	tenant, err := mastermodels.Tenants(
		qm.Select("id", "connection"),
		mastermodels.TenantWhere.Identifier.EQ(tenantIdentifier),
	).One(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
	}

	return tenant, nil
}

func (r *MainRepository) TenantGetConections() ([]*mastermodels.Tenant, error) {
	ctx := context.Background()

	tenants, err := mastermodels.Tenants(
		qm.Select("id", "name", "connection", "identifier"),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
	}

	var response []*mastermodels.Tenant
	for _, tenant := range tenants {
		connectionDecrypted, err := utils.Decrypt(tenant.Connection)
		if err != nil {
			return nil, schemas.ErrorResponse(500, "Error interno al obtener las connections", err)
		}

		tenant.Connection = connectionDecrypted
		response = append(response, tenant)
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

	tenant := &mastermodels.Tenant{
		Name:       tenantCreate.Name,
		Address:    tenantCreate.Address,
		Phone:      tenantCreate.Phone,
		Email:      tenantCreate.Email,
		CuitPDV:    tenantCreate.CuitPdv,
		Connection: connection,
		Identifier: identifier,
	}

	if err := tenant.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Tenant", schemas.Create)
	}

	user, err := mastermodels.Users(mastermodels.UserWhere.ID.EQ(userID)).One(ctx, tx)
	if err != nil {
		return 0, schemas.HandlerErrorDB(err, "User", schemas.Read)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO user_tenants (user_id, tenant_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", user.ID, tenant.ID)
	if err != nil {
		return 0, schemas.HandlerErrorDB(err, "UserTenant", schemas.Create)
	}

	memberAdmin := tenantmodels.Member{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Password:  "1",
		IsAdmin:   true,
		RoleID:    1,
	}
	if user.Address.Valid {
		memberAdmin.Address = user.Address
	}

	if err := database.PrepareDB(uri, memberAdmin); err != nil {
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

	tenant := &mastermodels.Tenant{
		Name:       tenantUserCreate.TenantCreate.Name,
		Address:    tenantUserCreate.TenantCreate.Address,
		Phone:      tenantUserCreate.TenantCreate.Phone,
		Email:      tenantUserCreate.TenantCreate.Email,
		CuitPDV:    tenantUserCreate.TenantCreate.CuitPdv,
		Connection: connection,
		Identifier: identifier,
		PlanID:     tenantUserCreate.TenantCreate.PlanID,
	}

	if err := tenant.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Tenant", schemas.Create)
	}

	user := &mastermodels.User{
		FirstName: tenantUserCreate.UserCreate.FirstName,
		LastName:  tenantUserCreate.UserCreate.LastName,
		Email:     tenantUserCreate.UserCreate.Email,
		Address:   null.StringFrom(tenantUserCreate.TenantCreate.Address),
		Username:  tenantUserCreate.UserCreate.Username,
	}

	if err := user.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "User", schemas.Create)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO user_tenants (user_id, tenant_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", user.ID, tenant.ID)
	if err != nil {
		return 0, schemas.HandlerErrorDB(err, "UserTenant", schemas.Create)
	}

	memberAdmin := tenantmodels.Member{
		FirstName: tenantUserCreate.UserCreate.FirstName,
		LastName:  tenantUserCreate.UserCreate.LastName,
		Username:  tenantUserCreate.UserCreate.Username,
		Email:     tenantUserCreate.UserCreate.Email,
		Password:  tenantUserCreate.UserCreate.Password,
		IsAdmin:   true,
		RoleID:    1,
	}
	if tenantUserCreate.TenantCreate.Address != "" {
		memberAdmin.Address = null.StringFrom(tenantUserCreate.TenantCreate.Address)
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

	var count int
	err = tx.QueryRowContext(ctx, "SELECT count(*) FROM user_tenants WHERE user_id = $1 AND tenant_id = $2", userID, tenantUpdate.ID).Scan(&count)
	if err != nil || count == 0 {
		return schemas.HandlerErrorDB(errors.New("user not in tenant"), "UserTenant", schemas.Read)
	}

	tenantIDInt, _ := strconv.ParseInt(tenantUpdate.ID, 10, 64)

	tenant, err := mastermodels.Tenants(mastermodels.TenantWhere.ID.EQ(tenantIDInt)).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
	}

	if tenantUpdate.Name != "" {
		tenant.Name = tenantUpdate.Name
	}
	if tenantUpdate.Address != "" {
		tenant.Address = tenantUpdate.Address
	}
	if tenantUpdate.Phone != "" {
		tenant.Phone = tenantUpdate.Phone
	}
	if tenantUpdate.Email != "" {
		tenant.Email = tenantUpdate.Email
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

	tenant, err := mastermodels.Tenants(mastermodels.TenantWhere.ID.EQ(updateExp.ID)).One(ctx, tx)
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

	tenant, err := mastermodels.Tenants(mastermodels.TenantWhere.ID.EQ(tenantID)).One(ctx, r.DB)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
	}

	tenant.AcceptedTerms = updateTerms.AcceptedTerms
	tenant.IP = null.StringFrom(updateTerms.IP)
	tenant.DateAccepted = null.TimeFrom(time.Now())

	if _, err := tenant.Update(ctx, r.DB, boil.Whitelist(
		mastermodels.TenantColumns.AcceptedTerms,
		mastermodels.TenantColumns.IP,
		mastermodels.TenantColumns.DateAccepted,
	)); err != nil {
		return schemas.HandlerErrorDB(err, "Tenant", schemas.Update)
	}

	return nil
}

func (r *MainRepository) TenantGetSettings(tenantID int64) (*schemas.TenantSettingsResponse, error) {
	ctx := context.Background()

	settings, err := mastermodels.SettingTenants(mastermodels.SettingTenantWhere.TenantID.EQ(tenantID)).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Configuraciones Tenant", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Configuraciones Tenant", schemas.Read)
	}

	res := &schemas.TenantSettingsResponse{}
	if settings.Logo.Valid {
		res.Logo = &settings.Logo.String
	}
	if settings.FrontPage.Valid {
		res.FrontPage = &settings.FrontPage.String
	}
	if settings.Title.Valid {
		res.Title = &settings.Title.String
	}
	if settings.Slogan.Valid {
		res.Slogan = &settings.Slogan.String
	}
	if settings.PrimaryColor.Valid {
		res.PrimaryColor = &settings.PrimaryColor.String
	}
	if settings.SecondaryColor.Valid {
		res.SecondaryColor = &settings.SecondaryColor.String
	}
	if settings.Phone.Valid {
		res.Phone = &settings.Phone.String
	}
	return res, nil
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
