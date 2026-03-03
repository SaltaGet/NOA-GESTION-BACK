package repository


import (
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
				return db.Select("id", "first_name", "last_name", "username").Unscoped()
			}).
			Preload("TypeExpense").
			Preload("PointSale").
			Where("point_sale_id = ?", *pointSaleID).
			First(&expenseOther, id).Error; err != nil {
			return nil, schemas.HandlerErrorGorm(err, "Otros egresos", schemas.Read)
		}
	} else {
		if err := r.DB.
			Preload("Member", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "first_name", "last_name", "username").Unscoped()
			}).
			Preload("TypeExpense").
			Preload("PointSale").
			First(&expenseOther, id).Error; err != nil {
			return nil, schemas.HandlerErrorGorm(err, "Otros egresos", schemas.Read)
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
			return db.Select("id", "first_name", "last_name", "username").Unscoped()
		}).
		Preload("TypeExpense").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&expensesOther).Error; err != nil {
		return nil, 0, schemas.HandlerErrorGorm(err, "Otros egresos", schemas.Read)
	}

	// Contar total
	var total int64
	countQuery := r.DB.Model(&models.ExpenseOther{}).
		Where("created_at BETWEEN ? AND ?", fromDate, toDate)

	if pointSaleID != nil {
		countQuery = countQuery.Where("point_sale_id = ?", *pointSaleID)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, schemas.HandlerErrorGorm(err, "Otros egresos", schemas.Read)
	}

	var expenseSchema []*schemas.ExpenseOtherResponse
	copier.Copy(&expenseSchema, &expensesOther)

	return expenseSchema, total, nil
}

func (r *ExpenseOtherRepository) ExpenseOtherCreate(memberID int64, pointSaleID *int64, expenseOtherCreate *schemas.ExpenseOtherCreate) (int64, error) {
	var expenseOtherSave models.ExpenseOther

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var typeExpense models.TypeExpense
		if err := tx.First(&typeExpense, expenseOtherCreate.TypeExpenseID).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Tipo de egreso", schemas.Read)
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
				return schemas.HandlerErrorGorm(err, "Otros egresos", schemas.Create)
			}

			expenseOtherSave = expenseOther
			return nil
		}

		var register models.CashRegister
		if err := tx.
			Where("is_close = ? AND point_sale_id = ?", false, pointSaleID).
			Order("hour_open DESC").
			First(&register).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Caja", schemas.Read)
		}

		expenseOther.PointSaleID = pointSaleID
		expenseOther.CashRegisterID = &register.ID

		if err := tx.Create(&expenseOther).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Otros egresos", schemas.Create)
		}

		expenseOtherSave = expenseOther
		return nil
	})

	if err != nil {
		return 0, err
	}

	return expenseOtherSave.ID, nil
}

func (r *ExpenseOtherRepository) ExpenseOtherUpdate(memberID int64, pointSaleID *int64, expenseOtherUpdate *schemas.ExpenseOtherUpdate) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var existingExpense models.ExpenseOther
		if pointSaleID == nil {
			if err := tx.
				Where("id = ?", expenseOtherUpdate.ID).
				First(&existingExpense).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Otros egresos", schemas.Read)
			}
		} else {
			if err := tx.
				Where("id = ? AND point_sale_id = ?", expenseOtherUpdate.ID, pointSaleID).
				First(&existingExpense).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Otros egresos", schemas.Read)
			}
		}

		var typeExpense models.TypeExpense
		if err := tx.Select("id").First(&typeExpense, expenseOtherUpdate.TypeExpenseID).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Tipo de egreso", schemas.Read)
		}

		existingExpense.Details = expenseOtherUpdate.Details
		existingExpense.Total = expenseOtherUpdate.Total
		existingExpense.PayMethod = expenseOtherUpdate.PayMethod
		existingExpense.MemberID = memberID
		existingExpense.TypeExpenseID = expenseOtherUpdate.TypeExpenseID

		if err := tx.Save(&existingExpense).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Otros egresos", schemas.Update)
		}

		return nil
	})

	return err
}

func (r *ExpenseOtherRepository) ExpenseOtherDelete(memberID int64, expenseOtherID int64, pointSaleID *int64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var existingExpense models.ExpenseOther
		if pointSaleID == nil {
			if err := tx.
				Where("id = ?", expenseOtherID).
				First(&existingExpense).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Otros egresos", schemas.Read)
			}
		} else {
			if err := tx.
				Where("id = ? AND point_sale_id = ?", expenseOtherID, pointSaleID).
				First(&existingExpense).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Otros egresos", schemas.Read)
			}
		}

		// Guardar estado anterior para auditoría
		
		if err := tx.Delete(&existingExpense).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Otros egresos", schemas.Delete)
		}
		
		return nil
	})

	return err
}
