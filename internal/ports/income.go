package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"


type IncomeRepository interface {
	IncomeGetByID(id string) (income *schemas.IncomeResponse, err error)
	IncomeGetAll(page, limit int) (incomes *[]schemas.IncomeDTO, err error)
	IncomeGetToday(page, limit int) (incomes *[]schemas.IncomeDTO, err error)
	IncomeCreate(incomeCreate *schemas.IncomeCreate) (id string, err error)
	IncomeUpdate(incomeUpdate *schemas.IncomeUpdate) (err error)
	IncomeDelete(id string) error
}

type IncomeService interface {
	IncomeGetByID(id string) (income *schemas.IncomeResponse, err error)
	IncomeGetAll(page, limit int) (incomes *[]schemas.IncomeDTO, err error)
	IncomeGetToday(page, limit int) (incomes *[]schemas.IncomeDTO, err error)
	IncomeCreate(incomeCreate *schemas.IncomeCreate) (id string, err error)
	IncomeUpdate(incomeUpdate *schemas.IncomeUpdate) (err error)
	IncomeDelete(id string) error
}
