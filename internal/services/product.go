package services

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (p *ProductService) ProductGetByID(id int64) (*schemas.ProductFullResponse, error) {
	product, err := p.ProductRepository.ProductGetByID(id)
	if err != nil {
		return nil, err
	}

	var productResponse schemas.ProductFullResponse

	productResponse.ID = product.ID
	productResponse.Code = product.Code
	productResponse.Name = product.Name
	productResponse.Description = product.Description
	productResponse.Category = schemas.CategoryResponse{
		ID:   product.Category.ID,
		Name: product.Category.Name,
	}
	productResponse.Price = product.Price

	if product.StockDeposit != nil {
		productResponse.StockDeposit = product.StockDeposit.Stock
	} else {
		productResponse.StockDeposit = 0
	}

	productResponse.Notifier = product.Notifier
	productResponse.MinAmount = product.MinAmount

	for _, stock := range product.StockPointSales {
		productResponse.StockPointSales = append(productResponse.StockPointSales, &schemas.PointSaleStockResponse{
			ID:        stock.PointSale.ID,
			Name:      stock.PointSale.Name,
			Stock:     stock.Stock,
			IsDeposit: stock.PointSale.IsDeposit,
		})
	}

	return &productResponse, nil
}

func (p *ProductService) ProductGetByCode(code string) (*schemas.ProductFullResponse, error) {
	product, err := p.ProductRepository.ProductGetByCode(code)
	if err != nil {
		return nil, err
	}

	var productResponse schemas.ProductFullResponse

	productResponse.ID = product.ID
	productResponse.Code = product.Code
	productResponse.Name = product.Name
	productResponse.Description = product.Description
	productResponse.Category = schemas.CategoryResponse{
		ID:   product.Category.ID,
		Name: product.Category.Name,
	}
	productResponse.Price = product.Price

	if product.StockDeposit != nil {
		productResponse.StockDeposit = product.StockDeposit.Stock
	} else {
		productResponse.StockDeposit = 0
	}

	productResponse.Notifier = product.Notifier
	productResponse.MinAmount = product.MinAmount

	for _, stock := range product.StockPointSales {
		productResponse.StockPointSales = append(productResponse.StockPointSales, &schemas.PointSaleStockResponse{
			ID:        stock.PointSale.ID,
			Name:      stock.PointSale.Name,
			Stock:     stock.Stock,
			IsDeposit: stock.PointSale.IsDeposit,
		})
	}

	return &productResponse, nil
}

func (p *ProductService) ProductGetByName(name string) ([]*schemas.ProductFullResponse, error) {
	products, err := p.ProductRepository.ProductGetByName(name)
	if err != nil {
		return nil, err
	}

	productsResponse := make([]*schemas.ProductFullResponse, len(products))

	for i, prod := range products {
		productsResponse[i] = &schemas.ProductFullResponse{
			ID:          prod.ID,
			Code:        prod.Code,
			Name:        prod.Name,
			Description: prod.Description,
			Category: schemas.CategoryResponse{
				ID:   prod.Category.ID,
				Name: prod.Category.Name,
			},
			Price:     prod.Price,
			Notifier:  prod.Notifier,
			MinAmount: prod.MinAmount,
		}
		if prod.StockDeposit != nil {
			productsResponse[i].StockDeposit = prod.StockDeposit.Stock
		} else {
			productsResponse[i].StockDeposit = 0
		}
		for _, stock := range prod.StockPointSales {
			productsResponse[i].StockPointSales = append(productsResponse[i].StockPointSales, &schemas.PointSaleStockResponse{
				ID:        stock.PointSale.ID,
				Name:      stock.PointSale.Name,
				Stock:     stock.Stock,
				IsDeposit: stock.PointSale.IsDeposit,
			})
		}
	}

	return productsResponse, nil
}

func (p *ProductService) ProductGetByCategoryID(categoryID int64) ([]*schemas.ProductFullResponse, error) {
	products, err := p.ProductRepository.ProductGetByCategoryID(categoryID)
	if err != nil {
		return nil, err
	}

	productsResponse := make([]*schemas.ProductFullResponse, len(products))

	for i, prod := range products {
		productsResponse[i] = &schemas.ProductFullResponse{
			ID:          prod.ID,
			Code:        prod.Code,
			Name:        prod.Name,
			Description: prod.Description,
			Category: schemas.CategoryResponse{
				ID:   prod.Category.ID,
				Name: prod.Category.Name,
			},
			Price:     prod.Price,
			Notifier:  prod.Notifier,
			MinAmount: prod.MinAmount,
		}
		if prod.StockDeposit != nil {
			productsResponse[i].StockDeposit = prod.StockDeposit.Stock
		} else {
			productsResponse[i].StockDeposit = 0
		}
		for _, stock := range prod.StockPointSales {
			productsResponse[i].StockPointSales = append(productsResponse[i].StockPointSales, &schemas.PointSaleStockResponse{
				ID:        stock.PointSale.ID,
				Name:      stock.PointSale.Name,
				Stock:     stock.Stock,
				IsDeposit: stock.PointSale.IsDeposit,
			})
		}
	}

	return productsResponse, nil
}

func (p *ProductService) ProductGetAll(page, limit int) ([]*schemas.ProductFullResponse, int64, error) {
	products, total, err := p.ProductRepository.ProductGetAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	productsResponse := make([]*schemas.ProductFullResponse, len(products))

	for i, prod := range products {
		productsResponse[i] = &schemas.ProductFullResponse{
			ID:          prod.ID,
			Code:        prod.Code,
			Name:        prod.Name,
			Description: prod.Description,
			Category: schemas.CategoryResponse{
				ID:   prod.Category.ID,
				Name: prod.Category.Name,
			},
			Price:     prod.Price,
			Notifier:  prod.Notifier,
			MinAmount: prod.MinAmount,
		}
		if prod.StockDeposit != nil {
			productsResponse[i].StockDeposit = prod.StockDeposit.Stock
		} else {
			productsResponse[i].StockDeposit = 0
		}
		for _, stock := range prod.StockPointSales {
			productsResponse[i].StockPointSales = append(productsResponse[i].StockPointSales, &schemas.PointSaleStockResponse{
				ID:        stock.PointSale.ID,
				Name:      stock.PointSale.Name,
				Stock:     stock.Stock,
				IsDeposit: stock.PointSale.IsDeposit,
			})
		}
	}

	return productsResponse, total, nil
}

func (p *ProductService) ProductCreate(memberID int64, productCreate *schemas.ProductCreate, plan *schemas.PlanResponseDTO) (int64, error) {
	return p.ProductRepository.ProductCreate(memberID, productCreate, plan)
}

func (p *ProductService) ProductUpdate(memberID int64, productUpdate *schemas.ProductUpdate) error {
	return p.ProductRepository.ProductUpdate(memberID, productUpdate)
}
func (p *ProductService) ProductPriceUpdate(memberID int64, productUpdate *schemas.ListPriceUpdate) error {
	return p.ProductRepository.ProductPriceUpdate(memberID, productUpdate)
}

func (p *ProductService) ProductDelete(memberID int64, id int64) error {
	return p.ProductRepository.ProductDelete(memberID, id)

}
