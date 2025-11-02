package services

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
)

func (r *ResumeService) ResumeExpenseCreate(resume *models.ResumeExpenseCreate) (string, error) {
	id, err := r.ResumeExpenseRepository.ResumeExpenseCreate(resume)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *ResumeService) ResumeExpenseGetByDateBetween(fromDate string, toDate string) (*[]models.ResumeExpense, error) {
	resumes, err := r.ResumeExpenseRepository.ResumeExpenseGetByDateBetween(fromDate, toDate)
	if err != nil {
		return nil, err
	}

	return resumes, nil
}

func (r *ResumeService) ResumeExpenseGetByID(id string) (*models.ResumeExpense, error) {
	resumes, err := r.ResumeExpenseRepository.ResumeExpenseGetByID(id)
	if err != nil {
		return nil, err
	}

	return resumes, nil
}

func (r *ResumeService) ResumeExpenseUpdate(resume *models.ResumeExpenseUpdate) error {
	err := r.ResumeExpenseRepository.ResumeExpenseUpdate(resume)
	if err != nil {
		return err
	}

	return nil
}


// INCOME

func (r *ResumeService) ResumeIncomeCreate(resume *models.ResumeIncomeCreate) (string, error) {
	id, err := r.ResumeIncomeRepository.ResumeIncomeCreate(resume)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *ResumeService) ResumeIncomeGetByDateBetween(fromDate string, toDate string) (*[]models.ResumeIncome, error) {
	incomes, err := r.ResumeIncomeRepository.ResumeIncomeGetByDateBetween(fromDate, toDate)
	if err != nil {
		return nil, err
	}

	return incomes, nil
}

func (r *ResumeService) ResumeIncomeGetByID(id string) (*models.ResumeIncome, error) {
	income, err := r.ResumeIncomeRepository.ResumeIncomeGetByID(id)
	if err != nil {
		return nil, err
	}

	return income, nil
}

func (r *ResumeService) ResumeIncomeUpdate(resume *models.ResumeIncomeUpdate) error {
	err := r.ResumeIncomeRepository.ResumeIncomeUpdate(resume)
	if err != nil {
		return err
	}

	return nil
}
