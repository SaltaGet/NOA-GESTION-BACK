package grpc_serv

import (
	"context"

	pb "github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
)

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
			Id: tenant.ID,
			Name: tenant.Name,
			Identifier: tenant.Identifier,
			Address: tenant.Address,
			Phone: tenant.Phone,
			Email: tenant.Email,
			SettingTenant: &pb.SettingTenant{
				Id: tenant.Setting.ID,
				Logo: tenant.Setting.Logo,
				FrontPage: tenant.Setting.FrontPage,
				Title: tenant.Setting.Title,
				Slogan: tenant.Setting.Slogan,
				PrimaryColor: tenant.Setting.PrimaryColor,
				SecondaryColor: tenant.Setting.SecondaryColor,
			},
			TokenMp: tenant.Credentials.AccessTokenMP,
			TokenEmail: tenant.Credentials.TokenEmail,
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
		Id:         tenant.ID,
		Name:       tenant.Name,
		Identifier: tenant.Identifier,
		Address:    tenant.Address,
		Phone:      tenant.Phone,
		Email:      tenant.Email,
		SettingTenant: &pb.SettingTenant{
			Id: tenant.Setting.ID,
			Logo: tenant.Setting.Logo,
			FrontPage: tenant.Setting.FrontPage,
			Title: tenant.Setting.Title,
			Slogan: tenant.Setting.Slogan,
			PrimaryColor: tenant.Setting.PrimaryColor,
			SecondaryColor: tenant.Setting.SecondaryColor,
		},
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