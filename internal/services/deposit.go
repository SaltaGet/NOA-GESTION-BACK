package services

import (
	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/utils"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func mapToDepositResponse(product *boilmodels.Product) *schemas.DepositResponse {
	if product == nil {
		return nil
	}

	var desc *string
	if product.Description.Valid {
		desc = &product.Description.String
	}

	var primaryImage *string
	if product.PrimaryImage.Valid {
		primaryImage = &product.PrimaryImage.String
	}

	var secImages string
	if product.SecondaryImages.Valid {
		secImages = product.SecondaryImages.String
	}

	catID := int64(0)
	catName := ""
	if product.R != nil && product.R.Category != nil {
		catID = product.R.Category.ID
		catName = product.R.Category.Name
	}

	price, _ := product.Price.Float64()
	var stock float64 = 0
	if product.R != nil && len(product.R.Deposits) > 0 {
		stock, _ = product.R.Deposits[0].Stock.Float64()
	}

	return &schemas.DepositResponse{
		ID:             product.ID,
		Code:           product.Code,
		Description:    desc,
		Name:           product.Name,
		PrimaryImage:   primaryImage,
		SecondaryImage: utils.SplitStrings(&secImages),
		Category: schemas.CategoryResponse{
			ID:   catID,
			Name: catName,
		},
		Price: price,
		Stock: stock,
	}
}

func (s *DepositService) DepositGetByID(id int64) (*schemas.DepositResponse, error) {
	product, err := s.DepositRepository.DepositGetByID(id)
	if err != nil {
		return nil, err
	}

	return mapToDepositResponse(product), nil
}

func (s *DepositService) DepositGetByCode(code string) (*schemas.DepositResponse, error) {
	product, err := s.DepositRepository.DepositGetByCode(code)
	if err != nil {
		return nil, err
	}

	return mapToDepositResponse(product), nil
}

func (s *DepositService) DepositGetByName(name string) ([]*schemas.DepositResponse, error) {
	products, err := s.DepositRepository.DepositGetByName(name)
	if err != nil {
		return nil, err
	}

	productsResponse := make([]*schemas.DepositResponse, len(products))

	for i, prod := range products {
		productsResponse[i] = mapToDepositResponse(prod)
	}

	return productsResponse, nil
}

func (s *DepositService) DepositGetAll(page, limit int) ([]*schemas.DepositResponse, int64, error) {
	products, total, err := s.DepositRepository.DepositGetAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	productsResponse := make([]*schemas.DepositResponse, len(products))

	for i, prod := range products {
		productsResponse[i] = mapToDepositResponse(prod)

	}

	return productsResponse, total, nil
}

func (s *DepositService) DepositUpdateStock(memberID int64, updateStock schemas.DepositUpdateStock) error {
	return s.DepositRepository.DepositUpdateStock(memberID, updateStock)
}
