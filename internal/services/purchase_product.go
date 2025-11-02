package services

// import (
// 	"errors"

// 	"github.com/DanielChachagua/GestionCar/pkg/models"
// 	"github.com/DanielChachagua/GestionCar/pkg/repositories"
// 	"gorm.io/gorm"
// )

// func (p *PurchaseProductService) PurchaseProductCreate(purchaseOrder *models.PurchaseProductCreate, workplace string) (string, error) {
// 	id, err := p.PurchaseProductRepository.PurchaseProductCreate(purchaseOrder, workplace)
// 	if err != nil {
// 		return "", models.ErrorResponse(500, "Error al actualizar cliente", err)
// 	}
// 	return id, nil
// }

// func (p *PurchaseProductService) PurchaseProductUpdate(purchaseOrder *models.PurchaseProductUpdate, workplace string) error {
// 	err := p.PurchaseProductRepository.UpdatePurchaseElement(purchaseOrder, workplace)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return models.ErrorResponse(404, "Empleado no encontrado", err)
// 		}
// 		return models.ErrorResponse(500, "Error al actualizar cliente", err)
// 	}
// 	return nil
// }

// func (p *PurchaseProductService) PurchaseProductDelete(id string, workplace string) error {
// 	err := p.PurchaseProductRepository.DeletePurchaseElementByID(id, workplace)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return models.ErrorResponse(404, "Empleado no encontrado", err)
// 		}
// 		return models.ErrorResponse(500, "Error al actualizar cliente", err)
// 	}
// 	return nil
// }

// func (p *PurchaseProductService) PurchaseProductGetByID(id string, workplace string) (*models.PurchaseProductLaundry, *models.PurchasePartWorkshop, error) {
// 	purchaseOrderLaundry, purchaseOrderWorkshop, err := p.PurchaseProductRepository.GetPurchaseElementByID(id, workplace)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil, models.ErrorResponse(404, "Empleado no encontrado", err)
// 		}
// 		return nil, nil, models.ErrorResponse(500, "Error al actualizar cliente", err)
// 	}
// 	return purchaseOrderLaundry, purchaseOrderWorkshop, nil
// }

// func (p *PurchaseProductService) PurchaseProductGetAllByPurhcaseID(purchaseID string, workplace string) (*[]models.PurchaseProductLaundry, *[]models.PurchasePartWorkshop, error) {
// 	purchaseOrderLaundry, purchaseOrderWorkshop, err := p.PurchaseProductRepository.GetPurchaseElementByPurchaseID(purchaseID, workplace)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil, models.ErrorResponse(404, "Empleado no encontrado", err)
// 		}
// 		return nil, nil, models.ErrorResponse(500, "Error al actualizar cliente", err)
// 	}
// 	return purchaseOrderLaundry, purchaseOrderWorkshop, nil
// }