package grpc_serv

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/ports/grpc_ports"

// type GrpcAuthService struct {
// 	GrpcAuthRepository grpc_ports.GrpcAuthRepository
// }

type GrpcTenantService struct {
	GrpcTenantRepository grpc_ports.GrpcTenantRepository
}

type GrpcProductService struct {
	GrpcProductRepository grpc_ports.GrpcProductRepository
}

type GrpcCategoryService struct {
	GrpcCategoryRepository grpc_ports.GrpcCategoryRepository
}

type GrpcMPService struct {
	GrpcMPRepository grpc_ports.GrpcMPRepository
}