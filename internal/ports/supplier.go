package ports

import "github.com/DanielChachagua/GestionCar/pkg/models"

type SupplierService interface {
	SupplierGetByID(id string) (suplier *models.Supplier, err error)
	SupplierGetByName(name string) (suplier *[]models.Supplier, err error)
	SupplierGetAll() (suppliers *[]models.Supplier, err error)
	SupplierCreate(supplier *models.SupplierCreate) (id string, err error)
	SupplierUpdate(supplierUpdate *models.SupplierUpdate) (err error)
	SupplierDelete(id string) (err error)
}

type SupplierRepository interface {
	SupplierGetByID(id string) (suplier *models.Supplier, err error)
	SupplierGetByName(name string) (suplier *[]models.Supplier, err error)
	SupplierGetAll() (suppliers *[]models.Supplier, err error)
	SupplierCreate(supplier *models.SupplierCreate) (id string, err error)
	SupplierUpdate(supplierUpdate *models.SupplierUpdate) (err error)
	SupplierDelete(id string) (err error)
}
