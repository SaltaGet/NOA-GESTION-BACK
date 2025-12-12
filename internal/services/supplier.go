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

func (s *SupplierService) SupplierCreate(memberID int64, supplierCreate *schemas.SupplierCreate) (int64, error) {
	return s.SupplierRepository.SupplierCreate(memberID, supplierCreate)
}

func (s *SupplierService) SupplierUpdate(memberID int64, supplierUpdate *schemas.SupplierUpdate) error {
	return s.SupplierRepository.SupplierUpdate(memberID, supplierUpdate)
}

func (s *SupplierService) SupplierDelete(memberID int64, id int64) error {
	return s.SupplierRepository.SupplierDelete(memberID, id)
}
