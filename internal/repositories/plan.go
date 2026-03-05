package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
)

func mapToTenantResponse(t *boilmodels.Tenant) schemas.TenantResponse {
	if t == nil {
		return schemas.TenantResponse{}
	}

	res := schemas.TenantResponse{
		ID:            t.ID,
		Name:          t.Name,
		SchemaName:    t.SchemaName, // Verify field
		Connection:    t.Connection,
		IsActive:      t.IsActive,
		AcceptedTerms: t.AcceptedTerms,     // Keep boolean default mapping straight
		PlanID:        int64(t.PlanID.Int), // Handling nulls
	}

	if t.Expiration.Valid {
		ex := t.Expiration.Time.Format("2006-01-02")
		res.Expiration = &ex
	}

	if t.Phone.Valid {
		res.Phone = &t.Phone.String
	}

	if t.TypeID.Valid {
		res.TypeID = int64(t.TypeID.Int)
	}

	if t.CreatedAt.Valid {
		res.CreatedAt = t.CreatedAt.Time
	}

	return res
}

func mapToPlanResponse(p *boilmodels.Plan) *schemas.PlanResponse {
	if p == nil {
		return nil
	}

	mounthly, _ := p.PriceMounthly.Big.Float64()
	yearly, _ := p.PriceYearly.Big.Float64()

	res := &schemas.PlanResponse{
		ID:              p.ID,
		Name:            p.Name,
		PriceMounthly:   mounthly,
		PriceYearly:     yearly,
		Description:     p.Description.String,
		Features:        p.Features.String,
		AmountPointSale: int64(p.AmountPointSale.Int),
		AmountMember:    int64(p.AmountMember.Int),
		AmountProduct:   int64(p.AmountProduct.Int),
	}

	if p.R != nil && len(p.R.Tenants) > 0 {
		for _, tenant := range p.R.Tenants {
			res.Tenants = append(res.Tenants, mapToTenantResponse(tenant))
		}
	}

	return res
}

func mapToPlanResponseDTO(p *boilmodels.Plan) *schemas.PlanResponseDTO {
	if p == nil {
		return nil
	}

	mounthly, _ := p.PriceMounthly.Big.Float64()
	yearly, _ := p.PriceYearly.Big.Float64()

	res := &schemas.PlanResponseDTO{
		ID:              p.ID,
		Name:            p.Name,
		PriceMounthly:   mounthly,
		PriceYearly:     yearly,
		Description:     p.Description.String,
		Features:        p.Features.String,
		AmountPointSale: int64(p.AmountPointSale.Int),
		AmountMember:    int64(p.AmountMember.Int),
		AmountProduct:   int64(p.AmountProduct.Int),
	}

	return res
}

func (r *MainRepository) PlanCreate(adminID int64, planCreate *schemas.PlanCreate) (int64, error) {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", adminID)); err != nil {
		return 0, err
	}

	plan := &boilmodels.Plan{
		Name:            planCreate.Name,
		PriceMounthly:   types.NewNullDecimal(types.NewDecimal(fmt.Sprintf("%.4f", planCreate.PriceMounthly))),
		PriceYearly:     types.NewNullDecimal(types.NewDecimal(fmt.Sprintf("%.4f", planCreate.PriceYearly))),
		Description:     null.StringFrom(planCreate.Description),
		Features:        null.StringFrom(planCreate.Features),
		AmountPointSale: null.IntFrom(int(planCreate.AmountPointSale)),
		AmountMember:    null.IntFrom(int(planCreate.AmountMember)),
		AmountProduct:   null.IntFrom(int(planCreate.AmountProduct)),
	}

	if err := plan.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Plan", schemas.Create)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return plan.ID, nil
}

func (r *MainRepository) PlanUpdate(adminID int64, planUpdate *schemas.PlanUpdate) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", adminID)); err != nil {
		return err
	}

	plan, err := boilmodels.Plans(boilmodels.PlanWhere.ID.EQ(planUpdate.ID)).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schemas.ErrorResponse(404, "El plan no encopntrado", err)
		}
		return schemas.ErrorResponse(500, "Error al buscar el plan", err)
	}

	plan.Name = planUpdate.Name
	plan.PriceMounthly = types.NewNullDecimal(types.NewDecimal(fmt.Sprintf("%.4f", planUpdate.PriceMounthly)))
	plan.PriceYearly = types.NewNullDecimal(types.NewDecimal(fmt.Sprintf("%.4f", planUpdate.PriceYearly)))
	plan.Description = null.StringFrom(planUpdate.Description)
	plan.Features = null.StringFrom(planUpdate.Features)
	plan.AmountPointSale = null.IntFrom(int(planUpdate.AmountPointSale))
	plan.AmountMember = null.IntFrom(int(planUpdate.AmountMember))
	plan.AmountProduct = null.IntFrom(int(planUpdate.AmountProduct))

	if _, err := plan.Update(ctx, tx, boil.Infer()); err != nil {
		// Unique error simulation check can be added or parsed from db error msg inside schemas.IsDuplicateError
		if schemas.IsDuplicateError(err) {
			return schemas.ErrorResponse(409, "El plan "+plan.Name+" ya existe", err)
		}
		return schemas.ErrorResponse(500, "Error al actualizar el plan", err)
	}

	return tx.Commit()
}

func (r *MainRepository) PlanGetAll() ([]*schemas.PlanResponseDTO, error) {
	ctx := context.Background()

	plans, err := boilmodels.Plans().All(ctx, r.DB)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al obtener los planes", err)
	}

	var plansResponse []*schemas.PlanResponseDTO
	for _, p := range plans {
		plansResponse = append(plansResponse, mapToPlanResponseDTO(p))
	}

	return plansResponse, nil
}

func (r *MainRepository) PlanGetByID(id int64) (*schemas.PlanResponse, error) {
	ctx := context.Background()

	plan, err := boilmodels.Plans(
		boilmodels.PlanWhere.ID.EQ(id),
		qm.Load(boilmodels.PlanRels.Tenants),
	).One(ctx, r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.ErrorResponse(404, "El plan no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error al obtener el plan", err)
	}

	return mapToPlanResponse(plan), nil
}
