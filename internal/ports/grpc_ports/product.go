package grpc_ports

import (
	"context"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
)

type GrpcProductRepository interface {
	ProductGetByCode(code string) (*tenant.Product, error)
	ProductList(req *pb.ListProductsRequest) ([]*tenant.Product, int64, error)
	SaveUrlImage(req *pb.SaveImageRequest) (error)
	ProductGetByID(productId int64) (*tenant.Product, error)
	ValidateProducts(ctx context.Context, req *pb.ProductValidateRequest) ([]tenant.Product, error)
}

type GrpcProductService interface {
	ProductGetByCode(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error)
	ProductList(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error)
	SaveUrlImage(ctx context.Context, req *pb.SaveImageRequest) (*pb.SaveImageResponse, error)
	ProductGetByID(ctx context.Context, req *pb.ProductRequest) (*pb.Product, error)
	ValidateProducts(ctx context.Context, req *pb.ProductValidateRequest) (*pb.ProductValidateResponse, error)
}