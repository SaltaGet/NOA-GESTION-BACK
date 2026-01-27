package grpc_repo

import (
	"context"
	"errors"
	// "math"
	"strconv"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// func (s *GrpcMPRepository) SyncPurchasePayment(ctx context.Context, req *pb.DataInfoPay) error {
// 	return s.DB.Transaction(func(tx *gorm.DB) error {
// 		newPay := &models.IncomeEcommerce{
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
// 		if err := tx.Create(newPay).Error; err != nil {
// 			return status.Errorf(codes.Internal, "error interno: %v", err)
// 		}

// 		for _, item := range req.AdditionalInfo.Items {
// 			id, err := strconv.ParseInt(item.Id, 10, 64)
// 			if err != nil {
// 				return status.Errorf(codes.InvalidArgument, "error interno id: %v", err)
// 			}
// 			amount, err := strconv.ParseFloat(item.Quantity, 64)
// 			if err != nil {
// 				return status.Errorf(codes.InvalidArgument, "error interno amount: %v", err)
// 			}
// 			price, err := strconv.ParseFloat(item.UnitPrice, 64)
// 			if err != nil {
// 				return status.Errorf(codes.InvalidArgument, "error interno price: %v", err)
// 			}

// 			var productPrice models.Product
// 			if err := tx.Select("price").
// 				Where("id = ?", item.Id).First(&productPrice).Error; err != nil {
// 				if errors.Is(err, gorm.ErrRecordNotFound) {
// 					return status.Errorf(codes.NotFound, "%s", err.Error())
// 				}
// 				return status.Errorf(codes.Internal, "error interno product: %v", err)
// 			}

// 			if math.Abs(productPrice.Price - price) > 100 {
// 				return status.Errorf(codes.InvalidArgument, "error interno: %v", err)
// 			}

// 			var stock models.Deposit
// 			if err := tx.
// 				Where("product_id = ?", item.Id).
// 				First(&stock).Error; err != nil {
// 				if errors.Is(err, gorm.ErrRecordNotFound) {
// 					return status.Errorf(codes.NotFound, "%s", err.Error())
// 				}
// 				return status.Errorf(codes.Internal, "error interno: %v", err)
// 			}

// 			if stock.Stock < amount {
// 				return status.Errorf(codes.InvalidArgument, "stock insuficiente para el producto %s (disponible: %.2f, requerido: %g)", item.Id, stock.Stock, amount)
// 			}

// 			stock.Stock -= amount
// 			if err := tx.Save(&stock).Error; err != nil {
// 				return status.Errorf(codes.Internal, "error interno: %v", err)
// 			}

// 			var priceCost models.ExpenseBuyItem
// 			if err := tx.
// 				Select("price").
// 				Where("product_id = ?", id).
// 				Order("created_at DESC").
// 				First(&priceCost).Error; err != nil {
// 				if errors.Is(err, gorm.ErrRecordNotFound) {
// 					priceCost.Price = productPrice.Price
// 				}
// 			}

// 			// subtotalItem := amount * productPrice.Price
// 			// totalItem := 0.0

// 			newItem := &models.IncomeEcommerceItem{
// 				IncomeEcommerceID: newPay.ID,
// 				ProductID:         id,
// 				Amount:            amount,
// 				Price_Cost:        priceCost.Price,
// 				Price:             price,
// 				Subtotal:          price * amount,
// 				Total:             price * amount,
// 			}
// 			if err := tx.Create(newItem).Error; err != nil {
// 				return status.Errorf(codes.Internal, "error interno: %v", err)
// 			}
// 		}

// 		return nil
// 	})
// }


func (s *GrpcMPRepository) SyncPurchasePayment(ctx context.Context, req *pb.DataInfoPay) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var existingPay models.IncomeEcommerce
		// Buscamos si ya existe el pago por su referencia única
		err := tx.Where("external_reference = ?", req.ExternalReference).First(&existingPay).Error
		isNew := errors.Is(err, gorm.ErrRecordNotFound)

		if err != nil && !isNew {
			return status.Errorf(codes.Internal, "error al buscar pago existente: %v", err)
		}

		// Definimos los estados
		isApprovedOrPending := req.Status == "approved" || req.Status == "pending"
		isRejected := req.Status == "rejected" || req.Status == "cancelled"

		// Determinamos la acción de stock basándonos en la transición de estados
		// 1. Si es nuevo y viene aprobado/pendiente -> Descontar
		// 2. Si ya existía, no estaba rechazado y ahora SI lo rechazan -> Reestablecer
		// 3. Si ya existía y el estado previo no era "fuerte" y ahora sí -> Ya se descontó, no hacer nada
		
		shouldDecreaseStock := isNew && isApprovedOrPending
		shouldRestoreStock := !isNew && (existingPay.Status == "approved" || existingPay.Status == "pending") && isRejected

		newPay := &models.IncomeEcommerce{
			PaymentID:         req.Id,
			ExternalReference: req.ExternalReference,
			Status:            req.Status,
			DeliveryStatus:    "pendiente",
			Total:             req.TransactionDetails.NetReceivedAmount,
			DateCreated:       req.DateCreated,
			DateApproved:      req.DateApproved,
			TransactionAmount: req.TransactionDetails.TotalPaidAmount,
			NetReceivedAmount: req.TransactionDetails.NetReceivedAmount,
			PayerFirstName:    req.AdditionalInfo.Payer.FirstName,
			PayerLastName:     req.AdditionalInfo.Payer.LastName,
			PayerEmail:        req.Payer.Email,
			PayMethod:         req.PaymentMethod.Type,
			OperationType:     req.OperationType,
			Message:           req.Message,
		}

		if isNew {
			if err := tx.Create(newPay).Error; err != nil {
				return status.Errorf(codes.Internal, "error al crear: %v", err)
			}
		} else {
			newPay.ID = existingPay.ID // Mantener el ID original
			if err := tx.Model(&existingPay).Updates(newPay).Error; err != nil {
				return status.Errorf(codes.Internal, "error al actualizar: %v", err)
			}
		}

		// PROCESAMIENTO DE ITEMS Y STOCK
		for _, item := range req.AdditionalInfo.Items {
			id, _ := strconv.ParseInt(item.Id, 10, 64)
			amount, _ := strconv.ParseFloat(item.Quantity, 64)
			price, _ := strconv.ParseFloat(item.UnitPrice, 64)

			// Solo tocamos el stock si es una creación aprobada o una cancelación de algo previo aprobado
			if shouldDecreaseStock || shouldRestoreStock {
				var stock models.Deposit
				if err := tx.Where("product_id = ?", item.Id).First(&stock).Error; err != nil {
					return status.Errorf(codes.NotFound, "stock no encontrado")
				}

				if shouldDecreaseStock {
					if stock.Stock < amount {
						return status.Errorf(codes.InvalidArgument, "stock insuficiente")
					}
					stock.Stock -= amount
				} else if shouldRestoreStock {
					stock.Stock += amount
				}

				if err := tx.Save(&stock).Error; err != nil {
					return status.Errorf(codes.Internal, "error actualizando stock: %v", err)
				}
			}

			// Solo creamos los items si el registro es nuevo
			if isNew {
				// (Aquí va tu lógica existente de buscar priceCost...)
				newItem := &models.IncomeEcommerceItem{
					IncomeEcommerceID: newPay.ID,
					ProductID:         id,
					Amount:            amount,
					Price:             price,
					Subtotal:          price * amount,
					Total:             price * amount,
				}
				if err := tx.Create(newItem).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}