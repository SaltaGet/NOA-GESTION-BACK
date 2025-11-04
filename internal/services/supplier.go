package services

import (
"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (s *SupplierService) SupplierCreate(supplier *schemas.SupplierCreate) (string, error) {
	id, err := s.SupplierRepository.SupplierCreate(supplier)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *SupplierService) SupplierGetAll() (*[]schemas.Supplier, error) {
	suppliers, err := s.SupplierRepository.SupplierGetAll()
	if err != nil {
		return nil, err
	}
	return suppliers,nil
}

func (s *SupplierService) SupplierGetByID(id string) (*schemas.Supplier, error) {
	supplier, err := s.SupplierRepository.SupplierGetByID(id)
	if err != nil {
		return nil, err
	}
	return supplier, nil
}

func (s *SupplierService) SupplierGetByName(name string) (*[]schemas.Supplier, error) {
	supplier, err := s.SupplierRepository.SupplierGetByName(name)
	if err != nil {
		return nil, err
	}
	return supplier, nil
}

func (s *SupplierService) SupplierDelete(id string) error {
	err := s.SupplierRepository.SupplierDelete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *SupplierService) SupplierUpdate(supplierUpdate *schemas.SupplierUpdate) error {
	err := s.SupplierRepository.SupplierUpdate(supplierUpdate)
	if err != nil {
		return err
	}
	return nil
}