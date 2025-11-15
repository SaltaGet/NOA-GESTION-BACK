package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// IncomeOtherGetByID obtiene un ingreso por ID
func (r *IncomeOtherRepository) IncomeOtherGetByID(id int64) (*schemas.IncomeOtherResponse, error) {
	var incomeOther models.IncomeOther

	if err := r.DB.
		Preload("Member", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "username")
		}).
		Preload("TypeIncome").
		First(&incomeOther, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Ingreso no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error al obtener el ingreso", err)
	}

	var incomeSchema schemas.IncomeOtherResponse
	copier.Copy(&incomeSchema, &incomeOther)

	return &incomeSchema, nil
}

// IncomeOtherGetByDate obtiene ingresos por rango de fechas
func (r *IncomeOtherRepository) IncomeOtherGetByDate(pointSaleID *int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.IncomeOtherResponse, int64, error) {
	var incomesOther []*models.IncomeOther

	offset := (page - 1) * limit

	query := r.DB.Where("created_at BETWEEN ? AND ?", fromDate, toDate)
	
	// Si se proporciona pointSaleID, filtrar por punto de venta
	if pointSaleID != nil {
		query = query.Where("point_sale_id = ?", *pointSaleID)
	}

	if err := query.
		Preload("Member", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "username")
		}).
		Preload("TypeIncome").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&incomesOther).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "Error al obtener los ingresos", err)
	}

	// Contar total
	var total int64
	countQuery := r.DB.Model(&models.IncomeOther{}).
		Where("created_at BETWEEN ? AND ?", fromDate, toDate)
	
	if pointSaleID != nil {
		countQuery = countQuery.Where("point_sale_id = ?", *pointSaleID)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "Error al contar los ingresos", err)
	}

	var incomeSchema []*schemas.IncomeOtherResponse
	copier.Copy(&incomeSchema, &incomesOther)

	return incomeSchema, total, nil
}

// IncomeOtherCreate crea un nuevo ingreso
func (r *IncomeOtherRepository) IncomeOtherCreate(memberID, pointSaleID int64, incomeOtherCreate *schemas.IncomeOtherCreate) (int64, error) {
	var incomeOtherID int64
	
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el tipo de ingreso existe
		var typeIncome models.TypeIncome
		if err := tx.First(&typeIncome, incomeOtherCreate.TypeIncomeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(400, fmt.Sprintf("El tipo de ingreso %d no existe", incomeOtherCreate.TypeIncomeID), err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el tipo de ingreso", err)
		}

		// Buscar caja abierta para el punto de venta
		var register models.CashRegister
		var cashRegisterID *int64
		
		if err := tx.
			Where("is_close = ? AND point_sale_id = ?", false, pointSaleID).
			Order("hour_open DESC").
			First(&register).Error; err != nil {
			// Si no hay caja abierta, cashRegisterID será nil
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(500, "Error al obtener la apertura de caja", err)
			}
		} else {
			cashRegisterID = &register.ID
		}

		// Crear el ingreso
		incomeOther := models.IncomeOther{
			PointSaleID:    &pointSaleID,
			MemberID:       &memberID,
			CashRegisterID: cashRegisterID,
			Total:          incomeOtherCreate.Total,
			TypeIncomeID:   incomeOtherCreate.TypeIncomeID,
			Details:        incomeOtherCreate.Details,
			MethodIncome:   incomeOtherCreate.MethodIncome,
		}

		if err := tx.Create(&incomeOther).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al crear el ingreso", err)
		}

		incomeOtherID = incomeOther.ID
		return nil
	})

	if err != nil {
		return 0, err
	}

	return incomeOtherID, nil
}

// IncomeOtherUpdate actualiza un ingreso existente
func (r *IncomeOtherRepository) IncomeOtherUpdate(memberID, pointSaleID int64, incomeOtherUpdate *schemas.IncomeOtherUpdate) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el ingreso existe y pertenece al punto de venta
		var existingIncome models.IncomeOther
		if err := tx.
			Where("id = ? AND point_sale_id = ?", incomeOtherUpdate.ID, pointSaleID).
			First(&existingIncome).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Ingreso no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el ingreso", err)
		}

		// Verificar que el tipo de ingreso existe
		var typeIncome models.TypeIncome
		if err := tx.First(&typeIncome, incomeOtherUpdate.TypeIncomeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(400, fmt.Sprintf("El tipo de ingreso %d no existe", incomeOtherUpdate.TypeIncomeID), err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el tipo de ingreso", err)
		}

		// Actualizar campos
		existingIncome.Total = incomeOtherUpdate.Total
		existingIncome.TypeIncomeID = incomeOtherUpdate.TypeIncomeID
		existingIncome.Details = incomeOtherUpdate.Details
		existingIncome.MethodIncome = incomeOtherUpdate.MethodIncome

		if err := tx.Save(&existingIncome).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al actualizar el ingreso", err)
		}

		return nil
	})
}

// IncomeOtherDelete elimina un ingreso
func (r *IncomeOtherRepository) IncomeOtherDelete(incomeOtherID, pointSaleID int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el ingreso existe y pertenece al punto de venta
		var existingIncome models.IncomeOther
		if err := tx.
			Where("id = ? AND point_sale_id = ?", incomeOtherID, pointSaleID).
			First(&existingIncome).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Ingreso no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el ingreso", err)
		}

		// Eliminar el ingreso
		if err := tx.Delete(&existingIncome).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al eliminar el ingreso", err)
		}

		return nil
	})
}