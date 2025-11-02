package repositories

import (
	"errors"

	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *PurchaseOrderRepository) PurchaseOrderGetByID(id string) (*models.PurchaseOrder, error) {
	var purchaseOrder models.PurchaseOrder
	if err := r.DB.Where("id = ?", id).First(&purchaseOrder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrorResponse(404, "Compra no encontrada", err)
		}
		return nil, models.ErrorResponse(500, "Error interno al obtener compra", err)
	}
	return &purchaseOrder, nil

}

func (r *PurchaseOrderRepository) PurchaseOrderGetAll() (*[]models.PurchaseOrder, error) {
	var purchaseOrders []models.PurchaseOrder
	if err := r.DB.Find(&purchaseOrders).Error; err != nil {
		return nil, models.ErrorResponse(500, "Error interno al obtener compras", err)
	}
	return &purchaseOrders, nil

}

func (r *PurchaseOrderRepository) PurchaseOrderCreate(purchaseOrder *models.PurchaseOrderCreate) (string, error) {
	newID := uuid.NewString()
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&models.PurchaseOrder{
			ID:          newID,
			OrderNumber: purchaseOrder.OrderNumber,
			OrderDate:   purchaseOrder.OrderDate,
			Amount:      purchaseOrder.Amount,
			SupplierID:  purchaseOrder.SupplierID,
		}).Error; err != nil {
			return models.ErrorResponse(500, "Error interno al crear compra", err)
		}
		for _, element := range purchaseOrder.PurchaseProductCreates {
			if err := tx.Create(&models.PurchaseProduct{
				ID:              uuid.NewString(),
				ProductID:       element.ProductID,
				PurchaseOrderID: newID,
				ExpiredAt:       element.ExpiredAt,
				UnitPrice:       element.UnitPrice,
				Quantity:        element.Quantity,
				TotalPrice:      element.UnitPrice * float32(element.Quantity),
			}).Error; err != nil {
				return models.ErrorResponse(500, "Error interno al crear productos de compra", err)
			}
		}
		return nil
	})
	if err != nil {
		return "", models.ErrorResponse(500, "Error interno al crear compra", err)
	}
	return newID, nil
}

func (r *PurchaseOrderRepository) PurchaseOrderUpdate(purchaseOrder *models.PurchaseOrderUpdate) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", purchaseOrder.ID).Updates(&models.PurchaseOrder{
			OrderNumber: purchaseOrder.OrderNumber,
			OrderDate:   purchaseOrder.OrderDate,
			Amount:      purchaseOrder.Amount,
			SupplierID:  purchaseOrder.SupplierID,
		}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return models.ErrorResponse(404, "Compra no encontrada", err)
			}
			return models.ErrorResponse(500, "Error interno al actualizar compra", err)
		}

		var existingProducts []models.PurchaseProduct
		if err := tx.Where("purchase_order_id = ?", purchaseOrder.ID).Find(&existingProducts).Error; err != nil {
			return models.ErrorResponse(500, "Error interno al buscar productos de compra", err)
		}
		existingIDs := map[string]bool{}
		for _, p := range existingProducts {
			existingIDs[p.ID] = true
		}

		receivedIDs := map[string]bool{}
		for _, prod := range purchaseOrder.PurchaseProductUpdates {
			receivedIDs[prod.ID] = true
		}

		for _, p := range existingProducts {
			if !receivedIDs[p.ID] {
				if err := tx.Delete(&models.PurchaseProduct{}, "id = ?", p.ID).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return models.ErrorResponse(404, "Producto de la orden de compra no encontrado", err)
					}
					return models.ErrorResponse(500, "Error interno al eliminar producto de compra", err)
				}
			}
		}

		for _, prod := range purchaseOrder.PurchaseProductUpdates {
			if prod.ID == "" || !existingIDs[prod.ID] {
				newProd := models.PurchaseProduct{
					ID:              uuid.NewString(),
					ProductID:       prod.ProductID,
					PurchaseOrderID: purchaseOrder.ID,
					ExpiredAt:       prod.ExpiredAt,
					UnitPrice:       prod.UnitPrice,
					Quantity:        prod.Quantity,
					TotalPrice:      prod.UnitPrice * float32(prod.Quantity),
				}
				if err := tx.Create(&newProd).Error; err != nil {
					return models.ErrorResponse(500, "Error interno al crear producto de compra", err)
				}
			} else {
				if err := tx.Model(&models.PurchaseProduct{}).
					Where("id = ?", prod.ID).
					Updates(map[string]interface{}{
						"product_id":  prod.ProductID,
						"expired_at":  prod.ExpiredAt,
						"unit_price":  prod.UnitPrice,
						"quantity":    prod.Quantity,
						"total_price": prod.UnitPrice * float32(prod.Quantity),
					}).Error; err != nil {
					return models.ErrorResponse(500, "Error interno al actualizar producto de compra", err)
				}
			}
		}
		return nil

	})
}

func (r *PurchaseOrderRepository) PurchaseOrderDelete(id string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("purchase_order_id = ?", id).Delete(&models.PurchaseProduct{}).Error; err != nil {
			return models.ErrorResponse(500, "Error interno al eliminar productos de compra", err)
		}
		if err := tx.Where("id = ?", id).Delete(&models.PurchaseOrder{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return models.ErrorResponse(404, "Compra no encontrada", err)
			}
			return models.ErrorResponse(500, "Error interno al eliminar compra", err)
		}
		return nil
	})
}

