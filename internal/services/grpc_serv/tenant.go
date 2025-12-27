package grpc_serv

import (
	"context"

	pb "github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// // TenantGRPCServer debe cumplir con la interfaz generada
// type GRPCTenantService struct {
// 	pb.UnimplementedTenantServiceServer // Recomendado para compatibilidad futura
//     // Aquí podrías inyectar tu repositorio de base de datos
//     // Repo domain.TenantRepository
// }

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