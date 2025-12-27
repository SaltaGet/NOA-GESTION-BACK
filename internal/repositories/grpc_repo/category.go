package grpc_repo

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/models"

func (r *GrpcCategoryRepository) CategoryGetAll() ([]*models.Category, error) {
	var categories []*models.Category

	err := r.DB.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}