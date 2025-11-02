package ports

import "github.com/DanielChachagua/GestionCar/pkg/models"

type IncomeRepository interface {
	IncomeGetByID(id string) (income *models.IncomeResponse, err error)
	IncomeGetAll(page, limit int) (incomes *[]models.IncomeDTO, err error)
	IncomeGetToday(page, limit int) (incomes *[]models.IncomeDTO, err error)
	IncomeCreate(incomeCreate *models.IncomeCreate) (id string, err error)
	IncomeUpdate(incomeUpdate *models.IncomeUpdate) (err error)
	IncomeDelete(id string) error
}

type IncomeService interface {
	IncomeGetByID(id string) (income *models.IncomeResponse, err error)
	IncomeGetAll(page, limit int) (incomes *[]models.IncomeDTO, err error)
	IncomeGetToday(page, limit int) (incomes *[]models.IncomeDTO, err error)
	IncomeCreate(incomeCreate *models.IncomeCreate) (id string, err error)
	IncomeUpdate(incomeUpdate *models.IncomeUpdate) (err error)
	IncomeDelete(id string) error
}
