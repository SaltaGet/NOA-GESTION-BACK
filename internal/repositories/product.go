package repositories

import (
	"errors"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *ProductRepository) ProductGetByID(id string) (*schemas.Product, error) {
	var product schemas.Product
	if err := r.DB.Where("id = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Producto no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al buscar producto", err)
	}
	return &product, nil

}

func (r *ProductRepository) ProductGetByIdentifier(identifier string) (*[]schemas.Product, error) {
	var product []schemas.Product
	if err := r.DB.Where("identifier LIKE ?", "%"+identifier+"%").Find(&product).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al buscar producto", err)
	}
	return &product, nil
}

func (r *ProductRepository) ProductGetByName(name string) (*[]schemas.Product, error) {
	var products []schemas.Product
	if err := r.DB.Where("name LIKE ?", "%"+name+"%").Find(&products).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al buscar productos", err)
	}
	return &products, nil
}

func (r *ProductRepository) ProductGetAll() (*[]schemas.Product, error) {
	var products []schemas.Product
	if err := r.DB.Find(&products).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al buscar productos", err)
	}
	return &products, nil

}

func (r *ProductRepository) ProductCreate(element *schemas.ProductCreate) (string, error) {
	newID := uuid.NewString()
	if err := r.DB.Create(&schemas.Product{
		ID:         newID,
		Identifier: element.Identifier,
		Name:       element.Name,
		Stock:      0,
	}).Error; err != nil {
		return "", schemas.ErrorResponse(500, "Error interno al crear producto", err)
	}
	return newID, nil

}

func (r *ProductRepository) ProductUpdate(element *schemas.ProductUpdate) error {
	if err := r.DB.Model(&schemas.Product{}).Where("id = ?", element.ID).Updates(&schemas.Product{
		Identifier: element.Identifier,
		Name:       element.Name,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "Producto no encontrado", err)
		}
		return schemas.ErrorResponse(500, "Error interno al actualizar producto", err)
	}
	return nil

}

func (r *ProductRepository) UpdateStock(stockUpdate *schemas.StockUpdate) error {
	if err := r.DB.Model(&schemas.Product{}).
		Where("id = ?", stockUpdate.ID).
		Update("stock", stockUpdate.Stock).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "Producto no encontrado", err)
		}
		return schemas.ErrorResponse(500, "Error interno al actualizar stock", err)
	}
	return nil

}

func (r *ProductRepository) AddToStock(stockUpdate *schemas.StockUpdate) error {
	if err := r.DB.Model(&schemas.Product{}).
		Where("id = ?", stockUpdate.ID).
		UpdateColumn("stock", gorm.Expr("stock + ?", stockUpdate.Stock)).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Producto no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error interno al actualizar stock", err)
		}

	return nil
}

func (r *ProductRepository) SubtractFromStockToStock(stockUpdate *schemas.StockUpdate) error {
	if err := r.DB.Model(&schemas.Product{}).
		Where("id = ?", stockUpdate.ID).
		UpdateColumn("stock", gorm.Expr("stock - ?", stockUpdate.Stock)).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Producto no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error interno al actualizar stock", err)
		}

	return nil
}

func (r *ProductRepository) ProductDelete(id string) error {
	if err := r.DB.Where("id = ?", id).Delete(&schemas.Product{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "Producto no encontrado", err)
			
		}
		return schemas.ErrorResponse(500, "Error interno al eliminar producto", err)
	}
	return nil
}
