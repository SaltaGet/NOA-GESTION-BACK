package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/utils"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
)

func mapToModelProduct(c *boilmodels.Product) *models.Product {
	if c == nil {
		return nil
	}

	p := &models.Product{
		ID:         c.ID,
		Code:       c.Code,
		Name:       c.Name,
		CategoryID: c.CategoryID,
		IsVisible:  c.IsVisible,
		Notifier:   c.Notifier,
	}

	if c.Description.Valid {
		p.Description = &c.Description.String
	}
	if c.PrimaryImage.Valid {
		p.PrimaryImage = &c.PrimaryImage.String
	}
	if c.SecondaryImages.Valid {
		p.SecondaryImages = &c.SecondaryImages.String
	}

	if !c.Price.IsZero() {
		val, _ := c.Price.Big.Float64()
		p.Price = val
	}
	if !c.MinAmount.IsZero() {
		val, _ := c.MinAmount.Big.Float64()
		p.MinAmount = val
	}

	if c.CreatedAt.Valid {
		p.CreatedAt = c.CreatedAt.Time
	}
	if c.UpdatedAt.Valid {
		p.UpdatedAt = c.UpdatedAt.Time
	}

	if c.R != nil {
		if c.R.Category != nil {
			p.Category = *mapToModelCategory(c.R.Category)
		}
		if len(c.R.Deposits) > 0 {
			dp := c.R.Deposits[0]
			stk, _ := dp.Stock.Big.Float64()
			p.StockDeposit = &models.Deposit{
				ID:        dp.ID,
				ProductID: dp.ProductID,
				Stock:     stk,
			}
		}
		if len(c.R.StockPointSales) > 0 {
			for _, sp := range c.R.StockPointSales {
				stk, _ := sp.Stock.Big.Float64()
				msps := &models.StockPointSale{
					ID:          sp.ID,
					ProductID:   sp.ProductID,
					PointSaleID: sp.PointSaleID,
					Stock:       stk,
				}
				if sp.R != nil && sp.R.PointSale != nil {
					msps.PointSale = models.PointSale{
						ID:        sp.R.PointSale.ID,
						Name:      sp.R.PointSale.Name,
						IsDeposit: sp.R.PointSale.IsDeposit,
					}
				}
				p.StockPointSales = append(p.StockPointSales, msps)
			}
		}
	}

	return p
}

func getProductPreloads() []qm.QueryMod {
	return []qm.QueryMod{
		qm.Load(boilmodels.ProductRels.Category),
		qm.Load(boilmodels.ProductRels.Deposits),
		qm.Load(boilmodels.ProductRels.StockPointSales),
		qm.Load(fmt.Sprintf("%s.%s", boilmodels.ProductRels.StockPointSales, boilmodels.StockPointSaleRels.PointSale)),
	}
}

func floatToDecimal(val float64) types.NullDecimal {
	d := types.NewNullDecimal(types.NewDecimal(fmt.Sprintf("%f", val)))
	return d
}

func (r *ProductRepository) ProductGetByID(id int64) (*models.Product, error) {
	ctx := context.Background()
	mods := getProductPreloads()
	mods = append(mods, boilmodels.ProductWhere.ID.EQ(id))

	product, err := boilmodels.Products(mods...).One(ctx, r.DB)
	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	return mapToModelProduct(product), nil
}

func (r *ProductRepository) ProductGetByCode(code string) (*models.Product, error) {
	ctx := context.Background()
	mods := getProductPreloads()
	mods = append(mods, boilmodels.ProductWhere.Code.EQ(code))

	product, err := boilmodels.Products(mods...).One(ctx, r.DB)
	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	return mapToModelProduct(product), nil
}

func (r *ProductRepository) ProductGetByCategoryID(categoryID int64) ([]*models.Product, error) {
	ctx := context.Background()
	mods := getProductPreloads()
	mods = append(mods, boilmodels.ProductWhere.CategoryID.EQ(categoryID))

	boilProducts, err := boilmodels.Products(mods...).All(ctx, r.DB)
	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	var products []*models.Product
	for _, p := range boilProducts {
		products = append(products, mapToModelProduct(p))
	}
	return products, nil
}

func (r *ProductRepository) ProductGetByName(name string) ([]*models.Product, error) {
	ctx := context.Background()
	mods := getProductPreloads()
	if name != "" {
		mods = append(mods, boilmodels.ProductWhere.Name.ILIKE("%"+name+"%"))
	}

	boilProducts, err := boilmodels.Products(mods...).All(ctx, r.DB)
	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	var allProducts []*models.Product
	for _, p := range boilProducts {
		allProducts = append(allProducts, mapToModelProduct(p))
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

	sort.Slice(scored, func(i, j int) bool {
		if scored[i].Score != scored[j].Score {
			return scored[i].Score > scored[j].Score
		}
		return scored[i].Length < scored[j].Length
	})

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

func (r *ProductRepository) ProductGetAll(page, limit int, isVisible *bool) ([]*models.Product, int64, error) {
	ctx := context.Background()

	mods := getProductPreloads()
	if isVisible != nil {
		mods = append(mods, boilmodels.ProductWhere.IsVisible.EQ(*isVisible))
	}

	count, err := boilmodels.Products(mods...).Count(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	offset := (page - 1) * limit
	mods = append(mods, qm.Offset(offset), qm.Limit(limit))

	boilProducts, err := boilmodels.Products(mods...).All(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	var products []*models.Product
	for _, p := range boilProducts {
		products = append(products, mapToModelProduct(p))
	}

	return products, count, nil
}

func (r *ProductRepository) ProductGetByCodeToQR(code string) (*models.Product, error) {
	ctx := context.Background()
	p, err := boilmodels.Products(
		qm.Select(boilmodels.ProductColumns.Code, boilmodels.ProductColumns.Name),
		boilmodels.ProductWhere.Code.EQ(code),
	).One(ctx, r.DB)
	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}
	return mapToModelProduct(p), nil
}

func (r *ProductRepository) ProductCount() (int64, error) {
	ctx := context.Background()
	count, err := boilmodels.Products().Count(ctx, r.DB)
	if err != nil {
		return 0, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}
	return count, nil
}

func (r *ProductRepository) ProductInsertToExcel(memberID int64, products []models.Product) ([]map[string]string, error) {
	ctx := context.Background()
	rejected := make([]map[string]string, 0)

	for _, product := range products {
		err := func() error {
			tx, err := r.DB.BeginTx(ctx, nil)
			if err != nil {
				return err
			}
			defer tx.Rollback()

			if product.Category.Name != "" {
				c, err := boilmodels.Categories(boilmodels.CategoryWhere.Name.EQ(strings.ToLower(product.Category.Name))).One(ctx, tx)
				if err != nil && !errors.Is(err, sql.ErrNoRows) {
					return err
				}
				if errors.Is(err, sql.ErrNoRows) {
					c = &boilmodels.Category{Name: strings.ToLower(product.Category.Name)}
					if err := c.Insert(ctx, tx, boil.Infer()); err != nil {
						return err
					}
				}
				product.CategoryID = c.ID
			} else {
				product.CategoryID = 1
			}

			p := boilmodels.Product{
				Name:        product.Name,
				Code:        product.Code,
				CategoryID:  product.CategoryID,
				Description: null.StringFromPtr(product.Description),
				Price:       floatToDecimal(product.Price),
				MinAmount:   floatToDecimal(product.MinAmount).Decimal,
			}

			if err := p.Insert(ctx, tx, boil.Infer()); err != nil {
				return schemas.HandlerErrorDB(err, "Producto", schemas.Create)
			}

			product.ID = p.ID

			if product.StockDeposit != nil && product.StockDeposit.Stock <= 0 {
				return schemas.ErrorResponse(400, "stock no puede ser menor o igual a 0", fmt.Errorf("stock no puede ser menor o igual a 0"))
			}

			stk := 0.0
			if product.StockDeposit != nil {
				stk = product.StockDeposit.Stock
			}
			dep := boilmodels.Deposit{
				ProductID: p.ID,
				Stock:     floatToDecimal(stk),
			}
			if err := dep.Insert(ctx, tx, boil.Infer()); err != nil {
				return schemas.HandlerErrorDB(err, "Producto", schemas.Create)
			}

			return tx.Commit()
		}()

		if err != nil {
			rejected = append(rejected, map[string]string{"code": product.Code, "name": product.Name})
		}
	}

	return rejected, nil
}

func (r *ProductRepository) ProductCreate(memberID int64, productCreate *schemas.ProductCreate, plan *schemas.PlanResponseDTO) (int64, error) {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return 0, err
	}

	countTotal, err := boilmodels.Products().Count(ctx, tx)
	if err != nil {
		return 0, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	if countTotal >= plan.AmountProduct {
		return 0, schemas.ErrorResponse(400, "el plan actual no permite crear más productos", fmt.Errorf("el plan actual no permite crear más productos"))
	}

	exists, err := boilmodels.Categories(boilmodels.CategoryWhere.ID.EQ(productCreate.CategoryID)).Exists(ctx, tx)
	if err != nil || !exists {
		if !exists {
			err = sql.ErrNoRows
		}
		return 0, schemas.HandlerErrorDB(err, "Categoria", schemas.Read)
	}

	p := boilmodels.Product{
		Name:        productCreate.Name,
		Code:        productCreate.Code,
		Description: null.StringFromPtr(productCreate.Description),
		CategoryID:  productCreate.CategoryID,
		Notifier:    productCreate.Notifier,
		MinAmount:   floatToDecimal(productCreate.MinAmount).Decimal,
	}

	if productCreate.Price != nil {
		p.Price = floatToDecimal(*productCreate.Price)
	}

	if err := p.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Producto", schemas.Create)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return p.ID, nil
}

func (r *ProductRepository) ProductUpdate(memberID int64, product *schemas.ProductUpdate) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	p, err := boilmodels.FindProduct(ctx, tx, product.ID)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	if product.Price != nil {
		p.Price = floatToDecimal(*product.Price)
	}

	p.Code = product.Code
	p.Name = product.Name
	p.Description = null.StringFromPtr(&product.Description)
	p.CategoryID = product.CategoryID
	p.Notifier = product.Notifier
	p.MinAmount = floatToDecimal(product.MinAmount).Decimal

	if _, err := p.Update(ctx, tx, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Producto", schemas.Update)
	}

	return tx.Commit()
}

func (r *ProductRepository) ProductPriceUpdate(memberID int64, product *schemas.ListPriceUpdate) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	for _, pup := range product.ListProductPriceUpdate {
		p, err := boilmodels.FindProduct(ctx, tx, pup.ID)
		if err != nil {
			return schemas.HandlerErrorDB(err, "Producto", schemas.Read)
		}

		p.Price = floatToDecimal(pup.Price)
		rowsAff, err := p.Update(ctx, tx, boil.Whitelist(boilmodels.ProductColumns.Price, boilmodels.ProductColumns.UpdatedAt))
		if err != nil {
			return schemas.HandlerErrorDB(err, "Producto", schemas.Update)
		}

		if rowsAff == 0 {
			return schemas.ErrorResponse(404, fmt.Sprintf("producto %d no encontrado", pup.ID), fmt.Errorf("producto %d no encontrado", pup.ID))
		}
	}

	return tx.Commit()
}

func (r *ProductRepository) ProductDelete(memberID int64, id int64) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	p, err := boilmodels.FindProduct(ctx, tx, id)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	_, err = boilmodels.StockPointSales(boilmodels.StockPointSaleWhere.ProductID.EQ(id)).DeleteAll(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Producto", schemas.Delete)
	}

	_, err = p.Delete(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Producto", schemas.Delete)
	}

	return tx.Commit()
}

func (r *ProductRepository) ValidateProductImages(productValidateImage schemas.ProductValidateImage, plan *schemas.PlanResponseDTO) error {
	ctx := context.Background()
	p, err := boilmodels.FindProduct(ctx, r.DB, productValidateImage.ProductID)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	if productValidateImage.PrimaryImage == "keep" && (!p.PrimaryImage.Valid || p.PrimaryImage.String == "") {
		return schemas.ErrorResponse(400, "la imagen princial es obligatoria", fmt.Errorf("la imagen princial es obligatoria"))
	}

	count := 0
	if p.SecondaryImages.Valid && p.SecondaryImages.String != "" {
		count += len(strings.Split(p.SecondaryImages.String, ","))
	}

	if len(productValidateImage.SecondaryImage.KeepUUIDs) > count || len(productValidateImage.SecondaryImage.RemoveUUIDs) > count {
		message := fmt.Sprintf("tienes %d imagenes secundarias, no puedes retener o eliminar mas de las que tienes", count)
		return schemas.ErrorResponse(400, message, fmt.Errorf("%s", message))
	}

	sum := int(*productValidateImage.SecondaryImage.Add) - len(productValidateImage.SecondaryImage.RemoveUUIDs) + len(productValidateImage.SecondaryImage.KeepUUIDs)
	typePrimary := productValidateImage.PrimaryImage
	existPrimary := utils.Ternary(typePrimary == "set", 1, 0)

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
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, prod := range productUpdate.ListProductVisibilityUpdate {
		p, err := boilmodels.FindProduct(ctx, tx, prod.ProductID)
		if err != nil {
			return schemas.HandlerErrorDB(err, "Producto", schemas.Update)
		}

		p.IsVisible = prod.Visibility
		rowsAff, err := p.Update(ctx, tx, boil.Whitelist(boilmodels.ProductColumns.IsVisible, boilmodels.ProductColumns.UpdatedAt))

		if err != nil {
			return schemas.HandlerErrorDB(err, "Producto", schemas.Update)
		}

		if rowsAff == 0 {
			return schemas.ErrorResponse(404, fmt.Sprintf("producto con ID %d no encontrado", prod.ProductID), fmt.Errorf("producto con ID %d no encontrado", prod.ProductID))
		}
	}

	return tx.Commit()
}
