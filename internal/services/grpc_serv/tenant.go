package grpc_serv

import (
	"context"

	pb "github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/aarondl/null/v8"
)

func nullStringPtr(ns null.String) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func (s *GrpcTenantService) ListTenants(ctx context.Context) (*pb.ListTenantsResponse, error) {
	tenants, err := s.GrpcTenantRepository.ListTenants()
	if err != nil {
		return nil, err
	}

	// Convertir a proto
	// return &pb.ListTenantsResponse{Tenants: tenants}, nil
	var protoTenants []*pb.TenantResponse
	for _, tenant := range tenants {
		pbTenant := &pb.TenantResponse{
			Id:            tenant.ID,
			Name:          tenant.Name,
			Identifier:    tenant.Identifier,
			Address:       tenant.Address,
			Phone:         tenant.Phone,
			Email:         tenant.Email,
			SettingTenant: &pb.SettingTenant{},
		}

		if tenant.R != nil && tenant.R.SettingTenant != nil {
			pbTenant.SettingTenant = &pb.SettingTenant{
				Id:             tenant.R.SettingTenant.ID,
				Logo:           nullStringPtr(tenant.R.SettingTenant.Logo),
				FrontPage:      nullStringPtr(tenant.R.SettingTenant.FrontPage),
				Title:          nullStringPtr(tenant.R.SettingTenant.Title),
				Slogan:         nullStringPtr(tenant.R.SettingTenant.Slogan),
				PrimaryColor:   nullStringPtr(tenant.R.SettingTenant.PrimaryColor),
				SecondaryColor: nullStringPtr(tenant.R.SettingTenant.SecondaryColor),
			}
		}

		if tenant.R != nil && tenant.R.Credential != nil {
			pbTenant.TokenMp = nullStringPtr(tenant.R.Credential.AccessTokenMP)
			pbTenant.TokenEmail = nullStringPtr(tenant.R.Credential.TokenEmail)
		}

		protoTenants = append(protoTenants, pbTenant)
	}

	return &pb.ListTenantsResponse{Tenants: protoTenants}, nil
}

func (s *GrpcTenantService) GetTenant(req *pb.TenantRequest) (*pb.TenantResponse, error) {
	tenant, err := s.GrpcTenantRepository.GetTenant(req)
	if err != nil {
		return nil, err
	}

	tenantResponse := &pb.TenantResponse{
		Id:            tenant.ID,
		Name:          tenant.Name,
		Identifier:    tenant.Identifier,
		Address:       tenant.Address,
		Phone:         tenant.Phone,
		Email:         tenant.Email,
		SettingTenant: &pb.SettingTenant{},
	}

	if tenant.R != nil && tenant.R.SettingTenant != nil {
		tenantResponse.SettingTenant = &pb.SettingTenant{
			Id:             tenant.R.SettingTenant.ID,
			Logo:           nullStringPtr(tenant.R.SettingTenant.Logo),
			FrontPage:      nullStringPtr(tenant.R.SettingTenant.FrontPage),
			Title:          nullStringPtr(tenant.R.SettingTenant.Title),
			Slogan:         nullStringPtr(tenant.R.SettingTenant.Slogan),
			PrimaryColor:   nullStringPtr(tenant.R.SettingTenant.PrimaryColor),
			SecondaryColor: nullStringPtr(tenant.R.SettingTenant.SecondaryColor),
		}
	}

	return tenantResponse, nil
}

func (s *GrpcTenantService) UpdateImageSetting(ctx context.Context, req *pb.TenantRequestImageSetting) (*pb.TenantUpdateImageResponse, error) {
	resp, err := s.GrpcTenantRepository.UpdateImageSetting(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
