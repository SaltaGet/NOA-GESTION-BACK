package grpc_serv

import (
	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
)

func (s *GrpcCategoryService) CategoryGetAll() (*pb.ListCategoriesResponse, error) {
	categoriesModel, err := s.GrpcCategoryRepository.CategoryGetAll()
	if err != nil {
		return nil, err
	}

	var categoriesSchema []*pb.Category
	for _, categoryModel := range categoriesModel {
		categorySchema := pb.Category{
			Id:   categoryModel.ID,
			Name: categoryModel.Name,
		}
		categoriesSchema = append(categoriesSchema, &categorySchema)
	}

	return &pb.ListCategoriesResponse{
		Categories: categoriesSchema,
	}, nil
}