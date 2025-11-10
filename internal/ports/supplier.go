package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

type SupplierService interface {
	SupplierGetByID(id int64) (suplier *schemas.SupplierResponse, err error)
	SupplierGetAll(limit, page int, search *map[string]string) ([]*schemas.SupplierResponseDTO, int64, error)
	SupplierCreate(supplierCreate *schemas.SupplierCreate) (int64, error)
	SupplierUpdate(supplierUpdate *schemas.SupplierUpdate) error
	SupplierDelete(id int64) error
}

type SupplierRepository interface {
	SupplierGetByID(id int64) (suplier *schemas.SupplierResponse, err error)
	SupplierGetAll(limit, page int, search *map[string]string) ([]*schemas.SupplierResponseDTO, int64, error)
	SupplierCreate(supplierCreate *schemas.SupplierCreate) (int64, error)
	SupplierUpdate(supplierUpdate *schemas.SupplierUpdate) error
	SupplierDelete(id int64) error
}
