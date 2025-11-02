package ports

import "github.com/DanielChachagua/GestionCar/pkg/models"

type PurchaseOrderService interface {
	PurchaseOrderGetByID(id string) (client *models.PurchaseOrder, err error)
	PurchaseOrderGetAll() (clients *[]models.PurchaseOrder, err error)
	PurchaseOrderCreate(purchaseOrderCreate *models.PurchaseOrderCreate) (id string, err error)
	PurchaseOrderUpdate(purchaseOrderUpdate *models.PurchaseOrderUpdate) (err error)
	PurchaseOrderDelete(id string) (err error)
}

type PurchaseOrderRepository interface {
	PurchaseOrderGetByID(id string) (client *models.PurchaseOrder, err error)
	PurchaseOrderGetAll() (clients *[]models.PurchaseOrder, err error)
	PurchaseOrderCreate(purchaseOrderCreate *models.PurchaseOrderCreate) (id string, err error)
	PurchaseOrderUpdate(purchaseOrderUpdate *models.PurchaseOrderUpdate) (err error)
	PurchaseOrderDelete(id string) (err error)
}
