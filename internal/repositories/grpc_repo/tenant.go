package grpc_repo

import (
	"context"
	"errors"

	// "time"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *GrpcMainRepository) ListTenants() ([]models.Tenant, error) {
	var tenants []models.Tenant
	if err := r.DB.Find(&tenants).Error; err != nil {
		return nil, status.Error(codes.Internal, "Error al obtener los tenants")
	}

	return tenants, nil
}

func (r *GrpcMainRepository) GetTenant(req *pb.TenantRequest) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := r.DB.Preload("Setting").Where("identifier = ?", req.Identifier).First(&tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Tenant no encontrado")
		}
		return nil, status.Error(codes.Internal, "Error al obtener el tenant")
	}

	if !tenant.IsActive {
		return nil, status.Error(codes.PermissionDenied, "El tenant no esta activo")
	}

	// if tenant.Expiration.Before(time.Now()) {
	// 	return nil, status.Error(codes.PermissionDenied, "El tenant ha expirado")
	// }

	return &tenant, nil
}

func (r *GrpcMainRepository) UpdateImageSetting(ctx context.Context, req *pb.TenantRequestImageSetting) (*pb.TenantUpdateImageResponse, error) {
	resp := &pb.TenantUpdateImageResponse{}

	var tenant models.Tenant
	if err := r.DB.Select("id", "is_active").Where("identifier = ?", req.TenantIdentifier).First(&tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Tenant no encontrado")
		}
		return nil, status.Error(codes.Internal, "Error al obtener el tenant")
	}

	if !tenant.IsActive {
		return nil, status.Error(codes.PermissionDenied, "El tenant no esta activo")
	}

	var setting models.SettingTenant
	if err := r.DB.Where("tenant_id = ?", tenant.ID).First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Setting no encontrado")
		}
		return nil, status.Error(codes.Internal, "Error al obtener el setting")
	}
	// if tenant.Expiration.Before(time.Now()) {
	// 	return nil, status.Error(codes.PermissionDenied, "El tenant ha expirado")
	// }

	resp.LogoUuid = utils.Ternary(req.LogoUuid == nil, nil, setting.Logo)
	resp.FrontPageUuid = utils.Ternary(req.FrontPageUuid == nil, nil, setting.FrontPage)
	updateData := models.SettingTenant{
		Logo:      req.LogoUuid,
		FrontPage: req.FrontPageUuid,
	}

	if err := r.DB.Model(&models.SettingTenant{}).Where("tenant_id = ?", tenant.ID).Updates(updateData).Error; err != nil {
		return nil, status.Error(codes.Internal, "Error al actualizar el tenant")
	}

	return resp, nil
}
