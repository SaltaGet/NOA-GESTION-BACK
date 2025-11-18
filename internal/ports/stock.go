package ports

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type StockRepository interface {
	StockGetByID(id, pointID int64) (*schemas.ProductStockFullResponse, error)
	StockGetByCode(code string, pointID int64) (*schemas.ProductStockFullResponse, error)
	StockGetByCategoryID(categoryID, pointID int64) ([]*schemas.ProductStockFullResponse, error)
	StockGetByName(name string, pointID int64) ([]*schemas.ProductStockFullResponse, error)
	StockGetAll(page, limit int, pointID int64) ([]*schemas.ProductStockFullResponse, int64, error)
}

type StockService interface {
	StockGetByID(id, pointID int64) (*schemas.ProductStockFullResponse, error)
	StockGetByCode(code string, pointID int64) (*schemas.ProductStockFullResponse, error)
	StockGetByCategoryID(categoryID, pointID int64) ([]*schemas.ProductStockFullResponse, error)
	StockGetByName(name string, pointID int64) ([]*schemas.ProductStockFullResponse, error)
	StockGetAll(page, limit int, pointID int64) ([]*schemas.ProductStockFullResponse, int64, error)
}
