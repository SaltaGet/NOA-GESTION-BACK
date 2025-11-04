package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"


type ProductRepository interface {
	ProductGetByID(id string) (product *schemas.Product, err error)
	ProductGetByName(name string) (clients *[]schemas.Product, err error)
	ProductGetByIdentifier(name string) (clients *[]schemas.Product, err error)
	ProductGetAll() (products *[]schemas.Product, err error)
	ProductCreate(productCreate *schemas.ProductCreate) (id string, err error)
	ProductUpdate(productUpdate *schemas.ProductUpdate) (err error)
	AddToStock(stockUpdate *schemas.StockUpdate) (err error)
	SubtractFromStockToStock(stockUpdate *schemas.StockUpdate) (err error)
	UpdateStock(stockUpdate *schemas.StockUpdate) (err error)
	ProductDelete(id string) (err error)
}

type ProductService interface {
	ProductGetByID(id string) (product *schemas.Product, err error)
	ProductGetByName(name string) (clients *[]schemas.Product, err error)
	ProductGetByIdentifier(name string) (clients *[]schemas.Product, err error)
	ProductGetAll() (products *[]schemas.Product, err error)
	ProductCreate(productCreate *schemas.ProductCreate) (id string, err error)
	ProductUpdate(productUpdate *schemas.ProductUpdate) (err error)
	ProductUpdateStock(productUpdate *schemas.StockUpdate) (err error)
	ProductDelete(id string) (err error)
}
