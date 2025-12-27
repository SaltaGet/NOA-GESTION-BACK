package server

import (
	"context"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
)

func (h *GrpcTenantServer) ListTenants(ctx context.Context, req *pb.ListTenantsRequest) (*pb.ListTenantsResponse, error) {
    tenants, err := h.GrpcTenantService.ListTenants(ctx)
    if err != nil {
        return nil, err 
    }

    return tenants, nil
}