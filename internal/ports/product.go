package ports

import "github.com/DanielChachagua/GestionCar/pkg/models"

type ProductRepository interface {
	ProductGetByID(id string) (product *models.Product, err error)
	ProductGetByName(name string) (clients *[]models.Product, err error)
	ProductGetByIdentifier(name string) (clients *[]models.Product, err error)
	ProductGetAll() (products *[]models.Product, err error)
	ProductCreate(productCreate *models.ProductCreate) (id string, err error)
	ProductUpdate(productUpdate *models.ProductUpdate) (err error)
	AddToStock(stockUpdate *models.StockUpdate) (err error)
	SubtractFromStockToStock(stockUpdate *models.StockUpdate) (err error)
	UpdateStock(stockUpdate *models.StockUpdate) (err error)
	ProductDelete(id string) (err error)
}

type ProductService interface {
	ProductGetByID(id string) (product *models.Product, err error)
	ProductGetByName(name string) (clients *[]models.Product, err error)
	ProductGetByIdentifier(name string) (clients *[]models.Product, err error)
	ProductGetAll() (products *[]models.Product, err error)
	ProductCreate(productCreate *models.ProductCreate) (id string, err error)
	ProductUpdate(productUpdate *models.ProductUpdate) (err error)
	ProductUpdateStock(productUpdate *models.StockUpdate) (err error)
	ProductDelete(id string) (err error)
}
