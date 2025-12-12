package repositories

import (
	"errors"
	"sort"
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func (r *StockRepository) StockGetByID(id, pointID int64) (*schemas.ProductStockFullResponse, error) {
	var pointSale models.PointSale
	if err := r.DB.Select("id", "is_deposit").First(&pointSale, pointID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "punto de venta no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "error al obtener el punto de venta", err)
	}

	query := r.DB.Model(&models.Product{}).Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	})

	if pointSale.IsDeposit {
		query = query.
			Preload("StockDeposit", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "product_id", "stock")
			})
	} else {
		query = query.
			Preload("StockPointSales", func(db *gorm.DB) *gorm.DB {
				return db.Select("product_id", "stock").Where("point_sale_id = ?", pointID)
			})
	}

	var product models.Product
	if err := query.Where("id = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "producto no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "error al obtener el producto", err)
	}

	var productSchema schemas.ProductStockFullResponse
	_ = copier.Copy(&productSchema, &product)

	if pointSale.IsDeposit {
		productSchema.Stock = product.StockDeposit.Stock
	} else {
		if len(product.StockPointSales) > 0 {
			productSchema.Stock = product.StockPointSales[0].Stock
		} else {
			productSchema.Stock = 0
		}
	}

	return &productSchema, nil
}

func (r *StockRepository) StockGetByCode(code string, pointID int64) (*schemas.ProductStockFullResponse, error) {
	var pointSale models.PointSale
	if err := r.DB.Select("id", "is_deposit").First(&pointSale, pointID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "punto de venta no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "error al obtener el punto de venta", err)
	}

	query := r.DB.Model(&models.Product{}).Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	})

	if pointSale.IsDeposit {
		query = query.
			Preload("StockDeposit", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "product_id", "stock")
			})
	} else {
		query = query.
			Preload("StockPointSales", func(db *gorm.DB) *gorm.DB {
				return db.Select("product_id", "stock").Where("point_sale_id = ?", pointID)
			})
	}

	var product models.Product
	if err := query.Where("code = ?", code).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "producto no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "error al obtener el producto", err)
	}

	var productSchema schemas.ProductStockFullResponse
	_ = copier.Copy(&productSchema, &product)

	if pointSale.IsDeposit {
		productSchema.Stock = product.StockDeposit.Stock
	} else {
		if len(product.StockPointSales) > 0 {
			productSchema.Stock = product.StockPointSales[0].Stock
		} else {
			productSchema.Stock = 0
		}
	}

	return &productSchema, nil
}

func (r *StockRepository) StockGetByCategoryID(categoryID, pointID int64) ([]*schemas.ProductStockFullResponse, error) {
	var pointSale models.PointSale
	if err := r.DB.Select("id", "is_deposit").First(&pointSale, pointID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "punto de venta no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "error al obtener el punto de venta", err)
	}

	var products []*schemas.ProductStockFullResponseCategory

	baseSelect := []string{
		"products.id",
		"products.code",
		"products.name",
		"products.description",
		"products.price",
		"s.stock",
		"categories.id AS category_id",
		"categories.name AS category_name",
	}

	query := r.DB.Model(&models.Product{}).
		Select(baseSelect).
		Joins("INNER JOIN categories ON categories.id = products.category_id")

	if pointSale.IsDeposit {
		query = query.Joins("INNER JOIN deposits s ON s.product_id = products.id")
	} else {
		query = query.Joins("INNER JOIN stock_point_sales s ON s.product_id = products.id AND s.point_sale_id = ?", pointID)
	}

	if err := query.Where("products.category_id = ?", categoryID).Scan(&products).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "error al obtener productos", err)
	}

	var result []*schemas.ProductStockFullResponse
	for _, p := range products {
		result = append(result, &schemas.ProductStockFullResponse{
			ID:           p.ID,
			Code:         p.Code,
			Name:         p.Name,
			Description:  p.Description,
			Price:        p.Price,
			Stock:        p.Stock,
			Category: schemas.CategoryResponseStock{
				ID:   p.CategoryID,
				Name: p.CategoryName,
			},
		})
	}

	return result, nil
}

// func (r *StockRepository) StockGetByName(name string, pointID int64) ([]*schemas.ProductStockFullResponse, error) {
// 	var pointSale models.PointSale
// 	if err := r.DB.Select("id", "is_deposit").First(&pointSale, pointID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, schemas.ErrorResponse(404, "punto de venta no encontrado", err)
// 		}
// 		return nil, schemas.ErrorResponse(500, "error al obtener el punto de venta", err)
// 	}

// 	var products []*schemas.ProductStockFullResponseCategory

// 	baseSelect := []string{
// 		"products.id",
// 		"products.code",
// 		"products.name",
// 		"products.description",
// 		"products.price",
// 		"s.stock",
// 		"categories.id AS category_id",
// 		"categories.name AS category_name",
// 	}

// 	query := r.DB.Model(&models.Product{}).
// 		Select(baseSelect).
// 		Joins("INNER JOIN categories ON categories.id = products.category_id")

// 	if pointSale.IsDeposit {
// 		query = query.Joins("INNER JOIN deposits s ON s.product_id = products.id")
// 	} else {
// 		query = query.Joins("INNER JOIN stock_point_sales s ON s.product_id = products.id AND s.point_sale_id = ?", pointID)
// 	}

// 	if err := query.Limit(10).Where("products.name LIKE ?", "%"+name+"%").Scan(&products).Error; err != nil {
// 		return nil, schemas.ErrorResponse(500, "error al obtener productos", err)
// 	}

// 	var result []*schemas.ProductStockFullResponse
// 	for _, p := range products {
// 		result = append(result, &schemas.ProductStockFullResponse{
// 			ID:           p.ID,
// 			Code:         p.Code,
// 			Name:         p.Name,
// 			Description:  p.Description,
// 			Price:        p.Price,
// 			Stock:        p.Stock,
// 			Category: schemas.CategoryResponseStock{
// 				ID:   p.CategoryID,
// 				Name: p.CategoryName,
// 			},
// 		})
// 	}

// 	return result, nil
// }

func (r *StockRepository) StockGetByName(name string, pointID int64) ([]*schemas.ProductStockFullResponse, error) {
    var pointSale models.PointSale
    if err := r.DB.Select("id", "is_deposit").First(&pointSale, pointID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, schemas.ErrorResponse(404, "punto de venta no encontrado", err)
        }
        return nil, schemas.ErrorResponse(500, "error al obtener el punto de venta", err)
    }

    var products []*schemas.ProductStockFullResponseCategory

    baseSelect := []string{
        "products.id",
        "products.code",
        "products.name",
        "products.description",
        "products.price",
        "s.stock",
        "categories.id AS category_id",
        "categories.name AS category_name",
    }

    query := r.DB.Model(&models.Product{}).
        Select(baseSelect).
        Joins("INNER JOIN categories ON categories.id = products.category_id")

    if pointSale.IsDeposit {
        query = query.Joins("INNER JOIN deposits s ON s.product_id = products.id")
    } else {
        query = query.Joins("INNER JOIN stock_point_sales s ON s.product_id = products.id AND s.point_sale_id = ?", pointID)
    }

    // Traer más resultados para luego filtrar y ordenar
    if err := query.Where("products.name LIKE ?", "%"+name+"%").Scan(&products).Error; err != nil {
        return nil, schemas.ErrorResponse(500, "error al obtener productos", err)
    }

    // Si no hay búsqueda, retornar los primeros 10
    if strings.TrimSpace(name) == "" {
        result := make([]*schemas.ProductStockFullResponse, 0, 10)
        for i, p := range products {
            if i >= 10 {
                break
            }
            result = append(result, &schemas.ProductStockFullResponse{
                ID:          p.ID,
                Code:        p.Code,
                Name:        p.Name,
                Description: p.Description,
                Price:       p.Price,
                Stock:       p.Stock,
                Category: schemas.CategoryResponseStock{
                    ID:   p.CategoryID,
                    Name: p.CategoryName,
                },
            })
        }
        return result, nil
    }

    // Calcular relevancia para cada producto
    scored := make([]schemas.ProductStockWithScore, 0)
    lowerSearch := strings.ToLower(strings.TrimSpace(name))

    for _, p := range products {
        lowerName := strings.ToLower(p.Name)
        score := models.CalculateRelevance(lowerSearch, lowerName)

        if score > 0 {
            scored = append(scored, schemas.ProductStockWithScore{
                Product: &schemas.ProductStockFullResponse{
                    ID:          p.ID,
                    Code:        p.Code,
                    Name:        p.Name,
                    Description: p.Description,
                    Price:       p.Price,
                    Stock:       p.Stock,
                    Category: schemas.CategoryResponseStock{
                        ID:   p.CategoryID,
                        Name: p.CategoryName,
                    },
                },
                Score:  score,
                Length: len(p.Name),
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
    result := make([]*schemas.ProductStockFullResponse, 0, limit)
    for i, ps := range scored {
        if i >= limit {
            break
        }
        result = append(result, ps.Product)
    }

    return result, nil
}

func (r *StockRepository) StockGetAll(page, limit int, pointID int64) ([]*schemas.ProductStockFullResponse, int64, error) {
	var total int64
	offset := (page - 1) * limit

	// verificar punto de venta
	var pointSale models.PointSale
	if err := r.DB.Select("id", "is_deposit").First(&pointSale, pointID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, schemas.ErrorResponse(404, "punto de venta no encontrado", err)
		}
		return nil, 0, schemas.ErrorResponse(500, "error al obtener el punto de venta", err)
	}

	var products []*schemas.ProductStockFullResponseCategory

	baseSelect := []string{
		"products.id",
		"products.code",
		"products.name",
		"products.description",
		"products.price",
		"s.stock",
		"categories.id AS category_id",
		"categories.name AS category_name",
	}

	// QUERY BASE (sin limit ni offset)
	baseQuery := r.DB.Model(&models.Product{}).
		Joins("INNER JOIN categories ON categories.id = products.category_id")

	// join según el tipo de stock
	if pointSale.IsDeposit {
		baseQuery = baseQuery.Joins("INNER JOIN deposits s ON s.product_id = products.id")
	} else {
		baseQuery = baseQuery.Joins("INNER JOIN stock_point_sales s ON s.product_id = products.id AND s.point_sale_id = ?", pointID)
	}

	// COUNT REAL (MISMO JOIN)
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "error al contar productos", err)
	}

	// QUERY FINAL (paginada)
	query := baseQuery.
		Select(baseSelect).
		Offset(offset).
		Limit(limit)

	if err := query.Scan(&products).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "error al obtener productos", err)
	}

	var result []*schemas.ProductStockFullResponse
	for _, p := range products {
		result = append(result, &schemas.ProductStockFullResponse{
			ID:           p.ID,
			Code:         p.Code,
			Name:         p.Name,
			Description:  p.Description,
			Price:        p.Price,
			Stock:        p.Stock,
			Category: schemas.CategoryResponseStock{
				ID:   p.CategoryID,
				Name: p.CategoryName,
			},
		})
	}

	return result, total, nil
}

