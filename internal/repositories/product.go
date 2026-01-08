package repositories

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
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
			return db.Select("id", "product_id", "stock")
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
	var allProducts []*models.Product

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

// func (r *ProductRepository) ProductGetAll(page, limit int, isVisible bool) ([]*models.Product, int64, error) {
// 	var products []*models.Product
// 	var total int64
// 	offset := (page - 1) * limit

// 	if err := r.DB.
// 		Model(&models.Product{}).
// 		Count(&total).Error; err != nil {
// 		return nil, 0, schemas.ErrorResponse(500, "error al contar productos", err)
// 	}

// 	if err := r.DB.
// 		Preload("Category", func(db *gorm.DB) *gorm.DB {
// 			return db.Select("id", "name")
// 		}).
// 		Preload("StockPointSales", func(db *gorm.DB) *gorm.DB {
// 			return db.Select("product_id", "stock", "point_sale_id")
// 		}).
// 		Preload("StockPointSales.PointSale", func(db *gorm.DB) *gorm.DB {
// 			return db.Select("id", "name", "is_deposit")
// 		}).
// 		Preload("StockDeposit", func(db *gorm.DB) *gorm.DB {
// 			return db.Select("id", "product_id", "stock")
// 		}).
// 		Where("is_visible = ?", true).
// 		Offset(offset).
// 		Limit(limit).
// 		Find(&products).Error; err != nil {
// 		return nil, 0, schemas.ErrorResponse(500, "error al obtener productos", err)
// 	}

// 	return products, total, nil
// }
func (r *ProductRepository) ProductGetAll(page, limit int, isVisible *bool) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64
	offset := (page - 1) * limit

	// 1. Iniciamos la query base
	query := r.DB.Model(&models.Product{})

	// 2. Filtro condicional: Si isVisible es true, aplicamos el filtro.
	// Si es false, no entramos aquí y la query traerá visibles e invisibles.
	if isVisible != nil {
		query = query.Where("is_visible = ?", isVisible)
	}

	// 3. Contar el total basado en si se aplicó el filtro o no
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "error al contar productos", err)
	}

	// 4. Ejecutar la búsqueda con los Preloads
	if err := query.
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

func (r *ProductRepository) ProductGetByCodeToQR(code string) (*models.Product, error) {
	var product *models.Product

	if err := r.DB.
		Select("code", "name").
		Where("code = ?", code).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "producto no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "error al obtener el producto", err)
	}

	return product, nil
}

func (r *ProductRepository) ProductCount() (int64, error) {
	var count int64
	if err := r.DB.Model(&models.Product{}).Count(&count).Error; err != nil {
		return 0, schemas.ErrorResponse(500, "error al contar productos", err)
	}
	return count, nil
}

func (r *ProductRepository) ProductInsertToExcel(memberID int64, products []models.Product) ([]map[string]string, error) {
	rejected := make([]map[string]string, 0)
	for _, product := range products {
		err := r.DB.Transaction(func(tx *gorm.DB) error {
			if product.Category.Name != "" {
				if err := tx.FirstOrCreate(&product.Category, models.Category{Name: product.Category.Name}).Error; err != nil {
					return err
				}
				product.CategoryID = product.Category.ID
			} else {
				product.CategoryID = 1
			}

			if err := tx.Create(&product).Error; err != nil {
				return err
			}

			if product.StockDeposit.Stock <= 0 {
				return errors.New("stock no puede ser menor o igual a 0")
			}

			if err := tx.Create(&models.Deposit{ProductID: product.ID, Stock: product.StockDeposit.Stock}).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			rejected = append(rejected, map[string]string{"code": product.Code, "name": product.Name})
		}
	}

	return rejected, nil
}

func (r *ProductRepository) ProductCreate(memberID int64, productCreate *schemas.ProductCreate, plan *schemas.PlanResponseDTO) (int64, error) {
	var productSave models.Product

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var countTotal int64
		if err := tx.Model(&models.Product{}).Count(&countTotal).Error; err != nil {
			return schemas.ErrorResponse(500, "error al contar productos", err)
		}

		if countTotal >= plan.AmountProduct {
			return schemas.ErrorResponse(400, "el plan actual no permite crear más productos", fmt.Errorf("el plan actual no permite crear más productos"))
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

func (r *ProductRepository) ValidateProductImages(productValidateImage schemas.ProductValidateImage, plan *schemas.PlanResponseDTO) error {
	var product models.Product
	if err := r.DB.First(&product, productValidateImage.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "producto no encontrado", err)
		}
		return schemas.ErrorResponse(500, "error al obtener productos", err)
	}

	if productValidateImage.PrimaryImage == "keep" && product.PrimaryImage == nil {
		return schemas.ErrorResponse(400, "la imagen princial es obligatoria", nil)
	}

	var count int = 0
	if product.SecondaryImages != nil {
		count += len(strings.Split(*product.SecondaryImages, ","))
	}

	if len(productValidateImage.SecondaryImage.KeepUUIDs) > count || len(productValidateImage.SecondaryImage.RemoveUUIDs) > count {
		message := fmt.Sprintf("tienes %d imagenes secundarias, no puedes retener o eliminar mas de las que tienes", count)
		return schemas.ErrorResponse(400, message, fmt.Errorf("%s", message))
	}

	var sum = int(*productValidateImage.SecondaryImage.Add) - len(productValidateImage.SecondaryImage.RemoveUUIDs) + len(productValidateImage.SecondaryImage.KeepUUIDs)
	var typePrimary string = productValidateImage.PrimaryImage
	var existPrimary int = utils.Ternary(typePrimary == "set", 1, 0)

	for _, module := range plan.Modules {
		if module.Name == "ecommerce" {
			if (sum + existPrimary) <= int(module.AmountImagesPerProduct) {
				return nil
			} else {
				return schemas.ErrorResponse(400, "la cantidad máxima de imágenes por productos es de "+strconv.Itoa(int(plan.Modules[0].AmountImagesPerProduct))+"", fmt.Errorf("la cantidad máxima de imágenes por productos es de %d", int(plan.Modules[0].AmountImagesPerProduct)))
			}
		}
	}

	return schemas.ErrorResponse(400, "no existe módulo ecommerce para el tenant", errors.New("no existe módulo ecommerce para el tenant"))
}

func (r *ProductRepository) ProductUpdateVisibility(productUpdate *schemas.ListVisibilityUpdate) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		for _, prod := range productUpdate.ListProductVisibilityUpdate {
			// Ejecutamos el update
			result := tx.Model(&models.Product{}).
				Where("id = ?", prod.ProductID).
				Update("is_visible", prod.Visibility) // Asegúrate que el nombre de la columna sea correcto

			// 1. Verificar errores de base de datos (conexión, sintaxis, etc.)
			if result.Error != nil {
				return schemas.ErrorResponse(500, "error al editar el producto", result.Error)
			}

			// 2. Verificar si el producto realmente existía
			if result.RowsAffected == 0 {
				return schemas.ErrorResponse(404, fmt.Sprintf("producto con ID %d no encontrado", prod.ProductID), nil)
			}
		}

		return nil
	})
}