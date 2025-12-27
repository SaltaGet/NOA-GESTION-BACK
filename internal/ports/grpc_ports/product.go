package grpc_ports

import (
	"context"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
)

type GrpcProductRepository interface {
	ProductGetByCode(code string) (*models.Product, error)
	ProductList(page, pageSize int32, categoryID *int32, search *string, sortBy int32) ([]*models.Product, int64, error)
}

type GrpcProductService interface {
	ProductGetByCode(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error)
	ProductList(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error)
}