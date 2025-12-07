package ports

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type ProductRepository interface {
	ProductGetByID(id int64) (*models.Product, error)
	ProductGetByCode(code string) (*models.Product, error)
	ProductGetByCategoryID(categoryID int64) ([]*models.Product, error)
	ProductGetByName(name string) ([]*models.Product, error)
	ProductGetAll(page, limit int) ([]*models.Product, int64, error)
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