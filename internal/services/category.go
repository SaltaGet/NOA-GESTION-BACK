package services

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
)

func (s *CategoryService) CategoryGetByID(id int64) (*schemas.CategoryResponse, error) {
	category, err := s.CategoryRepository.CategoryGetByID(id)
	if err != nil {
		return nil, err
	}

	var categoryResponse schemas.CategoryResponse
	_ = copier.Copy(&categoryResponse, &category)

	return &categoryResponse, nil
}

func (s *CategoryService) CategoryGetAll() ([]*schemas.CategoryResponse, error) {
	categories, err := s.CategoryRepository.CategoryGetAll()
	if err != nil {
		return nil, err
	}

	var categoriesResponse []*schemas.CategoryResponse
	_ = copier.Copy(&categoriesResponse, &categories)

	return categoriesResponse, nil
}

func (s *CategoryService) CategoryCreate(categoryCreate *schemas.CategoryCreate) (int64, error) {
	return s.CategoryRepository.CategoryCreate(categoryCreate)
}

func (s *CategoryService) CategoryUpdate(categoryUpdate *schemas.CategoryUpdate) error {
	return s.CategoryRepository.CategoryUpdate(categoryUpdate)
}

func (s *CategoryService) CategoryDelete(id int64) error {
	return s.CategoryRepository.CategoryDelete(id)
}