package repositories

import (
	"errors"
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"gorm.io/gorm"
)

func (r *ProductRepository) ProductGetByID(id int64) (*models.Product, error) {
	var product models.Product
	if err := r.DB.
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("StockPointSales", func(db *gorm.DB) *gorm.DB {
			return db.Select("product_id", "stock", "point_sale_id")
		}).
		Preload("StockPointSales.PointSale", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "is_deposit")
		}).
		Preload("StockDeposit", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "product_id", "stock")
		}).
		First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "producto no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "error al obtener el producto", err)
	}

	return &product, nil
}

func (r *ProductRepository) ProductGetByCode(code string) (*models.Product, error) {
	var product *models.Product

	if err := r.DB.
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("StockPointSales", func(db *gorm.DB) *gorm.DB {
			return db.Select("product_id", "stock", "point_sale_id")
		}).
		Preload("StockPointSales.PointSale", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "is_deposit")
		}).
		Preload("StockDeposit", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "product_id", "stock", "")
		}).
		Where("code = ?", code).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "producto no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "error al obtener el producto", err)
	}

	return product, nil
}

func (r *ProductRepository) ProductGetByCategoryID(categoryID int64) ([]*models.Product, error) {
	var products []*models.Product

	if err := r.DB.
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("StockPointSales", func(db *gorm.DB) *gorm.DB {
			return db.Select("product_id", "stock", "point_sale_id")
		}).
		Preload("StockPointSales.PointSale", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "is_deposit")
		}).
		Preload("StockDeposit", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "product_id", "stock")
		}).
		Where("category_id = ?", categoryID).Find(&products).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "error al obtener productos", err)
	}

	return products, nil
}

func (r *ProductRepository) ProductGetByName(name string) ([]*models.Product, error) {
	var products []*models.Product

	if err := r.DB.
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("StockPointSales", func(db *gorm.DB) *gorm.DB {
			return db.Select("product_id", "stock", "point_sale_id")
		}).
		Preload("StockPointSales.PointSale", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "is_deposit")
		}).
		Preload("StockDeposit", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "product_id", "stock")
		}).
		Where("name LIKE ?", "%"+name+"%").Find(&products).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "error al obtener productos", err)
	}

	return products, nil
}

func (r *ProductRepository) ProductGetAll(page, limit int) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64
	offset := (page - 1) * limit

	if err := r.DB.
		Model(&models.Product{}).
		Count(&total).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "error al contar productos", err)
	}

	if err := r.DB.
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("StockPointSales", func(db *gorm.DB) *gorm.DB {
			return db.Select("product_id", "stock", "point_sale_id")
		}).
		Preload("StockPointSales.PointSale", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "is_deposit")
		}).
		Preload("StockDeposit", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "product_id", "stock")
		}).
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "error al obtener productos", err)
	}

	return products, total, nil
}

// func (r *ProductRepository) ProductCreate(productCreate *schemas.ProductCreate, plan *schemas.PlanResponseDTO) (int64, error) {
// 	var countTotal int64
// 	if err := r.DB.Model(&models.Product{}).Count(&countTotal).Error; err != nil {
// 		return 0, schemas.ErrorResponse(500, "error al contar productos", err)
// 	}

// 	if countTotal >= plan.AmountProduct {
// 		return 0, schemas.ErrorResponse(400, "el plan actual no permite crear más productos", nil)
// 	}

// 	var product models.Product
// 	var category models.Category
// 	if err := r.DB.First(&category, productCreate.CategoryID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return 0, schemas.ErrorResponse(404, "categoria no encontrada", err)
// 		}
// 		return 0, schemas.ErrorResponse(500, "error al obtener la categoria", err)
// 	}

// 	product.Name = productCreate.Name
// 	product.Code = productCreate.Code
// 	product.Description = productCreate.Description
// 	if productCreate.Price != nil {
// 		product.Price = *productCreate.Price
// 	}
// 	product.CategoryID = productCreate.CategoryID
// 	product.Notifier = productCreate.Notifier
// 	product.MinAmount = productCreate.MinAmount

// 	if err := r.DB.Create(&product).Error; err != nil {
// 		if schemas.IsDuplicateError(err) {
// 			return 0, schemas.ErrorResponse(400, "el producto de codigo "+product.Code+" ya existe", err)
// 		}
// 		return 0, schemas.ErrorResponse(500, "error al crear el producto", err)
// 	}

// 	return product.ID, nil
// }

// func (r *ProductRepository) ProductUpdate(product *schemas.ProductUpdate) error {
// 	var p models.Product
// 	if err := r.DB.First(&p, product.ID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return schemas.ErrorResponse(404, "producto no encontrado", err)
// 		}
// 		return schemas.ErrorResponse(500, "error al obtener el producto", err)
// 	}

// 	if product.Price != nil {
// 		p.Price = *product.Price
// 	}

// 	updates := map[string]any{
// 		"code":        product.Code,
// 		"name":        product.Name,
// 		"description": &product.Description,
// 		"category_id": product.CategoryID,
// 		"price":       p.Price,
// 		"notifier":    product.Notifier,
// 		"min_amount":  product.MinAmount,
// 	}

// 	if err := r.DB.Model(&p).Updates(updates).Error; err != nil {
// 		if schemas.IsDuplicateError(err) {
// 			return schemas.ErrorResponse(400, "el producto de código "+product.Code+" ya existe", err)
// 		}
// 		return schemas.ErrorResponse(500, "error al actualizar el producto", err)
// 	}

// 	return nil
// }

// func (r *ProductRepository) ProductPriceUpdate(product *schemas.ListPriceUpdate) error {
// 	return r.DB.Transaction(func(tx *gorm.DB) error {
// 		for _, p := range product.ListProductPriceUpdate {
// 			res := tx.Model(&models.Product{}).
// 				Where("id = ?", p.ID).
// 				Update("price", p.Price)

// 			if res.Error != nil {
// 				return schemas.ErrorResponse(500, "error al actualizar el producto", res.Error)
// 			}

// 			if res.RowsAffected == 0 {
// 				return schemas.ErrorResponse(404, fmt.Sprintf("producto %d no encontrado", p.ID), fmt.Errorf("producto %d no encontrado", p.ID))
// 			}
// 		}
// 		return nil
// 	})
// }

// func (r *ProductRepository) ProductDelete(id int64) error {
// 	return r.DB.Transaction(func(tx *gorm.DB) error {
// 		if err := tx.Where("product_id = ?", id).Delete(&models.StockPointSale{}).Error; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return schemas.ErrorResponse(404, "producto no encontrado", err)
// 			}
// 			return schemas.ErrorResponse(500, "error al eliminar el producto", err)
// 		}

// 		if err := tx.Where("id = ?", id).Delete(&models.Product{}).Error; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return schemas.ErrorResponse(404, "producto no encontrado", err)
// 			}
// 			return schemas.ErrorResponse(500, "error al eliminar el producto", err)
// 		}
// 		return nil
// 	})
// }

// ProductCreate crea un nuevo producto con auditoría
func (r *ProductRepository) ProductCreate(memberID int64, productCreate *schemas.ProductCreate, plan *schemas.PlanResponseDTO) (int64, error) {
	var productSave models.Product

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var countTotal int64
		if err := tx.Model(&models.Product{}).Count(&countTotal).Error; err != nil {
			return schemas.ErrorResponse(500, "error al contar productos", err)
		}

		if countTotal >= plan.AmountProduct {
			return schemas.ErrorResponse(400, "el plan actual no permite crear más productos", nil)
		}

		var category models.Category
		if err := tx.First(&category, productCreate.CategoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "categoria no encontrada", err)
			}
			return schemas.ErrorResponse(500, "error al obtener la categoria", err)
		}

		product := models.Product{
			Name:        productCreate.Name,
			Code:        productCreate.Code,
			Description: productCreate.Description,
			CategoryID:  productCreate.CategoryID,
			Notifier:    productCreate.Notifier,
			MinAmount:   productCreate.MinAmount,
		}

		if productCreate.Price != nil {
			product.Price = *productCreate.Price
		}

		if err := tx.Create(&product).Error; err != nil {
			if schemas.IsDuplicateError(err) {
				return schemas.ErrorResponse(400, "el producto de codigo "+product.Code+" ya existe", err)
			}
			return schemas.ErrorResponse(500, "error al crear el producto", err)
		}

		productSave = product

		// Guardar auditoría

		return nil
	})

	if err != nil {
		return 0, err
	}

	go database.SaveAuditAsync(r.DB, models.AuditLog{
		MemberID: memberID,
		Method:   "create",
		Path:     "product",
	}, nil, productSave)

	return productSave.ID, nil
}

// ProductUpdate actualiza un producto con auditoría
func (r *ProductRepository) ProductUpdate(memberID int64, product *schemas.ProductUpdate) error {
	var saveProduct, updatedProduct models.Product
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var p models.Product
		if err := tx.First(&p, product.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "producto no encontrado", err)
			}
			return schemas.ErrorResponse(500, "error al obtener el producto", err)
		}

		// Guardar estado anterior
		saveProduct = p

		if product.Price != nil {
			p.Price = *product.Price
		}

		updates := map[string]any{
			"code":        product.Code,
			"name":        product.Name,
			"description": &product.Description,
			"category_id": product.CategoryID,
			"price":       p.Price,
			"notifier":    product.Notifier,
			"min_amount":  product.MinAmount,
		}

		if err := tx.Model(&p).Updates(updates).Error; err != nil {
			if schemas.IsDuplicateError(err) {
				return schemas.ErrorResponse(400, "el producto de código "+product.Code+" ya existe", err)
			}
			return schemas.ErrorResponse(500, "error al actualizar el producto", err)
		}

		tx.First(&updatedProduct, p.ID)
		return nil
	})

	if err == nil {
		// Guardar auditoría
		go database.SaveAuditAsync(r.DB, models.AuditLog{
			MemberID: memberID,
			Method:   "update",
			Path:     "product",
		}, saveProduct, updatedProduct)
	}

	return err
}

// ProductPriceUpdate actualiza los precios de múltiples productos con auditoría
func (r *ProductRepository) ProductPriceUpdate(memberID int64, product *schemas.ListPriceUpdate) error {
	var saveProduct, updatedProduct models.Product
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		for _, p := range product.ListProductPriceUpdate {
			// Obtener el estado anterior del producto
			var oldProduct models.Product
			if err := tx.First(&oldProduct, p.ID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return schemas.ErrorResponse(404, fmt.Sprintf("producto %d no encontrado", p.ID), err)
				}
				return schemas.ErrorResponse(500, fmt.Sprintf("error al obtener el producto %d", p.ID), err)
			}

			// Guardar estado anterior
			saveProduct = oldProduct

			res := tx.Model(&models.Product{}).
				Where("id = ?", p.ID).
				Update("price", p.Price)

			if res.Error != nil {
				return schemas.ErrorResponse(500, "error al actualizar el producto", res.Error)
			}

			if res.RowsAffected == 0 {
				return schemas.ErrorResponse(404, fmt.Sprintf("producto %d no encontrado", p.ID), fmt.Errorf("producto %d no encontrado", p.ID))
			}

			tx.First(&updatedProduct, p.ID)
		}
		return nil
	})

	if err == nil {
		go database.SaveAuditAsync(r.DB, models.AuditLog{
			MemberID: memberID,
			Method:   "update",
			Path:     "product",
		}, saveProduct, updatedProduct)
	}

	return err
}

// ProductDelete elimina un producto con auditoría
func (r *ProductRepository) ProductDelete(memberID int64, id int64) error {
	var productSave models.Product
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Obtener el producto antes de eliminarlo
		var product models.Product
		if err := tx.First(&product, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "producto no encontrado", err)
			}
			return schemas.ErrorResponse(500, "error al obtener el producto", err)
		}

		// Guardar estado del producto

		productSave = product

		// Obtener los stocks del punto de venta asociados
		var stockPointSales []models.StockPointSale
		if err := tx.Where("product_id = ?", id).Find(&stockPointSales).Error; err != nil {
			return schemas.ErrorResponse(500, "error al obtener stocks del producto", err)
		}

		// Eliminar stocks del punto de venta
		if err := tx.Where("product_id = ?", id).Delete(&models.StockPointSale{}).Error; err != nil {
			return schemas.ErrorResponse(500, "error al eliminar stocks del producto", err)
		}

		// Eliminar el producto
		if err := tx.Where("id = ?", id).Delete(&models.Product{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "producto no encontrado", err)
			}
			return schemas.ErrorResponse(500, "error al eliminar el producto", err)
		}

		return nil
	})

	if err == nil {
		// Guardar auditoría del producto eliminado
		go database.SaveAuditAsync(r.DB, models.AuditLog{
			MemberID: memberID,
			Method:   "delete",
			Path:     "product",
		}, productSave, nil)
	}

	return err
}
