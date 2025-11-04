package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"


type ExpenseRepository interface {
	ExpenseGetByID(id string) (expense *schemas.ExpenseResponse, err error)
	ExpenseGetAll(page, limit int) (expenses *[]schemas.ExpenseDTO, err error)
	ExpenseGetToday(page, limit int) (expenses *[]schemas.ExpenseDTO, err error)
	ExpenseCreate(expenseCreate *schemas.ExpenseCreate) (id string, err error)
	ExpenseUpdate(expenseUpdate *schemas.ExpenseUpdate) (err error)
	ExpenseDelete(id string) error
}

type ExpenseService interface {
	ExpenseGetByID(id string) (expense *schemas.ExpenseResponse, err error)
	ExpenseGetAll(page, limit int) (expenses *[]schemas.ExpenseDTO, err error)
	ExpenseGetToday(page, limit int) (expenses *[]schemas.ExpenseDTO, err error)
	ExpenseCreate(expenseCreate *schemas.ExpenseCreate) (id string, err error)
	ExpenseUpdate(expenseUpdate *schemas.ExpenseUpdate) (err error)
	ExpenseDelete(id string) error
}