package server

import (
	"context"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	grpc_cache "github.com/SaltaGet/NOA-GESTION-BACK/internal/cache/grpc"
)

func (s *GrpcCategoryServer) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
	deps := grpc_cache.GetGrpcContainerFromContext(ctx)
	categories, err := deps.Services.GrpcCategoryService.CategoryGetAll()
	if err != nil {
		return nil, err
	}

	return categories, nil
}