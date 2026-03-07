package grpc_ports

import (
	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
)


type GrpcCategoryRepository interface {
	CategoryGetAll() ([]*tenant.Category, error)
}

type GrpcCategoryService interface {
	CategoryGetAll() (*pb.ListCategoriesResponse, error)
}