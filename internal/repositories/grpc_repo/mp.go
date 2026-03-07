package grpc_repo

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GrpcMPRepository) SyncPurchasePayment(ctx context.Context, req *pb.DataInfoPay) error {
	ctx = boil.WithDebug(ctx, true)

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return status.Errorf(codes.Internal, "no se pudo iniciar transacción: %v", err)
	}

	defer tx.Rollback()

	// 1. Buscar si ya existe el pago
	existingPay, err := tenant.IncomeEcommerces(
		qm.Where("external_reference = ?", req.ExternalReference),
	).One(ctx, tx)

	isNew := errors.Is(err, sql.ErrNoRows)
	if err != nil && !isNew {
		return status.Errorf(codes.Internal, "error al buscar pago existente: %v", err)
	}

	// Definimos los estados
	isApprovedOrPending := req.Status == "approved" || req.Status == "pending"
	isRejected := req.Status == "rejected" || req.Status == "cancelled"

	var shouldDecreaseStock bool
	var shouldRestoreStock bool

	// Mapeo de datos al modelo de SQLBoiler
	newPay := &tenant.IncomeEcommerce{
		PaymentID:         req.Id,
		ExternalReference: req.ExternalReference,
		Status:            req.Status,
		DeliveryStatus:    "pendiente",
		Total:             types.NewDecimal(decimal.New(0, 0).SetFloat64(req.TransactionDetails.NetReceivedAmount)),
		DateCreated:       req.DateCreated,
		DateApproved:      req.DateApproved,
		TransactionAmount: types.NewDecimal(decimal.New(0, 0).SetFloat64(req.TransactionDetails.TotalPaidAmount)),
		NetReceivedAmount: types.NewDecimal(decimal.New(0, 0).SetFloat64(req.TransactionDetails.NetReceivedAmount)),
		PayerFirstName:    req.AdditionalInfo.Payer.FirstName,
		PayerLastName:     req.AdditionalInfo.Payer.LastName,
		PayerEmail:        req.Payer.Email,
		PayMethod:         req.PaymentMethod.Type,
		OperationType:     req.OperationType,
		Message:           null.NewString(req.Message, true),
	}

	if isNew {
		shouldDecreaseStock = isApprovedOrPending
		// Insertar nuevo registro
		if err := newPay.Insert(ctx, tx, boil.Infer()); err != nil {
			return status.Errorf(codes.Internal, "error al crear: %v", err)
		}
	} else {
		shouldRestoreStock = (existingPay.Status == "approved" || existingPay.Status == "pending") && isRejected
		newPay.ID = existingPay.ID // Mantener ID para el Update

		if _, err := newPay.Update(ctx, tx, boil.Infer()); err != nil {
			return status.Errorf(codes.Internal, "error al actualizar: %v", err)
		}
	}

	// PROCESAMIENTO DE ITEMS Y STOCK
	for _, item := range req.AdditionalInfo.Items {
		id, _ := strconv.ParseInt(item.Id, 10, 64)
		qty, _ := strconv.ParseFloat(item.Quantity, 64)
		price, _ := strconv.ParseFloat(item.UnitPrice, 64)

		if shouldDecreaseStock || shouldRestoreStock {
			// Buscar stock (Deposit)
			stock, err := tenant.Deposits(qm.Where("product_id = ?", id)).One(ctx, tx)
			if err != nil {
				return status.Errorf(codes.NotFound, "stock no encontrado para producto %d", id)
			}

			stockValue, _ := stock.Stock.Float64()

			if shouldDecreaseStock {
				if stockValue < qty {
					return status.Errorf(codes.InvalidArgument, "stock insuficiente para producto %d", id)
				}
				stock.Stock = types.NewDecimal(new(decimal.Big).SetFloat64(stockValue - qty))
			} else if shouldRestoreStock {
				stock.Stock = types.NewDecimal(new(decimal.Big).SetFloat64(stockValue + qty))
			}

			// Guardar cambio de stock (solo la columna stock para ser eficientes)
			if _, err := stock.Update(ctx, tx, boil.Whitelist(tenant.DepositColumns.Stock)); err != nil {
				return status.Errorf(codes.Internal, "error actualizando stock: %v", err)
			}
		}

		if isNew {
			newItem := &tenant.IncomeEcommerceItem{
				IncomeEcommerceID: newPay.ID,
				ProductID:         id,
				Amount:            types.NewDecimal(decimal.New(0, 0).SetFloat64(qty)),
				Price:             types.NewDecimal(decimal.New(0, 0).SetFloat64(price)),
				Subtotal:          types.NewDecimal(decimal.New(0, 0).SetFloat64(price * qty)),
				Total:             types.NewDecimal(decimal.New(0, 0).SetFloat64(price * qty)),
			}
			if err := newItem.Insert(ctx, tx, boil.Infer()); err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return status.Errorf(codes.Internal, "error confirmando transacción: %v", err)
	}

	return nil
}

// import (
// 	"context"
// 	"errors"
// 	// "math"
// 	"strconv"

// 	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// 	"gorm.io/gorm"
// )

// func (s *GrpcMPRepository) SyncPurchasePayment(ctx context.Context, req *pb.DataInfoPay) error {
// 	return s.DB.Transaction(func(tx *gorm.DB) error {
// 		var existingPay tenant.IncomeEcommerce
// 		// Buscamos si ya existe el pago por su referencia única
// 		err := tx.Where("external_reference = ?", req.ExternalReference).First(&existingPay).Error
// 		isNew := errors.Is(err, gorm.ErrRecordNotFound)

// 		if err != nil && !isNew {
// 			return status.Errorf(codes.Internal, "error al buscar pago existente: %v", err)
// 		}

// 		// Definimos los estados
// 		isApprovedOrPending := req.Status == "approved" || req.Status == "pending"
// 		isRejected := req.Status == "rejected" || req.Status == "cancelled"

// 		// Determinamos la acción de stock basándonos en la transición de estados
// 		// 1. Si es nuevo y viene aprobado/pendiente -> Descontar
// 		// 2. Si ya existía, no estaba rechazado y ahora SI lo rechazan -> Reestablecer
// 		// 3. Si ya existía y el estado previo no era "fuerte" y ahora sí -> Ya se descontó, no hacer nada

// 		shouldDecreaseStock := isNew && isApprovedOrPending
// 		shouldRestoreStock := !isNew && (existingPay.Status == "approved" || existingPay.Status == "pending") && isRejected

// 		newPay := &tenant.IncomeEcommerce{
// 			PaymentID:         req.Id,
// 			ExternalReference: req.ExternalReference,
// 			Status:            req.Status,
// 			DeliveryStatus:    "pendiente",
// 			Total:             req.TransactionDetails.NetReceivedAmount,
// 			DateCreated:       req.DateCreated,
// 			DateApproved:      req.DateApproved,
// 			TransactionAmount: req.TransactionDetails.TotalPaidAmount,
// 			NetReceivedAmount: req.TransactionDetails.NetReceivedAmount,
// 			PayerFirstName:    req.AdditionalInfo.Payer.FirstName,
// 			PayerLastName:     req.AdditionalInfo.Payer.LastName,
// 			PayerEmail:        req.Payer.Email,
// 			PayMethod:         req.PaymentMethod.Type,
// 			OperationType:     req.OperationType,
// 			Message:           req.Message,
// 		}

// 		if isNew {
// 			if err := tx.Create(newPay).Error; err != nil {
// 				return status.Errorf(codes.Internal, "error al crear: %v", err)
// 			}
// 		} else {
// 			newPay.ID = existingPay.ID // Mantener el ID original
// 			if err := tx.Model(&existingPay).Updates(newPay).Error; err != nil {
// 				return status.Errorf(codes.Internal, "error al actualizar: %v", err)
// 			}
// 		}

// 		// PROCESAMIENTO DE ITEMS Y STOCK
// 		for _, item := range req.AdditionalInfo.Items {
// 			id, _ := strconv.ParseInt(item.Id, 10, 64)
// 			amount, _ := strconv.ParseFloat(item.Quantity, 64)
// 			price, _ := strconv.ParseFloat(item.UnitPrice, 64)

// 			// Solo tocamos el stock si es una creación aprobada o una cancelación de algo previo aprobado
// 			if shouldDecreaseStock || shouldRestoreStock {
// 				var stock tenant.Deposit
// 				if err := tx.Where("product_id = ?", item.Id).First(&stock).Error; err != nil {
// 					return status.Errorf(codes.NotFound, "stock no encontrado")
// 				}

// 				if shouldDecreaseStock {
// 					if stock.Stock < amount {
// 						return status.Errorf(codes.InvalidArgument, "stock insuficiente")
// 					}
// 					stock.Stock -= amount
// 				} else if shouldRestoreStock {
// 					stock.Stock += amount
// 				}

// 				if err := tx.Save(&stock).Error; err != nil {
// 					return status.Errorf(codes.Internal, "error actualizando stock: %v", err)
// 				}
// 			}

// 			// Solo creamos los items si el registro es nuevo
// 			if isNew {
// 				// (Aquí va tu lógica existente de buscar priceCost...)
// 				newItem := &tenant.IncomeEcommerceItem{
// 					IncomeEcommerceID: newPay.ID,
// 					ProductID:         id,
// 					Amount:            amount,
// 					Price:             price,
// 					Subtotal:          price * amount,
// 					Total:             price * amount,
// 				}
// 				if err := tx.Create(newItem).Error; err != nil {
// 					return err
// 				}
// 			}
// 		}

// 		return nil
// 	})
// }
