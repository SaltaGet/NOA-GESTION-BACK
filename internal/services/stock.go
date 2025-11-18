package services

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (p *StockService) StockGetByID(id, pointID int64) (*schemas.ProductStockFullResponse, error) {
	return p.StockRepository.StockGetByID(id, pointID)
}

func (p *StockService) StockGetByCode(code string, pointID int64) (*schemas.ProductStockFullResponse, error) {
	return p.StockRepository.StockGetByCode(code, pointID)
}

func (p *StockService) StockGetByName(name string, pointID int64) ([]*schemas.ProductStockFullResponse, error) {
	return p.StockRepository.StockGetByName(name, pointID)
}

func (p *StockService) StockGetByCategoryID(categoryID, pointID int64) ([]*schemas.ProductStockFullResponse, error) {
	return p.StockRepository.StockGetByCategoryID(categoryID, pointID)
}

func (p *StockService) StockGetAll(page, limit int, pointID int64) ([]*schemas.ProductStockFullResponse, int64, error) {
	return p.StockRepository.StockGetAll(page, limit, pointID)
}
