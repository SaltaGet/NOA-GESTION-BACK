package services

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
)

func (e *EmployeeService) EmployeeGetByID(id string) (*models.Employee, error) {
	employee, err := e.EmployeeRepository.EmployeeGetByID(id)
	if err != nil {
		return nil, err
	}
	return employee, nil
}

func (e *EmployeeService) EmployeeGetByName(name string) (*[]models.Employee, error) {
	employees, err := e.EmployeeRepository.EmployeeGetByName(name)
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (e *EmployeeService) EmployeeGetAll() (*[]models.Employee, error) {
	employees, err := e.EmployeeRepository.EmployeeGetAll()
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (e *EmployeeService) EmployeeCreate(employee *models.EmployeeCreate) (string, error) {
	id, err := e.EmployeeRepository.EmployeeCreate(employee)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (e *EmployeeService) EmployeeUpdate(employee *models.EmployeeUpdate) error {
	err := e.EmployeeRepository.EmployeeUpdate(employee)
	if err != nil {
		return err
	}
	return nil
}

func (e *EmployeeService) EmployeeDelete(id string) error {
	err := e.EmployeeRepository.EmployeeDelete(id)
	if err != nil {
		return err
	}
	return nil
}