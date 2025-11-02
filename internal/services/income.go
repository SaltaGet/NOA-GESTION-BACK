package services

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
)

func (i *IncomeService) IncomeGetByID(id string) (*models.IncomeResponse, error) {
	income, err := i.IncomeRepository.IncomeGetByID(id)
	if err != nil {
		return nil, err
	}

	return income, nil
}

func (i *IncomeService) IncomeGetAll(page, limit int) (*[]models.IncomeDTO, error) {
	incomes, err := i.IncomeRepository.IncomeGetAll(page, limit)
	
	if err != nil {
		return nil, err
	}

	return incomes, nil
}

func (i *IncomeService) IncomeGetToday(page, limit int) (*[]models.IncomeDTO, error) {
	incomes, err := i.IncomeRepository.IncomeGetToday(page, limit)
	
	if err != nil {
		return nil, err
	}

	return incomes, nil
}

func (i *IncomeService) IncomeCreate(incomeCreate *models.IncomeCreate) (string, error) {
	id, err := i.IncomeRepository.IncomeCreate(incomeCreate)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (i *IncomeService) IncomeUpdate(incomeUpdate *models.IncomeUpdate) error {
	err := i.IncomeRepository.IncomeUpdate(incomeUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (i *IncomeService) IncomeDelete(id string) error {
	err := i.IncomeRepository.IncomeDelete(id)
	if err != nil {
		return err
	}
	return nil
}