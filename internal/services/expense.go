package services

import (
"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (e *ExpenseService) ExpenseGetByID(id string) (*schemas.ExpenseResponse, error) {
	expense, err := e.ExpenseRepository.ExpenseGetByID(id)
	if err != nil {
		return nil, err
	}

	return expense, nil
}

func (e *ExpenseService) ExpenseGetAll(page, limit int) (*[]schemas.ExpenseDTO, error) {
	expenses, err := e.ExpenseRepository.ExpenseGetAll(page, limit)
	
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func (e *ExpenseService) ExpenseGetToday(page, limit int) (*[]schemas.ExpenseDTO, error) {
	expenses, err := e.ExpenseRepository.ExpenseGetToday(page, limit)
	
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func (e *ExpenseService) ExpenseCreate(expense *schemas.ExpenseCreate) (string, error) {
	id, err := e.ExpenseRepository.ExpenseCreate(expense)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (e *ExpenseService) ExpenseUpdate(expense *schemas.ExpenseUpdate) error {
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