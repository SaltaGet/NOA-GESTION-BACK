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

func mapToTypeIncomeResponse(c *boilmodels.TypeIncome) schemas.TypeIncomeResponse {
	if c == nil {
		return schemas.TypeIncomeResponse{}
	}
	return schemas.TypeIncomeResponse{
		ID:   c.ID,
		Name: c.Name,
	}
}

func mapToIncomeOtherResponse(i *boilmodels.IncomeOther) *schemas.IncomeOtherResponse {
	if i == nil {
		return nil
	}

	total, _ := i.Total.Big.Float64()

	res := &schemas.IncomeOtherResponse{
		ID:           i.ID,
		Total:        total,
		MethodIncome: i.MethodIncome,
	}

	if i.Details.Valid {
		res.Details = &i.Details.String
	}
	if i.CashRegisterID.Valid {
		res.CashRegisterID = &i.CashRegisterID.Int64
	}
	if i.CreatedAt.Valid {
		res.CreatedAt = i.CreatedAt.Time
	}

	if i.R != nil {
		if i.R.Member != nil {
			res.Member = mapToMemberSimpleDTO(i.R.Member)
		}
		if i.R.TypeIncome != nil {
			res.TypeIncome = mapToTypeIncomeResponse(i.R.TypeIncome)
		}
		if i.R.PointSale != nil {
			p := schemas.PointSaleResponse{
				ID:        i.R.PointSale.ID,
				Name:      i.R.PointSale.Name,
				Number:    int64(i.R.PointSale.Number),
				IsDeposit: i.R.PointSale.IsDeposit,
				IsMain:    i.R.PointSale.IsMain,
			}
			if i.R.PointSale.Description.Valid {
				p.Description = &i.R.PointSale.Description.String
			}
			res.PointSale = &p
		}
	}

	return res
}

func (r *IncomeOtherRepository) IncomeOtherGetByID(id int64, pointSaleId *int64) (*schemas.IncomeOtherResponse, error) {
	ctx := context.Background()

	var qms []qm.QueryMod
	qms = append(qms,
		boilmodels.IncomeOtherWhere.ID.EQ(id),
		qm.Load(boilmodels.IncomeOtherRels.Member),
		qm.Load(boilmodels.IncomeOtherRels.TypeIncome),
	)

	if pointSaleId != nil {
		qms = append(qms,
			boilmodels.IncomeOtherWhere.PointSaleID.EQ(null.Int64From(*pointSaleId)),
			qm.Load(boilmodels.IncomeOtherRels.PointSale),
		)
	}

	incomeOther, err := boilmodels.IncomeOthers(qms...).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Otros ingresos", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Otros ingresos", schemas.Read)
	}

	return mapToIncomeOtherResponse(incomeOther), nil
}

func (r *IncomeOtherRepository) IncomeOtherGetByDate(pointSaleID *int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.IncomeOtherResponse, int64, error) {
	ctx := context.Background()
	offset := (page - 1) * limit

	var qms []qm.QueryMod
	qms = append(qms,
		boilmodels.IncomeOtherWhere.CreatedAt.GTE(null.TimeFrom(fromDate)),
		boilmodels.IncomeOtherWhere.CreatedAt.LTE(null.TimeFrom(toDate)),
	)

	if pointSaleID != nil {
		qms = append(qms, boilmodels.IncomeOtherWhere.PointSaleID.EQ(null.Int64From(*pointSaleID)))
	}

	total, err := boilmodels.IncomeOthers(qms...).Count(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Otros ingresos", schemas.Read)
	}

	qms = append(qms,
		qm.Load(boilmodels.IncomeOtherRels.Member),
		qm.Load(boilmodels.IncomeOtherRels.TypeIncome),
		qm.OrderBy("created_at DESC"),
		qm.Offset(offset),
		qm.Limit(limit),
	)

	if pointSaleID == nil {
		qms = append(qms, qm.Load(boilmodels.IncomeOtherRels.PointSale))
	}

	incomesOther, err := boilmodels.IncomeOthers(qms...).All(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Otros ingresos", schemas.Read)
	}

	var results []*schemas.IncomeOtherResponse
	for _, i := range incomesOther {
		results = append(results, mapToIncomeOtherResponse(i))
	}

	return results, total, nil
}

func (r *IncomeOtherRepository) IncomeOtherCreate(memberID int64, pointSaleID *int64, incomeOtherCreate *schemas.IncomeOtherCreate) (int64, error) {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return 0, err
	}

	typeIncomeExists, err := boilmodels.TypeIncomes(boilmodels.TypeIncomeWhere.ID.EQ(incomeOtherCreate.TypeIncomeID)).Exists(ctx, tx)
	if err != nil || !typeIncomeExists {
		return 0, schemas.HandlerErrorDB(sql.ErrNoRows, "Tipo de ingreso", schemas.Read)
	}

	totalDec := types.NewNullDecimal(types.NewDecimal(fmt.Sprintf("%.4f", incomeOtherCreate.Total)))

	incomeOther := &boilmodels.IncomeOther{
		MemberID:     null.Int64From(memberID),
		Total:        totalDec,
		TypeIncomeID: incomeOtherCreate.TypeIncomeID,
		Details:      null.StringFromPtr(incomeOtherCreate.Details),
		MethodIncome: incomeOtherCreate.MethodIncome,
	}

	if pointSaleID == nil {
		if err := incomeOther.Insert(ctx, tx, boil.Infer()); err != nil {
			return 0, schemas.HandlerErrorDB(err, "Otros ingresos", schemas.Create)
		}
		if err := tx.Commit(); err != nil {
			return 0, err
		}
		return incomeOther.ID, nil
	}

	incomeOther.PointSaleID = null.Int64From(*pointSaleID)

	register, err := boilmodels.CashRegisters(
		boilmodels.CashRegisterWhere.IsClose.EQ(false),
		boilmodels.CashRegisterWhere.PointSaleID.EQ(*pointSaleID),
		qm.OrderBy("hour_open DESC"),
	).One(ctx, tx)

	if err != nil {
		return 0, schemas.HandlerErrorDB(err, "Apertura de caja", schemas.Read)
	}

	incomeOther.CashRegisterID = null.Int64From(register.ID)

	if err := incomeOther.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Otros ingresos", schemas.Create)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return incomeOther.ID, nil
}

func (r *IncomeOtherRepository) IncomeOtherUpdate(memberID int64, pointSaleID *int64, incomeOtherUpdate *schemas.IncomeOtherUpdate) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	var qms []qm.QueryMod
	qms = append(qms, boilmodels.IncomeOtherWhere.ID.EQ(incomeOtherUpdate.ID))
	if pointSaleID != nil {
		qms = append(qms, boilmodels.IncomeOtherWhere.PointSaleID.EQ(null.Int64From(*pointSaleID)))
	}

	existingIncome, err := boilmodels.IncomeOthers(qms...).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Otros ingresos", schemas.Read)
	}

	typeIncomeExists, err := boilmodels.TypeIncomes(boilmodels.TypeIncomeWhere.ID.EQ(incomeOtherUpdate.TypeIncomeID)).Exists(ctx, tx)
	if err != nil || !typeIncomeExists {
		return schemas.HandlerErrorDB(sql.ErrNoRows, "Tipo de ingreso", schemas.Read)
	}

	totalDec := types.NewNullDecimal(types.NewDecimal(fmt.Sprintf("%.4f", incomeOtherUpdate.Total)))

	existingIncome.Total = totalDec
	existingIncome.TypeIncomeID = incomeOtherUpdate.TypeIncomeID
	existingIncome.Details = null.StringFromPtr(incomeOtherUpdate.Details)
	existingIncome.MethodIncome = incomeOtherUpdate.MethodIncome

	if _, err := existingIncome.Update(ctx, tx, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Otros ingresos", schemas.Update)
	}

	return tx.Commit()
}

func (r *IncomeOtherRepository) IncomeOtherDelete(memberID int64, incomeOtherID int64, pointSaleID *int64) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	var qms []qm.QueryMod
	qms = append(qms, boilmodels.IncomeOtherWhere.ID.EQ(incomeOtherID))
	if pointSaleID != nil {
		qms = append(qms, boilmodels.IncomeOtherWhere.PointSaleID.EQ(null.Int64From(*pointSaleID)))
	}

	existingIncome, err := boilmodels.IncomeOthers(qms...).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Otros ingresos", schemas.Read)
	}

	if _, err := existingIncome.Delete(ctx, tx, false); err != nil {
		return schemas.HandlerErrorDB(err, "Otros ingresos", schemas.Delete)
	}

	return tx.Commit()
}
