package repositories

import (
	"errors"

	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *PurchaseProductRepository) PurchaseProdcutGetByID(id string) (*models.PurchaseProduct, error) {
	var purchaseProduct models.PurchaseProduct
	if err := r.DB.Where("id = ?", id).First(&purchaseProduct).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrorResponse(404, "Producto de compra no encontrada", err)
		}
		return nil, models.ErrorResponse(500, "Error interno al obtener producto de compra", err)
	}
	return &purchaseProduct, nil
}

func (r *PurchaseProductRepository) PurchaseProductGetByPurchaseID(purchaseID string) (*[]models.PurchaseProduct, error) {
	var purchaseProduct []models.PurchaseProduct
	if err := r.DB.Where("purchase_order_id = ?", purchaseID).Find(&purchaseProduct).Error; err != nil {
		return nil, models.ErrorResponse(500, "Error interno al obtener productos de compra", err)
	}
	return &purchaseProduct, nil
}

func (r *PurchaseProductRepository) PurchaseProductGetAll() ([]models.PurchaseProduct, error) {
	var purchaseProducts []models.PurchaseProduct
	if err := r.DB.Find(&purchaseProducts).Error; err != nil {
		return nil, models.ErrorResponse(500, "Error interno al obtener productos de compra", err)
	}
	return purchaseProducts, nil
}

func (r *PurchaseProductRepository) PurchaseProductCreate(element *models.PurchaseProductCreate) (string, error) {
	newID := uuid.NewString()
	if err := r.DB.Create(&models.PurchaseProduct{
		ID:         newID,
		ProductID:  element.ProductID,
		ExpiredAt:  element.ExpiredAt,
		UnitPrice:  element.UnitPrice,
		Quantity:   element.Quantity,
		TotalPrice: element.UnitPrice * float32(element.Quantity),
	}).Error; err != nil {
		return "", models.ErrorResponse(500, "Error interno al crear producto de compra", err)
	}
	return newID, nil

}

func (r *PurchaseProductRepository) PurchaseProductUpdate(element *models.PurchaseProductUpdate) error {
	if err := r.DB.Where("id = ?", element.ID).Updates(&models.PurchaseProduct{
		ProductID:  element.ProductID,
		ExpiredAt:  element.ExpiredAt,
		UnitPrice:  element.UnitPrice,
		Quantity:   element.Quantity,
		TotalPrice: element.UnitPrice * float32(element.Quantity),
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrorResponse(404, "Producto de compra no encontrado", err)
		}
		return models.ErrorResponse(500, "Error interno al actualizar producto de compra", err)
	}
	return nil

}

func (r *PurchaseProductRepository) PurchaseProductDelete(id string) error {
	var purchaseProduct models.PurchaseProduct
	if err := r.DB.Where("id = ?", id).Delete(&purchaseProduct).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrorResponse(404, "Producto de compra no encontrado", err)
		}
		return models.ErrorResponse(500, "Error interno al eliminar producto de compra", err)
	}

	return nil
}
