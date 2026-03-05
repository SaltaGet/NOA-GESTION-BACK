package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
)

func mapToMovementStockResponse(m *boilmodels.MovementStock) *schemas.MovementStockResponse {
	if m == nil {
		return nil
	}

	amt, _ := m.Amount.Big.Float64()

	res := &schemas.MovementStockResponse{
		ID:          m.ID,
		Amount:      amt,
		FromID:      m.FromID.Int64,
		FromType:    m.FromType.String,
		ToID:        m.ToID.Int64,
		ToType:      m.ToType.String,
		IgnoreStock: m.IgnoreStock,
	}

	if m.CreatedAt.Valid {
		res.CreatedAt = m.CreatedAt.Time
	}

	if m.R != nil {
		if m.R.Member != nil {
			member := mapToMemberSimpleDTO(m.R.Member)
			if member != nil {
				res.Member = *member
			}
		}

		if m.R.Product != nil {
			pPrice, _ := m.R.Product.Price.Big.Float64()
			res.Product = schemas.ProductSimpleResponseDTO{
				ID:    m.R.Product.ID,
				Code:  m.R.Product.Code,
				Name:  m.R.Product.Name,
				Price: pPrice,
			}
		}
	}

	return res
}

func mapToMovementStockResponseDTO(m *boilmodels.MovementStock) schemas.MovementStockResponseDTO {
	if m == nil {
		return schemas.MovementStockResponseDTO{}
	}

	amt, _ := m.Amount.Big.Float64()

	res := schemas.MovementStockResponseDTO{
		ID:          m.ID,
		Amount:      amt,
		FromID:      m.FromID.Int64,
		FromType:    m.FromType.String,
		ToID:        m.ToID.Int64,
		ToType:      m.ToType.String,
		IgnoreStock: m.IgnoreStock,
	}

	if m.CreatedAt.Valid {
		res.CreatedAt = m.CreatedAt.Time
	}

	if m.R != nil {
		if m.R.Member != nil {
			member := mapToMemberSimpleDTO(m.R.Member)
			if member != nil {
				res.Member = *member
			}
		}

		if m.R.Product != nil {
			pPrice, _ := m.R.Product.Price.Big.Float64()
			res.Product = schemas.ProductSimpleResponseDTO{
				ID:    m.R.Product.ID,
				Code:  m.R.Product.Code,
				Name:  m.R.Product.Name,
				Price: pPrice,
			}
		}
	}

	return res
}

func (r *MovementStockRepository) MovementStockGetByID(id int64) (*schemas.MovementStockResponse, error) {
	ctx := context.Background()

	movement, err := boilmodels.MovementStocks(
		boilmodels.MovementStockWhere.ID.EQ(id),
		qm.Load(boilmodels.MovementStockRels.Member),
		qm.Load(boilmodels.MovementStockRels.Product),
		qm.Load(qm.Rels(boilmodels.MovementStockRels.Product, boilmodels.ProductRels.Category)),
	).One(ctx, r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Movimiento de stock", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Movimiento de stock", schemas.Read)
	}

	return mapToMovementStockResponse(movement), nil
}

func (r *MovementStockRepository) MovementStockGetByDate(page, limit int, fromDate, toDate time.Time) ([]schemas.MovementStockResponseDTO, int64, error) {
	ctx := context.Background()
	offset := (page - 1) * limit

	var qms []qm.QueryMod
	// If the original GORM query just listed all but we passed fromDate, it seems fromDate/toDate were ignored in GORM?
	// Checking previous code: `r.DB.Offset(offset).Limit(limit).Order("created_at desc").Find(&movements)`
	// Indeed, fromDate, toDate were ignored. But we'll apply them if they are valid.
	if !fromDate.IsZero() && !toDate.IsZero() {
		qms = append(qms,
			boilmodels.MovementStockWhere.CreatedAt.GTE(null.TimeFrom(fromDate)),
			boilmodels.MovementStockWhere.CreatedAt.LTE(null.TimeFrom(toDate)),
		)
	}

	total, err := boilmodels.MovementStocks(qms...).Count(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Movimiento de stock", schemas.Read)
	}

	qms = append(qms,
		qm.Load(boilmodels.MovementStockRels.Member),
		qm.Load(boilmodels.MovementStockRels.Product),
		qm.Offset(offset),
		qm.Limit(limit),
		qm.OrderBy("created_at DESC"),
	)

	movements, err := boilmodels.MovementStocks(qms...).All(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Movimiento de stock", schemas.Read)
	}

	var results []schemas.MovementStockResponseDTO
	for _, m := range movements {
		results = append(results, mapToMovementStockResponseDTO(m))
	}

	return results, total, nil
}

func (r *MovementStockRepository) MoveStockList(memberID int64, input []schemas.MovementStockList) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	if len(input) == 0 {
		return schemas.ErrorResponse(400, "no hay movimientos para procesar", fmt.Errorf("lista vacía"))
	}

	for _, movementList := range input {
		product, err := boilmodels.Products(boilmodels.ProductWhere.ID.EQ(movementList.ProductID)).One(ctx, tx)
		if err != nil {
			return schemas.HandlerErrorDB(err, "Producto", schemas.Read)
		}

		prodPrice, _ := product.Price.Big.Float64()
		if prodPrice <= 0.0 {
			return schemas.ErrorResponse(400, fmt.Sprintf("no se puede editar el producto %d sin precio", movementList.ProductID),
				fmt.Errorf("no se puede editar un producto sin precio"))
		}

		if len(movementList.MovementStockItem) == 0 {
			return schemas.ErrorResponse(400, fmt.Sprintf("no hay movimientos para el producto %d", movementList.ProductID),
				fmt.Errorf("lista de movimientos vacía"))
		}

		for idx, item := range movementList.MovementStockItem {
			if err := r.processSingleMovement(ctx, tx, memberID, movementList.ProductID, &item, idx); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *MovementStockRepository) processSingleMovement(
	ctx context.Context,
	tx *sql.Tx,
	memberID int64,
	productID int64,
	item *schemas.MovementStockItem,
	index int,
) error {
	var fromID, toID int64

	ignoreStock := item.IgnoreStock != nil && *item.IgnoreStock

	switch item.FromType {
	case "deposit":
		fromID = 100

		deposit, err := boilmodels.Deposits(boilmodels.DepositWhere.ProductID.EQ(productID)).One(ctx, tx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				deposit = &boilmodels.Deposit{ProductID: productID, Stock: 0}
				if err := deposit.Insert(ctx, tx, boil.Infer()); err != nil {
					return schemas.HandlerErrorDB(err, "Deposito", schemas.Create)
				}
			} else {
				return schemas.HandlerErrorDB(err, "Deposito", schemas.Read)
			}
		}

		if !ignoreStock {
			if deposit.Stock < item.Amount {
				return schemas.ErrorResponse(400,
					fmt.Sprintf("stock insuficiente en depósito (%d)", index+1),
					fmt.Errorf("actual %.2f < requerido %.2f", deposit.Stock, item.Amount))
			}
		}

		deposit.Stock -= item.Amount
		if _, err := deposit.Update(ctx, tx, boil.Whitelist(boilmodels.DepositColumns.Stock, boilmodels.DepositColumns.UpdatedAt)); err != nil {
			return schemas.HandlerErrorDB(err, "Deposito", schemas.Update)
		}

	case "point_sale":
		ps, err := boilmodels.PointSales(qm.Select(boilmodels.PointSaleColumns.ID, boilmodels.PointSaleColumns.IsDeposit), boilmodels.PointSaleWhere.ID.EQ(item.FromID)).One(ctx, tx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return schemas.ErrorResponse(404, fmt.Sprintf("punto de venta %d no encontrado (%d)", item.FromID, index+1), err)
			}
			return schemas.HandlerErrorDB(err, "Punto de venta", schemas.Read)
		}

		if ps.IsDeposit {
			return schemas.ErrorResponse(400,
				fmt.Sprintf("no se puede usar un punto de venta depósito como origen (%d)", index+1),
				fmt.Errorf("point_sale es depósito"))
		}

		fromID = item.FromID

		oldPS, err := boilmodels.StockPointSales(
			boilmodels.StockPointSaleWhere.ProductID.EQ(productID),
			boilmodels.StockPointSaleWhere.PointSaleID.EQ(item.FromID),
		).One(ctx, tx)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				oldPS = &boilmodels.StockPointSale{ProductID: productID, PointSaleID: item.FromID, Stock: 0}
				if err := oldPS.Insert(ctx, tx, boil.Infer()); err != nil {
					return schemas.HandlerErrorDB(err, "Stock Punto de venta", schemas.Create)
				}
			} else {
				return schemas.HandlerErrorDB(err, "Stock Punto de venta", schemas.Read)
			}
		}

		if !ignoreStock {
			if oldPS.Stock < item.Amount {
				return schemas.ErrorResponse(400,
					fmt.Sprintf("stock insuficiente en point_sale origen (%d)", index+1),
					fmt.Errorf("actual %.2f < requerido %.2f", oldPS.Stock, item.Amount))
			}
		}

		oldPS.Stock -= item.Amount
		if _, err := oldPS.Update(ctx, tx, boil.Whitelist(boilmodels.StockPointSaleColumns.Stock, boilmodels.StockPointSaleColumns.UpdatedAt)); err != nil {
			return schemas.HandlerErrorDB(err, "Stock Punto de venta", schemas.Update)
		}

	default:
		return schemas.ErrorResponse(400,
			fmt.Sprintf("tipo de origen inválido (%d)", index+1),
			fmt.Errorf("FromType='%s'", item.FromType))
	}

	switch item.ToType {
	case "deposit":
		toID = 100

		oldDeposit, err := boilmodels.Deposits(boilmodels.DepositWhere.ProductID.EQ(productID)).One(ctx, tx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				oldDeposit = &boilmodels.Deposit{ProductID: productID, Stock: 0}
				if err := oldDeposit.Insert(ctx, tx, boil.Infer()); err != nil {
					return schemas.HandlerErrorDB(err, "Deposito", schemas.Create)
				}
			} else {
				return schemas.HandlerErrorDB(err, "Deposito", schemas.Read)
			}
		}

		oldDeposit.Stock += item.Amount
		if _, err := oldDeposit.Update(ctx, tx, boil.Whitelist(boilmodels.DepositColumns.Stock, boilmodels.DepositColumns.UpdatedAt)); err != nil {
			return schemas.HandlerErrorDB(err, "Deposito", schemas.Update)
		}

	case "point_sale":
		ps, err := boilmodels.PointSales(qm.Select(boilmodels.PointSaleColumns.ID, boilmodels.PointSaleColumns.IsDeposit), boilmodels.PointSaleWhere.ID.EQ(item.ToID)).One(ctx, tx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return schemas.ErrorResponse(404, fmt.Sprintf("punto de venta %d no encontrado (%d)", item.ToID, index+1), err)
			}
			return schemas.HandlerErrorDB(err, "Punto de venta", schemas.Read)
		}

		if ps.IsDeposit {
			return schemas.ErrorResponse(400,
				fmt.Sprintf("no se puede usar un punto de venta depósito como destino (%d)", index+1),
				fmt.Errorf("point_sale es depósito"))
		}

		toID = item.ToID

		oldPS, err := boilmodels.StockPointSales(
			boilmodels.StockPointSaleWhere.ProductID.EQ(productID),
			boilmodels.StockPointSaleWhere.PointSaleID.EQ(item.ToID),
		).One(ctx, tx)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				oldPS = &boilmodels.StockPointSale{ProductID: productID, PointSaleID: item.ToID, Stock: 0}
				if err := oldPS.Insert(ctx, tx, boil.Infer()); err != nil {
					return schemas.HandlerErrorDB(err, "Stock Punto de venta", schemas.Create)
				}
			} else {
				return schemas.HandlerErrorDB(err, "Stock Punto de venta", schemas.Read)
			}
		}

		oldPS.Stock += item.Amount
		if _, err := oldPS.Update(ctx, tx, boil.Whitelist(boilmodels.StockPointSaleColumns.Stock, boilmodels.StockPointSaleColumns.UpdatedAt)); err != nil {
			return schemas.HandlerErrorDB(err, "Stock Punto de venta", schemas.Update)
		}

	default:
		return schemas.ErrorResponse(400,
			fmt.Sprintf("tipo de destino inválido (%d)", index+1),
			fmt.Errorf("ToType='%s'", item.ToType))
	}

	amtDec := types.NewNullDecimal(types.NewDecimal(fmt.Sprintf("%.4f", item.Amount)))

	movement := &boilmodels.MovementStock{
		MemberID:    null.Int64From(memberID),
		ProductID:   null.Int64From(productID),
		Amount:      amtDec,
		FromID:      null.Int64From(fromID),
		FromType:    null.StringFrom(item.FromType),
		ToID:        null.Int64From(toID),
		ToType:      null.StringFrom(item.ToType),
		IgnoreStock: ignoreStock,
	}

	if err := movement.Insert(ctx, tx, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Movimiento", schemas.Create)
	}

	return nil
}
