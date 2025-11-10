package repositories

import (
	"errors"
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// SupplierGetByID obtiene un proveedor por ID
func (r *SupplierRepository) SupplierGetByID(id int64) (*schemas.SupplierResponse, error) {
	var supplier models.Supplier

	if err := r.DB.First(&supplier, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Proveedor no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error al obtener el proveedor", err)
	}

	var supplierSchema schemas.SupplierResponse
	copier.Copy(&supplierSchema, &supplier)

	return &supplierSchema, nil
}

// SupplierGetAll obtiene todos los proveedores con paginación y búsqueda
func (r *SupplierRepository) SupplierGetAll(limit, page int, search *map[string]string) ([]*schemas.SupplierResponseDTO, int64, error) {
	var suppliers []*models.Supplier
	var total int64

	offset := (page - 1) * limit

	query := r.DB.Model(&models.Supplier{})

	// Aplicar filtros de búsqueda si existen
	if search != nil {
		for key, value := range *search {
			switch key {
			case "name":
				query = query.Where("name ILIKE ?", "%"+value+"%")
			case "company_name":
				query = query.Where("company_name ILIKE ?", "%"+value+"%")
			case "identifier":
				query = query.Where("identifier ILIKE ?", "%"+value+"%")
			case "email":
				query = query.Where("email ILIKE ?", "%"+value+"%")
			case "phone":
				query = query.Where("phone ILIKE ?", "%"+value+"%")
			}
		}
	}

	// Contar total de registros
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "Error al contar los proveedores", err)
	}

	// Obtener registros con paginación
	if err := query.
		Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&suppliers).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "Error al obtener los proveedores", err)
	}

	var suppliersSchema []*schemas.SupplierResponseDTO
	copier.Copy(&suppliersSchema, &suppliers)

	return suppliersSchema, total, nil
}

// SupplierCreate crea un nuevo proveedor
func (r *SupplierRepository) SupplierCreate(supplierCreate *schemas.SupplierCreate) (int64, error) {
	var supplierID int64

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Crear proveedor
		supplier := models.Supplier{
			Name:        supplierCreate.Name,
			CompanyName: supplierCreate.CompanyName,
			Identifier:  supplierCreate.Identifier,
			Address:     supplierCreate.Address,
			DebtLimit:   supplierCreate.DebtLimit,
			Email:       supplierCreate.Email,
			Phone:       supplierCreate.Phone,
		}

		if err := tx.Create(&supplier).Error; err != nil {
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
				if strings.Contains(err.Error(), "identifier") {
					return schemas.ErrorResponse(400, "El CUIT ya existe", err)
				}
				if strings.Contains(err.Error(), "email") {
					return schemas.ErrorResponse(400, "El email ya existe", err)
				}
				return schemas.ErrorResponse(400, "El proveedor ya existe", err)
			}
			return schemas.ErrorResponse(500, "Error al crear el proveedor", err)
		}

		supplierID = supplier.ID
		return nil
	})

	if err != nil {
		return 0, err
	}

	return supplierID, nil
}

// SupplierUpdate actualiza un proveedor existente
func (r *SupplierRepository) SupplierUpdate(supplierUpdate *schemas.SupplierUpdate) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el proveedor existe
		var existingSupplier models.Supplier
		if err := tx.First(&existingSupplier, supplierUpdate.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Proveedor no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el proveedor", err)
		}

		// Actualizar campos
		existingSupplier.Name = supplierUpdate.Name
		existingSupplier.CompanyName = supplierUpdate.CompanyName
		existingSupplier.Identifier = supplierUpdate.Identifier
		existingSupplier.Address = supplierUpdate.Address
		existingSupplier.DebtLimit = supplierUpdate.DebtLimit
		existingSupplier.Email = supplierUpdate.Email
		existingSupplier.Phone = supplierUpdate.Phone

		if err := tx.Save(&existingSupplier).Error; err != nil {
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
				if strings.Contains(err.Error(), "identifier") {
					return schemas.ErrorResponse(400, "El CUIT ya existe", err)
				}
				if strings.Contains(err.Error(), "email") {
					return schemas.ErrorResponse(400, "El email ya existe", err)
				}
				return schemas.ErrorResponse(400, "Error de duplicación", err)
			}
			return schemas.ErrorResponse(500, "Error al actualizar el proveedor", err)
		}

		return nil
	})
}

// SupplierDelete elimina un proveedor (soft delete)
func (r *SupplierRepository) SupplierDelete(id int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el proveedor existe
		var supplier models.Supplier
		if err := tx.First(&supplier, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Proveedor no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el proveedor", err)
		}

		// Verificar si el proveedor tiene compras asociadas
		var expenseCount int64
		if err := tx.Model(&models.ExpenseBuy{}).Where("supplier_id = ?", id).Count(&expenseCount).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al verificar compras asociadas", err)
		}

		if expenseCount > 0 {
			return schemas.ErrorResponse(400, "No se puede eliminar el proveedor porque tiene compras asociadas", nil)
		}

		// Soft delete
		if err := tx.Delete(&supplier).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al eliminar el proveedor", err)
		}

		return nil
	})
}