package repositories

import (
	"context"
	"fmt"
	"time"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
)

func mapToMemberSimpleDTO(c *boilmodels.Member) *schemas.MemberSimpleDTO {
	if c == nil {
		return nil
	}
	return &schemas.MemberSimpleDTO{
		ID:        c.ID,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Username:  c.Username,
	}
}

func (r *CashRegisterRepository) CashRegisterExistOpen(pointSaleID int64) (bool, error) {
	ctx := context.Background()
	exists, err := boilmodels.CashRegisters(
		boilmodels.CashRegisterWhere.IsClose.EQ(false),
		boilmodels.CashRegisterWhere.PointSaleID.EQ(pointSaleID),
	).Exists(ctx, r.DB)
	if err != nil {
		return false, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
	}
	return exists, nil
}

func (r *CashRegisterRepository) CashRegisterGetByID(pointSaleID, id int64) (*schemas.CashRegisterFullResponse, error) {
	ctx := context.Background()

	register, err := boilmodels.CashRegisters(
		boilmodels.CashRegisterWhere.ID.EQ(id),
		boilmodels.CashRegisterWhere.PointSaleID.EQ(pointSaleID),
		qm.Load(boilmodels.CashRegisterRels.MemberOpen),
		qm.Load(boilmodels.CashRegisterRels.MemberClose),
	).One(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
	}

	var openAmount float64
	if register.OpenAmount.Valid {
		openAmount, _ = register.OpenAmount.Big.Float64()
	}

	res := &schemas.CashRegisterFullResponse{
		ID:         register.ID,
		IsClose:    register.IsClose,
		OpenAmount: openAmount,
	}

	if register.CreatedAt.Valid {
		res.CreatedAt = register.CreatedAt.Time
	}
	if register.HourOpen.Valid {
		res.HourOpen = register.HourOpen.Time
	}
	if register.CloseAmount.Valid {
		cAmount, _ := register.CloseAmount.Big.Float64()
		res.CloseAmount = &cAmount
	}
	if register.HourClose.Valid {
		hClose := register.HourClose.Time
		res.HourClose = &hClose
	}

	if register.R != nil {
		if register.R.MemberOpen != nil {
			m := mapToMemberSimpleDTO(register.R.MemberOpen)
			res.MemberOpen = *m
		}
		if register.R.MemberClose != nil {
			res.MemberClose = mapToMemberSimpleDTO(register.R.MemberClose)
		}
	}

	// 1. IncomeSales
	incomeSales, err := boilmodels.IncomeSales(
		boilmodels.IncomeSaleWhere.CashRegisterID.EQ(null.Int64From(id)),
		boilmodels.IncomeSaleWhere.PointSaleID.EQ(pointSaleID),
		qm.Load(boilmodels.IncomeSaleRels.IncomeSaleItems),
		qm.Load(qm.Rels(boilmodels.IncomeSaleRels.IncomeSaleItems, boilmodels.IncomeSaleItemRels.Product)),
		qm.Load(boilmodels.IncomeSaleRels.PayIncomes),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
	}

	var incomes []*schemas.IncomeSaleSimpleResponse
	for _, in := range incomeSales {
		total, _ := in.Total.Big.Float64()
		iso := &schemas.IncomeSaleSimpleResponse{
			ID:       in.ID,
			Total:    total,
			IsBudget: in.IsBudget,
		}
		if in.CreatedAt.Valid {
			iso.CreatedAt = in.CreatedAt.Time
		}
		if in.InvoiceID.Valid {
			iso.InvoiceID = &in.InvoiceID.Int64
		}

		if in.R != nil {
			for _, item := range in.R.IncomeSaleItems {
				itemAmt, _ := item.Amount.Big.Float64()
				itemTotal, _ := item.Total.Big.Float64()
				iResp := schemas.IncomeSaleItemResponseDTO{
					ID:     item.ID,
					Amount: itemAmt,
					Total:  itemTotal,
				}
				if item.CreatedAt.Valid {
					iResp.CreatedAt = item.CreatedAt.Time
				}
				if item.R != nil && item.R.Product != nil {
					pPrice, _ := item.R.Product.Price.Big.Float64()
					iResp.Product = schemas.ProductSimpleResponseDTO{
						ID:    item.R.Product.ID,
						Code:  item.R.Product.Code,
						Name:  item.R.Product.Name,
						Price: pPrice,
					}
				}
				iso.Items = append(iso.Items, iResp)
			}

			for _, pay := range in.R.PayIncomes {
				val, _ := pay.Total.Big.Float64()
				pResp := schemas.PayResponse{
					ID:        pay.ID,
					Total:     val,
					MethodPay: pay.MethodPay,
				}
				if pay.CreatedAt.Valid {
					pResp.CreatedAt = pay.CreatedAt.Time.Format(time.RFC3339)
				}
				iso.Pay = append(iso.Pay, pResp)
			}
		}

		incomes = append(incomes, iso)
	}

	// 2. IncomeOthers
	incomeOthers, err := boilmodels.IncomeOthers(
		boilmodels.IncomeOtherWhere.CashRegisterID.EQ(null.Int64From(id)),
		boilmodels.IncomeOtherWhere.PointSaleID.EQ(null.Int64From(pointSaleID)),
		qm.Load(boilmodels.IncomeOtherRels.Member),
		qm.Load(boilmodels.IncomeOtherRels.TypeIncome),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
	}

	var others []*schemas.IncomeOtherResponse
	for _, o := range incomeOthers {
		val, _ := o.Total.Big.Float64()
		or := &schemas.IncomeOtherResponse{
			ID:           o.ID,
			Total:        val,
			MethodIncome: o.MethodIncome,
		}
		if o.Details.Valid {
			or.Details = &o.Details.String
		}
		if o.CreatedAt.Valid {
			or.CreatedAt = o.CreatedAt.Time
		}
		if o.R != nil {
			if o.R.Member != nil {
				or.Member = mapToMemberSimpleDTO(o.R.Member)
			}
			if o.R.TypeIncome != nil {
				or.TypeIncome = schemas.TypeIncomeResponse{
					ID:   o.R.TypeIncome.ID,
					Name: o.R.TypeIncome.Name,
				}
			}
		}
		others = append(others, or)
	}

	// 3. ExpenseOthers
	expenseOthers, err := boilmodels.ExpenseOthers(
		boilmodels.ExpenseOtherWhere.CashRegisterID.EQ(null.Int64From(id)),
		boilmodels.ExpenseOtherWhere.PointSaleID.EQ(null.Int64From(pointSaleID)),
		qm.Load(boilmodels.ExpenseOtherRels.Member),
		qm.Load(boilmodels.ExpenseOtherRels.TypeExpense),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
	}

	var expOthers []*schemas.ExpenseOtherResponse
	for _, e := range expenseOthers {
		val, _ := e.Total.Big.Float64()
		eo := &schemas.ExpenseOtherResponse{
			ID:        e.ID,
			Total:     val,
			PayMethod: e.PayMethod.String,
		}
		if e.Details.Valid {
			eo.Details = &e.Details.String
		}
		if e.CreatedAt.Valid {
			eo.CreatedAt = e.CreatedAt.Time
		}
		if e.R != nil {
			if e.R.Member != nil {
				eo.Member = mapToMemberSimpleDTO(e.R.Member)
			}
			if e.R.TypeExpense != nil {
				eo.TypeExpense = schemas.TypeExpenseResponse{
					ID:   e.R.TypeExpense.ID,
					Name: e.R.TypeExpense.Name,
				}
			}
		}
		expOthers = append(expOthers, eo)
	}

	res.IncomeSale = incomes
	res.IncomeOther = others
	res.ExpenseOther = expOthers

	return res, nil
}

func (r *CashRegisterRepository) CashRegisterOpen(pointSaleID int64, userID int64, amountOpen schemas.CashRegisterOpen) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", userID)); err != nil {
		return err
	}

	exists, err := boilmodels.CashRegisters(
		boilmodels.CashRegisterWhere.IsClose.EQ(false),
		boilmodels.CashRegisterWhere.PointSaleID.EQ(pointSaleID),
	).Exists(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Caja", schemas.Read)
	}

	if exists {
		return schemas.ErrorResponse(400, "ya existe una apertura de caja, antes de continuar cierre la caja", fmt.Errorf("ya existe una apertura de caja, antes de continuar cerrar"))
	}

	openAmtDec := types.NewNullDecimal(types.NewDecimal(fmt.Sprintf("%.2f", amountOpen.OpenAmount)))
	registerOpen := boilmodels.CashRegister{
		PointSaleID:  pointSaleID,
		MemberOpenID: userID,
		OpenAmount:   openAmtDec,
		HourOpen:     null.TimeFrom(time.Now().UTC()),
	}

	if err := registerOpen.Insert(ctx, tx, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Caja", schemas.Create)
	}

	return tx.Commit()
}

func (r *CashRegisterRepository) CashRegisterClose(pointSaleID int64, userID int64, amountOpen schemas.CashRegisterClose) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", userID)); err != nil {
		return err
	}

	register, err := boilmodels.CashRegisters(
		boilmodels.CashRegisterWhere.IsClose.EQ(false),
		boilmodels.CashRegisterWhere.PointSaleID.EQ(pointSaleID),
		qm.OrderBy("hour_open DESC"),
	).One(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Caja", schemas.Read)
	}

	_, err = boilmodels.Members(
		boilmodels.MemberWhere.ID.EQ(userID),
		qm.Load(boilmodels.MemberRels.Role),
	).One(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Caja", schemas.Read)
	}

	now := time.Now().UTC()
	closeAmtDec := types.NewNullDecimal(types.NewDecimal(fmt.Sprintf("%.2f", amountOpen.CloseAmount)))

	register.CloseAmount = closeAmtDec
	register.IsClose = true
	register.HourClose = null.TimeFrom(now)
	register.MemberCloseID = null.Int64From(userID)

	if _, err := register.Update(ctx, tx, boil.Whitelist(
		boilmodels.CashRegisterColumns.CloseAmount,
		boilmodels.CashRegisterColumns.IsClose,
		boilmodels.CashRegisterColumns.HourClose,
		boilmodels.CashRegisterColumns.MemberCloseID,
		boilmodels.CashRegisterColumns.UpdatedAt,
	)); err != nil {
		return schemas.HandlerErrorDB(err, "Caja", schemas.Update)
	}

	return tx.Commit()
}

type total struct {
	Cash  float64 `boil:"cash"`
	Other float64 `boil:"other"`
}

func (r *CashRegisterRepository) CashRegisterInform(pointSaleID int64, userID int64, fromDate, toDate time.Time) ([]*schemas.CashRegisterInformResponse, error) {
	ctx := context.Background()

	registers, err := boilmodels.CashRegisters(
		boilmodels.CashRegisterWhere.PointSaleID.EQ(pointSaleID),
		boilmodels.CashRegisterWhere.CreatedAt.GTE(null.TimeFrom(fromDate)),
		boilmodels.CashRegisterWhere.CreatedAt.LTE(null.TimeFrom(toDate)),
		qm.Load(boilmodels.CashRegisterRels.MemberOpen),
		qm.Load(boilmodels.CashRegisterRels.MemberClose),
		qm.OrderBy("created_at DESC"),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
	}

	var results []*schemas.CashRegisterInformResponse
	for _, register := range registers {
		var openAmount float64
		if register.OpenAmount.Valid {
			openAmount, _ = register.OpenAmount.Big.Float64()
		}

		res := &schemas.CashRegisterInformResponse{
			ID:         register.ID,
			OpenAmount: openAmount,
			IsClose:    register.IsClose,
		}
		if register.CreatedAt.Valid {
			res.CreatedAt = register.CreatedAt.Time
		}
		if register.HourOpen.Valid {
			res.HourOpen = register.HourOpen.Time
		}
		if register.CloseAmount.Valid {
			cAmount, _ := register.CloseAmount.Big.Float64()
			res.CloseAmount = &cAmount
		}
		if register.HourClose.Valid {
			hClose := register.HourClose.Time
			res.HourClose = &hClose
		}

		if register.R != nil {
			if register.R.MemberOpen != nil {
				res.MemberOpen = *mapToMemberSimpleDTO(register.R.MemberOpen)
			}
			if register.R.MemberClose != nil {
				res.MemberClose = mapToMemberSimpleDTO(register.R.MemberClose)
			}
		}

		// Calculate aggregations
		var incomes total
		err = queries.Raw(`
			SELECT 
			COALESCE(SUM(CASE WHEN method_pay = 'cash' THEN total ELSE 0 END), 0) AS cash,
			COALESCE(SUM(CASE WHEN method_pay <> 'cash' AND method_pay <> 'credit' THEN total ELSE 0 END), 0) AS other
			FROM pay_incomes WHERE cash_register_id = $1 AND delete_at IS NULL
		`, register.ID).Bind(ctx, r.DB, &incomes)

		if err != nil {
			return nil, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
		}

		var incomeOther total
		err = queries.Raw(`
			SELECT 
			COALESCE(SUM(CASE WHEN method_income = 'cash' THEN total ELSE 0 END), 0) AS cash,
			COALESCE(SUM(CASE WHEN method_income <> 'cash' AND method_income <> 'credit' THEN total ELSE 0 END), 0) AS other
			FROM income_others WHERE cash_register_id = $1 AND delete_at IS NULL
		`, register.ID).Bind(ctx, r.DB, &incomeOther)

		if err != nil {
			return nil, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
		}

		var expenseBuy total
		err = queries.Raw(`
			SELECT 
			COALESCE(SUM(CASE WHEN method_pay = 'cash' THEN total ELSE 0 END), 0) AS cash,
			COALESCE(SUM(CASE WHEN method_pay <> 'cash' AND method_pay <> 'credit' THEN total ELSE 0 END), 0) AS other
			FROM pay_expense_buys WHERE cash_register_id = $1 AND delete_at IS NULL
		`, register.ID).Bind(ctx, r.DB, &expenseBuy)

		if err != nil {
			return nil, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
		}

		var expenseOther total
		err = queries.Raw(`
			SELECT 
			COALESCE(SUM(CASE WHEN pay_method = 'cash' THEN total ELSE 0 END), 0) AS cash,
			COALESCE(SUM(CASE WHEN pay_method <> 'cash' AND pay_method <> 'credit' THEN total ELSE 0 END), 0) AS other
			FROM expense_others WHERE cash_register_id = $1 AND delete_at IS NULL
		`, register.ID).Bind(ctx, r.DB, &expenseOther)

		if err != nil {
			return nil, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
		}

		totalIncomesCash := incomes.Cash + incomeOther.Cash
		totalIncomesOthers := incomes.Other + incomeOther.Other
		totalExpenseCash := expenseBuy.Cash + expenseOther.Cash
		totalExpenseOther := expenseBuy.Other + expenseOther.Other

		res.TotalIncomeCash = &totalIncomesCash
		res.TotalIncomeOthers = &totalIncomesOthers
		res.TotalExpenseCash = &totalExpenseCash
		res.TotalExpenseOthers = &totalExpenseOther

		results = append(results, res)
	}

	return results, nil
}
