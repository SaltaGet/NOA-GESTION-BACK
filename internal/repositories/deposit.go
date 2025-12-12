package repositories

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
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
	var allProducts []*models.Product

	if err := r.DB.
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("StockDeposit").
		Where("name LIKE ?", "%"+name+"%").
		Find(&allProducts).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "error al obtener productos", err)
	}

	if strings.TrimSpace(name) == "" {
		if len(allProducts) > 10 {
			return allProducts[:10], nil
		}
		return allProducts, nil
	}

	scored := make([]models.ProductWithScore, 0)
	lowerSearch := strings.ToLower(strings.TrimSpace(name))

	for _, product := range allProducts {
		lowerName := strings.ToLower(product.Name)
		score := models.CalculateRelevance(lowerSearch, lowerName)

		if score > 0 {
			scored = append(scored, models.ProductWithScore{
				Product: product,
				Score:   score,
				Length:  len(product.Name),
			})
		}
	}

	// Ordenar según los criterios especificados
	sort.Slice(scored, func(i, j int) bool {
		// Si los scores son diferentes, ordenar por score (descendente)
		if scored[i].Score != scored[j].Score {
			return scored[i].Score > scored[j].Score
		}
		// Si los scores son iguales, ordenar por longitud (ascendente - más corto primero)
		return scored[i].Length < scored[j].Length
	})

	// Limitar a 10 resultados
	limit := 10
	products := make([]*models.Product, 0, limit)
	for i, ps := range scored {
		if i >= limit {
			break
		}
		products = append(products, ps.Product)
	}

	return products, nil
}

// func (r *DepositRepository) DepositGetByName(name string) ([]*models.Product, error) {
// 	var products []*models.Product
// 	if err := r.DB.Preload("Category").Preload("StockDeposit").Where("name LIKE ?", "%"+name+"%").Find(&products).Error; err != nil {
// 		return nil, schemas.ErrorResponse(500, "error al obtener productos", err)
// 	}
// 	return products, nil
// }

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

// func (r *DepositRepository) DepositUpdateStock(updateStock schemas.DepositUpdateStock) error {
// 	return r.DB.Transaction(func(tx *gorm.DB) error {
// 		var product models.Product
// 		if err := tx. Select("id").Where("id = ?", updateStock.ProductID).First(&product).Error; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return schemas.ErrorResponse(404, "producto no encontrado", err)
// 			}
// 			return schemas.ErrorResponse(500, "error al obtener el producto", err)
// 		}

// 		var deposit models.Deposit

// 		if err := tx.Where("product_id = ?", updateStock.ProductID).FirstOrCreate(&deposit, &models.Deposit{ProductID: updateStock.ProductID}).Error; err != nil {
// 			return schemas.ErrorResponse(500, "error al actualizar el stock", err)
// 		}

// 		stock := *updateStock.Stock
// 		switch updateStock.Method {
// 		case "add":
// 			deposit.Stock += stock
// 		case "subtract":
// 			if deposit.Stock < stock{
// 				return schemas.ErrorResponse(400, "stock insuficiente", fmt.Errorf("stock insuficiente: %.2f", stock))
// 			}
// 			deposit.Stock -= stock
// 		case "set":
// 			deposit.Stock = stock
// 		default:
// 			return schemas.ErrorResponse(400, "metodo de actualizacion no valido", fmt.Errorf("metodo de actualizacion no valido"))
// 		}

// 		if err := tx.Save(&deposit).Error; err != nil {
// 			return schemas.ErrorResponse(500, "error al actualizar el stock", err)
// 		}

// 		return nil
// 	})
// }

func (r *DepositRepository) DepositUpdateStock(memberID int64, updateStock schemas.DepositUpdateStock) error {
	var saveDeposit any
	var finalDeposit any
	var isNewRecord bool
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var product models.Product
		if err := tx.Select("id").Where("id = ?", updateStock.ProductID).First(&product).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "producto no encontrado", err)
			}
			return schemas.ErrorResponse(500, "error al obtener el producto", err)
		}

		var deposit models.Deposit
		// Verificar si el registro existe
		if err := tx.Where("product_id = ?", updateStock.ProductID).First(&deposit).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Si no existe, crear uno nuevo
				deposit = models.Deposit{ProductID: updateStock.ProductID, Stock: 0}
				isNewRecord = true
			} else {
				return schemas.ErrorResponse(500, "error al obtener el stock", err)
			}
		}

		// Guardar estado anterior para auditoría
		saveDeposit = deposit

		stock := *updateStock.Stock
		switch updateStock.Method {
		case "add":
			deposit.Stock += stock
		case "subtract":
			if deposit.Stock < stock {
				return schemas.ErrorResponse(400, "stock insuficiente", fmt.Errorf("stock insuficiente: %.2f", stock))
			}
			deposit.Stock -= stock
		case "set":
			deposit.Stock = stock
		default:
			return schemas.ErrorResponse(400, "metodo de actualizacion no valido", fmt.Errorf("metodo de actualizacion no valido"))
		}

		if err := tx.Save(&deposit).Error; err != nil {
			return schemas.ErrorResponse(500, "error al actualizar el stock", err)
		}

		// Guardar auditoría

		finalDeposit = deposit

		return nil
	})

	if err == nil {
		method := "update"
		path := "deposit"

		if isNewRecord {
			method = "create"
			path = "deposit"
		}

		go database.SaveAuditAsync(r.DB, models.AuditLog{
			MemberID: memberID,
			Method:   method,
			Path:     path,
		}, saveDeposit, finalDeposit)
	}

	return err
}
