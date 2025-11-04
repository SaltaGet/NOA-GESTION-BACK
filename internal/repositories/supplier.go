package repositories

import (
	"errors"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *SupplierRepository) SupplierGetByID(id string) (*schemas.Supplier, error) {
	var supplier schemas.Supplier
	if err := r.DB.Where("id = ?", id).First(&supplier).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Proveedor no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al buscar proveedor", err)
	}
	return &supplier, nil
}

func (r *SupplierRepository) SupplierGetAll() (*[]schemas.Supplier, error) {
	var suppliers []schemas.Supplier
	if err := r.DB.Find(&suppliers).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al buscar proveedores", err)
	}
	return &suppliers, nil
}

func (r *SupplierRepository) SupplierGetByName(name string) (*[]schemas.Supplier, error) {
	var supplier []schemas.Supplier
	if err := r.DB.Where("name LIKE ?", "%"+name+"%").Find(&supplier).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al buscar proveedores", err)
	}
	return &supplier, nil
}

func (r *SupplierRepository) SupplierCreate(supplierCreate *schemas.SupplierCreate) (string, error) {
	var supplierID string
	supplier := schemas.Supplier{
		ID:      uuid.NewString(),
		Name:    supplierCreate.Name,
		Address: supplierCreate.Address,
		Phone:   supplierCreate.Phone,
		Email:   supplierCreate.Email,
	}
	if err := r.DB.Create(&supplier).Error; err != nil {
		return "", schemas.ErrorResponse(500, "Error al crear proveedor", err)
	}
	supplierID = supplier.ID
	return supplierID, nil
}

func (r *SupplierRepository) SupplierUpdate(supplierUpdate *schemas.SupplierUpdate) error {
	var supplierLaundry schemas.Supplier
	if err := r.DB.Where("id = ?", supplierUpdate.ID).First(&supplierLaundry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			return schemas.ErrorResponse(404, "Proveedor no encontrado", err)
		}
		return schemas.ErrorResponse(500, "Error interno al buscar proveedor", err)
	}
	supplierLaundry.Name = supplierUpdate.Name
	supplierLaundry.Address = supplierUpdate.Address
	supplierLaundry.Phone = supplierUpdate.Phone
	supplierLaundry.Email = supplierUpdate.Email
	if err := r.DB.Save(&supplierLaundry).Error; err != nil {
		return schemas.ErrorResponse(500, "Error al actualizar proveedor", err)
	}

	return nil
}

func (r *SupplierRepository) SupplierDelete(id string) error {
	var supplier schemas.Supplier
	if err := r.DB.Where("id = ?", id).Delete(&supplier).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "Proveedor no encontrado", err)
		}
		return schemas.ErrorResponse(500, "Error interno al eliminar proveedor", err)
	}
	return nil
}
