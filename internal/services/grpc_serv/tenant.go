package grpc_serv

import (
	"context"

	pb "github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GrpcTenantService) ListTenants(ctx context.Context) (*pb.ListTenantsResponse, error) {
	tenants, err := s.GrpcTenantRepository.ListTenants()
	if err != nil {
		return nil, err
	}

	// Convertir a proto
	// return &pb.ListTenantsResponse{Tenants: tenants}, nil
	var protoTenants []*pb.Tenant
	for _, tenant := range tenants {
		var expirationProto *timestamppb.Timestamp

    // 2. SOLO si la fecha en BD NO es nil, hacemos la conversión
    if tenant.Expiration != nil {
        expirationProto = timestamppb.New(*tenant.Expiration)
    }

    // 3. Construimos el mensaje Proto de forma segura
    pbTenant := &pb.Tenant{
        Identifier: tenant.Identifier,
        IsActive:   tenant.IsActive,
        Expiration: expirationProto, // Aquí pasamos el puntero (sea nil o tenga valor)
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
			Logo: *tenant.Setting.Logo,
			FrontPage: *tenant.Setting.FrontPage,
			Title: *tenant.Setting.Title,
			Slogan: *tenant.Setting.Slogan,
			PrimaryColor: *tenant.Setting.PrimaryColor,
			SecondaryColor: *tenant.Setting.SecondaryColor,
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