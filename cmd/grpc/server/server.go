package server

import (
	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/ports/grpc_ports"
)

type GrpcTenantServer struct {
	pb.UnimplementedTenantServiceServer
	GrpcTenantService grpc_ports.GrpcTenantService
}

type GrpcProductServer struct {
	pb.UnimplementedProductServiceServer
}

type GrpcCategoryServer struct {
	pb.UnimplementedCategoryServiceServer
}