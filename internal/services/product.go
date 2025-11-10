package services

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (p *ProductService) ProductGetByID(id int64) (*schemas.ProductFullResponse, error) {
	return p.ProductRepository.ProductGetByID(id)
}

func (p *ProductService) ProductGetByCode(code string) (*schemas.ProductFullResponse, error) {
	return p.ProductRepository.ProductGetByCode(code)
}

func (p *ProductService) ProductGetByName(name string) ([]*schemas.ProductFullResponse, error) {
	return p.ProductRepository.ProductGetByName(name)
}

func (p *ProductService) ProductGetByCategoryID(categoryID int64) ([]*schemas.ProductFullResponse, error) {
	return p.ProductRepository.ProductGetByCategoryID(categoryID)
}

func (p *ProductService) ProductGetAll(page, limit int) ([]*schemas.ProductFullResponse, int64, error) {
	return p.ProductRepository.ProductGetAll(page, limit)
}

func (p *ProductService) ProductCreate(productCreate *schemas.ProductCreate) (int64, error) {
	return p.ProductRepository.ProductCreate(productCreate)
}

func (p *ProductService) ProductUpdate(productUpdate *schemas.ProductUpdate) error {
	return p.ProductRepository.ProductUpdate(productUpdate)
}
func (p *ProductService) ProductPriceUpdate(productUpdate *schemas.ListPriceUpdate) error {
	return p.ProductRepository.ProductPriceUpdate(productUpdate)
}

func (p *ProductService) ProductDelete(id int64) error {
	return p.ProductRepository.ProductDelete(id)

}
