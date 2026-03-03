package repository


import (
	"fmt"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// IncomeOtherGetByID obtiene un ingreso por ID
func (r *IncomeOtherRepository) IncomeOtherGetByID(id int64, pointSaleId *int64) (*schemas.IncomeOtherResponse, error) {
	var incomeOther models.IncomeOther

	if pointSaleId != nil {
		if err := r.DB.
			Preload("Member", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "first_name", "last_name", "username").Unscoped()
			}).
			Preload("TypeIncome").
			Preload("PointSale").
			Where("point_sale_id = ?", *pointSaleId).
			First(&incomeOther, id).Error; err != nil {
			return nil, schemas.HandlerErrorGorm(err, "Otros ingresos", schemas.Read)
		}
	} else {
		if err := r.DB.
			Preload("Member", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "first_name", "last_name", "username").Unscoped()
			}).
			Preload("TypeIncome").
			First(&incomeOther, id).Error; err != nil {
			return nil, schemas.HandlerErrorGorm(err, "Otros ingresos", schemas.Read)
		}
	}

	var incomeSchema schemas.IncomeOtherResponse
	copier.Copy(&incomeSchema, &incomeOther)

	return &incomeSchema, nil
}

func (r *IncomeOtherRepository) IncomeOtherGetByDate(pointSaleID *int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.IncomeOtherResponse, int64, error) {
	var incomesOther []*models.IncomeOther

	offset := (page - 1) * limit

	query := r.DB.Where("created_at BETWEEN ? AND ?", fromDate, toDate)

	// Si se proporciona pointSaleID, filtrar por punto de venta
	if pointSaleID != nil {
		query = query.Where("point_sale_id = ?", *pointSaleID)
	} else {
		query = query.Preload("PointSale")
	}

	if err := query.
		Preload("Member", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "username").Unscoped()
		}).
		Preload("TypeIncome").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&incomesOther).Error; err != nil {
		return nil, 0, schemas.HandlerErrorGorm(err, "Otros ingresos", schemas.Read)
	}

	// Contar total
	var total int64
	countQuery := r.DB.Model(&models.IncomeOther{}).
		Where("created_at BETWEEN ? AND ?", fromDate, toDate)

	if pointSaleID != nil {
		countQuery = countQuery.Where("point_sale_id = ?", *pointSaleID)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, schemas.HandlerErrorGorm(err, "Otros ingresos", schemas.Read)
	}

	var incomeSchema []*schemas.IncomeOtherResponse
	copier.Copy(&incomeSchema, &incomesOther)

	return incomeSchema, total, nil
}

func (r *IncomeOtherRepository) IncomeOtherCreate(memberID int64, pointSaleID *int64, incomeOtherCreate *schemas.IncomeOtherCreate) (int64, error) {
	var incomeOtherSave models.IncomeOther

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var typeIncome models.TypeIncome
		if err := tx.First(&typeIncome, incomeOtherCreate.TypeIncomeID).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Tipo de ingreso", schemas.Read)
		}

		incomeOther := models.IncomeOther{
			PointSaleID:  pointSaleID,
			MemberID:     &memberID,
			Total:        incomeOtherCreate.Total,
			TypeIncomeID: incomeOtherCreate.TypeIncomeID,
			Details:      incomeOtherCreate.Details,
			MethodIncome: incomeOtherCreate.MethodIncome,
		}

		if pointSaleID == nil {
			if err := tx.Create(&incomeOther).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Otros ingresos", schemas.Create)
			}

			incomeOtherSave = incomeOther
			return nil
		}

		var register models.CashRegister
		if err := tx.
			Where("is_close = ? AND point_sale_id = ?", false, pointSaleID).
			Order("hour_open DESC").
			First(&register).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Apertura de caja", schemas.Read)
		}

		incomeOther.CashRegisterID = &register.ID

		if err := tx.Create(&incomeOther).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Otros ingresos", schemas.Create)
		}

		incomeOtherSave = incomeOther
		return nil
	})

	if err != nil {
		return 0, err
	}

	return incomeOtherSave.ID, nil
}

// IncomeOtherUpdate actualiza un ingreso existente
func (r *IncomeOtherRepository) IncomeOtherUpdate(memberID int64, pointSaleID *int64, incomeOtherUpdate *schemas.IncomeOtherUpdate) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var existingIncome models.IncomeOther
		if pointSaleID != nil {
			if err := tx.
				Where("id = ? AND point_sale_id = ?", incomeOtherUpdate.ID, pointSaleID).
				First(&existingIncome).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Otros ingresos", schemas.Read)
			}
		} else {
			if err := tx.
				Where("id = ?", incomeOtherUpdate.ID).
				First(&existingIncome).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Otros ingresos", schemas.Read)
			}
		}

		var typeIncome models.TypeIncome
		if err := tx.First(&typeIncome, incomeOtherUpdate.TypeIncomeID).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Tipo de ingreso", schemas.Read)
		}

		existingIncome.Total = incomeOtherUpdate.Total
		existingIncome.TypeIncomeID = incomeOtherUpdate.TypeIncomeID
		existingIncome.Details = incomeOtherUpdate.Details
		existingIncome.MethodIncome = incomeOtherUpdate.MethodIncome

		if err := tx.Save(&existingIncome).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Otros ingresos", schemas.Update)
		}

		return nil
	})

	return err
}

func (r *IncomeOtherRepository) IncomeOtherDelete(memberID int64, incomeOtherID int64, pointSaleID *int64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var existingIncome models.IncomeOther
		if pointSaleID != nil {
			if err := tx.
				Where("id = ? AND point_sale_id = ?", incomeOtherID, pointSaleID).
				First(&existingIncome).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Otros ingresos", schemas.Read)
			}
		} else {
			if err := tx.
				Where("id = ?", incomeOtherID).
				First(&existingIncome).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Otros ingresos", schemas.Read)
			}
		}

		if err := tx.Delete(&existingIncome).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Otros ingresos", schemas.Delete)
		}
		return nil
	})

	return err
}
