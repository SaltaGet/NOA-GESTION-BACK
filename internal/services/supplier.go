package services

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
)

func (s *SupplierService) SupplierCreate(supplier *models.SupplierCreate) (string, error) {
	id, err := s.SupplierRepository.SupplierCreate(supplier)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *SupplierService) SupplierGetAll() (*[]models.Supplier, error) {
	suppliers, err := s.SupplierRepository.SupplierGetAll()
	if err != nil {
		return nil, err
	}
	return suppliers,nil
}

func (s *SupplierService) SupplierGetByID(id string) (*models.Supplier, error) {
	supplier, err := s.SupplierRepository.SupplierGetByID(id)
	if err != nil {
		return nil, err
	}
	return supplier, nil
}

func (s *SupplierService) SupplierGetByName(name string) (*[]models.Supplier, error) {
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

func (s *SupplierService) SupplierUpdate(supplierUpdate *models.SupplierUpdate) error {
	err := s.SupplierRepository.SupplierUpdate(supplierUpdate)
	if err != nil {
		return err
	}
	return nil
}