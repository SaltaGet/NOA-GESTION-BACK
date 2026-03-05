package repositories

import (
	"context"
	"database/sql"
	"errors"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func mapToEcommerceResponse(c *boilmodels.IncomeEcommerce) *schemas.EcommerceResponse {
	if c == nil {
		return nil
	}

	tAmount, _ := c.TransactionAmount.Big.Float64()
	nAmount, _ := c.NetReceivedAmount.Big.Float64()
	total, _ := c.Total.Big.Float64()

	res := &schemas.EcommerceResponse{
		ID:                c.ID,
		PaymentID:         c.PaymentID,
		ExternalReference: c.ExternalReference,
		Status:            c.Status,
		Total:             total,
		DeliveryStatus:    c.DeliveryStatus,
		DateCreated:       c.DateCreated,
		DateApproved:      c.DateApproved,
		TransactionAmount: tAmount,
		NetReceivedAmount: nAmount,
		PayerFirstName:    c.PayerFirstName,
		PayerLastName:     c.PayerLastName,
		PayerEmail:        c.PayerEmail,
		PayMethod:         c.PayMethod,
		OperationType:     c.OperationType,
	}

	if c.DeliveryID.Valid {
		res.DeliveryID = &c.DeliveryID.String
	}
	if c.Message.Valid {
		res.Message = c.Message.String
	}

	if c.R != nil {
		for _, item := range c.R.IncomeEcommerceItems {
			iAmt, _ := item.Amount.Big.Float64()
			iPriceCost, _ := item.PriceCost.Big.Float64()
			iPrice, _ := item.Price.Big.Float64()
			iDiscount, _ := item.Discount.Big.Float64()
			iSubtotal, _ := item.Subtotal.Big.Float64()
			iTotal, _ := item.Total.Big.Float64()

			ier := schemas.EcommerceItemResponse{
				Amount:       iAmt,
				Price_Cost:   iPriceCost,
				Price:        iPrice,
				Discount:     iDiscount,
				TypeDiscount: item.TypeDiscount,
				Subtotal:     iSubtotal,
				Total:        iTotal,
			}
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
	}

	return res
}

func mapToEcommerceResponseDTO(c *boilmodels.IncomeEcommerce) schemas.EcommerceResponseDTO {
	total, _ := c.Total.Big.Float64()
	return schemas.EcommerceResponseDTO{
		ID:                c.ID,
		ExternalReference: c.ExternalReference,
		Status:            c.Status,
		Total:             total,
		DateCreated:       c.DateCreated,
		PayerEmail:        c.PayerEmail,
	}
}

func (er *EcommerceRepository) GetByID(id int64) (*schemas.EcommerceResponse, error) {
	ctx := context.Background()

	e, err := boilmodels.IncomeEcommerces(
		boilmodels.IncomeEcommerceWhere.ID.EQ(id),
		qm.Load(boilmodels.IncomeEcommerceRels.IncomeEcommerceItems),
		qm.Load(qm.Rels(boilmodels.IncomeEcommerceRels.IncomeEcommerceItems, boilmodels.IncomeEcommerceItemRels.Product)),
	).One(ctx, er.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Compra ecommerce", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Compra ecommerce", schemas.Read)
	}

	return mapToEcommerceResponse(e), nil
}

func (er *EcommerceRepository) GetByReference(reference string) (*schemas.EcommerceResponse, error) {
	ctx := context.Background()

	e, err := boilmodels.IncomeEcommerces(
		boilmodels.IncomeEcommerceWhere.ExternalReference.EQ(reference),
		qm.Load(boilmodels.IncomeEcommerceRels.IncomeEcommerceItems),
		qm.Load(qm.Rels(boilmodels.IncomeEcommerceRels.IncomeEcommerceItems, boilmodels.IncomeEcommerceItemRels.Product)),
	).One(ctx, er.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Compra ecommerce", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Compra ecommerce", schemas.Read)
	}

	return mapToEcommerceResponse(e), nil
}

func (er *EcommerceRepository) GetAll(page, limit int, status *string) ([]schemas.EcommerceResponseDTO, error) {
	ctx := context.Background()
	offset := (page - 1) * limit

	var qms []qm.QueryMod
	if status != nil {
		qms = append(qms, boilmodels.IncomeEcommerceWhere.Status.EQ(*status))
	}
	qms = append(qms,
		qm.Select(
			boilmodels.IncomeEcommerceColumns.ID,
			boilmodels.IncomeEcommerceColumns.ExternalReference,
			boilmodels.IncomeEcommerceColumns.Status,
			boilmodels.IncomeEcommerceColumns.Total,
			boilmodels.IncomeEcommerceColumns.DateCreated,
			boilmodels.IncomeEcommerceColumns.PayerEmail,
		),
		qm.Offset(offset),
		qm.Limit(limit),
	)

	ecommerces, err := boilmodels.IncomeEcommerces(qms...).All(ctx, er.DB)
	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Compra ecommerce", schemas.Read)
	}

	var ecommerceResponses []schemas.EcommerceResponseDTO
	for _, e := range ecommerces {
		ecommerceResponses = append(ecommerceResponses, mapToEcommerceResponseDTO(e))
	}

	return ecommerceResponses, nil
}

func (er *EcommerceRepository) UpdateStatus(update *schemas.EcommerceStatusUpdate) error {
	ctx := context.Background()

	ecommerce, err := boilmodels.IncomeEcommerces(
		boilmodels.IncomeEcommerceWhere.ID.EQ(update.ID),
	).One(ctx, er.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schemas.ErrorResponse(404, "compra electrónica no encontrada", errors.New("compra de ecomerce no encontrada"))
		}
		return schemas.HandlerErrorDB(err, "Compra ecommerce", schemas.Update)
	}

	ecommerce.Status = update.NewStatus
	if _, err := ecommerce.Update(ctx, er.DB, boil.Whitelist(boilmodels.IncomeEcommerceColumns.Status, boilmodels.IncomeEcommerceColumns.UpdatedAt)); err != nil {
		return schemas.HandlerErrorDB(err, "Compra ecommerce", schemas.Update)
	}

	return nil
}
