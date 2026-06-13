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
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
)

func mapToIncomeSaleResponse(i *boilmodels.IncomeSale) *schemas.IncomeSaleResponse {
	if i == nil {
		return nil
	}

	subtotal, _ := i.Subtotal.Float64()
	discount, _ := i.Discount.Float64()
	total, _ := i.Total.Float64()

	res := &schemas.IncomeSaleResponse{
		ID:       i.ID,
		SubTotal: subtotal,
		Discount: discount,
		Type:     i.Type,
		Total:    total,
		IsBudget: i.IsBudget,
	}

	if i.InvoiceID.Valid {
		res.InvoiceID = &i.InvoiceID.Int64
	}
	res.CreatedAt = i.CreatedAt

	if i.R != nil {
		// Client relation
		if i.R.Client != nil {
			var clientDebt float64 = 0 // Optional placeholder
			res.Client = mapToClientResponseDTO(i.R.Client, &clientDebt)
		}

		// Member relation
		if i.R.Member != nil {
			res.Member = *mapToMemberSimpleDTO(i.R.Member)
		}

		for _, item := range i.R.IncomeSaleItems {
			iAmt, _ := item.Amount.Big.Float64()
			iPrice, _ := item.Price.Big.Float64()
			iDiscount, _ := item.Discount.Big.Float64()
			iSubtotal, _ := item.Subtotal.Big.Float64()
			iTotal, _ := item.Total.Big.Float64()

			ier := schemas.IncomeSaleItemResponse{
				ID:       item.ID,
				Amount:   iAmt,
				Price:    iPrice,
				Discount: iDiscount,
				SubTotal: iSubtotal,
				Total:    iTotal,
			}
			if item.TypeDiscount != "" {
				ier.TypeDiscount = item.TypeDiscount
			}
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
			res.Items = append(res.Items, ier)
		}

		for _, pay := range i.R.PayIncomes {
			pTotal, _ := pay.Total.Big.Float64()
			pr := schemas.PayResponse{
				ID:        pay.ID,
				Total:     pTotal,
				MethodPay: pay.MethodPay,
			}
			pr.CreatedAt = pay.CreatedAt
			res.Pay = append(res.Pay, pr)
		}
	}

	return res
}

func mapToIncomeSaleResponseDTO(i *boilmodels.IncomeSale) *schemas.IncomeSaleResponseDTO {
	if i == nil {
		return nil
	}

	total, _ := i.Total.Big.Float64()

	res := &schemas.IncomeSaleResponseDTO{
		ID:    i.ID,
		Total: total,
	}

	if i.InvoiceID.Valid {
		res.InvoiceID = &i.InvoiceID.Int64
	}
	res.CreatedAt = i.CreatedAt

	if i.R != nil {
		if i.R.Client != nil {
			var clientDebt float64 = 0 // Placeholder
			res.Client = *mapToClientSimpleDTO(i.R.Client, &clientDebt)
		}

		if i.R.Member != nil {
			res.Member = *mapToMemberSimpleDTO(i.R.Member)
		}

		for _, pay := range i.R.PayIncomes {
			pTotal, _ := pay.Total.Big.Float64()
			pr := schemas.PayResponse{
				ID:        pay.ID,
				Total:     pTotal,
				MethodPay: pay.MethodPay,
			}
			pr.CreatedAt = pay.CreatedAt
			res.Pay = append(res.Pay, pr)
		}
	}

	return res
}

func mapToMemberSimpleDTO(m *boilmodels.Member) *schemas.MemberSimpleDTO {
	if m == nil {
		return nil
	}
	res := &schemas.MemberSimpleDTO{
		ID: m.ID,
		FirstName: m.FirstName,
		LastName: m.LastName,
		Username: m.Username,
	}
	return res
}

func mapToClientSimpleDTO(c *boilmodels.Client, debt *float64) *schemas.ClientSimpleDTO {
	if c == nil {
		return nil
	}
	res := &schemas.ClientSimpleDTO{
		ID:        c.ID,
		FirstName: c.FirstName,
		LastName:  c.LastName,
	}
	if c.CompanyName.Valid {
		res.CompanyName = &c.CompanyName.String
	}
	return res
}

func (i *IncomeSaleRepository) IncomeSaleGetByID(pointSaleID, id int64) (*schemas.IncomeSaleResponse, error) {
	ctx := context.Background()

	incomeSaleModel, err := boilmodels.IncomeSales(
		boilmodels.IncomeSaleWhere.ID.EQ(id),
		qm.Load(boilmodels.IncomeSaleRels.Member),
		qm.Load(boilmodels.IncomeSaleRels.Client),
		qm.Load(boilmodels.IncomeSaleRels.IncomeSaleItems),
		qm.Load(qm.Rels(boilmodels.IncomeSaleRels.IncomeSaleItems, boilmodels.IncomeSaleItemRels.Product)),
		qm.Load(boilmodels.IncomeSaleRels.PayIncomes),
	).One(ctx, i.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Ingreso de venta", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Ingreso de venta", schemas.Read)
	}

	return mapToIncomeSaleResponse(incomeSaleModel), nil
}

func (i *IncomeSaleRepository) IncomeSaleGetByDate(pointSaleID int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.IncomeSaleResponseDTO, int64, error) {
	ctx := context.Background()
	offSet := (page - 1) * limit

	qms := []qm.QueryMod{
		boilmodels.IncomeSaleWhere.PointSaleID.EQ(pointSaleID),
		boilmodels.IncomeSaleWhere.CreatedAt.GTE(fromDate),
		boilmodels.IncomeSaleWhere.CreatedAt.LTE(toDate),
	}

	total, err := boilmodels.IncomeSales(qms...).Count(ctx, i.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Ingreso de venta", schemas.Read)
	}

	qms = append(qms,
		qm.Load(boilmodels.IncomeSaleRels.Member),
		qm.Load(boilmodels.IncomeSaleRels.Client),
		qm.Load(boilmodels.IncomeSaleRels.PayIncomes),
		qm.OrderBy("created_at DESC"),
		qm.Offset(offSet),
		qm.Limit(limit),
	)

	incomeSaleModel, err := boilmodels.IncomeSales(qms...).All(ctx, i.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Ingreso de venta", schemas.Read)
	}

	var results []*schemas.IncomeSaleResponseDTO
	for _, inc := range incomeSaleModel {
		results = append(results, mapToIncomeSaleResponseDTO(inc))
	}

	return results, total, nil
}

func (i *IncomeSaleRepository) IncomeSaleCreate(memberID, pointSaleID int64, incomeSaleCreate *schemas.IncomeSaleCreate) (int64, error) {
	ctx := context.Background()
	tx, err := i.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return 0, err
	}

	register, err := boilmodels.CashRegisters(
		boilmodels.CashRegisterWhere.IsClose.EQ(false),
		boilmodels.CashRegisterWhere.PointSaleID.EQ(pointSaleID),
		qm.OrderBy("hour_open DESC"),
	).One(ctx, tx)

	if err != nil {
		return 0, schemas.HandlerErrorDB(err, "Ingreso de venta", schemas.Read)
	}

	var isDeposit bool
	pointSale, err := boilmodels.PointSales(
		qm.Select(boilmodels.PointSaleColumns.IsDeposit),
		boilmodels.PointSaleWhere.ID.EQ(pointSaleID),
	).One(ctx, tx)
	if err != nil {
		return 0, schemas.HandlerErrorDB(err, "Ingreso de venta", schemas.Read)
	}
	isDeposit = pointSale.IsDeposit

	clientExists, err := boilmodels.Clients(boilmodels.ClientWhere.ID.EQ(incomeSaleCreate.ClientID)).Exists(ctx, tx)
	if err != nil || !clientExists {
		return 0, schemas.HandlerErrorDB(sql.ErrNoRows, "Cliente", schemas.Read)
	}

	var newItems []*boilmodels.IncomeSaleItem
	subtotal := 0.0

	for _, item := range incomeSaleCreate.Items {
		product, err := boilmodels.Products(
			qm.Select(boilmodels.ProductColumns.Price),
			boilmodels.ProductWhere.ID.EQ(item.ProductID),
		).One(ctx, tx)

		if err != nil {
			return 0, schemas.HandlerErrorDB(err, "Producto", schemas.Read)
		}
		productPrice, _ := product.Price.Big.Float64()

		if isDeposit {
			stock, err := boilmodels.Deposits(boilmodels.DepositWhere.ProductID.EQ(item.ProductID)).One(ctx, tx)
			if err != nil {
				return 0, schemas.HandlerErrorDB(err, "Stock Depósito", schemas.Read)
			}
			stockF, _ := stock.Stock.Big.Float64()
			if stockF < float64(item.Amount) {
				return 0, schemas.ErrorResponse(
					400,
					fmt.Sprintf("stock insuficiente para el producto %d (disponible: %.2f, requerido: %v)", item.ProductID, stockF, item.Amount),
					fmt.Errorf("stock insuficiente"),
				)
			}
			stock.Stock = types.NewDecimal(decimal.New(0, 0).SetFloat64(stockF - float64(item.Amount)))
			if _, err := stock.Update(ctx, tx, boil.Whitelist(boilmodels.DepositColumns.Stock, boilmodels.DepositColumns.UpdatedAt)); err != nil {
				return 0, schemas.HandlerErrorDB(err, "Stock Depósito", schemas.Update)
			}
		} else {
			stock, err := boilmodels.StockPointSales(
				boilmodels.StockPointSaleWhere.PointSaleID.EQ(pointSaleID),
				boilmodels.StockPointSaleWhere.ProductID.EQ(item.ProductID),
			).One(ctx, tx)

			if err != nil {
				return 0, schemas.HandlerErrorDB(err, "Stock Punto de Venta", schemas.Read)
			}
			stockF, _ := stock.Stock.Big.Float64()
			if stockF < float64(item.Amount) {
				return 0, schemas.ErrorResponse(
					400,
					fmt.Sprintf("stock insuficiente para el producto %d (disponible: %.2f, requerido: %v)", item.ProductID, stockF, item.Amount),
					fmt.Errorf("stock insuficiente"),
				)
			}
			stock.Stock = types.NewDecimal(decimal.New(0, 0).SetFloat64(stockF - float64(item.Amount)))
			if _, err := stock.Update(ctx, tx, boil.Whitelist(boilmodels.StockPointSaleColumns.Stock, boilmodels.StockPointSaleColumns.UpdatedAt)); err != nil {
				return 0, schemas.HandlerErrorDB(err, "Stock Punto de Venta", schemas.Update)
			}
		}

		priceCostExp, err := boilmodels.ExpenseBuyItems(
			qm.Select(boilmodels.ExpenseBuyItemColumns.Price),
			boilmodels.ExpenseBuyItemWhere.ProductID.EQ(item.ProductID),
			qm.OrderBy("created_at DESC"),
		).One(ctx, tx)

		priceCostVal := productPrice
		if err == nil {
			priceCostVal, _ = priceCostExp.Price.Big.Float64()
		}

		subtotalItem := item.Amount * productPrice
		totalItem := 0.0

		if item.Discount > 0 {
			if item.TypeDiscount == "amount" {
				totalItem = subtotalItem - item.Discount
			} else if item.TypeDiscount == "percent" {
				totalItem = subtotalItem - (subtotalItem * item.Discount / 100)
			}
		} else {
			totalItem = subtotalItem
		}

		newItems = append(newItems, &boilmodels.IncomeSaleItem{
			ProductID:    null.Int64From(item.ProductID),
			Amount:       types.NewDecimal(decimal.New(0, 0).SetFloat64(item.Amount)),
			Price:        types.NewDecimal(decimal.New(0, 0).SetFloat64(productPrice)),
			PriceCost:    types.NewDecimal(decimal.New(0, 0).SetFloat64(priceCostVal)),
			Subtotal:     types.NewDecimal(decimal.New(0, 0).SetFloat64(subtotalItem)),
			Discount:     types.NewDecimal(decimal.New(0, 0).SetFloat64(item.Discount)),
			TypeDiscount: item.TypeDiscount,
			Total:        types.NewDecimal(decimal.New(0, 0).SetFloat64(totalItem)),
		})

		subtotal += totalItem
	}

	totalIncome := 0.0
	if incomeSaleCreate.Discount > 0 {
		if incomeSaleCreate.Type == "amount" {
			totalIncome = subtotal - incomeSaleCreate.Discount
		} else if incomeSaleCreate.Type == "percent" {
			totalIncome = subtotal - (subtotal * incomeSaleCreate.Discount / 100)
		}
	} else {
		totalIncome = subtotal
	}

	income := &boilmodels.IncomeSale{
		PointSaleID:    pointSaleID,
		MemberID:       memberID,
		ClientID:       incomeSaleCreate.ClientID,
		CashRegisterID: register.ID,
		Subtotal:       types.NewDecimal(decimal.New(0, 0).SetFloat64(subtotal)),
		Discount:       types.NewDecimal(decimal.New(0, 0).SetFloat64(incomeSaleCreate.Discount)),
		Type:           incomeSaleCreate.Type,
		Total:          types.NewDecimal(decimal.New(0, 0).SetFloat64(totalIncome)),
		IsBudget:       *incomeSaleCreate.IsBudget,
	}

	if err := income.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.ErrorResponse(500, "Error al crear el ingreso", err)
	}
	incomeSaleID := income.ID

	for _, item := range newItems {
		item.IncomeSaleID = null.Int64From(incomeSaleID)
		if err := item.Insert(ctx, tx, boil.Infer()); err != nil {
			return 0, schemas.ErrorResponse(500, "Error al crear items del ingreso", err)
		}
	}

	totalPay := 0.0
	for _, pay := range incomeSaleCreate.Pay {
		totalPay += pay.Total

		newPay := boilmodels.PayIncome{
			IncomeSaleID:   incomeSaleID,
			CashRegisterID: null.Int64From(register.ID),
			ClientID:       null.Int64From(incomeSaleCreate.ClientID),
			Total:          types.NewDecimal(decimal.New(0, 0).SetFloat64(pay.Total)),
			MethodPay:      pay.MethodPay,
		}

		if pay.MethodPay == "credit" {
			newPay.CashRegisterID = null.NewInt64(0, false)
		}

		if err := newPay.Insert(ctx, tx, boil.Infer()); err != nil {
			return 0, schemas.HandlerErrorDB(err, "Pago de Ingreso", schemas.Create)
		}
	}

	if math.Abs(totalPay-totalIncome) > 1 {
		message := fmt.Sprintf("la diferencia entre la suma de pagos (%.2f) y el total del ingreso (%.2f)", totalPay, totalIncome)
		return 0, schemas.ErrorResponse(400, message, fmt.Errorf("%s", message))
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return incomeSaleID, nil
}

func (i *IncomeSaleRepository) IncomeSaleUpdate(memberID, pointSaleID int64, incomeSaleUpdate *schemas.IncomeSaleUpdate) error {
	ctx := context.Background()
	tx, err := i.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	existingIncome, err := boilmodels.IncomeSales(
		boilmodels.IncomeSaleWhere.ID.EQ(incomeSaleUpdate.ID),
		boilmodels.IncomeSaleWhere.PointSaleID.EQ(pointSaleID),
	).One(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Ingreso de venta", schemas.Read)
	}

	oldItems, err := boilmodels.IncomeSaleItems(boilmodels.IncomeSaleItemWhere.IncomeSaleID.EQ(null.Int64From(incomeSaleUpdate.ID))).All(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Items de Ingreso", schemas.Read)
	}

	for _, oldItem := range oldItems {
		oldItemAmt, _ := oldItem.Amount.Big.Float64()
		var isDeposit bool
		pointSale, _ := boilmodels.PointSales(qm.Select(boilmodels.PointSaleColumns.IsDeposit), boilmodels.PointSaleWhere.ID.EQ(pointSaleID)).One(ctx, tx)
		if pointSale != nil {
			isDeposit = pointSale.IsDeposit
		}

		if isDeposit {
			deposit, err := boilmodels.Deposits(boilmodels.DepositWhere.ProductID.EQ(oldItem.ProductID.Int64)).One(ctx, tx)
			if err == nil {
				stockF, _ := deposit.Stock.Big.Float64()
				deposit.Stock = types.NewDecimal(decimal.New(0, 0).SetFloat64(stockF + oldItemAmt))
				deposit.Update(ctx, tx, boil.Whitelist(boilmodels.DepositColumns.Stock, boilmodels.DepositColumns.UpdatedAt))
			}
		} else {
			stockPointSale, err := boilmodels.StockPointSales(boilmodels.StockPointSaleWhere.PointSaleID.EQ(pointSaleID), boilmodels.StockPointSaleWhere.ProductID.EQ(oldItem.ProductID.Int64)).One(ctx, tx)
			if err == nil {
				stockF, _ := stockPointSale.Stock.Big.Float64()
				stockPointSale.Stock = types.NewDecimal(decimal.New(0, 0).SetFloat64(stockF + oldItemAmt))
				stockPointSale.Update(ctx, tx, boil.Whitelist(boilmodels.StockPointSaleColumns.Stock, boilmodels.StockPointSaleColumns.UpdatedAt))
			}
		}
	}

	if _, err := boilmodels.IncomeSaleItems(boilmodels.IncomeSaleItemWhere.IncomeSaleID.EQ(null.Int64From(incomeSaleUpdate.ID))).DeleteAll(ctx, tx); err != nil {
		return schemas.HandlerErrorDB(err, "Items de Ingreso", schemas.Delete)
	}

	var isDeposit bool
	pointSale, err := boilmodels.PointSales(qm.Select(boilmodels.PointSaleColumns.IsDeposit), boilmodels.PointSaleWhere.ID.EQ(pointSaleID)).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Punto de Venta", schemas.Read)
	}
	isDeposit = pointSale.IsDeposit

	clientExists, err := boilmodels.Clients(boilmodels.ClientWhere.ID.EQ(incomeSaleUpdate.ClientID)).Exists(ctx, tx)
	if err != nil || !clientExists {
		return schemas.HandlerErrorDB(sql.ErrNoRows, "Cliente", schemas.Read)
	}

	var newIncomeSaleItems []*boilmodels.IncomeSaleItem
	subtotal := 0.0

	for _, item := range incomeSaleUpdate.Items {
		product, err := boilmodels.Products(qm.Select(boilmodels.ProductColumns.Price), boilmodels.ProductWhere.ID.EQ(item.ProductID)).One(ctx, tx)
		if err != nil {
			return schemas.HandlerErrorDB(err, "Producto", schemas.Read)
		}
		productPrice, _ := product.Price.Big.Float64()

		if isDeposit {
			stock, err := boilmodels.Deposits(boilmodels.DepositWhere.ProductID.EQ(item.ProductID)).One(ctx, tx)
			if err != nil {
				return schemas.HandlerErrorDB(err, "Stock Depósito", schemas.Read)
			}
			stockF, _ := stock.Stock.Big.Float64()
			if stockF < float64(item.Amount) {
				return schemas.ErrorResponse(400, "stock insuficiente", fmt.Errorf("stock insuficiente"))
			}
			stock.Stock = types.NewDecimal(decimal.New(0, 0).SetFloat64(stockF - float64(item.Amount)))
			stock.Update(ctx, tx, boil.Whitelist(boilmodels.DepositColumns.Stock, boilmodels.DepositColumns.UpdatedAt))
		} else {
			stock, err := boilmodels.StockPointSales(boilmodels.StockPointSaleWhere.PointSaleID.EQ(pointSaleID), boilmodels.StockPointSaleWhere.ProductID.EQ(item.ProductID)).One(ctx, tx)
			if err != nil {
				return schemas.HandlerErrorDB(err, "Stock Punto de Venta", schemas.Read)
			}
			stockF, _ := stock.Stock.Big.Float64()
			if stockF < float64(item.Amount) {
				return schemas.ErrorResponse(400, "stock insuficiente", fmt.Errorf("stock insuficiente"))
			}
			stock.Stock = types.NewDecimal(decimal.New(0, 0).SetFloat64(stockF - float64(item.Amount)))
			stock.Update(ctx, tx, boil.Whitelist(boilmodels.StockPointSaleColumns.Stock, boilmodels.StockPointSaleColumns.UpdatedAt))
		}

		priceCostExp, err := boilmodels.ExpenseBuyItems(
			qm.Select(boilmodels.ExpenseBuyItemColumns.Price),
			boilmodels.ExpenseBuyItemWhere.ProductID.EQ(item.ProductID),
			qm.OrderBy("created_at DESC"),
		).One(ctx, tx)

		priceCostVal := productPrice
		if err == nil {
			priceCostVal, _ = priceCostExp.Price.Big.Float64()
		}

		subtotalItem := item.Amount * productPrice
		totalItem := 0.0

		if item.Discount > 0 {
			if item.TypeDiscount == "amount" {
				totalItem = subtotalItem - item.Discount
			} else if item.TypeDiscount == "percent" {
				totalItem = subtotalItem - (subtotalItem * item.Discount / 100)
			}
		} else {
			totalItem = subtotalItem
		}

		newIncomeSaleItems = append(newIncomeSaleItems, &boilmodels.IncomeSaleItem{
			IncomeSaleID: null.Int64From(incomeSaleUpdate.ID),
			ProductID:    null.Int64From(item.ProductID),
			Amount:       types.NewDecimal(decimal.New(0, 0).SetFloat64(item.Amount)),
			Price:        types.NewDecimal(decimal.New(0, 0).SetFloat64(productPrice)),
			PriceCost:    types.NewDecimal(decimal.New(0, 0).SetFloat64(priceCostVal)),
			Subtotal:     types.NewDecimal(decimal.New(0, 0).SetFloat64(subtotalItem)),
			Discount:     types.NewDecimal(decimal.New(0, 0).SetFloat64(item.Discount)),
			TypeDiscount: item.TypeDiscount,
			Total:        types.NewDecimal(decimal.New(0, 0).SetFloat64(totalItem)),
		})
		subtotal += totalItem
	}

	totalIncome := 0.0
	if incomeSaleUpdate.Discount > 0 {
		if incomeSaleUpdate.Type == "amount" {
			totalIncome = subtotal - incomeSaleUpdate.Discount
		} else if incomeSaleUpdate.Type == "percent" {
			totalIncome = subtotal - (subtotal * incomeSaleUpdate.Discount / 100)
		}
	} else {
		totalIncome = subtotal
	}

	existingIncome.ClientID = incomeSaleUpdate.ClientID
	existingIncome.Subtotal = types.NewDecimal(decimal.New(0, 0).SetFloat64(subtotal))
	existingIncome.Discount = types.NewDecimal(decimal.New(0, 0).SetFloat64(incomeSaleUpdate.Discount))
	existingIncome.Type = incomeSaleUpdate.Type
	existingIncome.Total = types.NewDecimal(decimal.New(0, 0).SetFloat64(totalIncome))
	existingIncome.IsBudget = incomeSaleUpdate.IsBudget

	if _, err := existingIncome.Update(ctx, tx, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Ingreso de venta", schemas.Update)
	}

	for _, item := range newIncomeSaleItems {
		if err := item.Insert(ctx, tx, boil.Infer()); err != nil {
			return schemas.HandlerErrorDB(err, "Items de Ingreso", schemas.Create)
		}
	}

	if _, err := boilmodels.PayIncomes(boilmodels.PayIncomeWhere.IncomeSaleID.EQ(incomeSaleUpdate.ID)).DeleteAll(ctx, tx); err != nil {
		return schemas.HandlerErrorDB(err, "Pagos de Ingreso", schemas.Delete)
	}

	totalPay := 0.0
	for _, pay := range incomeSaleUpdate.Pay {
		totalPay += pay.Total

		newPay := boilmodels.PayIncome{
			IncomeSaleID:   incomeSaleUpdate.ID,
			CashRegisterID: null.Int64From(existingIncome.CashRegisterID),
			ClientID:       null.Int64From(incomeSaleUpdate.ClientID),
			Total:          types.NewDecimal(decimal.New(0, 0).SetFloat64(pay.Total)),
			MethodPay:      pay.MethodPay,
		}

		if pay.MethodPay == "credit" {
			newPay.CashRegisterID = null.NewInt64(0, false)
		}

		if err := newPay.Insert(ctx, tx, boil.Infer()); err != nil {
			return schemas.HandlerErrorDB(err, "Pagos de Ingreso", schemas.Create)
		}
	}

	if math.Abs(totalPay-totalIncome) > 1 {
		message := fmt.Sprintf("la diferencia entre la suma de pagos (%.2f) y el total del ingreso (%.2f)", totalPay, totalIncome)
		return schemas.ErrorResponse(400, message, fmt.Errorf("%s", message))
	}

	return tx.Commit()
}

func (i *IncomeSaleRepository) IncomeSaleDelete(memberID, incomeSaleID, pointSaleID int64) error {
	ctx := context.Background()
	tx, err := i.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	existingIncome, err := boilmodels.IncomeSales(
		boilmodels.IncomeSaleWhere.ID.EQ(incomeSaleID),
		boilmodels.IncomeSaleWhere.PointSaleID.EQ(pointSaleID),
	).One(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Ingreso de venta", schemas.Read)
	}

	items, err := boilmodels.IncomeSaleItems(boilmodels.IncomeSaleItemWhere.IncomeSaleID.EQ(null.Int64From(incomeSaleID))).All(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Items de Ingreso", schemas.Read)
	}

	var isDeposit bool
	pointSale, err := boilmodels.PointSales(qm.Select(boilmodels.PointSaleColumns.IsDeposit), boilmodels.PointSaleWhere.ID.EQ(pointSaleID)).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Punto de Venta", schemas.Read)
	}
	isDeposit = pointSale.IsDeposit

	for _, item := range items {
		itemAmt, _ := item.Amount.Big.Float64()
		if isDeposit {
			deposit, err := boilmodels.Deposits(boilmodels.DepositWhere.ProductID.EQ(item.ProductID.Int64)).One(ctx, tx)
			if err == nil {
				stockF, _ := deposit.Stock.Big.Float64()
				deposit.Stock = types.NewDecimal(decimal.New(0, 0).SetFloat64(stockF + itemAmt))
				deposit.Update(ctx, tx, boil.Whitelist(boilmodels.DepositColumns.Stock, boilmodels.DepositColumns.UpdatedAt))
			}
		} else {
			stockPointSale, err := boilmodels.StockPointSales(boilmodels.StockPointSaleWhere.PointSaleID.EQ(pointSaleID), boilmodels.StockPointSaleWhere.ProductID.EQ(item.ProductID.Int64)).One(ctx, tx)
			if err == nil {
				stockF, _ := stockPointSale.Stock.Big.Float64()
				stockPointSale.Stock = types.NewDecimal(decimal.New(0, 0).SetFloat64(stockF + itemAmt))
				stockPointSale.Update(ctx, tx, boil.Whitelist(boilmodels.StockPointSaleColumns.Stock, boilmodels.StockPointSaleColumns.UpdatedAt))
			}
		}
	}

	if _, err := boilmodels.PayIncomes(boilmodels.PayIncomeWhere.IncomeSaleID.EQ(incomeSaleID)).DeleteAll(ctx, tx); err != nil {
		return schemas.HandlerErrorDB(err, "Pagos de Ingreso", schemas.Delete)
	}

	if _, err := boilmodels.IncomeSaleItems(boilmodels.IncomeSaleItemWhere.IncomeSaleID.EQ(null.Int64From(incomeSaleID))).DeleteAll(ctx, tx); err != nil {
		return schemas.HandlerErrorDB(err, "Items de Ingreso", schemas.Delete)
	}

	if _, err := existingIncome.Delete(ctx, tx); err != nil {
		return schemas.HandlerErrorDB(err, "Ingreso de venta", schemas.Delete)
	}

	return tx.Commit()
}
