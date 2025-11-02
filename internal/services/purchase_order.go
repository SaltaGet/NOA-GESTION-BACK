package services

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
)

func (p *PurchaseOrderService) PurchaseOrderGetByID(id string) (*models.PurchaseOrder, error) {
	purchaseOrder, err := p.PurchaseOrderRepository.PurchaseOrderGetByID(id)
	if err != nil {
		return nil, err
	}
	return purchaseOrder, nil
}

func (p *PurchaseOrderService) PurchaseOrderGetAll() (*[]models.PurchaseOrder, error) {
	purchaseOrder, err := p.PurchaseOrderRepository.PurchaseOrderGetAll()
	if err != nil {
		return nil, err
	}
	return purchaseOrder, nil
}

func (p *PurchaseOrderService) PurchaseOrderCreate(purchaseOrder *models.PurchaseOrderCreate) (string, error) {
	id, err := p.PurchaseOrderRepository.PurchaseOrderCreate(purchaseOrder)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (p *PurchaseOrderService) PurchaseOrderUpdate(purchaseOrder *models.PurchaseOrderUpdate) error {
	err := p.PurchaseOrderRepository.PurchaseOrderUpdate(purchaseOrder)
	if err != nil {
		return err
	}
	return nil
}

func (p *PurchaseOrderService) PurchaseOrderDelete(id string) error {
	err := p.PurchaseOrderRepository.PurchaseOrderDelete(id)
	if err != nil {
		return err
	}
	return nil
}