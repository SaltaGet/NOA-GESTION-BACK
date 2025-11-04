package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"


type SupplierService interface {
	SupplierGetByID(id string) (suplier *schemas.Supplier, err error)
	SupplierGetByName(name string) (suplier *[]schemas.Supplier, err error)
	SupplierGetAll() (suppliers *[]schemas.Supplier, err error)
	SupplierCreate(supplier *schemas.SupplierCreate) (id string, err error)
	SupplierUpdate(supplierUpdate *schemas.SupplierUpdate) (err error)
	SupplierDelete(id string) (err error)
}

type SupplierRepository interface {
	SupplierGetByID(id string) (suplier *schemas.Supplier, err error)
	SupplierGetByName(name string) (suplier *[]schemas.Supplier, err error)
	SupplierGetAll() (suppliers *[]schemas.Supplier, err error)
	SupplierCreate(supplier *schemas.SupplierCreate) (id string, err error)
	SupplierUpdate(supplierUpdate *schemas.SupplierUpdate) (err error)
	SupplierDelete(id string) (err error)
}
