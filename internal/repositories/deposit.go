package repositories

import (
	"errors"
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"gorm.io/gorm"
)

func (r *DepositRepository) DepositGetByID(id int64) (*models.Product, error) {
	var product models.Product
	if err := r.DB.Preload("Category").Preload("StockDeposit").Where("id = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "producto no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "error al obtener el producto", err)
	}
	return &product, nil
}

func (r *DepositRepository) DepositGetByCode(code string) (*models.Product, error) {
	var product models.Product
	if err := r.DB.Preload("Category").Preload("StockDeposit").Where("code = ?", code).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "producto no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "error al obtener el producto", err)
	}
	return &product, nil
}

func (r *DepositRepository) DepositGetByName(name string) ([]*models.Product, error) {
	var products []*models.Product
	if err := r.DB.Preload("Category").Preload("StockDeposit").Where("name LIKE ?", "%"+name+"%").Find(&products).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "error al obtener productos", err)
	}
	return products, nil
}

func (r *DepositRepository) DepositGetAll(page, limit int) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64
	if err := r.DB.Preload("Category").Preload("StockDeposit").Offset((page - 1) * limit).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "error al obtener productos", err)
	}
	if err := r.DB.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "error al contar productos", err)
	}
	return products, total, nil
}

func (r *DepositRepository) DepositUpdateStock(updateStock schemas.DepositUpdateStock) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var product models.Product
		if err := tx. Select("id").Where("id = ?", updateStock.ProductID).First(&product).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "producto no encontrado", err)
			}
			return schemas.ErrorResponse(500, "error al obtener el producto", err)
		}	

		var deposit models.Deposit

		if err := tx.Where("product_id = ?", updateStock.ProductID).FirstOrCreate(&deposit, &models.Deposit{ProductID: updateStock.ProductID}).Error; err != nil {
			return schemas.ErrorResponse(500, "error al actualizar el stock", err)
		}

		switch updateStock.Method {
		case "add":
			deposit.Stock += updateStock.Stock
		case "subtract":
			if deposit.Stock < updateStock.Stock {
				return schemas.ErrorResponse(400, "stock insuficiente", fmt.Errorf("stock insuficiente: %.2f", updateStock.Stock))
			}
			deposit.Stock -= updateStock.Stock
		case "set":
			deposit.Stock = updateStock.Stock
		default:
			return schemas.ErrorResponse(400, "metodo de actualizacion no valido", fmt.Errorf("metodo de actualizacion no valido"))
		}

		if err := tx.Save(&deposit).Error; err != nil {
			return schemas.ErrorResponse(500, "error al actualizar el stock", err)
		}

		return nil
	})
}
