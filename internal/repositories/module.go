package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/master"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
)

func mapToModuleResponse(m *boilmodels.Module) *schemas.ModuleResponse {
	if m == nil {
		return nil
	}

	price, _ := m.PriceMonthly.Big.Float64()
	priceYearly, _ := m.PriceYearly.Big.Float64()

	res := &schemas.ModuleResponse{
		ID:                     m.ID,
		Name:                   m.Name,
		PriceMonthly:           price,
		PriceYearly:            priceYearly,
		Description:            m.Description,
		Features:               m.Features,
		AmountImagesPerProduct: int32(m.AmountImagesPerProduct),
	}

	return res
}

func (r *MainRepository) ModuleGet(id int64) (*schemas.ModuleResponse, error) {
	ctx := context.Background()

	module, err := boilmodels.Modules(boilmodels.ModuleWhere.ID.EQ(id)).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Modulo", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Modulo", schemas.Read)
	}

	return mapToModuleResponse(module), nil
}

func (r *MainRepository) ModuleGetAll() ([]schemas.ModuleResponse, error) {
	ctx := context.Background()

	modules, err := boilmodels.Modules().All(ctx, r.DB)
	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Modulo", schemas.Read)
	}

	var modulesResponse []schemas.ModuleResponse
	for _, m := range modules {
		res := mapToModuleResponse(m)
		if res != nil {
			modulesResponse = append(modulesResponse, *res)
		}
	}

	return modulesResponse, nil
}

func (r *MainRepository) ModuleCreate(moduleCreate *schemas.ModuleCreate) (int64, error) {
	ctx := context.Background()

	desc := ""
	if moduleCreate.Description != nil {
		desc = *moduleCreate.Description
	}

	newModule := &boilmodels.Module{
		Name:                   moduleCreate.Name,
		AmountImagesPerProduct: int(moduleCreate.AmountImagesPerProduct),
		PriceMonthly:           types.NewDecimal(decimal.New(0, 0).SetFloat64(moduleCreate.PriceMonthly)),
		PriceYearly:            types.NewDecimal(decimal.New(0, 0).SetFloat64(moduleCreate.PriceYearly)),
		Description:            desc,
		Features:               moduleCreate.Features,
	}

	if err := newModule.Insert(ctx, r.DB, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Modulo", schemas.Create)
	}

	return newModule.ID, nil
}

func (r *MainRepository) ModuleUpdate(moduleUpdate *schemas.ModuleUpdate) error {
	ctx := context.Background()

	module, err := boilmodels.Modules(boilmodels.ModuleWhere.ID.EQ(moduleUpdate.ID)).One(ctx, r.DB)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Modulo", schemas.Read)
	}

	module.Name = moduleUpdate.Name
	module.AmountImagesPerProduct = int(moduleUpdate.AmountImagesPerProduct)
	module.PriceMonthly = types.NewDecimal(decimal.New(0, 0).SetFloat64(moduleUpdate.PriceMonthly))
	module.PriceYearly = types.NewDecimal(decimal.New(0, 0).SetFloat64(moduleUpdate.PriceYearly))
	desc := ""
	if moduleUpdate.Description != nil {
		desc = *moduleUpdate.Description
	}
	module.Description = desc
	module.Features = moduleUpdate.Features

	if _, err := module.Update(ctx, r.DB, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Modulo", schemas.Update)
	}

	return nil
}

func (r *MainRepository) ModuleDelete(id int64) error {
	// Original code returns nil explicitly for testing or avoiding deletes.
	// We retain it as is.
	return nil
}

func (r *MainRepository) ModuleAddTenant(moduleAddTenant *schemas.ModuleAddTenant) error {
	ctx := context.Background()

	tenantExists, err := boilmodels.Tenants(boilmodels.TenantWhere.ID.EQ(moduleAddTenant.TenantID)).Exists(ctx, r.DB)
	if err != nil || !tenantExists {
		return schemas.HandlerErrorDB(sql.ErrNoRows, "Tenant", schemas.Read)
	}

	moduleExists, err := boilmodels.Modules(boilmodels.ModuleWhere.ID.EQ(moduleAddTenant.ModuleID)).Exists(ctx, r.DB)
	if err != nil || !moduleExists {
		return schemas.HandlerErrorDB(sql.ErrNoRows, "Modulo", schemas.Read)
	}

	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		loc = time.UTC
	}

	exp, err := time.ParseInLocation("2006-01-02", moduleAddTenant.Expiration, loc)
	if err != nil {
		return schemas.ErrorResponse(422, "Formato de fecha inválido, debe ser YYYY-MM-DD", err)
	}

	newModuleAdd := &boilmodels.TenantModule{
		ModuleID:   moduleAddTenant.ModuleID,
		TenantID:   moduleAddTenant.TenantID,
		Expiration: null.TimeFrom(exp),
	}

	err = newModuleAdd.Upsert(ctx, r.DB, true, []string{"tenant_id", "module_id"}, boil.Infer(), boil.Infer())
	if err != nil {
		return schemas.HandlerErrorDB(err, "Modulo", schemas.Create)
	}

	return nil
}

func (r *MainRepository) ModuleGetByTenantID(tenantID int64) ([]schemas.ModuleResponseDTO, error) {
	ctx := context.Background()

	tenantModules, err := boilmodels.TenantModules(
		boilmodels.TenantModuleWhere.TenantID.EQ(tenantID),
		qm.Load(boilmodels.TenantModuleRels.Module),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Modulo", schemas.Read)
	}

	modulesResponse := make([]schemas.ModuleResponseDTO, 0, len(tenantModules))

	for _, tm := range tenantModules {
		if tm.R != nil && tm.R.Module != nil {
			var expire *time.Time
			if tm.Expiration.Valid {
				t := tm.Expiration.Time
				expire = &t
			}
			modulesResponse = append(modulesResponse, schemas.ModuleResponseDTO{
				ID:                     tm.R.Module.ID,
				Name:                   tm.R.Module.Name,
				AmountImagesPerProduct: int32(tm.R.Module.AmountImagesPerProduct),
				Expiration:             expire,
				AcceptTerms:            tm.AcceptedTerms,
			})
		}
	}

	return modulesResponse, nil
}
