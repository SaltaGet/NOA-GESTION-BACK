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

func (r *CashRegisterRepository) CashRegisterExistOpen(pointSaleID int64) (bool, error) {
	var existCashRegisterOpen float64
	if err := r.DB.
		Model(&models.CashRegister{}).
		Select("count(*)").
		Where("is_close = ? AND point_sale_id = ?", false, pointSaleID).
		Scan(&existCashRegisterOpen).Error; err != nil {
		return false, schemas.ErrorResponse(500, "error al contar aperturas de caja", err)
	}

	if existCashRegisterOpen > 0 {
		return true, nil
	}

	return false, nil
}

func (r *CashRegisterRepository) CashRegisterGetByID(pointSaleID, id int64) (*schemas.CashRegisterFullResponse, error) {
	var register models.CashRegister
	if err := r.DB.
		Preload("MemberOpen").
		Preload("MemberClose").
		Where("id = ? AND point_sale_id = ?", id, pointSaleID).
		First(&register).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "caja no encontrada", err)
		}
		return nil, schemas.ErrorResponse(500, "error obtener caja", err)
	}

	var cashRegisterResponse schemas.CashRegisterFullResponse
	_ = copier.Copy(&cashRegisterResponse, &register)

	var incomesModel []models.IncomeSale
	if err := r.DB.Select("id", "total", "is_budget", "created_at").
		Preload("Items.Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "code", "name")
		}).
		Preload("Pay", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "amount", "method_pay", "income_sale_id")
		}).
		Where("cash_register_id = ? AND point_sale_id = ?", id, pointSaleID).
		Find(&incomesModel).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "error al obtener ingresos de caja", err)
	}
	var incomes []schemas.IncomeSaleSimpleResponse
	_ = copier.Copy(&incomes, &incomesModel)

	var incomeOtherModel []models.IncomeOther
	if err := r.DB.Select("id", "total", "details", "method_income", "created_at").
		Preload("Member", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "username")
		}).
		Preload("TypeIncome", func(db *gorm.DB) *gorm.DB {
			return db.Select("name")
		}).
		Where("cash_register_id = ? AND point_sale_id = ?", id, pointSaleID).
		Find(&incomeOtherModel).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "error al obtener ingresos de cancha de caja", err)
	}

	var incomeOther []schemas.IncomeOtherResponse
	_ = copier.Copy(&incomeOther, &incomeOtherModel)

	var expenseBuyModel []models.ExpenseBuy
	if err := r.DB.Select("id", "total", "details", "method_pay", "created_at").
		Preload("Supplier", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "company_name")
		}).
		Preload("PayExpenseBuy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "amount", "method_pay")
		}).
		Where("cash_register_id = ? AND point_sale_id = ?", id, pointSaleID).
		Find(&expenseBuyModel).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "error al obtener egresos de compra de caja", err)
	}

	var expenseBuyResponseSimple []schemas.ExpenseBuyResponseSimple
	_ = copier.Copy(&expenseBuyResponseSimple, &expenseBuyModel)

	var expensesOtherModel []models.ExpenseOther
	if err := r.DB.Select("id", "total", "register_id", "description", "pay_method", "created_at").
		Where("cash_register_id = ? AND point_sale_id = ?", id, pointSaleID).
		Find(&expensesOtherModel).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "error al obtener egresos de caja", err)
	}
	var expenseOtherResponse []schemas.ExpenseOtherResponse
	_ = copier.Copy(&expenseOtherResponse, &expensesOtherModel)

	cashRegisterResponse.IncomeSale = &incomes
	cashRegisterResponse.IncomeOther = &incomeOther
	cashRegisterResponse.ExpenseBuy = &expenseBuyResponseSimple
	cashRegisterResponse.ExpenseOther = &expenseOtherResponse

	return &cashRegisterResponse, nil
}

func (r *CashRegisterRepository) CashRegisterOpen(pointSaleID int64, userID int64, amountOpen schemas.CashRegisterOpen) error {
	var existRegisterOpen float64

	if err := r.DB.
		Model(&models.CashRegister{}).
		Select("count(*)").
		Where("is_close = ? AND point_sale_id = ?", false, pointSaleID).
		Scan(&existRegisterOpen).Error; err != nil {
		return schemas.ErrorResponse(500, "error al contar aperturas de caja", err)
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

	if err := r.DB.Create(&registerOpen).Error; err != nil {
		return schemas.ErrorResponse(500, "error al registrar la apertura de caja", err)
	}

	return nil
}

func (r *CashRegisterRepository) CashRegisterClose(pointSaleID int64, userID int64, amountOpen schemas.CashRegisterClose) error {
	var register models.CashRegister
	if err := r.DB.
		Where("is_close = ? AND point_sale_id = ?", false, pointSaleID).
		Order("hour_open DESC").
		First(&register).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "No se encontraron aperturas de caja", err)
		}
		return schemas.ErrorResponse(500, "error al obtener la apertura de caja", err)
	}

	var member models.Member
	if err := r.DB.Preload("Role").First(&member, userID).Error; err != nil {
		return schemas.ErrorResponse(404, "usuario no encontrado", err)
	}

	if member.Role.Name != "admin" && member.ID != register.MemberOpenID {
		return schemas.ErrorResponse(403, "no tienes permiso para cerrar la caja, solo el creador o el admin", fmt.Errorf("no tienes permiso para cerrar la caja, solo el creador o el admin"))
	}

	now := time.Now().UTC()
	register.CloseAmount = &amountOpen.CloseAmount
	register.IsClose = true
	register.HourClose = &now
	register.MemberCloseID = &userID

	if err := r.DB.Save(&register).Error; err != nil {
		return schemas.ErrorResponse(500, "error al cerrar la caja", err)
	}

	return nil
}

func (r *CashRegisterRepository) CashRegisterInform(pointSaleID int64, userID int64, fromDate, toDate time.Time) ([]*schemas.CashRegisterInformResponse, error) {
	var registers []*models.CashRegister
	if err := r.DB.
		Preload("MemberOpen").
		Preload("MemberClose").
		Where("point_sale_id = ? AND created_at >= ? AND created_at <= ?", pointSaleID, fromDate, toDate).
		Order("created_at DESC").
		Find(&registers).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "error al obtener aperturas de caja", err)
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
			return nil, schemas.ErrorResponse(500, "error al obtener ingresos por ventas", err)
		}

		var incomeOther total
		if err := r.DB.Model(&models.IncomeOther{}).
			Select(`
			SUM(CASE WHEN method_income = 'cash' THEN total ELSE 0 END) AS cash,
			SUM(CASE WHEN method_income <> 'cash' AND method_income <> 'credit' THEN total ELSE 0 END) AS other
		`).
			Where("cash_register_id = ?", register.ID).
			Find(&incomeOther).Error; err != nil {
			return nil, schemas.ErrorResponse(500, "error al obtener otros ingresos", err)
		}

		var expenseBuy total
		if err := r.DB.Model(&models.PayExpenseBuy{}).
			Select(`
			SUM(CASE WHEN method_pay = 'cash' THEN total ELSE 0 END) AS cash,
			SUM(CASE WHEN method_pay <> 'cash' AND method_pay <> 'credit' THEN total ELSE 0 END) AS other
		`).
			Where("cash_register_id = ?", register.ID).
			Scan(&expenseBuy).Error; err != nil {
			return nil, schemas.ErrorResponse(500, "error al obtener egresos por compras", err)
		}

		var expenseOther total
		if err := r.DB.Model(&models.PayExpenseOther{}).
			Select(`
			SUM(CASE WHEN method_pay = 'cash' THEN total ELSE 0 END) AS cash,
			SUM(CASE WHEN method_pay <> 'cash' AND method_pay <> 'credit' THEN total ELSE 0 END) AS other
		`).
			Where("cash_register_id = ?", register.ID).
			Scan(&expenseOther).Error; err != nil {
				return nil, schemas.ErrorResponse(500, "error al obtener otros egresos", err)
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
