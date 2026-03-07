package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
	"github.com/aarondl/null/v8"
)

func mapToExpenseBuyResponse(e *boilmodels.ExpenseBuy) *schemas.ExpenseBuyResponse {
	if e == nil {
		return nil
	}

	subtotal, _ := e.Subtotal.Big.Float64()
	discount, _ := e.Discount.Big.Float64()
	total, _ := e.Total.Big.Float64()

	res := &schemas.ExpenseBuyResponse{
		ID:       e.ID,
		Subtotal: subtotal,
		Discount: discount,
		Total:    total,
	}

	res.TypeDiscount = e.TypeDiscount
	res.Details = &e.Details.String
	res.CreatedAt = e.CreatedAt

	if e.R != nil {
		if e.R.Member != nil {
			res.Member = schemas.MemberSimpleDTO{
				ID:        e.R.Member.ID,
				Username:  e.R.Member.Username,
				FirstName: e.R.Member.FirstName,
				LastName:  e.R.Member.LastName,
			}
		}

		if e.R.Supplier != nil {
			res.Supplier.ID = e.R.Supplier.ID
			res.Supplier.Name = e.R.Supplier.Name
			res.Supplier.CompanyName = e.R.Supplier.CompanyName
		}

		for _, item := range e.R.ExpenseBuyItems {
			iAmt, _ := item.Amount.Big.Float64()
			iPrice, _ := item.Price.Big.Float64()
			iDiscount, _ := item.Discount.Big.Float64()
			iSubtotal, _ := item.Subtotal.Big.Float64()
			iTotal, _ := item.Total.Big.Float64()

			ier := schemas.ExpenseBuyItemResponse{
				ID:       item.ID,
				Amount:   iAmt,
				Price:    iPrice,
				Discount: iDiscount,
				Subtotal: iSubtotal,
				Total:    iTotal,
			}
			ier.TypeDiscount = item.TypeDiscount
			ier.CreatedAt = item.CreatedAt

			if item.R != nil && item.R.Product != nil {
				pPrice, _ := item.R.Product.Price.Big.Float64()
				ier.Product = schemas.ProductSimpleResponseDTO{
					ID:    item.R.Product.ID,
					Code:  item.R.Product.Code,
					Name:  item.R.Product.Name,
					Price: pPrice,
				}
			}
			res.ExpenseBuyItem = append(res.ExpenseBuyItem, ier)
		}

		for _, pay := range e.R.PayExpenseBuys {
			pTotal, _ := pay.Total.Big.Float64()
			pr := schemas.PayExpenseBuyResponse{
				ID:        pay.ID,
				Total:     pTotal,
				MethodPay: pay.MethodPay,
			}
			res.PayExpenseBuy = append(res.PayExpenseBuy, pr)
		}
	}

	return res
}

func mapToExpenseBuyResponseSimple(e *boilmodels.ExpenseBuy) *schemas.ExpenseBuyResponseSimple {
	if e == nil {
		return nil
	}

	subtotal, _ := e.Subtotal.Big.Float64()
	discount, _ := e.Discount.Big.Float64()
	total, _ := e.Total.Big.Float64()

	res := &schemas.ExpenseBuyResponseSimple{
		ID:       e.ID,
		Subtotal: subtotal,
		Discount: discount,
		Total:    total,
	}

	res.TypeDiscount = e.TypeDiscount
	res.Description = &e.Details.String
	res.CreatedAt = e.CreatedAt

	if e.R != nil {
		if e.R.Supplier != nil {
			res.Supplier = schemas.SupplierResponseDTO{
				ID:          e.R.Supplier.ID,
				Name:        e.R.Supplier.Name,
				CompanyName: e.R.Supplier.CompanyName,
			}
		}

		for _, pay := range e.R.PayExpenseBuys {
			pTotal, _ := pay.Total.Big.Float64()
			pr := schemas.PayExpenseBuyResponse{
				ID:        pay.ID,
				Total:     pTotal,
				MethodPay: pay.MethodPay,
			}
			res.PayExpenseBuy = append(res.PayExpenseBuy, pr)
		}
	}

	return res
}

func (r *ExpenseBuyRepository) ExpenseBuyGetByID(id int64) (*schemas.ExpenseBuyResponse, error) {
	ctx := context.Background()

	e, err := boilmodels.ExpenseBuys(
		boilmodels.ExpenseBuyWhere.ID.EQ(id),
		qm.Load(boilmodels.ExpenseBuyRels.Member),
		qm.Load(boilmodels.ExpenseBuyRels.Supplier),
		qm.Load(boilmodels.ExpenseBuyRels.ExpenseBuyItems),
		qm.Load(qm.Rels(boilmodels.ExpenseBuyRels.ExpenseBuyItems, boilmodels.ExpenseBuyItemRels.Product)),
		qm.Load(boilmodels.ExpenseBuyRels.PayExpenseBuys),
	).One(ctx, r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Egreso de compras", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Egreso de compras", schemas.Read)
	}

	return mapToExpenseBuyResponse(e), nil
}

func (r *ExpenseBuyRepository) ExpenseBuyGetByDate(fromDate, toDate time.Time, page, limit int) ([]*schemas.ExpenseBuyResponseSimple, int64, error) {
	ctx := context.Background()
	offset := (page - 1) * limit

	qms := []qm.QueryMod{
		boilmodels.ExpenseBuyWhere.CreatedAt.GTE(fromDate),
		boilmodels.ExpenseBuyWhere.CreatedAt.LTE(toDate),
	}

	total, err := boilmodels.ExpenseBuys(qms...).Count(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Egreso de compras", schemas.Read)
	}

	qms = append(qms,
		qm.Load(boilmodels.ExpenseBuyRels.Supplier),
		qm.Load(boilmodels.ExpenseBuyRels.PayExpenseBuys),
		qm.OrderBy("created_at DESC"),
		qm.Offset(offset),
		qm.Limit(limit),
	)

	boilExpenses, err := boilmodels.ExpenseBuys(qms...).All(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Egreso de compras", schemas.Read)
	}

	var results []*schemas.ExpenseBuyResponseSimple
	for _, e := range boilExpenses {
		results = append(results, mapToExpenseBuyResponseSimple(e))
	}

	return results, total, nil
}

func (r *ExpenseBuyRepository) ExpenseBuyCreate(memberID int64, expenseBuyCreate *schemas.ExpenseBuyCreate) (int64, error) {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return 0, err
	}

	supplierExists, err := boilmodels.Suppliers(
		boilmodels.SupplierWhere.ID.EQ(expenseBuyCreate.SupplierID),
	).Exists(ctx, tx)

	if err != nil || !supplierExists {
		return 0, schemas.HandlerErrorDB(sql.ErrNoRows, "Proveedor", schemas.Read)
	}

	var newItems []*boilmodels.ExpenseBuyItem
	total := 0.0

	for _, item := range expenseBuyCreate.ExpenseBuyItem {
		if item.Amount <= 0 {
			return 0, schemas.ErrorResponse(400, fmt.Sprintf("La cantidad para el producto %d no es válida", item.ProductID), fmt.Errorf("la cantidad para el producto %d no es válida", item.ProductID))
		}

		productExists, err := boilmodels.Products(boilmodels.ProductWhere.ID.EQ(item.ProductID)).Exists(ctx, tx)
		if err != nil || !productExists {
			return 0, schemas.HandlerErrorDB(sql.ErrNoRows, "Producto", schemas.Read)
		}

		deposit, err := boilmodels.Deposits(boilmodels.DepositWhere.ProductID.EQ(item.ProductID)).One(ctx, tx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				deposit = &boilmodels.Deposit{ProductID: item.ProductID, Stock: types.NewDecimal(decimal.New(0, 0))}
				if err := deposit.Insert(ctx, tx, boil.Infer()); err != nil {
					return 0, schemas.HandlerErrorDB(err, "Stock", schemas.Create)
				}
			} else {
				return 0, schemas.HandlerErrorDB(err, "Stock", schemas.Read)
			}
		}

		deposit.Stock.Big = deposit.Stock.Big.Add(deposit.Stock.Big, decimal.New(0, 0).SetFloat64(item.Amount))

		if _, err := deposit.Update(ctx, tx, boil.Whitelist(boilmodels.DepositColumns.Stock, boilmodels.DepositColumns.UpdatedAt)); err != nil {
			return 0, schemas.HandlerErrorDB(err, "Stock", schemas.Update)
		}

		subtotalItem := item.Amount * item.Price
		totalItem := 0.0
		if item.Discount > 0 {
			if item.TypeDiscount == "percent" {
				totalItem = subtotalItem - (subtotalItem * item.Discount / 100)
			} else {
				totalItem = subtotalItem - item.Discount
			}
		} else {
			totalItem = subtotalItem
		}

		total += totalItem

		newItem := &boilmodels.ExpenseBuyItem{
			ProductID:    item.ProductID,
			Amount:       types.NewDecimal(decimal.New(0, 0).SetFloat64(item.Amount)),
			Price:        types.NewDecimal(decimal.New(0, 0).SetFloat64(item.Price)),
			Discount:     types.NewDecimal(decimal.New(0, 0).SetFloat64(item.Discount)),
			TypeDiscount: item.TypeDiscount,
			Subtotal:     types.NewDecimal(decimal.New(0, 0).SetFloat64(subtotalItem)),
			Total:        types.NewDecimal(decimal.New(0, 0).SetFloat64(totalItem)),
		}
		newItems = append(newItems, newItem)
	}

	totalExpense := 0.0
	if expenseBuyCreate.Discount > 0 {
		if expenseBuyCreate.TypeDiscount == "percent" {
			totalExpense = total - (total * expenseBuyCreate.Discount / 100)
		} else {
			totalExpense = total - expenseBuyCreate.Discount
		}
	} else {
		totalExpense = total
	}

	subtotalDec := types.NewDecimal(decimal.New(0, 0).SetFloat64(total))
	discountDec := types.NewDecimal(decimal.New(0, 0).SetFloat64(expenseBuyCreate.Discount))
	totalExpenseDec := types.NewDecimal(decimal.New(0, 0).SetFloat64(totalExpense))

	expenseBuy := &boilmodels.ExpenseBuy{
		MemberID:     memberID,
		SupplierID:   expenseBuyCreate.SupplierID,
		Details:      null.StringFromPtr(expenseBuyCreate.Details),
		Subtotal:     subtotalDec,
		Discount:     discountDec,
		TypeDiscount: null.StringFrom(expenseBuyCreate.TypeDiscount).String,
		Total:        totalExpenseDec,
	}

	if err := expenseBuy.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Egreso de compras", schemas.Create)
	}

	for _, item := range newItems {
		item.ExpenseBuyID = expenseBuy.ID
		if err := item.Insert(ctx, tx, boil.Infer()); err != nil {
			return 0, schemas.HandlerErrorDB(err, "Items de egreso de compras", schemas.Create)
		}
	}

	totalPay := 0.0
	for _, pay := range expenseBuyCreate.PayExpenseBuy {
		totalPay += pay.Total
		pDec := types.NewDecimal(decimal.New(0, 0).SetFloat64(pay.Total))
		newPay := boilmodels.PayExpenseBuy{
			ExpenseBuyID: expenseBuy.ID,
			Total:        pDec,
			MethodPay:    null.StringFrom(pay.MethodPay).String,
		}
		if err := newPay.Insert(ctx, tx, boil.Infer()); err != nil {
			return 0, schemas.HandlerErrorDB(err, "Pagos de egreso de compras", schemas.Create)
		}
	}

	if math.Abs(totalPay-totalExpense) > 1 {
		message := fmt.Sprintf("la diferencia entre la suma de pagos (%.2f) y el total del egreso (%.2f)", totalPay, totalExpense)
		return 0, schemas.ErrorResponse(400, message, fmt.Errorf("%s", message))
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return expenseBuy.ID, nil
}

func (r *ExpenseBuyRepository) ExpenseBuyUpdate(memberID int64, expenseBuyUpdate *schemas.ExpenseBuyUpdate) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	existingExpense, err := boilmodels.ExpenseBuys(boilmodels.ExpenseBuyWhere.ID.EQ(expenseBuyUpdate.ID)).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Egreso de compras", schemas.Read)
	}

	supplierExists, err := boilmodels.Suppliers(boilmodels.SupplierWhere.ID.EQ(expenseBuyUpdate.SupplierID)).Exists(ctx, tx)
	if err != nil || !supplierExists {
		return schemas.HandlerErrorDB(sql.ErrNoRows, "Proveedor", schemas.Read)
	}

	oldItems, err := boilmodels.ExpenseBuyItems(boilmodels.ExpenseBuyItemWhere.ExpenseBuyID.EQ(expenseBuyUpdate.ID)).All(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Egreso de compras", schemas.Read)
	}

	for _, oldItem := range oldItems {
		deposit, err := boilmodels.Deposits(boilmodels.DepositWhere.ProductID.EQ(oldItem.ProductID)).One(ctx, tx)
		if err == nil {
			amt, _ := oldItem.Amount.Big.Float64()
			deposit.Stock.Big = deposit.Stock.Big.Sub(deposit.Stock.Big, decimal.New(0, 0).SetFloat64(amt))
			if _, err := deposit.Update(ctx, tx, boil.Whitelist(boilmodels.DepositColumns.Stock, boilmodels.DepositColumns.UpdatedAt)); err != nil {
				return schemas.HandlerErrorDB(err, "Stock", schemas.Update)
			}
		}
	}

	if _, err := boilmodels.ExpenseBuyItems(boilmodels.ExpenseBuyItemWhere.ExpenseBuyID.EQ(expenseBuyUpdate.ID)).DeleteAll(ctx, tx); err != nil {
		return schemas.HandlerErrorDB(err, "Items de egreso de compras", schemas.Delete)
	}

	var newItems []*boilmodels.ExpenseBuyItem
	total := 0.0

	for _, item := range expenseBuyUpdate.ExpenseBuyItem {
		if item.Amount <= 0 {
			return schemas.ErrorResponse(400, fmt.Sprintf("La cantidad para el producto %d no es válida", item.ProductID), fmt.Errorf("la cantidad para el producto %d no es válida", item.ProductID))
		}

		productExists, err := boilmodels.Products(boilmodels.ProductWhere.ID.EQ(item.ProductID)).Exists(ctx, tx)
		if err != nil || !productExists {
			return schemas.HandlerErrorDB(sql.ErrNoRows, "Producto", schemas.Read)
		}

		deposit, err := boilmodels.Deposits(boilmodels.DepositWhere.ProductID.EQ(item.ProductID)).One(ctx, tx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				deposit = &boilmodels.Deposit{ProductID: item.ProductID, Stock: types.NewDecimal(decimal.New(0, 0))}
				if err := deposit.Insert(ctx, tx, boil.Infer()); err != nil {
					return schemas.HandlerErrorDB(err, "Stock Depósito", schemas.Create)
				}
			} else {
				return schemas.HandlerErrorDB(err, "Stock Depósito", schemas.Read)
			}
		}

		deposit.Stock.Big = deposit.Stock.Big.Add(deposit.Stock.Big, decimal.New(0, 0).SetFloat64(item.Amount))
		if _, err := deposit.Update(ctx, tx, boil.Whitelist(boilmodels.DepositColumns.Stock, boilmodels.DepositColumns.UpdatedAt)); err != nil {
			return schemas.HandlerErrorDB(err, "Stock Depósito", schemas.Update)
		}

		subtotalItem := item.Amount * item.Price
		totalItem := 0.0
		if item.Discount > 0 {
			if item.TypeDiscount == "percent" {
				totalItem = subtotalItem - (subtotalItem * item.Discount / 100)
			} else {
				totalItem = subtotalItem - item.Discount
			}
		} else {
			totalItem = subtotalItem
		}

		newItem := &boilmodels.ExpenseBuyItem{
			ExpenseBuyID: expenseBuyUpdate.ID,
			ProductID:    item.ProductID,
			Amount:       types.NewDecimal(decimal.New(0, 0).SetFloat64(item.Amount)),
			Price:       types.NewDecimal(decimal.New(0, 0).SetFloat64(item.Price)),
			Subtotal:       types.NewDecimal(decimal.New(0, 0).SetFloat64(item.Discount)),
			Total:       types.NewDecimal(decimal.New(0, 0).SetFloat64(subtotalItem)),
			Discount:     types.NewDecimal(decimal.New(0, 0).SetFloat64(item.Discount)),
			TypeDiscount: null.StringFrom(item.TypeDiscount).String,
		}
		newItems = append(newItems, newItem)
		total += totalItem
	}

	totalExpense := 0.0
	if expenseBuyUpdate.Discount > 0 {
		if expenseBuyUpdate.Type == "percent" {
			totalExpense = total - (total * expenseBuyUpdate.Discount / 100)
		} else {
			totalExpense = total - expenseBuyUpdate.Discount
		}
	} else {
		totalExpense = total
	}

	existingExpense.SupplierID = expenseBuyUpdate.SupplierID
	existingExpense.Details = null.StringFromPtr(expenseBuyUpdate.Details)
	existingExpense.Subtotal = types.NewDecimal(decimal.New(0, 0).SetFloat64(total))
	existingExpense.Discount = types.NewDecimal(decimal.New(0, 0).SetFloat64(expenseBuyUpdate.Discount))
	existingExpense.TypeDiscount = null.StringFrom(expenseBuyUpdate.Type).String
	existingExpense.Total = types.NewDecimal(decimal.New(0, 0).SetFloat64(totalExpense))

	if _, err := existingExpense.Update(ctx, tx, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Egreso de compras", schemas.Update)
	}

	for _, item := range newItems {
		if err := item.Insert(ctx, tx, boil.Infer()); err != nil {
			return schemas.HandlerErrorDB(err, "Items de egreso de compras", schemas.Create)
		}
	}

	if _, err := boilmodels.PayExpenseBuys(boilmodels.PayExpenseBuyWhere.ExpenseBuyID.EQ(expenseBuyUpdate.ID)).DeleteAll(ctx, tx); err != nil {
		return schemas.HandlerErrorDB(err, "Pagos de egreso de compras", schemas.Delete)
	}

	totalPay := 0.0
	for _, pay := range expenseBuyUpdate.PayExpenseBuy {
		totalPay += pay.Total
		newPay := boilmodels.PayExpenseBuy{
			ExpenseBuyID: expenseBuyUpdate.ID,
			Total:       types.NewDecimal(decimal.New(0, 0).SetFloat64(pay.Total)),
			MethodPay:    null.StringFrom(pay.MethodPay).String,
		}
		if err := newPay.Insert(ctx, tx, boil.Infer()); err != nil {
			return schemas.HandlerErrorDB(err, "Pagos de egreso de compras", schemas.Create)
		}
	}

	if math.Abs(totalPay-totalExpense) > 1 {
		message := fmt.Sprintf("la diferencia entre la suma de pagos (%.2f) y el total del egreso (%.2f)", totalPay, totalExpense)
		return schemas.ErrorResponse(400, message, fmt.Errorf("%s", message))
	}

	return tx.Commit()
}

func (r *ExpenseBuyRepository) ExpenseBuyDelete(memberID int64, expenseBuyID int64) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	existingExpense, err := boilmodels.ExpenseBuys(boilmodels.ExpenseBuyWhere.ID.EQ(expenseBuyID)).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Compra", schemas.Read)
	}

	items, err := boilmodels.ExpenseBuyItems(boilmodels.ExpenseBuyItemWhere.ExpenseBuyID.EQ(expenseBuyID)).All(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Items de egreso de compras", schemas.Read)
	}

	for _, item := range items {
		deposit, err := boilmodels.Deposits(boilmodels.DepositWhere.ProductID.EQ(item.ProductID)).One(ctx, tx)
		if err != nil {
			continue
		}

		itemAmt, _ := item.Amount.Big.Float64()
		if deposit.Stock.Big.Cmp(item.Amount.Big) < 0 {
			return schemas.ErrorResponse(
				400,
				fmt.Sprintf("No se puede eliminar: stock insuficiente para el producto %d (disponible: %.2f, a revertir: %.2f)", item.ProductID, deposit.Stock, itemAmt),
				fmt.Errorf("stock insuficiente para revertir"),
			)
		}

		deposit.Stock.Big = deposit.Stock.Big.Sub(deposit.Stock.Big, item.Amount.Big)
		if _, err := deposit.Update(ctx, tx, boil.Whitelist(boilmodels.DepositColumns.Stock, boilmodels.DepositColumns.UpdatedAt)); err != nil {
			return schemas.HandlerErrorDB(err, "Stock", schemas.Update)
		}
	}

	if _, err := boilmodels.PayExpenseBuys(boilmodels.PayExpenseBuyWhere.ExpenseBuyID.EQ(expenseBuyID)).DeleteAll(ctx, tx); err != nil {
		return schemas.HandlerErrorDB(err, "Pagos de egreso de compras", schemas.Delete)
	}

	if _, err := boilmodels.ExpenseBuyItems(boilmodels.ExpenseBuyItemWhere.ExpenseBuyID.EQ(expenseBuyID)).DeleteAll(ctx, tx); err != nil {
		return schemas.HandlerErrorDB(err, "Items de egreso de compras", schemas.Delete)
	}

	if _, err := existingExpense.Delete(ctx, tx); err != nil {
		return schemas.HandlerErrorDB(err, "Egreso de compras", schemas.Delete)
	}

	return tx.Commit()
}
