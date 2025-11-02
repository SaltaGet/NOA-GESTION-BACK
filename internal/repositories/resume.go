package repositories

import (
	"encoding/json"
	"time"

	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/DanielChachagua/GestionCar/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EXPENSE
func (r *ResumeRepository) ResumeExpenseGetByDateBetween(fromDate string, toDate string) (*[]models.ResumeExpense, error) {
	return nil, nil
}

func (r *ResumeRepository) ResumeExpenseGetByID(id string) (*models.ResumeExpense, error) {
	return nil, nil
}

func (r *ResumeRepository) ResumeExpenseCreate(resume *models.ResumeExpenseCreate) (string, error) {
	err := r.DB.Create(&resume).Error
	if err != nil {
		return "", err
	}
	return "", nil
}

func (r *ResumeRepository) ResumeExpenseUpdate(resume *models.ResumeExpenseUpdate) error {
	return nil
}

// INCOME

func (r *ResumeRepository) ResumeIncomeGetByDateBetween(fromDate string, toDate string) (*[]models.ResumeIncome, error) {
	return nil, nil
}

func (r *ResumeRepository) ResumeIncomeGetByID(id string) (*models.ResumeIncome, error) {
	return nil, nil
}

func (r *ResumeRepository) ResumeIncomeCreate(resume *models.ResumeIncomeCreate) (string, error) {
	err := r.DB.Create(&resume).Error
	if err != nil {
		return "", err
	}
	return "", nil
}

func (r *ResumeRepository) ResumeIncomeUpdate(resume *models.ResumeIncomeUpdate) error {
	return nil
}





// RESUME

func GenerateResume(db *gorm.DB, firstday, lastDay time.Time) error {
	tx := db.Begin()

	var incomes []models.Income
	if err := tx.Where("created_at BETWEEN ? AND ?", firstday, lastDay).Order("created_at ASC").Find(&incomes).Error; err != nil {
		return models.ErrorResponse(500, "Error interno al buscar ingresos", err)
	}

	var expenses []models.Expense
	if err := tx.Where("created_at BETWEEN ? AND ?", firstday, lastDay).Order("created_at ASC").Find(&expenses).Error; err != nil {
		return models.ErrorResponse(500, "Error interno al buscar egresos", err)
	}

	jsonIncome, err := json.Marshal(incomes)
	if err != nil {
		return models.ErrorResponse(500, "Error al serializar los ingresos", err)
	}

	jsonExpense, err := json.Marshal(expenses)
	if err != nil {
		return models.ErrorResponse(500, "Error al serializar los egresos", err)
	}

	incomesCompressed, err := utils.CompressToBase64Bytes(jsonIncome)
	if err != nil {
		return models.ErrorResponse(500, "Error al comprimir los ingresos", err)
	}

	expensesCompressed, err := utils.CompressToBase64Bytes(jsonExpense)
	if err != nil {
		return models.ErrorResponse(500, "Error al comprimir los egresos", err)
	}

	if err := tx.Create(&models.ResumeIncome{
		ID:   uuid.NewString(),
		Data: incomesCompressed,
		Date: firstday,
	}).Error; err != nil {
		tx.Rollback()
		return models.ErrorResponse(500, "Error interno al crear el resumen de ingresos", err)
	}

	if err := tx.Create(&models.ResumeExpense{
		ID:   uuid.NewString(),
		Data: expensesCompressed,
		Date: firstday,
	}).Error; err != nil {
		tx.Rollback()
		return models.ErrorResponse(500, "Error interno al crear el resumen de egresos", err)
	}

	if err := tx.Where("created_at BETWEEN ? AND ?", firstday, lastDay).Delete(&models.Income{}).Error; err != nil {
		tx.Rollback()
		return models.ErrorResponse(500, "Error interno al eliminar los ingresos", err)
	}

	if err := tx.Where("created_at BETWEEN ? AND ?", firstday, lastDay).Delete(&models.Expense{}).Error; err != nil {
		tx.Rollback()
		return models.ErrorResponse(500, "Error interno al eliminar los egresos", err)
	}

	if err := tx.Commit().Error; err != nil {
		return models.ErrorResponse(500, "Error al confirmar la transacción", err)
	}
	return nil
}