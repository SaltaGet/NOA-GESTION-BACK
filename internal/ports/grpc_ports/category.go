package grpc_ports

import (
	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
)


type GrpcCategoryRepository interface {
	CategoryGetAll() ([]*models.Category, error)
}

type GrpcCategoryService interface {
	CategoryGetAll() (*pb.ListCategoriesResponse, error)
}