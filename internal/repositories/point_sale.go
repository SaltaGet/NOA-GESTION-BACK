package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	boiltenant "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

func mapToPointSaleResponse(p *boiltenant.PointSale) schemas.PointSaleResponse {
	if p == nil {
		return schemas.PointSaleResponse{}
	}

	res := schemas.PointSaleResponse{
		ID:        p.ID,
		Name:      p.Name,
		Number:    int64(p.Number),
		IsDeposit: p.IsDeposit,
		IsMain:    p.IsMain,
	}

	if p.Description.Valid {
		res.Description = &p.Description.String
	}

	return res
}

func (p *PointSaleRepository) PointSaleGetAllByMember(memberID int64) ([]schemas.PointSaleResponse, error) {
	ctx := context.Background()

	pointSales, err := boiltenant.PointSales(
		qm.Select(
			"point_sales."+boiltenant.PointSaleColumns.ID,
			"point_sales."+boiltenant.PointSaleColumns.Name,
			"point_sales."+boiltenant.PointSaleColumns.Description,
			"point_sales."+boiltenant.PointSaleColumns.IsDeposit,
			"point_sales."+boiltenant.PointSaleColumns.IsMain,
			"point_sales."+boiltenant.PointSaleColumns.Number,
		),
		qm.InnerJoin("member_point_sales mp ON mp.point_sale_id = point_sales.id"),
		qm.Where("mp.member_id = ?", memberID),
	).All(ctx, p.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Punto de venta", schemas.Read)
	}

	var response []schemas.PointSaleResponse
	for _, ps := range pointSales {
		response = append(response, mapToPointSaleResponse(ps))
	}

	return response, nil
}

func (p *PointSaleRepository) PointSaleGetAll() ([]schemas.PointSaleResponse, error) {
	ctx := context.Background()

	pointSales, err := boilmodels.PointSales(
		qm.Select(
			boilmodels.PointSaleColumns.ID,
			boilmodels.PointSaleColumns.Name,
			boilmodels.PointSaleColumns.Description,
			boilmodels.PointSaleColumns.IsDeposit,
			boilmodels.PointSaleColumns.IsMain,
			boilmodels.PointSaleColumns.Number,
		),
	).All(ctx, p.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Punto de venta", schemas.Read)
	}

	var response []schemas.PointSaleResponse
	for _, ps := range pointSales {
		response = append(response, mapToPointSaleResponse(ps))
	}

	return response, nil
}

func (p *PointSaleRepository) PointSaleGetByID(id int64) (*schemas.PointSaleResponse, error) {
	ctx := context.Background()

	ps, err := boilmodels.PointSales(
		qm.Select(
			boilmodels.PointSaleColumns.ID,
			boilmodels.PointSaleColumns.Name,
			boilmodels.PointSaleColumns.Description,
			boilmodels.PointSaleColumns.IsDeposit,
			boilmodels.PointSaleColumns.IsMain,
			boilmodels.PointSaleColumns.Number,
		),
		boilmodels.PointSaleWhere.ID.EQ(id),
	).One(ctx, p.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Punto de venta", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Punto de venta", schemas.Read)
	}

	res := mapToPointSaleResponse(ps)
	return &res, nil
}

func (p *PointSaleRepository) PointSaleCount() (int64, error) {
	ctx := context.Background()

	count, err := boilmodels.PointSales().Count(ctx, p.DB)
	if err != nil {
		return 0, schemas.HandlerErrorDB(err, "Punto de venta", schemas.Read)
	}

	return count, nil
}

func (p *PointSaleRepository) PointSaleCreate(memberID int64, pointSaleCreate *schemas.PointSaleCreate) (int64, error) {
	ctx := context.Background()
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return 0, err
	}

	hasMain, err := boilmodels.PointSales(boilmodels.PointSaleWhere.IsMain.EQ(true)).Exists(ctx, tx)
	if err != nil {
		return 0, schemas.HandlerErrorDB(err, "Punto de venta", schemas.Read)
	}

	ps := &boilmodels.PointSale{
		Name:      pointSaleCreate.Name,
		IsDeposit: *pointSaleCreate.IsDeposit,
		Number:    int64(pointSaleCreate.Number),
		IsMain:    !hasMain,
	}

	if pointSaleCreate.Description != nil {
		ps.Description = null.StringFromPtr(pointSaleCreate.Description)
	}

	if err := ps.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Punto de venta", schemas.Create)
	}

	membersAdmin, err := boilmodels.Members(boilmodels.MemberWhere.IsAdmin.EQ(true)).All(ctx, tx)
	if err != nil {
		return 0, schemas.HandlerErrorDB(err, "Punto de venta", schemas.Read)
	}

	if len(membersAdmin) > 0 {
		if err := ps.AddMembers(ctx, tx, false, membersAdmin...); err != nil {
			return 0, schemas.HandlerErrorDB(err, "Punto de venta", schemas.Update)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return ps.ID, nil
}

func (p *PointSaleRepository) PointSaleUpdate(memberID int64, pointSaleUpdate *schemas.PointSaleUpdate) error {
	ctx := context.Background()
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	ps, err := boilmodels.PointSales(boilmodels.PointSaleWhere.ID.EQ(pointSaleUpdate.ID)).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Punto de venta", schemas.Read)
	}

	ps.Name = pointSaleUpdate.Name
	ps.Number = int64(pointSaleUpdate.Number)
	if pointSaleUpdate.Description != nil {
		ps.Description = null.StringFromPtr(pointSaleUpdate.Description)
	}

	if !ps.IsDeposit && *pointSaleUpdate.IsDeposit {
		stockList, err := boilmodels.StockPointSales(boilmodels.StockPointSaleWhere.PointSaleID.EQ(ps.ID)).All(ctx, tx)
		if err != nil {
			return schemas.HandlerErrorDB(err, "Punto de venta", schemas.Read)
		}

		for _, s := range stockList {
			deposit, err := boilmodels.Deposits(boilmodels.DepositWhere.ProductID.EQ(s.ProductID)).One(ctx, tx)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					deposit = &boilmodels.Deposit{
						ProductID: s.ProductID,
						Stock:     s.Stock,
					}
					if err := deposit.Insert(ctx, tx, boil.Infer()); err != nil {
						return schemas.HandlerErrorDB(err, "Déposito", schemas.Create)
					}
				} else {
					return schemas.HandlerErrorDB(err, "Punto de venta", schemas.Read)
				}
			} else {
				f1, _ := deposit.Stock.Float64()
				f2, _ := s.Stock.Float64()
				deposit.Stock = floatToDecimal(f1 + f2)
				if _, err := deposit.Update(ctx, tx, boil.Whitelist(boilmodels.DepositColumns.Stock, boilmodels.DepositColumns.UpdatedAt)); err != nil {
					return schemas.HandlerErrorDB(err, "Déposito", schemas.Update)
				}
			}

			s.Stock = floatToDecimal(0)
			if _, err := s.Update(ctx, tx, boil.Whitelist(boilmodels.StockPointSaleColumns.Stock, boilmodels.StockPointSaleColumns.UpdatedAt)); err != nil {
				return schemas.HandlerErrorDB(err, "Stock punto de venta", schemas.Update)
			}
		}
	}

	ps.IsDeposit = *pointSaleUpdate.IsDeposit

	if _, err := ps.Update(ctx, tx, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Punto de venta", schemas.Update)
	}

	return tx.Commit()
}

func (p *PointSaleRepository) PointSaleUpdateMain(memberID int64, pointSaleUpdateMain *schemas.PointSaleUpdateMain) error {
	ctx := context.Background()
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	pointSaleOld, err := boilmodels.PointSales(
		qm.For("UPDATE"),
		boilmodels.PointSaleWhere.ID.EQ(pointSaleUpdateMain.ID),
	).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Punto de venta", schemas.Read)
	}

	pointSaleNew, err := boilmodels.PointSales(
		qm.For("UPDATE"),
		boilmodels.PointSaleWhere.ID.EQ(pointSaleUpdateMain.NewMain),
	).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Punto de venta", schemas.Read)
	}

	if !pointSaleOld.IsMain {
		return schemas.ErrorResponse(400, "El punto de venta indicado no es el principal actual", nil)
	}

	if pointSaleNew.IsMain {
		return schemas.ErrorResponse(400, "El nuevo punto de venta ya es el principal", nil)
	}

	pointSaleOld.IsMain = false
	if _, err := pointSaleOld.Update(ctx, tx, boil.Whitelist(boilmodels.PointSaleColumns.IsMain, boilmodels.PointSaleColumns.UpdatedAt)); err != nil {
		return schemas.HandlerErrorDB(err, "Punto de venta", schemas.Update)
	}

	pointSaleNew.IsMain = true
	if _, err := pointSaleNew.Update(ctx, tx, boil.Whitelist(boilmodels.PointSaleColumns.IsMain, boilmodels.PointSaleColumns.UpdatedAt)); err != nil {
		return schemas.HandlerErrorDB(err, "Punto de venta", schemas.Update)
	}

	return tx.Commit()
}
