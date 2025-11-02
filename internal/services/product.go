package services

import (
	"fmt"

	"github.com/DanielChachagua/GestionCar/pkg/models"
)

func (p *ProductService) ProductGetByID(id string) (*models.Product, error) {
	product, err := p.ProductRepository.ProductGetByID(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductService) ProductGetByIdentifier(identifier string) (*[]models.Product, error) {
	product, err := p.ProductRepository.ProductGetByIdentifier(identifier)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductService) ProductGetAll() (*[]models.Product, error) {
	products, err := p.ProductRepository.ProductGetAll()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *ProductService) ProductGetByName(name string) (*[]models.Product, error) {
	products, err := p.ProductRepository.ProductGetByName(name)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *ProductService) ProductCreate(product *models.ProductCreate) (string, error) {
	id, err := p.ProductRepository.ProductCreate(product)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (p *ProductService) ProductUpdate(product *models.ProductUpdate) error {
	err := p.ProductRepository.ProductUpdate(product)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductService) ProductUpdateStock(stock *models.StockUpdate) error {
	product, err := p.ProductRepository.ProductGetByID(stock.ID)
	if err != nil {
		return err
	}
	switch stock.Method {
	case "update":
		if stock.Stock < 0 {
			return models.ErrorResponse(400, "El stock no puede ser negativo", fmt.Errorf("el stock no puede ser negativo"))
		}
		return p.ProductRepository.UpdateStock(stock)
	case "add":
		if stock.Stock <= 0{
			return models.ErrorResponse(400, "El stock debe ser mayor a 0", fmt.Errorf("el stock debe ser mayor a 0"))
		}
		return p.ProductRepository.AddToStock(stock)
	case "subtract":
		if stock.Stock <= 0{
			return models.ErrorResponse(400, "El stock debe ser mayor a 0", fmt.Errorf("el stock debe ser mayor a 0"))
		}
		if product != nil && product.Stock < stock.Stock {
			return models.ErrorResponse(400, "El stock no puede ser negativo", fmt.Errorf("el stock no puede ser negativo"))
		}
		return p.ProductRepository.SubtractFromStockToStock(stock)
	
	default:
		return models.ErrorResponse(500, "Método de actualización no soportado", err)
	}
}

func (p *ProductService) ProductDelete(id string) error {
	err := p.ProductRepository.ProductDelete(id)
	if err != nil {
		return err
	}
	return nil
}