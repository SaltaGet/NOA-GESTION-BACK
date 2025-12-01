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

func (r *ExpenseOtherRepository) ExpenseOtherGetByID(id int64, pointSaleID *int64) (*schemas.ExpenseOtherResponse, error) {
	var expenseOther models.ExpenseOther

	if pointSaleID != nil {
		if err := r.DB.
			Preload("Member", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "first_name", "last_name", "username")
			}).
			Preload("TypeExpense").
			Where("point_sale_id = ?", *pointSaleID).
			First(&expenseOther, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, schemas.ErrorResponse(404, "Egreso no encontrado", err)
			}
			return nil, schemas.ErrorResponse(500, "Error al obtener el egreso", err)
		}
	} else {
		if err := r.DB.
			Preload("Member", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "first_name", "last_name", "username")
			}).
			Preload("TypeExpense").
			First(&expenseOther, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, schemas.ErrorResponse(404, "Egreso no encontrado", err)
			}
			return nil, schemas.ErrorResponse(500, "Error al obtener el egreso", err)
		}
	}

	var expenseSchema schemas.ExpenseOtherResponse
	copier.Copy(&expenseSchema, &expenseOther)

	return &expenseSchema, nil
}

func (r *ExpenseOtherRepository) ExpenseOtherGetByDate(pointSaleID *int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.ExpenseOtherResponse, int64, error) {
	var expensesOther []*models.ExpenseOther

	offset := (page - 1) * limit

	query := r.DB.Where("created_at BETWEEN ? AND ?", fromDate, toDate)

	if pointSaleID != nil {
		query = query.Where("point_sale_id = ?", *pointSaleID)
	} else {
		query = query.Preload("PointSale")
	}

	if err := query.
		Preload("Member", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "username")
		}).
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

	var expenseSchema []*schemas.ExpenseOtherResponse
	copier.Copy(&expenseSchema, &expensesOther)

	return expenseSchema, total, nil
}

func (r *ExpenseOtherRepository) ExpenseOtherCreate(memberID int64, pointSaleID *int64, expenseOtherCreate *schemas.ExpenseOtherCreate) (int64, error) {
	var expenseOtherID int64

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var typeExpense models.TypeExpense
		if err := tx.First(&typeExpense, expenseOtherCreate.TypeExpenseID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(400, fmt.Sprintf("El tipo de egreso %d no existe", expenseOtherCreate.TypeExpenseID), err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el tipo de egreso", err)
		}

		expenseOther := models.ExpenseOther{
			MemberID:      memberID,
			Details:       expenseOtherCreate.Details,
			TypeExpenseID: expenseOtherCreate.TypeExpenseID,
			Total:         expenseOtherCreate.Total,
			PayMethod:     expenseOtherCreate.PayMethod,
		}

		if pointSaleID == nil {
			if err := tx.Create(&expenseOther).Error; err != nil {
				return schemas.ErrorResponse(500, "Error al crear el egreso", err)
			}

			expenseOtherID = expenseOther.ID
			return nil
		}

		var register models.CashRegister
		if err := tx.
			Where("is_close = ? AND point_sale_id = ?", false, pointSaleID).
			Order("hour_open DESC").
			First(&register).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(400, "No hay una caja abierta para el punto de venta", err)
			}
			return schemas.ErrorResponse(500, "Error al obtener la apertura de caja", err)
		}

		expenseOther.PointSaleID = pointSaleID
		expenseOther.CashRegisterID = &register.ID

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

func (r *ExpenseOtherRepository) ExpenseOtherUpdate(memberID int64, pointSaleID *int64, expenseOtherUpdate *schemas.ExpenseOtherUpdate) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var existingExpense models.ExpenseOther
		if pointSaleID == nil {
			if err := tx.
				Where("id = ?", expenseOtherUpdate.ID).
				First(&existingExpense).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return schemas.ErrorResponse(404, "Egreso no encontrado", err)
				}
				return schemas.ErrorResponse(500, "Error al obtener el egreso", err)
			}
		} else {
			if err := tx.
				Where("id = ? AND point_sale_id = ?", expenseOtherUpdate.ID, pointSaleID).
				First(&existingExpense).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return schemas.ErrorResponse(404, "Egreso no encontrado", err)
				}
				return schemas.ErrorResponse(500, "Error al obtener el egreso", err)
			}
		}

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

func (r *ExpenseOtherRepository) ExpenseOtherDelete(expenseOtherID int64, pointSaleID *int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var existingExpense models.ExpenseOther
		if pointSaleID == nil {
			if err := tx.
				Where("id = ?", expenseOtherID).
				First(&existingExpense).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return schemas.ErrorResponse(404, "Egreso no encontrado", err)
				}
				return schemas.ErrorResponse(500, "Error al obtener el egreso", err)
			}
		} else {
			if err := tx.
				Where("id = ? AND point_sale_id = ?", expenseOtherID, pointSaleID).
				First(&existingExpense).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return schemas.ErrorResponse(404, "Egreso no encontrado", err)
				}
				return schemas.ErrorResponse(500, "Error al obtener el egreso", err)
			}
		}

		if err := tx.Delete(&existingExpense).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al eliminar el egreso", err)
		}

		return nil
	})
}
