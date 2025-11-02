package ports

import "github.com/DanielChachagua/GestionCar/pkg/models"

type EmployeeService interface {
	EmployeeGetByID(id string) (employee *models.Employee, err error)
	EmployeeGetByName(name string) (employees *[]models.Employee, err error)
	EmployeeGetAll() (employees *[]models.Employee, err error)
	EmployeeCreate(employeeCreate *models.EmployeeCreate) (id string, err error)
	EmployeeUpdate(employeeUpdate *models.EmployeeUpdate) (err error)
	EmployeeDelete(id string) error
}

type EmployeeRepository interface {
	EmployeeGetByID(id string) (employee *models.Employee, err error)
	EmployeeGetByName(name string) (employees *[]models.Employee, err error)
	EmployeeGetAll() (employees *[]models.Employee, err error)
	EmployeeCreate(employeeCreate *models.EmployeeCreate) (id string, err error)
	EmployeeUpdate(employeeUpdate *models.EmployeeUpdate) (err error)
	EmployeeDelete(id string) error
}
