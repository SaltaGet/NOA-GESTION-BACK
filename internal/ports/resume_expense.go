package ports

import "github.com/DanielChachagua/GestionCar/pkg/models"

type ResumeExpenseService interface {
	ResumeExpenseGetByID(id string) (resume *models.ResumeExpense, err error)
	ResumeExpenseGetByDateBetween(fromDate string, toDate string) (resumes *[]models.ResumeExpense, err error)
	ResumeExpenseCreate(resume *models.ResumeExpenseCreate) (id string, err error)
	ResumeExpenseUpdate(resume *models.ResumeExpenseUpdate) (err error)
}

type ResumeExpenseRepository interface {
	ResumeExpenseGetByID(id string) (resume *models.ResumeExpense, err error)
	ResumeExpenseGetByDateBetween(fromDate string, toDate string) (resumes *[]models.ResumeExpense, err error)
	ResumeExpenseCreate(resume *models.ResumeExpenseCreate) (id string, err error)
	ResumeExpenseUpdate(resume *models.ResumeExpenseUpdate) (err error)
}


