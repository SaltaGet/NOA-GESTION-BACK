package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

type SupplierService interface {
	SupplierGetByID(id int64) (suplier *schemas.SupplierResponse, err error)
	SupplierGetAll(limit, page int, search *map[string]string) ([]*schemas.SupplierResponseDTO, int64, error)
	SupplierCreate(memberID int64, supplierCreate *schemas.SupplierCreate) (int64, error)
	SupplierUpdate(memberID int64, supplierUpdate *schemas.SupplierUpdate) error
	SupplierDelete(memberID int64, id int64) error
}

type SupplierRepository interface {
	SupplierGetByID(id int64) (suplier *schemas.SupplierResponse, err error)
	SupplierGetAll(limit, page int, search *map[string]string) ([]*schemas.SupplierResponseDTO, int64, error)
	SupplierCreate(memberID int64, supplierCreate *schemas.SupplierCreate) (int64, error)
	SupplierUpdate(memberID int64, supplierUpdate *schemas.SupplierUpdate) error
	SupplierDelete(memberID int64, id int64) error
}
