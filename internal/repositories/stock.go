package repositories

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"strings"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/utils"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

func (r *StockRepository) getPointSale(ctx context.Context, pointID int64) (*boilmodels.PointSale, error) {
	pointSale, err := boilmodels.PointSales(
		qm.Select(
			boilmodels.PointSaleColumns.ID,
			boilmodels.PointSaleColumns.IsDeposit,
		),
		boilmodels.PointSaleWhere.ID.EQ(pointID),
	).One(ctx, r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Punto de Venta", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Punto de Venta", schemas.Read)
	}

	return pointSale, nil
}

func (r *StockRepository) buildStockQueryMods(pointSale *boilmodels.PointSale, pointID int64) []qm.QueryMod {
	queryMods := []qm.QueryMod{
		qm.Select(
			"products.id AS id",
			"products.code AS code",
			"products.name AS name",
			"products.description AS description",
			"products.price AS price",
			"s.stock AS stock",
			"categories.id AS category_id",
			"categories.name AS category_name",
			"products.primary_image AS primary_image",
			"products.secondary_images AS secondary_images",
		),
		qm.InnerJoin("categories ON categories.id = products.category_id"),
	}

	if pointSale.IsDeposit {
		queryMods = append(queryMods, qm.InnerJoin("deposits s ON s.product_id = products.id"))
	} else {
		queryMods = append(queryMods,
			qm.InnerJoin("stock_point_sales s ON s.product_id = products.id"),
			qm.Where("s.point_sale_id = ?", pointID),
		)
	}

	return queryMods
}

func (r *StockRepository) StockGetByID(id, pointID int64) (*schemas.ProductStockFullResponse, error) {
	ctx := context.Background()

	pointSale, err := r.getPointSale(ctx, pointID)
	if err != nil {
		return nil, err
	}

	queryMods := r.buildStockQueryMods(pointSale, pointID)
	queryMods = append(queryMods, qm.Where("products.id = ?", id))

	var product schemas.ProductStockFullResponseCategory
	if err := boilmodels.Products(queryMods...).Bind(ctx, r.DB, &product); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	secondaries := ""
	if product.SecondaryImages.Valid {
		secondaries = product.SecondaryImages.String
	}

	description := ""
	if product.Description.Valid {
		description = product.Description.String
	}

	var primaryImage *string
	if product.PrimaryImage.Valid {
		primaryImage = &product.PrimaryImage.String
	}

	return &schemas.ProductStockFullResponse{
		ID:             product.ID,
		Code:           product.Code,
		Name:           product.Name,
		Description:    description,
		Price:          product.Price,
		Stock:          product.Stock,
		PrimaryImage:   primaryImage,
		SecondaryImage: utils.SplitStrings(&secondaries),
		Category: schemas.CategoryResponseStock{
			ID:   product.CategoryID,
			Name: product.CategoryName,
		},
	}, nil
}

func (r *StockRepository) StockGetByCode(code string, pointID int64) (*schemas.ProductStockFullResponse, error) {
	ctx := context.Background()

	pointSale, err := r.getPointSale(ctx, pointID)
	if err != nil {
		return nil, err
	}

	queryMods := r.buildStockQueryMods(pointSale, pointID)
	queryMods = append(queryMods, qm.Where("products.code = ?", code))

	var product schemas.ProductStockFullResponseCategory
	if err := boilmodels.Products(queryMods...).Bind(ctx, r.DB, &product); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	secondaries := ""
	if product.SecondaryImages.Valid {
		secondaries = product.SecondaryImages.String
	}

	description := ""
	if product.Description.Valid {
		description = product.Description.String
	}

	var primaryImage *string
	if product.PrimaryImage.Valid {
		primaryImage = &product.PrimaryImage.String
	}

	return &schemas.ProductStockFullResponse{
		ID:             product.ID,
		Code:           product.Code,
		Name:           product.Name,
		Description:    description,
		Price:          product.Price,
		Stock:          product.Stock,
		PrimaryImage:   primaryImage,
		SecondaryImage: utils.SplitStrings(&secondaries),
		Category: schemas.CategoryResponseStock{
			ID:   product.CategoryID,
			Name: product.CategoryName,
		},
	}, nil
}

func (r *StockRepository) StockGetByCategoryID(categoryID, pointID int64) ([]*schemas.ProductStockFullResponse, error) {
	ctx := context.Background()

	pointSale, err := r.getPointSale(ctx, pointID)
	if err != nil {
		return nil, err
	}

	queryMods := r.buildStockQueryMods(pointSale, pointID)
	queryMods = append(queryMods, qm.Where("products.category_id = ?", categoryID))

	var products []*schemas.ProductStockFullResponseCategory
	if err := boilmodels.Products(queryMods...).Bind(ctx, r.DB, &products); err != nil {
		return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	var result []*schemas.ProductStockFullResponse
	for _, p := range products {
		secondaries := ""
		if p.SecondaryImages.Valid {
			secondaries = p.SecondaryImages.String
		}

		description := ""
		if p.Description.Valid {
			description = p.Description.String
		}

		var primaryImage *string
		if p.PrimaryImage.Valid {
			primaryImage = &p.PrimaryImage.String
		}

		result = append(result, &schemas.ProductStockFullResponse{
			ID:             p.ID,
			Code:           p.Code,
			Name:           p.Name,
			Description:    description,
			Price:          p.Price,
			Stock:          p.Stock,
			PrimaryImage:   primaryImage,
			SecondaryImage: utils.SplitStrings(&secondaries),
			Category: schemas.CategoryResponseStock{
				ID:   p.CategoryID,
				Name: p.CategoryName,
			},
		})
	}

	return result, nil
}

func calculateRelevance(search, name string) int {
	if search == "" {
		return 0
	}
	// Similar relevance algorithm
	idx := strings.Index(name, search)
	if idx == 0 {
		return 3
	} else if idx > 0 {
		return 2
	}

	// Check word boundaries
	words := strings.Fields(name)
	for _, w := range words {
		if strings.HasPrefix(w, search) {
			return 2
		}
	}

	return 0
}

func (r *StockRepository) StockGetByName(name string, pointID int64) ([]*schemas.ProductStockFullResponse, error) {
	ctx := context.Background()

	pointSale, err := r.getPointSale(ctx, pointID)
	if err != nil {
		return nil, err
	}

	queryMods := r.buildStockQueryMods(pointSale, pointID)

	var products []*schemas.ProductStockFullResponseCategory
	if err := boilmodels.Products(append(queryMods, qm.Where("products.name ILIKE ?", "%"+name+"%"))...).Bind(ctx, r.DB, &products); err != nil {
		return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	if strings.TrimSpace(name) == "" {
		result := make([]*schemas.ProductStockFullResponse, 0, 10)
		for i, p := range products {
			if i >= 10 {
				break
			}
			secondaries := ""
			if p.SecondaryImages.Valid {
				secondaries = p.SecondaryImages.String
			}

			description := ""
			if p.Description.Valid {
				description = p.Description.String
			}

			var primaryImage *string
			if p.PrimaryImage.Valid {
				primaryImage = &p.PrimaryImage.String
			}

			result = append(result, &schemas.ProductStockFullResponse{
				ID:             p.ID,
				Code:           p.Code,
				Name:           p.Name,
				Description:    description,
				Price:          p.Price,
				Stock:          p.Stock,
				PrimaryImage:   primaryImage,
				SecondaryImage: utils.SplitStrings(&secondaries),
				Category: schemas.CategoryResponseStock{
					ID:   p.CategoryID,
					Name: p.CategoryName,
				},
			})
		}
		return result, nil
	}

	scored := make([]schemas.ProductStockWithScore, 0)
	lowerSearch := strings.ToLower(strings.TrimSpace(name))

	for _, p := range products {
		lowerName := strings.ToLower(p.Name)
		score := calculateRelevance(lowerSearch, lowerName)

		if score > 0 {
			secondaries := ""
			if p.SecondaryImages.Valid {
				secondaries = p.SecondaryImages.String
			}

			description := ""
			if p.Description.Valid {
				description = p.Description.String
			}

			var primaryImage *string
			if p.PrimaryImage.Valid {
				primaryImage = &p.PrimaryImage.String
			}

			scored = append(scored, schemas.ProductStockWithScore{
				Product: &schemas.ProductStockFullResponse{
					ID:             p.ID,
					Code:           p.Code,
					Name:           p.Name,
					Description:    description,
					Price:          p.Price,
					Stock:          p.Stock,
					PrimaryImage:   primaryImage,
					SecondaryImage: utils.SplitStrings(&secondaries),
					Category: schemas.CategoryResponseStock{
						ID:   p.CategoryID,
						Name: p.CategoryName,
					},
				},
				Score:  float64(score),
				Length: len(p.Name),
			})
		}
	}

	sort.Slice(scored, func(i, j int) bool {
		if scored[i].Score != scored[j].Score {
			return scored[i].Score > scored[j].Score
		}
		return scored[i].Length < scored[j].Length
	})

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
	ctx := context.Background()
	offset := (page - 1) * limit

	pointSale, err := r.getPointSale(ctx, pointID)
	if err != nil {
		return nil, 0, err
	}

	queryMods := r.buildStockQueryMods(pointSale, pointID)

	// Since Count and Bind are two different things,
	// for Count we only need the joins avoiding Select
	countMods := []qm.QueryMod{
		qm.InnerJoin("categories ON categories.id = products.category_id"),
	}

	if pointSale.IsDeposit {
		countMods = append(countMods, qm.InnerJoin("deposits s ON s.product_id = products.id"))
	} else {
		countMods = append(countMods,
			qm.InnerJoin("stock_point_sales s ON s.product_id = products.id"),
			qm.Where("s.point_sale_id = ?", pointID),
		)
	}

	total, err := boilmodels.Products(countMods...).Count(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	queryMods = append(queryMods, qm.Offset(offset), qm.Limit(limit))

	var products []*schemas.ProductStockFullResponseCategory
	if err := boilmodels.Products(queryMods...).Bind(ctx, r.DB, &products); err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	var result []*schemas.ProductStockFullResponse
	for _, p := range products {
		secondaries := ""
		if p.SecondaryImages.Valid {
			secondaries = p.SecondaryImages.String
		}

		description := ""
		if p.Description.Valid {
			description = p.Description.String
		}

		var primaryImage *string
		if p.PrimaryImage.Valid {
			primaryImage = &p.PrimaryImage.String
		}

		result = append(result, &schemas.ProductStockFullResponse{
			ID:             p.ID,
			Code:           p.Code,
			Name:           p.Name,
			Description:    description,
			Price:          p.Price,
			Stock:          p.Stock,
			PrimaryImage:   primaryImage,
			SecondaryImage: utils.SplitStrings(&secondaries),
			Category: schemas.CategoryResponseStock{
				ID:   p.CategoryID,
				Name: p.CategoryName,
			},
		})
	}

	return result, total, nil
}
