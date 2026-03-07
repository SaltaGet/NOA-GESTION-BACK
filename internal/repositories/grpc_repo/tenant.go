package grpc_repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/master"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/utils"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *GrpcMainRepository) ListTenants() ([]master.Tenant, error) {
	var results []master.Tenant
	now := time.Now().UTC()
	targetModule := "ecommerce"

	tenants, err := master.Tenants(
		qm.InnerJoin("tenant_modules ON tenant_modules.tenant_id = tenants.id"),
		qm.InnerJoin("modules ON modules.id = tenant_modules.module_id"),
		qm.Where("modules.name = ?", targetModule),
		qm.Where("tenant_modules.expiration > ?", now),
		qm.Where("tenants.expiration > ?", now),
		qm.Where("tenants.is_active = ?", true),
		qm.GroupBy("tenants.id"),
		qm.Load(master.TenantRels.SettingTenant),
		qm.Load(master.TenantRels.Credential, qm.Select("tenant_id", "access_token_mp", "token_email")),
		qm.Load(master.TenantRels.TenantModules,
			qm.InnerJoin("modules ON modules.id = tenant_modules.module_id"),
			qm.Where("modules.name = ?", targetModule),
			qm.Where("tenant_modules.expiration > ?", now),
		),
	).All(context.Background(), r.DB)

	if err != nil {
		return nil, status.Error(codes.Internal, "Error al obtener los tenants")
	}

	for _, t := range tenants {
		results = append(results, *t)
	}

	return results, nil
}

func (r *GrpcMainRepository) GetTenant(req *pb.TenantRequest) (*master.Tenant, error) {
	tenant, err := master.Tenants(
		qm.Where("identifier = ?", req.Identifier),
		qm.Load(master.TenantRels.SettingTenant),
	).One(context.Background(), r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Tenant no encontrado")
		}
		return nil, status.Error(codes.Internal, "Error al obtener el tenant")
	}

	if !tenant.IsActive {
		return nil, status.Error(codes.PermissionDenied, "El tenant no esta activo")
	}

	return tenant, nil
}

func (r *GrpcMainRepository) UpdateImageSetting(ctx context.Context, req *pb.TenantRequestImageSetting) (*pb.TenantUpdateImageResponse, error) {
	resp := &pb.TenantUpdateImageResponse{}

	tenant, err := master.Tenants(
		qm.Select("id", "is_active"),
		qm.Where("identifier = ?", req.TenantIdentifier),
	).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Tenant no encontrado")
		}
		return nil, status.Error(codes.Internal, "Error al obtener el tenant")
	}

	if !tenant.IsActive {
		return nil, status.Error(codes.PermissionDenied, "El tenant no esta activo")
	}

	setting, err := master.SettingTenants(qm.Where("tenant_id = ?", tenant.ID)).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Setting no encontrado")
		}
		return nil, status.Error(codes.Internal, "Error al obtener el setting")
	}

	resp.LogoUuid = utils.Ternary(req.LogoUuid == nil, nil, setting.Logo.Ptr())
	resp.FrontPageUuid = utils.Ternary(req.FrontPageUuid == nil, nil, setting.FrontPage.Ptr())

	var updateCols []string
	if req.LogoUuid != nil {
		setting.Logo = null.StringFrom(*req.LogoUuid)
		updateCols = append(updateCols, master.SettingTenantColumns.Logo)
	}
	if req.FrontPageUuid != nil {
		setting.FrontPage = null.StringFrom(*req.FrontPageUuid)
		updateCols = append(updateCols, master.SettingTenantColumns.FrontPage)
	}

	if len(updateCols) > 0 {
		if _, err := setting.Update(ctx, r.DB, boil.Whitelist(updateCols...)); err != nil {
			return nil, status.Error(codes.Internal, "Error al actualizar el tenant")
		}
	}

	return resp, nil
}
