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

func mapToPointSaleResponse(c *boilmodels.PointSale) *schemas.PointSaleResponse {
	if c == nil {
		return nil
	}
	res := &schemas.PointSaleResponse{
		ID:        c.ID,
		Name:      c.Name,
		Number:    int64(c.Number),
		IsDeposit: c.IsDeposit,
		IsMain:    c.IsMain,
	}
	if c.Description.Valid {
		res.Description = &c.Description.String
	}
	return res
}

func mapToTypeExpenseResponse(c *boilmodels.TypeExpense) schemas.TypeExpenseResponse {
	if c == nil {
		return schemas.TypeExpenseResponse{}
	}
	return schemas.TypeExpenseResponse{
		ID:   c.ID,
		Name: c.Name,
	}
}

func mapToExpenseOtherResponse(e *boilmodels.ExpenseOther) *schemas.ExpenseOtherResponse {
	if e == nil {
		return nil
	}

	total, _ := e.Total.Big.Float64()

	res := &schemas.ExpenseOtherResponse{
		ID:        e.ID,
		Total:     total,
		PayMethod: e.PayMethod.String,
	}

	if e.Details.Valid {
		res.Details = &e.Details.String
	}
	if e.CashRegisterID.Valid {
		res.CashRegisterID = &e.CashRegisterID.Int64
	}
	if e.CreatedAt.Valid {
		res.CreatedAt = e.CreatedAt.Time
	}

	if e.R != nil {
		if e.R.Member != nil {
			res.Member = mapToMemberSimpleDTO(e.R.Member)
		}
		if e.R.PointSale != nil {
			res.PointSale = mapToPointSaleResponse(e.R.PointSale)
		}
		if e.R.TypeExpense != nil {
			res.TypeExpense = mapToTypeExpenseResponse(e.R.TypeExpense)
		}
	}

	return res
}

func (r *ExpenseOtherRepository) ExpenseOtherGetByID(id int64, pointSaleID *int64) (*schemas.ExpenseOtherResponse, error) {
	ctx := context.Background()

	var qms []qm.QueryMod
	qms = append(qms,
		boilmodels.ExpenseOtherWhere.ID.EQ(id),
		qm.Load(boilmodels.ExpenseOtherRels.Member),
		qm.Load(boilmodels.ExpenseOtherRels.TypeExpense),
		qm.Load(boilmodels.ExpenseOtherRels.PointSale),
	)

	if pointSaleID != nil {
		qms = append(qms, boilmodels.ExpenseOtherWhere.PointSaleID.EQ(null.Int64From(*pointSaleID)))
	}

	expenseOther, err := boilmodels.ExpenseOthers(qms...).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Otros egresos", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Otros egresos", schemas.Read)
	}

	return mapToExpenseOtherResponse(expenseOther), nil
}

func (r *ExpenseOtherRepository) ExpenseOtherGetByDate(pointSaleID *int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.ExpenseOtherResponse, int64, error) {
	ctx := context.Background()
	offset := (page - 1) * limit

	qms := []qm.QueryMod{
		boilmodels.ExpenseOtherWhere.CreatedAt.GTE(null.TimeFrom(fromDate)),
		boilmodels.ExpenseOtherWhere.CreatedAt.LTE(null.TimeFrom(toDate)),
	}

	if pointSaleID != nil {
		qms = append(qms, boilmodels.ExpenseOtherWhere.PointSaleID.EQ(null.Int64From(*pointSaleID)))
	}

	total, err := boilmodels.ExpenseOthers(qms...).Count(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Otros egresos", schemas.Read)
	}

	qms = append(qms,
		qm.Load(boilmodels.ExpenseOtherRels.Member),
		qm.Load(boilmodels.ExpenseOtherRels.TypeExpense),
		qm.Load(boilmodels.ExpenseOtherRels.PointSale),
		qm.OrderBy("created_at DESC"),
		qm.Offset(offset),
		qm.Limit(limit),
	)

	expensesOther, err := boilmodels.ExpenseOthers(qms...).All(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Otros egresos", schemas.Read)
	}

	var expenseSchema []*schemas.ExpenseOtherResponse
	for _, e := range expensesOther {
		expenseSchema = append(expenseSchema, mapToExpenseOtherResponse(e))
	}

	return expenseSchema, total, nil
}

func (r *ExpenseOtherRepository) ExpenseOtherCreate(memberID int64, pointSaleID *int64, expenseOtherCreate *schemas.ExpenseOtherCreate) (int64, error) {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return 0, err
	}

	typeExpenseExists, err := boilmodels.TypeExpenses(boilmodels.TypeExpenseWhere.ID.EQ(expenseOtherCreate.TypeExpenseID)).Exists(ctx, tx)
	if err != nil || !typeExpenseExists {
		return 0, schemas.HandlerErrorDB(sql.ErrNoRows, "Tipo de egreso", schemas.Read)
	}

	totalDec := types.NewNullDecimal(types.NewDecimal(fmt.Sprintf("%.4f", expenseOtherCreate.Total)))

	expenseOther := &boilmodels.ExpenseOther{
		MemberID:      null.Int64From(memberID),
		Details:       null.StringFromPtr(expenseOtherCreate.Details),
		TypeExpenseID: null.Int64From(expenseOtherCreate.TypeExpenseID),
		Total:         totalDec,
		PayMethod:     null.StringFrom(expenseOtherCreate.PayMethod),
	}

	if pointSaleID != nil {
		register, err := boilmodels.CashRegisters(
			boilmodels.CashRegisterWhere.IsClose.EQ(false),
			boilmodels.CashRegisterWhere.PointSaleID.EQ(*pointSaleID),
			qm.OrderBy("hour_open DESC"),
		).One(ctx, tx)

		if err != nil {
			return 0, schemas.HandlerErrorDB(err, "Caja", schemas.Read)
		}

		expenseOther.PointSaleID = null.Int64From(*pointSaleID)
		expenseOther.CashRegisterID = null.Int64From(register.ID)
	}

	if err := expenseOther.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Otros egresos", schemas.Create)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return expenseOther.ID, nil
}

func (r *ExpenseOtherRepository) ExpenseOtherUpdate(memberID int64, pointSaleID *int64, expenseOtherUpdate *schemas.ExpenseOtherUpdate) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	qms := []qm.QueryMod{boilmodels.ExpenseOtherWhere.ID.EQ(expenseOtherUpdate.ID)}
	if pointSaleID != nil {
		qms = append(qms, boilmodels.ExpenseOtherWhere.PointSaleID.EQ(null.Int64From(*pointSaleID)))
	}

	existingExpense, err := boilmodels.ExpenseOthers(qms...).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Otros egresos", schemas.Read)
	}

	typeExpenseExists, err := boilmodels.TypeExpenses(boilmodels.TypeExpenseWhere.ID.EQ(expenseOtherUpdate.TypeExpenseID)).Exists(ctx, tx)
	if err != nil || !typeExpenseExists {
		return schemas.HandlerErrorDB(sql.ErrNoRows, "Tipo de egreso", schemas.Read)
	}

	existingExpense.Details = null.StringFromPtr(expenseOtherUpdate.Details)
	existingExpense.Total = types.NewNullDecimal(types.NewDecimal(fmt.Sprintf("%.4f", expenseOtherUpdate.Total)))
	existingExpense.PayMethod = null.StringFrom(expenseOtherUpdate.PayMethod)
	existingExpense.MemberID = null.Int64From(memberID)
	existingExpense.TypeExpenseID = null.Int64From(expenseOtherUpdate.TypeExpenseID)

	if _, err := existingExpense.Update(ctx, tx, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Otros egresos", schemas.Update)
	}

	return tx.Commit()
}

func (r *ExpenseOtherRepository) ExpenseOtherDelete(memberID int64, expenseOtherID int64, pointSaleID *int64) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	qms := []qm.QueryMod{boilmodels.ExpenseOtherWhere.ID.EQ(expenseOtherID)}
	if pointSaleID != nil {
		qms = append(qms, boilmodels.ExpenseOtherWhere.PointSaleID.EQ(null.Int64From(*pointSaleID)))
	}

	existingExpense, err := boilmodels.ExpenseOthers(qms...).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Otros egresos", schemas.Read)
	}

	if _, err := existingExpense.Delete(ctx, tx, false); err != nil {
		return schemas.HandlerErrorDB(err, "Otros egresos", schemas.Delete)
	}

	return tx.Commit()
}
