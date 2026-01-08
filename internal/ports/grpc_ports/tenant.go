package grpc_ports

import (
	"context"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
)

type GrpcTenantRepository interface {
	// Define methods for tenant repository
	ListTenants() ([]models.Tenant, error)	
	GetTenant(req *pb.TenantRequest) (*models.Tenant, error)
	UpdateImageSetting(ctx context.Context, req *pb.TenantRequestImageSetting) (*pb.TenantUpdateImageResponse, error)
}

type GrpcTenantService interface {
	ListTenants(ctx context.Context) (*pb.ListTenantsResponse, error)	
	GetTenant(req *pb.TenantRequest) (*pb.TenantResponse, error)
	UpdateImageSetting(ctx context.Context, req *pb.TenantRequestImageSetting) (*pb.TenantUpdateImageResponse, error)
}