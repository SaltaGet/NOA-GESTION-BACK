package ports

import "github.com/DanielChachagua/GestionCar/pkg/models"

type ExpenseRepository interface {
	ExpenseGetByID(id string) (expense *models.ExpenseResponse, err error)
	ExpenseGetAll(page, limit int) (expenses *[]models.ExpenseDTO, err error)
	ExpenseGetToday(page, limit int) (expenses *[]models.ExpenseDTO, err error)
	ExpenseCreate(expenseCreate *models.ExpenseCreate) (id string, err error)
	ExpenseUpdate(expenseUpdate *models.ExpenseUpdate) (err error)
	ExpenseDelete(id string) error
}

type ExpenseService interface {
	ExpenseGetByID(id string) (expense *models.ExpenseResponse, err error)
	ExpenseGetAll(page, limit int) (expenses *[]models.ExpenseDTO, err error)
	ExpenseGetToday(page, limit int) (expenses *[]models.ExpenseDTO, err error)
	ExpenseCreate(expenseCreate *models.ExpenseCreate) (id string, err error)
	ExpenseUpdate(expenseUpdate *models.ExpenseUpdate) (err error)
	ExpenseDelete(id string) error
}