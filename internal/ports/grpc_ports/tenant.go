package grpc_ports

import (
	"context"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
)

type GrpcTenantRepository interface {
	// Define methods for tenant repository
	ListTenants() ([]models.Tenant, error)	
}

type GrpcTenantService interface {
	ListTenants(ctx context.Context) (*pb.ListTenantsResponse, error)	
}