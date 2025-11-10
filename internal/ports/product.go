package ports

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type ProductRepository interface {
	ProductGetByID(id int64) (*schemas.ProductFullResponse, error)
	ProductGetByCode(code string) (*schemas.ProductFullResponse, error)
	ProductGetByCategoryID(categoryID int64) ([]*schemas.ProductFullResponse, error)
	ProductGetByName(name string) ([]*schemas.ProductFullResponse, error)
	ProductGetAll(page, limit int) ([]*schemas.ProductFullResponse, int64, error)
	ProductCreate(productCreate *schemas.ProductCreate) (int64, error)
	ProductUpdate(productUpdate *schemas.ProductUpdate) error
	ProductPriceUpdate(productUpdate *schemas.ListPriceUpdate) error
	ProductDelete(id int64) error
}

type ProductService interface {
	ProductGetByID(id int64) (*schemas.ProductFullResponse, error)
	ProductGetByCode(code string) (*schemas.ProductFullResponse, error)
	ProductGetByName(name string) ([]*schemas.ProductFullResponse, error)
	ProductGetByCategoryID(categoryID int64) ([]*schemas.ProductFullResponse, error)
	ProductGetAll(page, limit int) ([]*schemas.ProductFullResponse, int64, error)
	ProductCreate(productCreate *schemas.ProductCreate) (int64, error)
	ProductUpdate(productUpdate *schemas.ProductUpdate) error
	ProductPriceUpdate(productUpdate *schemas.ListPriceUpdate) error
	ProductDelete(id int64) error
}