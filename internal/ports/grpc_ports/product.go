package grpc_ports

import (
	"context"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
)

type GrpcProductRepository interface {
	ProductGetByCode(code string) (*models.Product, error)
	ProductList(page, pageSize int32, categoryID *int32, search *string, sortBy int32) ([]*models.Product, int64, error)
	SaveUrlImage(req *pb.SaveImageRequest) (error)
	ProductGetByID(productId int64) (*models.Product, error)
}

type GrpcProductService interface {
	ProductGetByCode(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error)
	ProductList(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error)
	SaveUrlImage(ctx context.Context, req *pb.SaveImageRequest) (*pb.SaveImageResponse, error)
	ProductGetByID(ctx context.Context, req *pb.ProductRequest) (*pb.Product, error)
}