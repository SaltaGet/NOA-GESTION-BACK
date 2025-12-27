package server

import (
	"context"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	grpc_cache "github.com/SaltaGet/NOA-GESTION-BACK/internal/cache/grpc"
)



func (s *GrpcProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	deps := grpc_cache.GetGrpcContainerFromContext(ctx)
	prod, err := deps.Services.GrpcProductService.ProductGetByCode(ctx, req)
	if err != nil {
		return nil, err
	}

	return prod, nil
}

func (s *GrpcProductServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	deps := grpc_cache.GetGrpcContainerFromContext(ctx)
	return deps.Services.GrpcProductService.ProductList(ctx, req)
}