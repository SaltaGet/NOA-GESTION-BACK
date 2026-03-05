package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (r *DepositRepository) DepositGetByID(id int64) (*models.Product, error) {
	ctx := context.Background()

	p, err := boilmodels.Products(
		boilmodels.ProductWhere.ID.EQ(id),
		boilmodels.ProductWhere.DeleteAt.IsNull(),
		qm.Load(boilmodels.ProductRels.Category),
		qm.Load(boilmodels.ProductRels.Deposits),
	).One(ctx, r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	return mapToModelProduct(p), nil
}

func (r *DepositRepository) DepositGetByCode(code string) (*models.Product, error) {
	ctx := context.Background()

	p, err := boilmodels.Products(
		boilmodels.ProductWhere.Code.EQ(code),
		boilmodels.ProductWhere.DeleteAt.IsNull(),
		qm.Load(boilmodels.ProductRels.Category),
		qm.Load(boilmodels.ProductRels.Deposits),
	).One(ctx, r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	return mapToModelProduct(p), nil
}

func (r *DepositRepository) DepositGetByName(name string) ([]*models.Product, error) {
	ctx := context.Background()

	searchStr := "%" + name + "%"
	boilProducts, err := boilmodels.Products(
		boilmodels.ProductWhere.Name.ILIKE(searchStr),
		boilmodels.ProductWhere.DeleteAt.IsNull(),
		qm.Load(boilmodels.ProductRels.Category),
		qm.Load(boilmodels.ProductRels.Deposits),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	allProducts := make([]*models.Product, 0, len(boilProducts))
	for _, bp := range boilProducts {
		allProducts = append(allProducts, mapToModelProduct(bp))
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

func (r *DepositRepository) DepositGetAll(page, limit int) ([]*models.Product, int64, error) {
	ctx := context.Background()

	total, err := boilmodels.Products(boilmodels.ProductWhere.DeleteAt.IsNull()).Count(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	boilProducts, err := boilmodels.Products(
		boilmodels.ProductWhere.DeleteAt.IsNull(),
		qm.Load(boilmodels.ProductRels.Category),
		qm.Load(boilmodels.ProductRels.Deposits),
		qm.Limit(limit),
		qm.Offset((page-1)*limit),
	).All(ctx, r.DB)

	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	products := make([]*models.Product, 0, len(boilProducts))
	for _, bp := range boilProducts {
		products = append(products, mapToModelProduct(bp))
	}

	return products, total, nil
}

func (r *DepositRepository) DepositUpdateStock(memberID int64, updateStock schemas.DepositUpdateStock) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	exists, err := boilmodels.Products(
		boilmodels.ProductWhere.ID.EQ(updateStock.ProductID),
	).Exists(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}
	if !exists {
		return schemas.HandlerErrorDB(sql.ErrNoRows, "Producto", schemas.Read)
	}

	deposit, err := boilmodels.Deposits(
		boilmodels.DepositWhere.ProductID.EQ(updateStock.ProductID),
	).One(ctx, tx)

	isNewEntry := false
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			deposit = &boilmodels.Deposit{
				ProductID: updateStock.ProductID,
				Stock:     0,
			}
			isNewEntry = true
		} else {
			return schemas.HandlerErrorDB(err, "Stock", schemas.Read)
		}
	}

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

	if isNewEntry {
		if err := deposit.Insert(ctx, tx, boil.Infer()); err != nil {
			return schemas.HandlerErrorDB(err, "Stock", schemas.Create)
		}
	} else {
		if _, err := deposit.Update(ctx, tx, boil.Whitelist(boilmodels.DepositColumns.Stock, boilmodels.DepositColumns.UpdatedAt)); err != nil {
			return schemas.HandlerErrorDB(err, "Stock", schemas.Update)
		}
	}

	return tx.Commit()
}
