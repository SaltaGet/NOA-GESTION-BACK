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

	var incomes []models.Income
	if err := r.DB.
		Preload("Items").
		Preload("Items.Product").
		Where("register_id = ?", id).
		Find(&incomes).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "error al obtener ingresos de caja", err)
	}

	var expenses []models.Expense
	if err := r.DB.
		Where("register_id = ?", id).
		Find(&expenses).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "error al obtener egresos de caja", err)
	}

	var IncomeSportsCourts []models.IncomeSportsCourts
	if err := r.DB.
		Preload("SportsCourt").
		Where("partial_register_id = ? OR rest_register_id = ?", id, id).
		Find(&IncomeSportsCourts).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "error al obtener ingresos de cancha de caja", err)
	}

	return &register, nil
}

func (r *CashRegisterRepository) CashRegisterOpen(pointSaleID uint, userID uint, amountOpen schemas.RegisterOpen) error {
	var existRegisterOpen float64

	if err := r.DB.
		Model(&models.Register{}).
		Select("count(*)").
		Where("is_close = ? AND point_sale_id = ?", false, pointSaleID).
		Scan(&existRegisterOpen).Error; err != nil {
		return schemas.ErrorResponse(500, "error al contar aperturas de caja", err)
	}

	if existRegisterOpen > 0 {
		return schemas.ErrorResponse(400, "ya existe una apertura de caja, antes de continuar cierre la caja", fmt.Errorf("ya existe una apertura de caja, antes de continuar cerrar"))
	}

	registerOpen := models.Register{
		PointSaleID: pointSaleID,
		UserOpenID:  userID,
		OpenAmount:  amountOpen.OpenAmount,
		HourOpen:    time.Now().UTC(),
	}

	if err := r.DB.Create(&registerOpen).Error; err != nil {
		return schemas.ErrorResponse(500, "error al registrar la apertura de caja", err)
	}

	return nil
}

func (r *CashRegisterRepository) CashRegisterClose(pointSaleID uint, userID uint, amountOpen schemas.RegisterClose) error {
	var register models.Register
	if err := r.DB.
		Where("is_close = ? AND point_sale_id = ?", false, pointSaleID).
		Order("hour_open DESC").
		First(&register).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "No se encontraron aperturas de caja", err)
		}
		return schemas.ErrorResponse(500, "error al obtener la apertura de caja", err)
	}

	var user models.User
	if err := r.DB.Preload("Role").First(&user, userID).Error; err != nil {
		return schemas.ErrorResponse(404, "usuario no encontrado", err)
	}

	if user.Role.Name != "admin" && user.ID != register.UserOpenID {
		return schemas.ErrorResponse(403, "no tienes permiso para cerrar la caja, solo el creador o el admin", fmt.Errorf("no tienes permiso para cerrar la caja, solo el creador o el admin"))
	}

	// var totalsIncome schemas.Totals
	// if err := r.DB.
	// 	Model(&models.Income{}).
	// 	Select(`
	// 	COALESCE(SUM(CASE WHEN payment_method = 'efectivo' THEN COALESCE(total,0) ELSE 0	END),0) AS cash,
	// 	COALESCE(SUM(CASE WHEN payment_method IN ('tarjeta','transferencia') THEN COALESCE(total,0) ELSE 0 END),0) AS others
	// `).
	// 	Where("register_id = ?", register.ID).
	// 	Scan(&totalsIncome).Error; err != nil {
	// 	return schemas.ErrorResponse(500, "error al obtener ingresos por métodos", err)
	// }

	// var totalsIncomeCourts schemas.Totals
	// if err := r.DB.
	// 	Model(&models.IncomeSportsCourts{}).
	// 	Select(`
	// 	COALESCE(
	// 		SUM(
	// 			CASE WHEN partial_payment_method = 'efectivo' THEN COALESCE(partial_pay,0) ELSE 0	END +
	// 			CASE WHEN rest_payment_method = 'efectivo' THEN COALESCE(rest_pay,0) ELSE 0 END
	// 		),0
	// 	) AS cash,
	// 	COALESCE(
	//       SUM(
	//           CASE WHEN partial_payment_method IN ('tarjeta','transferencia') THEN COALESCE(partial_pay,0) ELSE 0 END +
	//           CASE WHEN rest_payment_method IN ('tarjeta','transferencia') THEN COALESCE(rest_pay,0) ELSE 0 END
	//       ), 0
	//   ) AS others
	// `).
	// 	Where("partial_register_id = ? OR rest_register_id = ?", register.ID, register.ID).
	// 	Scan(&totalsIncomeCourts).Error; err != nil {
	// 	return schemas.ErrorResponse(500, "error al obtener ingresos de canchas por métodos", err)
	// }

	// var totalsExpense schemas.Totals
	// if err := r.DB.
	// 	Model(&models.Expense{}).
	// 	Select(`
	// 	COALESCE(SUM(CASE WHEN payment_method = 'efectivo' THEN COALESCE(total,0) ELSE 0	END),0) AS cash,
	// 	COALESCE(SUM(CASE WHEN payment_method IN ('tarjeta','transferencia') THEN COALESCE(total,0) ELSE 0 END),0) AS others
	// `).
	// 	Where("register_id = ?", register.ID).
	// 	Scan(&totalsExpense).Error; err != nil {
	// 	return schemas.ErrorResponse(500, "error al obtener ingresos por métodos", err)
	// }

	now := time.Now().UTC()
	// totalIncomesCash := totalsIncome.Cash + totalsIncomeCourts.Cash
	// totalIncomesOthers := totalsIncome.Others + totalsIncomeCourts.Others
	// totalExpenseCash := totalsExpense.Cash + totalsExpense.Cash
	// totalExpenseOthers := totalsExpense.Others + totalsExpense.Others

	register.CloseAmount = &amountOpen.CloseAmount
	register.IsClose = true
	register.HourClose = &now
	register.UserCloseID = &userID
	// register.TotalIncomeCash = &totalIncomesCash
	// register.TotalIncomeOthers = &totalIncomesOthers
	// register.TotalExpenseCash = &totalExpenseCash
	// register.TotalExpenseOthers = &totalExpenseOthers

	if err := r.DB.Save(&register).Error; err != nil {
		return schemas.ErrorResponse(500, "error al cerrar la caja", err)
	}

	return nil
}

func (r *CashRegisterRepository) CashRegisterInform(pointSaleID uint, userID uint, fromDate, toDate time.Time) ([]*models.Register, error) {
	var registers []*models.Register
	if err := r.DB.
		Preload("UserOpen").
		Preload("UserClose").
		// Preload("PointSale").
		Where("point_sale_id = ? AND created_at >= ? AND created_at <= ?", pointSaleID, fromDate, toDate).
		Order("created_at DESC").
		Find(&registers).Error; err != nil {
		return nil, err
	}

	for _, register := range registers {
		var incomes []models.Income
		if err := r.DB.
			Preload("Items").
			Preload("Items.Product").
			Where("register_id = ?", register.ID).
			Find(&incomes).Error; err != nil {
			return nil, schemas.ErrorResponse(500, "error al obtener ingresos de caja", err)
		}

		var expenses []models.Expense
		if err := r.DB.
			Where("register_id = ?", register.ID).
			Find(&expenses).Error; err != nil {
			return nil, schemas.ErrorResponse(500, "error al obtener egresos de caja", err)
		}

		var incomeSportsCourts []models.IncomeSportsCourts
		if err := r.DB.
			Preload("SportsCourt").
			Where("partial_register_id = ? OR rest_register_id = ?", register.ID, register.ID).
			Find(&incomeSportsCourts).Error; err != nil {
			return nil, schemas.ErrorResponse(500, "error al obtener ingresos de canchas por métodos", err)
		}

		var totalIncomeCash, totalIncomeOther float64
		for _, income := range incomes {
			if income.PaymentMethod == "efectivo" {
				totalIncomeCash += income.Total
			} else {
				totalIncomeOther += income.Total
			}
		}

		var totalExpenseCash, totalExpenseOther float64
		for _, expense := range expenses {
			if expense.PaymentMethod == "efectivo" {
				totalExpenseCash += expense.Total
			} else {
				totalExpenseOther += expense.Total
			}
		}

		var totalIncomeSportsCourtsCash, totalIncomeSportsCourtsOther float64
		for _, income := range incomeSportsCourts {
			if income.PartialPaymentMethod == "efectivo" && income.PartialRegisterID == register.ID {
				totalIncomeSportsCourtsCash += income.PartialPay
			} else if income.PartialPaymentMethod != "efectivo" && income.PartialRegisterID == register.ID {
				totalIncomeSportsCourtsOther += income.PartialPay
			}
			if income.RestPay != nil && income.RestPaymentMethod != nil  && income.RestRegisterID == &register.ID {
				if *income.RestPaymentMethod == "efectivo" {
					totalIncomeSportsCourtsCash += *income.RestPay
				} else {
					totalIncomeSportsCourtsOther += *income.RestPay
				}
			}
		}

		totalIncomesCash := totalIncomeCash + totalIncomeSportsCourtsCash
		totalIncomesOthers := totalIncomeOther + totalIncomeSportsCourtsOther

		register.TotalIncomeCash = &totalIncomesCash
		register.TotalIncomeOthers = &totalIncomesOthers
		register.TotalExpenseCash = &totalExpenseCash
		register.TotalExpenseOthers = &totalExpenseOther

		register.Income = incomes
		register.Expense = expenses
		register.IncomeSportsCourts = incomeSportsCourts
	}

	return registers, nil
}
