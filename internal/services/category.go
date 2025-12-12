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

func (s *CategoryService) CategoryCreate(memberID int64, categoryCreate *schemas.CategoryCreate) (int64, error) {
	return s.CategoryRepository.CategoryCreate(memberID ,categoryCreate)
}

func (s *CategoryService) CategoryUpdate(memberID int64, categoryUpdate *schemas.CategoryUpdate) error {
	return s.CategoryRepository.CategoryUpdate(memberID ,categoryUpdate)
}

func (s *CategoryService) CategoryDelete(memberID, id int64) error {
	return s.CategoryRepository.CategoryDelete(memberID ,id)
}