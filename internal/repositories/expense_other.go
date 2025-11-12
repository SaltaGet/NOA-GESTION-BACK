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

// ExpenseOtherGetByID obtiene un egreso por ID
func (r *ExpenseOtherRepository) ExpenseOtherGetByID(id int64) (*schemas.ExpenseOtherResponse, error) {
	var expenseOther models.ExpenseOther

	if err := r.DB.
		Preload("Member", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "username")
		}).
		Preload("PointSale", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "address", "is_deposit")
		}).
		Preload("TypeExpense").
		First(&expenseOther, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Egreso no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error al obtener el egreso", err)
	}

	var expenseSchema schemas.ExpenseOtherResponse
	copier.Copy(&expenseSchema, &expenseOther)

	return &expenseSchema, nil
}

// ExpenseOtherGetByDate obtiene egresos por rango de fechas
func (r *ExpenseOtherRepository) ExpenseOtherGetByDate(pointSaleID *int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.ExpenseOtherResponseDTO, int64, error) {
	var expensesOther []*models.ExpenseOther

	offset := (page - 1) * limit

	query := r.DB.Where("created_at BETWEEN ? AND ?", fromDate, toDate)

	// Si se proporciona pointSaleID, filtrar por punto de venta
	if pointSaleID != nil {
		query = query.Where("point_sale_id = ?", *pointSaleID)
	}

	if err := query.
		Select("id", "cash_register_id", "details", "total", "pay_method", "created_at").
		Preload("TypeExpense").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&expensesOther).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "Error al obtener los egresos", err)
	}

	// Contar total
	var total int64
	countQuery := r.DB.Model(&models.ExpenseOther{}).
		Where("created_at BETWEEN ? AND ?", fromDate, toDate)

	if pointSaleID != nil {
		countQuery = countQuery.Where("point_sale_id = ?", *pointSaleID)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "Error al contar los egresos", err)
	}

	var expenseSchema []*schemas.ExpenseOtherResponseDTO
	copier.Copy(&expenseSchema, &expensesOther)

	return expenseSchema, total, nil
}

// ExpenseOtherCreate crea un nuevo egreso
func (r *ExpenseOtherRepository) ExpenseOtherCreate(memberID, pointSaleID int64, expenseOtherCreate *schemas.ExpenseOtherCreate) (int64, error) {
	var expenseOtherID int64

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el tipo de egreso existe
		var typeExpense models.TypeExpense
		if err := tx.First(&typeExpense, expenseOtherCreate.TypeExpenseID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(400, fmt.Sprintf("El tipo de egreso %d no existe", expenseOtherCreate.TypeExpenseID), err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el tipo de egreso", err)
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

		// Crear el egreso
		expenseOther := models.ExpenseOther{
			PointSaleID:    &pointSaleID,
			MemberID:       memberID,
			CashRegisterID: cashRegisterID,
			Details:    expenseOtherCreate.Details,
			TypeExpenseID:  expenseOtherCreate.TypeExpenseID,
			Total:          expenseOtherCreate.Total,
			PayMethod:      expenseOtherCreate.PayMethod,
		}

		if err := tx.Create(&expenseOther).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al crear el egreso", err)
		}

		expenseOtherID = expenseOther.ID
		return nil
	})

	if err != nil {
		return 0, err
	}

	return expenseOtherID, nil
}

// ExpenseOtherUpdate actualiza un egreso existente
func (r *ExpenseOtherRepository) ExpenseOtherUpdate(memberID, pointSaleID int64, expenseOtherUpdate *schemas.ExpenseOtherUpdate) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el egreso existe y pertenece al punto de venta
		var existingExpense models.ExpenseOther
		if err := tx.
			Where("id = ? AND point_sale_id = ?", expenseOtherUpdate.ID, pointSaleID).
			First(&existingExpense).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Egreso no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el egreso", err)
		}

		// Actualizar campos
		existingExpense.Details = expenseOtherUpdate.Details
		existingExpense.Total = expenseOtherUpdate.Total
		existingExpense.PayMethod = expenseOtherUpdate.PayMethod
		existingExpense.MemberID = memberID

		if err := tx.Save(&existingExpense).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al actualizar el egreso", err)
		}

		return nil
	})
}

// ExpenseOtherDelete elimina un egreso
func (r *ExpenseOtherRepository) ExpenseOtherDelete(expenseOtherID, pointSaleID int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el egreso existe y pertenece al punto de venta
		var existingExpense models.ExpenseOther
		if err := tx.
			Where("id = ? AND point_sale_id = ?", expenseOtherID, pointSaleID).
			First(&existingExpense).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Egreso no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el egreso", err)
		}

		// Eliminar el egreso
		if err := tx.Delete(&existingExpense).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al eliminar el egreso", err)
		}

		return nil
	})
}
