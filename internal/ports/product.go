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
	ProductGetByCodeToQR(code string) (*models.Product, error)
	ProductCreate(memberID int64, productCreate *schemas.ProductCreate, plan *schemas.PlanResponseDTO) (int64, error)
	ProductUpdate(memberID int64, productUpdate *schemas.ProductUpdate) error
	ProductPriceUpdate(memberID int64, productUpdate *schemas.ListPriceUpdate) error
	ProductDelete(memberID int64, id int64) error
}

type ProductService interface {
	ProductGetByID(id int64) (*schemas.ProductFullResponse, error)
	ProductGetByCode(code string) (*schemas.ProductFullResponse, error)
	ProductGetByName(name string) ([]*schemas.ProductFullResponse, error)
	ProductGetByCategoryID(categoryID int64) ([]*schemas.ProductFullResponse, error)
	ProductGetAll(page, limit int) ([]*schemas.ProductFullResponse, int64, error)
	ProductGenerateQR(code string) ([]byte, error)
	ProductCreate(memberID int64, productCreate *schemas.ProductCreate, plan *schemas.PlanResponseDTO) (int64, error)
	ProductUpdate(memberID int64, productUpdate *schemas.ProductUpdate) error
	ProductPriceUpdate(memberID int64, productUpdate *schemas.ListPriceUpdate) error
	ProductDelete(memberID int64, id int64) error
}