package repositories

import (
	"context"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/master"
	boiltenant "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

func (r *MainRepository) NotificationStock(tenantID int64) ([]*schemas.ProductSimpleResponse, error) {
	ctx := context.Background()

	_, err := boilmodels.Tenants(boilmodels.TenantWhere.ID.EQ(tenantID)).One(ctx, r.DB)
	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
	}

	dbTenant, err := database.GetTenantDB("", tenantID)
	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Tenant", schemas.Read)
	}

	// Products that have stock <= min_amount AND notifier = true
	products, err := boiltenant.Products(
		qm.InnerJoin("stock_deposits sd on sd.product_id = products.id"),
		qm.Where("sd.stock <= products.min_amount"),
		qm.Where("products.notifier = ?", true),
		qm.Load(boiltenant.ProductRels.Deposits),
	).All(ctx, dbTenant)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Product", schemas.Read)
	}

	var response []*schemas.ProductSimpleResponse
	for _, p := range products {
		pPrice, _ := p.Price.Big.Float64()
		minAmt, _ := p.MinAmount.Big.Float64()
		var sAmt float64
		if len(p.R.Deposits) > 0 {
			stk, _ := p.R.Deposits[0].Stock.Big.Float64()
			sAmt = stk
		}

		response = append(response, &schemas.ProductSimpleResponse{
			ID:        p.ID,
			Code:      p.Code,
			Name:      p.Name,
			Price:     pPrice,
			Stock:     sAmt,
			Notifier:  p.Notifier,
			MinAmount: minAmt,
		})
	}

	return response, nil
}
