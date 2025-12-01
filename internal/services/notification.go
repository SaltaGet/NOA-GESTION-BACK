package services

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
)

func (s *NotificationService) NotificationStock(tenantID int64) ([]*schemas.ProductSimpleResponse, error) {
	products, err := s.NotificationRepository.NotificationStock(tenantID)
	if err != nil {
		return nil, err
	}

	
	var productsResponse []*schemas.ProductSimpleResponse
	for i:=0; i<len(products); i++ {
		var productResponse schemas.ProductSimpleResponse
		_ = copier.Copy(&productResponse, &products[i])

		if products[i].StockDeposit != nil {
			productResponse.Stock = products[i].StockDeposit.Stock
		} else {
			productResponse.Stock = 0
		}
		productsResponse = append(productsResponse, &productResponse)
	}

	return productsResponse, nil
}