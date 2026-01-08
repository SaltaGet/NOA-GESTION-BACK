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

func (h *GrpcTenantServer) TenantGetIdentifier(ctx context.Context, req *pb.TenantRequest) (*pb.TenantResponse, error) {
	tenant, err := h.GrpcTenantService.GetTenant(req)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}

func (h *GrpcTenantServer) TenantUpdateImageSetting(ctx context.Context, req *pb.TenantRequestImageSetting) (*pb.TenantUpdateImageResponse, error) {
	resp, err := h.GrpcTenantService.UpdateImageSetting(ctx, req)
	if err != nil {
		return nil, err
	}

  return resp, nil
}
