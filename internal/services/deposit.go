package services

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

func (s *DepositService) DepositGetByID(id int64) (*schemas.DepositResponse, error) {
	product, err := s.DepositRepository.DepositGetByID(id)
	if err != nil {
		return nil, err
	}

	desc := product.Description
	productResponse := &schemas.DepositResponse{
		ID:          product.ID,
		Code:        product.Code,
		Description: desc,
		Name:        product.Name,
		Category: schemas.CategoryResponse{
			ID:   product.Category.ID,
			Name: product.Category.Name,
		},
		Price: product.Price,
	}

	if product.StockDeposit != nil {
		productResponse.Stock = product.StockDeposit.Stock
	} else {
		productResponse.Stock = 0
	}

	return productResponse, nil
}

func (s *DepositService) DepositGetByCode(code string) (*schemas.DepositResponse, error) {
	product, err := s.DepositRepository.DepositGetByCode(code)
	if err != nil {
		return nil, err
	}

	desc := product.Description
	productResponse := &schemas.DepositResponse{
		ID:          product.ID,
		Code:        product.Code,
		Description: desc,
		Name:        product.Name,
		Category: schemas.CategoryResponse{
			ID:   product.Category.ID,
			Name: product.Category.Name,
		},
		Price: product.Price,
	}

	if product.StockDeposit != nil {
		productResponse.Stock = product.StockDeposit.Stock
	} else {
		productResponse.Stock = 0
	}

	return productResponse, nil
}

func (s *DepositService) DepositGetByName(name string) ([]*schemas.DepositResponse, error) {
	products, err := s.DepositRepository.DepositGetByName(name)
	if err != nil {
		return nil, err
	}

	productsResponse := make([]*schemas.DepositResponse, len(products))

	for i, prod := range products {
		desc := prod.Description
		productsResponse[i] = &schemas.DepositResponse{
			ID:          prod.ID,
			Code:        prod.Code,
			Description: desc,
			Name:        prod.Name,
			Category: schemas.CategoryResponse{
				ID:   prod.Category.ID,
				Name: prod.Category.Name,
			},
			Price: prod.Price,
		}
		if prod.StockDeposit != nil {
			productsResponse[i].Stock = prod.StockDeposit.Stock
		} else {
			productsResponse[i].Stock = 0
		}
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
		desc := prod.Description
		productsResponse[i] = &schemas.DepositResponse{
			ID:          prod.ID,
			Code:        prod.Code,
			Description: desc,
			Name:        prod.Name,
			Category: schemas.CategoryResponse{
				ID:   prod.Category.ID,
				Name: prod.Category.Name,
			},
			Price: prod.Price,
			// Stock: utils.FloatDefault(&prod.StockDeposit.Stock, 0),
		}

		if prod.StockDeposit != nil {
			productsResponse[i].Stock = prod.StockDeposit.Stock
		} else {
			productsResponse[i].Stock = 0
		}

	}

	return productsResponse, total, nil
}

func (s *DepositService) DepositUpdateStock(updateStock schemas.DepositUpdateStock) error {
	return s.DepositRepository.DepositUpdateStock(updateStock)
}
