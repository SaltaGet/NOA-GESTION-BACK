package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
)

func (r *DepositRepository) DepositGetByID(id int64) (*boilmodels.Product, error) {
	ctx := context.Background()

	p, err := boilmodels.Products(
		boilmodels.ProductWhere.ID.EQ(id),
		qm.Load(boilmodels.ProductRels.Category),
		qm.Load(boilmodels.ProductRels.Deposits),
	).One(ctx, r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	return p, nil
}

func (r *DepositRepository) DepositGetByCode(code string) (*boilmodels.Product, error) {
	ctx := context.Background()

	p, err := boilmodels.Products(
		boilmodels.ProductWhere.Code.EQ(code),
		qm.Load(boilmodels.ProductRels.Category),
		qm.Load(boilmodels.ProductRels.Deposits),
	).One(ctx, r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	return p, nil
}

func (r *DepositRepository) DepositGetByName(name string) ([]*boilmodels.Product, error) {
	ctx := context.Background()

	searchStr := "%" + name + "%"
	boilProducts, err := boilmodels.Products(
		boilmodels.ProductWhere.Name.ILIKE(searchStr),
		qm.Load(boilmodels.ProductRels.Category),
		qm.Load(boilmodels.ProductRels.Deposits),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	allProducts := make([]*boilmodels.Product, 0, len(boilProducts))
	for _, bp := range boilProducts {
		allProducts = append(allProducts, bp)
	}

	if strings.TrimSpace(name) == "" {
		if len(allProducts) > 10 {
			return allProducts[:10], nil
		}
		return allProducts, nil
	}

	scored := make([]schemas.ProductWithScore, 0)
	lowerSearch := strings.ToLower(strings.TrimSpace(name))

	for _, product := range allProducts {
		lowerName := strings.ToLower(product.Name)
		score := schemas.CalculateRelevance(lowerSearch, lowerName)

		if score > 0 {
			scored = append(scored, schemas.ProductWithScore{
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
	products := make([]*boilmodels.Product, 0, limit)
	for i, ps := range scored {
		if i >= limit {
			break
		}
		products = append(products, ps.Product)
	}

	return products, nil
}

func (r *DepositRepository) DepositGetAll(page, limit int) ([]*boilmodels.Product, int64, error) {
	ctx := context.Background()

	total, err := boilmodels.Products().Count(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	boilProducts, err := boilmodels.Products(
		qm.Load(boilmodels.ProductRels.Category),
		qm.Load(boilmodels.ProductRels.Deposits),
		qm.Limit(limit),
		qm.Offset((page-1)*limit),
	).All(ctx, r.DB)

	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
	}

	return boilProducts, total, nil
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
				Stock:     types.NewDecimal(decimal.New(0, 0).SetFloat64(0)),
			}
			isNewEntry = true
		} else {
			return schemas.HandlerErrorDB(err, "Stock", schemas.Read)
		}
	}

	stock := *updateStock.Stock
	switch updateStock.Method {
	case "add":
		nuevoStock := decimal.New(0, 0).SetFloat64(stock)
		deposit.Stock.Big = deposit.Stock.Big.Add(deposit.Stock.Big, nuevoStock)
	case "subtract":
		f, _ := deposit.Stock.Big.Float64()
		if f < stock {
			return schemas.ErrorResponse(400, "stock insuficiente", fmt.Errorf("stock insuficiente: %.2f", stock))
		}
		deposit.Stock.Big = deposit.Stock.Big.Sub(deposit.Stock.Big, decimal.New(0, 0).SetFloat64(stock))
	case "set":
		deposit.Stock.Big = decimal.New(0, 0).SetFloat64(stock)
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
