package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
)

func (r *CashRegisterRepository) CashRegisterExistOpen(pointSaleID int64) (bool, error) {
	ctx := context.Background()
	exists, err := tenant.CashRegisters(
		tenant.CashRegisterWhere.IsClose.EQ(false),
		tenant.CashRegisterWhere.PointSaleID.EQ(pointSaleID),
	).Exists(ctx, r.DB)
	if err != nil {
		return false, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
	}
	return exists, nil
}

func (r *CashRegisterRepository) CashRegisterGetByID(pointSaleID, id int64) (*schemas.CashRegisterFullResponse, error) {
	ctx := context.Background()

	register, err := tenant.CashRegisters(
		tenant.CashRegisterWhere.ID.EQ(id),
		tenant.CashRegisterWhere.PointSaleID.EQ(pointSaleID),
		qm.Load(tenant.CashRegisterRels.MemberOpen),
		qm.Load(tenant.CashRegisterRels.MemberClose),
	).One(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
	}

	openAmount, _ := register.OpenAmount.Big.Float64()

	res := &schemas.CashRegisterFullResponse{
		ID:         register.ID,
		IsClose:    register.IsClose,
		OpenAmount: openAmount,
	}

	res.CreatedAt = register.CreatedAt
	res.HourOpen = register.HourOpen

	if register.CloseAmount.Big == nil {
		cAmount, _ := register.CloseAmount.Big.Float64()
		res.CloseAmount = &cAmount
	}
	if register.HourClose.Valid {
		hClose := register.HourClose.Time
		res.HourClose = &hClose
	}

	if register.R != nil {
		if register.R.MemberOpen != nil {
			m := register.R.MemberOpen
			res.MemberOpen.ID = m.ID
			res.MemberOpen.FirstName = m.FirstName
			res.MemberOpen.LastName = m.LastName
			res.MemberOpen.Username = m.Username
		}
		if register.R.MemberClose != nil {
			m := register.R.MemberClose
			res.MemberClose.ID = m.ID
			res.MemberClose.FirstName = m.FirstName
			res.MemberClose.LastName = m.LastName
			res.MemberClose.Username = m.Username
		}
	}

	// 1. IncomeSales
	incomeSales, err := tenant.IncomeSales(
		tenant.IncomeSaleWhere.CashRegisterID.EQ(id),
		tenant.IncomeSaleWhere.PointSaleID.EQ(pointSaleID),
		qm.Load(tenant.IncomeSaleRels.IncomeSaleItems),
		qm.Load(qm.Rels(tenant.IncomeSaleRels.IncomeSaleItems, tenant.IncomeSaleItemRels.Product)),
		qm.Load(tenant.IncomeSaleRels.PayIncomes),
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
		iso.CreatedAt = in.CreatedAt
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
				iResp.CreatedAt = item.CreatedAt
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
				pResp.CreatedAt = pay.CreatedAt
				iso.Pay = append(iso.Pay, pResp)
			}
		}

		incomes = append(incomes, iso)
	}

	// 2. IncomeOthers
	incomeOthers, err := tenant.IncomeOthers(
		tenant.IncomeOtherWhere.CashRegisterID.EQ(null.Int64From(id)),
		tenant.IncomeOtherWhere.PointSaleID.EQ(null.Int64From(pointSaleID)),
		qm.Load(tenant.IncomeOtherRels.Member),
		qm.Load(tenant.IncomeOtherRels.TypeIncome),
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
		or.CreatedAt = o.CreatedAt
		if o.R != nil {
			if o.R.Member != nil {
				or.Member.ID = o.R.Member.ID
				or.Member.FirstName = o.R.Member.FirstName
				or.Member.LastName = o.R.Member.LastName
				or.Member.Username = o.R.Member.Username
			}
			if o.R.TypeIncome != nil {
				or.TypeIncome.ID = o.R.TypeIncome.ID
				or.TypeIncome.Name = o.R.TypeIncome.Name
			}
		}
		others = append(others, or)
	}

	// 3. ExpenseOthers
	expenseOthers, err := tenant.ExpenseOthers(
		tenant.ExpenseOtherWhere.CashRegisterID.EQ(null.Int64From(id)),
		tenant.ExpenseOtherWhere.PointSaleID.EQ(null.Int64From(pointSaleID)),
		qm.Load(tenant.ExpenseOtherRels.Member),
		qm.Load(tenant.ExpenseOtherRels.TypeExpense),
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
		eo.CreatedAt = e.CreatedAt
		if e.R != nil {
			if e.R.Member != nil {
				eo.Member.ID = e.R.Member.ID
				eo.Member.FirstName = e.R.Member.FirstName
				eo.Member.LastName = e.R.Member.LastName
				eo.Member.Username = e.R.Member.Username
			}
			if e.R.TypeExpense != nil {
				eo.TypeExpense.ID = e.R.TypeExpense.ID
				eo.TypeExpense.Name = e.R.TypeExpense.Name
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

	exists, err := tenant.CashRegisters(
		tenant.CashRegisterWhere.IsClose.EQ(false),
		tenant.CashRegisterWhere.PointSaleID.EQ(pointSaleID),
	).Exists(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Caja", schemas.Read)
	}

	if exists {
		return schemas.ErrorResponse(400, "ya existe una apertura de caja, antes de continuar cierre la caja", fmt.Errorf("ya existe una apertura de caja, antes de continuar cerrar"))
	}

	// openAmtDec := types.NewNullDecimal(types.NewDecimal(fmt.Sprintf("%.2f", amountOpen.OpenAmount)))
	openAmtDec := types.NewDecimal(decimal.New(0, 0).SetFloat64(amountOpen.OpenAmount))
	registerOpen := tenant.CashRegister{
		PointSaleID:  pointSaleID,
		MemberOpenID: userID,
		OpenAmount:   openAmtDec,
		HourOpen:     time.Now().UTC(),
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

	register, err := tenant.CashRegisters(
		tenant.CashRegisterWhere.IsClose.EQ(false),
		tenant.CashRegisterWhere.PointSaleID.EQ(pointSaleID),
		qm.OrderBy("hour_open DESC"),
	).One(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Caja", schemas.Read)
	}

	_, err = tenant.Members(
		tenant.MemberWhere.ID.EQ(userID),
		qm.Load(tenant.MemberRels.Role),
	).One(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Caja", schemas.Read)
	}

	now := time.Now().UTC()
	closeAmtDec := types.NewNullDecimal(decimal.New(0, 0).SetFloat64(amountOpen.CloseAmount))

	register.CloseAmount = closeAmtDec
	register.IsClose = true
	register.HourClose = null.TimeFrom(now)
	register.MemberCloseID = null.Int64From(userID)

	if _, err := register.Update(ctx, tx, boil.Whitelist(
		tenant.CashRegisterColumns.CloseAmount,
		tenant.CashRegisterColumns.IsClose,
		tenant.CashRegisterColumns.HourClose,
		tenant.CashRegisterColumns.MemberCloseID,
		tenant.CashRegisterColumns.UpdatedAt,
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

	registers, err := tenant.CashRegisters(
		tenant.CashRegisterWhere.PointSaleID.EQ(pointSaleID),
		tenant.CashRegisterWhere.CreatedAt.GTE(fromDate),
		tenant.CashRegisterWhere.CreatedAt.LTE(toDate),
		qm.Load(tenant.CashRegisterRels.MemberOpen),
		qm.Load(tenant.CashRegisterRels.MemberClose),
		qm.OrderBy("created_at DESC"),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
	}

	var results []*schemas.CashRegisterInformResponse
	for _, register := range registers {
		var openAmount float64
		openAmount, _ = register.OpenAmount.Big.Float64()

		res := &schemas.CashRegisterInformResponse{
			ID:         register.ID,
			OpenAmount: openAmount,
			IsClose:    register.IsClose,
		}
		res.CreatedAt = register.CreatedAt
		res.HourOpen = register.HourOpen
		if register.CloseAmount.Big != nil {
			cAmount, _ := register.CloseAmount.Big.Float64()
			res.CloseAmount = &cAmount
		}
		if register.HourClose.Valid {
			hClose := register.HourClose.Time
			res.HourClose = &hClose
		}

		if register.R != nil {
			if register.R.MemberOpen != nil {
				res.MemberOpen.ID = register.R.MemberOpen.ID
				res.MemberOpen.FirstName = register.R.MemberOpen.FirstName
				res.MemberOpen.LastName = register.R.MemberOpen.LastName
				res.MemberOpen.Username = register.R.MemberOpen.Username
			}
			if register.R.MemberClose != nil {
				res.MemberClose.ID = register.R.MemberClose.ID
				res.MemberClose.FirstName = register.R.MemberClose.FirstName
				res.MemberClose.LastName = register.R.MemberClose.LastName
				res.MemberClose.Username = register.R.MemberClose.Username
			}
		}

		// Calculate aggregations
		var incomes total
		err = queries.Raw(`
			SELECT
			COALESCE(SUM(CASE WHEN method_pay = 'cash' THEN total ELSE 0 END), 0)::float8 AS cash,
			COALESCE(SUM(CASE WHEN method_pay <> 'cash' AND method_pay <> 'credit' THEN total ELSE 0 END), 0)::float8 AS other
			FROM pay_incomes WHERE cash_register_id = $1
		`, register.ID).Bind(ctx, r.DB, &incomes)

		if err != nil {
			return nil, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
		}

		var incomeOther total
		err = queries.Raw(`
			SELECT
			COALESCE(SUM(CASE WHEN method_income = 'cash' THEN total ELSE 0 END), 0)::float8 AS cash,
			COALESCE(SUM(CASE WHEN method_income <> 'cash' AND method_income <> 'credit' THEN total ELSE 0 END), 0)::float8 AS other
			FROM income_others WHERE cash_register_id = $1
		`, register.ID).Bind(ctx, r.DB, &incomeOther)

		if err != nil {
			return nil, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
		}

		var expenseBuy total
		err = queries.Raw(`
			SELECT
			COALESCE(SUM(CASE WHEN method_pay = 'cash' THEN total ELSE 0 END), 0)::float8 AS cash,
			COALESCE(SUM(CASE WHEN method_pay <> 'cash' AND method_pay <> 'credit' THEN total ELSE 0 END), 0)::float8 AS other
			FROM pay_expense_buys WHERE cash_register_id = $1
		`, register.ID).Bind(ctx, r.DB, &expenseBuy)

		if err != nil {
			return nil, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
		}

		var expenseOther total
		err = queries.Raw(`
			SELECT
			COALESCE(SUM(CASE WHEN pay_method = 'cash' THEN total ELSE 0 END), 0)::float8 AS cash,
			COALESCE(SUM(CASE WHEN pay_method <> 'cash' AND pay_method <> 'credit' THEN total ELSE 0 END), 0)::float8 AS other
			FROM expense_others WHERE cash_register_id = $1
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
