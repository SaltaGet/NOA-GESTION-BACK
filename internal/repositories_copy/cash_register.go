package repositories_copy

import (
	"fmt"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func (r *CashRegisterRepository) CashRegisterExistOpen(pointSaleID int64) (bool, error) {
	var existCashRegisterOpen float64
	if err := r.DB.
		Model(&models.CashRegister{}).
		Select("count(*)").
		Where("is_close = ? AND point_sale_id = ?", false, pointSaleID).
		Scan(&existCashRegisterOpen).Error; err != nil {
		return false, schemas.HandlerErrorGorm(err, "Caja", schemas.Read)
	}

	if existCashRegisterOpen > 0 {
		return true, nil
	}

	return false, nil
}

func (r *CashRegisterRepository) CashRegisterGetByID(pointSaleID, id int64) (*schemas.CashRegisterFullResponse, error) {
	var register models.CashRegister
	if err := r.DB.
		Preload("MemberOpen", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Preload("MemberClose", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Where("id = ? AND point_sale_id = ?", id, pointSaleID).
		First(&register).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Caja", schemas.Read)
	}

	var cashRegisterResponse schemas.CashRegisterFullResponse
	_ = copier.Copy(&cashRegisterResponse, &register)

	var incomesModel []models.IncomeSale
	if err := r.DB.Select("id", "total", "is_budget", "created_at").
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "income_sale_id", "amount", "total", "product_id", "created_at")
		}).
		Preload("Items.Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "code", "name", "price")
		}).
		Preload("Pay", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "total", "method_pay", "income_sale_id")
		}).
		Where("cash_register_id = ? AND point_sale_id = ?", id, pointSaleID).
		Find(&incomesModel).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Caja", schemas.Read)
	}
	var incomes []*schemas.IncomeSaleSimpleResponse
	_ = copier.Copy(&incomes, &incomesModel)

	var incomeOtherModel []models.IncomeOther
	if err := r.DB.Select("id", "total", "details", "method_income", "created_at", "type_income_id").
		Preload("Member", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "username").Unscoped()
		}).
		Preload("TypeIncome", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Where("cash_register_id = ? AND point_sale_id = ?", id, pointSaleID).
		Find(&incomeOtherModel).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Caja", schemas.Read)
	}

	var incomeOther []*schemas.IncomeOtherResponse
	_ = copier.Copy(&incomeOther, &incomeOtherModel)

	var expensesOtherModel []models.ExpenseOther
	if err := r.DB.Select("id", "total", "cash_register_id", "details", "pay_method", "created_at", "member_id").
		Preload("Member", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "username").Unscoped()
		}).
		Preload("TypeExpense").
		Where("cash_register_id = ? AND point_sale_id = ?", id, pointSaleID).
		Find(&expensesOtherModel).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Caja", schemas.Read)
	}
	var expenseOtherResponse []*schemas.ExpenseOtherResponse
	_ = copier.Copy(&expenseOtherResponse, &expensesOtherModel)

	cashRegisterResponse.IncomeSale = incomes
	cashRegisterResponse.IncomeOther = incomeOther
	// cashRegisterResponse.ExpenseBuy = &expenseBuyResponseSimple
	cashRegisterResponse.ExpenseOther = expenseOtherResponse

	return &cashRegisterResponse, nil
}

func (r *CashRegisterRepository) CashRegisterOpen(pointSaleID int64, userID int64, amountOpen schemas.CashRegisterOpen) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", userID)).Error; err != nil {
			return err
		}

		var existRegisterOpen float64
		if err := tx.
			Model(&models.CashRegister{}).
			Select("count(*)").
			Where("is_close = ? AND point_sale_id = ?", false, pointSaleID).
			Scan(&existRegisterOpen).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Caja", schemas.Read)
		}

		if existRegisterOpen > 0 {
			return schemas.ErrorResponse(400, "ya existe una apertura de caja, antes de continuar cierre la caja", fmt.Errorf("ya existe una apertura de caja, antes de continuar cerrar"))
		}

		registerOpen := models.CashRegister{
			PointSaleID:  pointSaleID,
			MemberOpenID: userID,
			OpenAmount:   amountOpen.OpenAmount,
			HourOpen:     time.Now().UTC(),
		}

		if err := tx.Create(&registerOpen).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Caja", schemas.Create)
		}

		return nil
	})
}

func (r *CashRegisterRepository) CashRegisterClose(pointSaleID int64, userID int64, amountOpen schemas.CashRegisterClose) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", userID)).Error; err != nil {
			return err
		}

		var register models.CashRegister
		if err := tx.
			Where("is_close = ? AND point_sale_id = ?", false, pointSaleID).
			Order("hour_open DESC").
			First(&register).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Caja", schemas.Read)
		}

		var member models.Member
		if err := tx.Preload("Role").First(&member, userID).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Caja", schemas.Read)
		}

		now := time.Now().UTC()
		register.CloseAmount = &amountOpen.CloseAmount
		register.IsClose = true
		register.HourClose = &now
		register.MemberCloseID = &userID

		if err := tx.Save(&register).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Caja", schemas.Update)
		}

		return nil
	})
}

func (r *CashRegisterRepository) CashRegisterInform(pointSaleID int64, userID int64, fromDate, toDate time.Time) ([]*schemas.CashRegisterInformResponse, error) {
	var registers []*models.CashRegister
	if err := r.DB.
		Preload("MemberOpen", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Preload("MemberClose", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Where("point_sale_id = ? AND created_at >= ? AND created_at <= ?", pointSaleID, fromDate, toDate).
		Order("created_at DESC").
		Find(&registers).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Caja", schemas.Read)
	}

	var cashRegisterInformResponse []*schemas.CashRegisterInformResponse
	_ = copier.Copy(&cashRegisterInformResponse, &registers)

	for _, register := range cashRegisterInformResponse {
		type total struct {
			Cash  float64 `json:"total"`
			Other float64 `json:"other"`
		}

		var incomes total
		if err := r.DB.Model(&models.PayIncome{}).
			Select(`
			SUM(CASE WHEN method_pay = 'cash' THEN total ELSE 0 END) AS cash,
			SUM(CASE WHEN method_pay <> 'cash' AND method_pay <> 'credit' THEN total ELSE 0 END) AS other
		`).
			Where("cash_register_id = ?", register.ID).
			Scan(&incomes).Error; err != nil {
			return nil, schemas.HandlerErrorGorm(err, "Caja", schemas.Read)
		}

		var incomeOther total
		if err := r.DB.Model(&models.IncomeOther{}).
			Select(`
			SUM(CASE WHEN method_income = 'cash' THEN total ELSE 0 END) AS cash,
			SUM(CASE WHEN method_income <> 'cash' AND method_income <> 'credit' THEN total ELSE 0 END) AS other
		`).
			Where("cash_register_id = ?", register.ID).
			Find(&incomeOther).Error; err != nil {
			return nil, schemas.HandlerErrorGorm(err, "Caja", schemas.Read)
		}

		var expenseBuy total
		if err := r.DB.Model(&models.PayExpenseBuy{}).
			Select(`
			SUM(CASE WHEN method_pay = 'cash' THEN total ELSE 0 END) AS cash,
			SUM(CASE WHEN method_pay <> 'cash' AND method_pay <> 'credit' THEN total ELSE 0 END) AS other
		`).
			Where("cash_register_id = ?", register.ID).
			Scan(&expenseBuy).Error; err != nil {
			return nil, schemas.HandlerErrorGorm(err, "Caja", schemas.Read)
		}

		var expenseOther total
		if err := r.DB.Model(&models.PayExpenseOther{}).
			Select(`
			SUM(CASE WHEN method_pay = 'cash' THEN total ELSE 0 END) AS cash,
			SUM(CASE WHEN method_pay <> 'cash' AND method_pay <> 'credit' THEN total ELSE 0 END) AS other
		`).
			Where("cash_register_id = ?", register.ID).
			Scan(&expenseOther).Error; err != nil {
			return nil, schemas.HandlerErrorGorm(err, "Caja", schemas.Read)
		}

		totalIncomesCash := incomes.Cash + incomeOther.Cash
		totalIncomesOthers := incomes.Other + incomeOther.Other
		totalExpenseCash := expenseBuy.Cash + expenseOther.Cash
		totalExpenseOther := expenseBuy.Other + expenseOther.Other

		register.TotalIncomeCash = &totalIncomesCash
		register.TotalIncomeOthers = &totalIncomesOthers
		register.TotalExpenseCash = &totalExpenseCash
		register.TotalExpenseOthers = &totalExpenseOther
	}

	return cashRegisterInformResponse, nil
}
