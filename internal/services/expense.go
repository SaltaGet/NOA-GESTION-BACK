package services

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
)

func (e *ExpenseService) ExpenseGetByID(id string) (*models.ExpenseResponse, error) {
	expense, err := e.ExpenseRepository.ExpenseGetByID(id)
	if err != nil {
		return nil, err
	}

	return expense, nil
}

func (e *ExpenseService) ExpenseGetAll(page, limit int) (*[]models.ExpenseDTO, error) {
	expenses, err := e.ExpenseRepository.ExpenseGetAll(page, limit)
	
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func (e *ExpenseService) ExpenseGetToday(page, limit int) (*[]models.ExpenseDTO, error) {
	expenses, err := e.ExpenseRepository.ExpenseGetToday(page, limit)
	
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func (e *ExpenseService) ExpenseCreate(expense *models.ExpenseCreate) (string, error) {
	id, err := e.ExpenseRepository.ExpenseCreate(expense)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (e *ExpenseService) ExpenseUpdate(expense *models.ExpenseUpdate) error {
	err := e.ExpenseRepository.ExpenseUpdate(expense)
	if err != nil {
		return err
	}
	return nil
}

func (e *ExpenseService) ExpenseDelete(id string) error {
	err := e.ExpenseRepository.ExpenseDelete(id)
	if err != nil {
		return err
	}
	return nil
}