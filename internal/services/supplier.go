package services

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (s *SupplierService) SupplierGetByID(id int64) (suplier *schemas.SupplierResponse, err error) {
	return s.SupplierRepository.SupplierGetByID(id)
}

func (s *SupplierService) SupplierGetAll(limit, page int, search *map[string]string) ([]*schemas.SupplierResponseDTO, int64, error) {
	return s.SupplierRepository.SupplierGetAll(limit, page, search)
}

func (s *SupplierService) SupplierCreate(supplierCreate *schemas.SupplierCreate) (int64, error) {
	return s.SupplierRepository.SupplierCreate(supplierCreate)
}

func (s *SupplierService) SupplierUpdate(supplierUpdate *schemas.SupplierUpdate) error {
	return s.SupplierRepository.SupplierUpdate(supplierUpdate)
}

func (s *SupplierService) SupplierDelete(id int64) error {
	return s.SupplierRepository.SupplierDelete(id)
}
