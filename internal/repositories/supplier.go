package repositories

import (
	"errors"

	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *SupplierRepository) SupplierGetByID(id string) (*models.Supplier, error) {
	var supplier models.Supplier
	if err := r.DB.Where("id = ?", id).First(&supplier).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrorResponse(404, "Proveedor no encontrado", err)
		}
		return nil, models.ErrorResponse(500, "Error interno al buscar proveedor", err)
	}
	return &supplier, nil
}

func (r *SupplierRepository) SupplierGetAll() (*[]models.Supplier, error) {
	var suppliers []models.Supplier
	if err := r.DB.Find(&suppliers).Error; err != nil {
		return nil, models.ErrorResponse(500, "Error interno al buscar proveedores", err)
	}
	return &suppliers, nil
}

func (r *SupplierRepository) SupplierGetByName(name string) (*[]models.Supplier, error) {
	var supplier []models.Supplier
	if err := r.DB.Where("name LIKE ?", "%"+name+"%").Find(&supplier).Error; err != nil {
		return nil, models.ErrorResponse(500, "Error interno al buscar proveedores", err)
	}
	return &supplier, nil
}

func (r *SupplierRepository) SupplierCreate(supplierCreate *models.SupplierCreate) (string, error) {
	var supplierID string
	supplier := models.Supplier{
		ID:      uuid.NewString(),
		Name:    supplierCreate.Name,
		Address: supplierCreate.Address,
		Phone:   supplierCreate.Phone,
		Email:   supplierCreate.Email,
	}
	if err := r.DB.Create(&supplier).Error; err != nil {
		return "", models.ErrorResponse(500, "Error al crear proveedor", err)
	}
	supplierID = supplier.ID
	return supplierID, nil
}

func (r *SupplierRepository) SupplierUpdate(supplierUpdate *models.SupplierUpdate) error {
	var supplierLaundry models.Supplier
	if err := r.DB.Where("id = ?", supplierUpdate.ID).First(&supplierLaundry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			return models.ErrorResponse(404, "Proveedor no encontrado", err)
		}
		return models.ErrorResponse(500, "Error interno al buscar proveedor", err)
	}
	supplierLaundry.Name = supplierUpdate.Name
	supplierLaundry.Address = supplierUpdate.Address
	supplierLaundry.Phone = supplierUpdate.Phone
	supplierLaundry.Email = supplierUpdate.Email
	if err := r.DB.Save(&supplierLaundry).Error; err != nil {
		return models.ErrorResponse(500, "Error al actualizar proveedor", err)
	}

	return nil
}

func (r *SupplierRepository) SupplierDelete(id string) error {
	var supplier models.Supplier
	if err := r.DB.Where("id = ?", id).Delete(&supplier).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrorResponse(404, "Proveedor no encontrado", err)
		}
		return models.ErrorResponse(500, "Error interno al eliminar proveedor", err)
	}
	return nil
}
