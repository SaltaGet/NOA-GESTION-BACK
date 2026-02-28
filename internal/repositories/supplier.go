package repositories

import (
	"errors"
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// SupplierGetByID obtiene un proveedor por ID
func (r *SupplierRepository) SupplierGetByID(id int64) (*schemas.SupplierResponse, error) {
	var supplier models.Supplier

	if err := r.DB.First(&supplier, id).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Proveedor", schemas.Read)
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
				query = query.Where("name LIKE ?", "%"+value+"%")
			case "company_name":
				query = query.Where("company_name LIKE ?", "%"+value+"%")
			case "identifier":
				query = query.Where("identifier LIKE ?", "%"+value+"%")
			case "email":
				query = query.Where("email LIKE ?", "%"+value+"%")
			}
		}
	}

	// Contar total de registros
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, schemas.HandlerErrorGorm(err, "Proveedor", schemas.Read)
	}

	// Obtener registros con paginación
	if err := query.
		Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&suppliers).Error; err != nil {
		return nil, 0, schemas.HandlerErrorGorm(err, "Proveedor", schemas.Read)
	}

	var suppliersSchema []*schemas.SupplierResponseDTO
	copier.Copy(&suppliersSchema, &suppliers)

	return suppliersSchema, total, nil
}

// SupplierCreate crea un nuevo proveedor con auditoría
func (r *SupplierRepository) SupplierCreate(memberID int64, supplierCreate *schemas.SupplierCreate) (int64, error) {
	var supplierSave models.Supplier
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}
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
			return schemas.HandlerErrorGorm(err, "Proveedor", schemas.Create)
		}

		supplierSave = supplier
		return nil
	})

	if err != nil {
		return 0, err
	}

	return supplierSave.ID, nil
}

// SupplierUpdate actualiza un proveedor existente con auditoría
func (r *SupplierRepository) SupplierUpdate(memberID int64, supplierUpdate *schemas.SupplierUpdate) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var existingSupplier models.Supplier
		if err := tx.First(&existingSupplier, supplierUpdate.ID).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Proveedor", schemas.Read)
		}

		existingSupplier.Name = supplierUpdate.Name
		existingSupplier.CompanyName = supplierUpdate.CompanyName
		existingSupplier.Identifier = supplierUpdate.Identifier
		existingSupplier.Address = supplierUpdate.Address
		existingSupplier.DebtLimit = supplierUpdate.DebtLimit
		existingSupplier.Email = supplierUpdate.Email
		existingSupplier.Phone = supplierUpdate.Phone

		if err := tx.Save(&existingSupplier).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Proveedor", schemas.Update)
		}

		return nil
	})

	return err
}

func (r *SupplierRepository) SupplierDelete(memberID int64, id int64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var supplier models.Supplier
		if err := tx.First(&supplier, id).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Proveedor", schemas.Read)
		}

		var expenseCount int64
		if err := tx.Model(&models.ExpenseBuy{}).Where("supplier_id = ?", id).Count(&expenseCount).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Egreso de compra", schemas.Read)
		}

		if expenseCount > 0 {
			return schemas.ErrorResponse(400, "No se puede eliminar el proveedor porque tiene compras asociadas", errors.New("No se puede eliminar el proveedor porque tiene compras asociadas"))
		}

		// Soft delete
		if err := tx.Delete(&supplier).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Proveedor", schemas.Delete)
		}
		return nil
	})

	return err
}
