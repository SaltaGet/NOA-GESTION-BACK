package repositories

import (
	"errors"

	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *EmployeeRepository) EmployeeGetByID(id string) (*models.Employee, error) {
	var employee models.Employee
	if err := r.DB.Where("id = ?", id).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrorResponse(404, "Empleado no encontrado", err)
		}
		return nil, models.ErrorResponse(500, "Error interno al buscar el empleado", err)
	}
	return &employee, nil
}

func (r *EmployeeRepository) EmployeeGetAll() (*[]models.Employee, error) {
	var employees []models.Employee
	if err := r.DB.Find(&employees).Error; err != nil {
		return nil, models.ErrorResponse(500, "Error interno al buscar los empleados", err)
	}
	return &employees, nil
}

func (r *EmployeeRepository) EmployeeGetByName(name string) (*[]models.Employee, error) {
	var employees []models.Employee
	if err := r.DB.Where("name LIKE ?", "%"+name+"%").Find(&employees).Error; err != nil {
		return nil, models.ErrorResponse(500, "Error interno al buscar los empleados", err)
	}
	return &employees, nil
}

func (r *EmployeeRepository)EmployeeCreate(employee *models.EmployeeCreate) (string, error) {
	newID := uuid.NewString()
	if err := r.DB.Create(&models.Employee{
		ID:      newID,
		Name:    employee.Name,
		Phone:   employee.Phone,
		Email:   employee.Email,
		Address: employee.Address,
	}).Error; err != nil {
		return "", models.ErrorResponse(500, "Error interno al crear el empleado", err)
	}
	return newID, nil
}

func (r *EmployeeRepository) EmployeeUpdate(employeeUpdate *models.EmployeeUpdate) error {
	var employee models.Employee
	if err := r.DB.Where("id = ?", employeeUpdate.ID).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrorResponse(404, "Empleado no encontrado", err)
		}
		return models.ErrorResponse(500, "Error interno al buscar el empleado", err)
	}
	employee.Name = employeeUpdate.Name
	employee.Phone = employeeUpdate.Phone
	employee.Email = employeeUpdate.Email
	employee.Address = employeeUpdate.Address
	if err := r.DB.Save(&employee).Error; err != nil {
		return models.ErrorResponse(500, "Error interno al actualizar el empleado", err)
	}
	return nil
}

func (r *EmployeeRepository) EmployeeDelete(id string) error {
	var employee models.Employee
	if err := r.DB.Where("id = ?", id).Delete(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrorResponse(404, "Empleado no encontrado", err)
		}
		return models.ErrorResponse(500, "Error interno al eliminar el empleado", err)
	}
	return nil
}

