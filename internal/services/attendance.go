package services

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
)

func (s *AttendanceService) AttendanceGetByID(id string) (*models.AttendanceDTO, error) {
	attendance, err := s.AttendanceRepository.AttendanceGetByID(id)
	if err != nil {
		return nil, err
	}
	return attendance, nil
}

func (s *AttendanceService) AttendanceGetAll() (*[]models.AttendanceDTO, error) {
	attendances, err := s.AttendanceRepository.AttendanceGetAll()
	if err != nil {
		return nil, err
	}
	return attendances, nil
}

func (s *AttendanceService) AttendanceGetByDate(date_start string, date_end string) (*[]models.AttendanceDTO, error) {
	attendances, err := s.AttendanceRepository.AttendanceGetByDate(date_start, date_end)
	if err != nil {
		return nil, err
	}
	return attendances, nil
}

func (s *AttendanceService) AttendanceGetByEmployeeID(employeeID string) (*[]models.AttendanceDTO, error) {
	attendances, err := s.AttendanceRepository.AttendanceGetByEmployeeID(employeeID)
	if err != nil {
		return nil, err
	}
	return attendances, nil
}

func (s *AttendanceService) AttendanceCreate(attendance *models.AttendanceCreate) (string, error) {
	id, err := s.AttendanceRepository.AttendanceCreate(attendance)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *AttendanceService) AttendanceUpdate(attendance *models.AttendanceUpdate) error {
	err := s.AttendanceRepository.AttendanceUpdate(attendance)
	if err != nil {
		return err
	}

	return nil
}

func (s *AttendanceService) AttendanceUpdatePay(listIDs []string) error {
	err := s.AttendanceRepository.AttendanceUpdatePay(listIDs)
	if err != nil {
		return err
	}

	return nil
}

func (s *AttendanceService) AttendancePay(listIDs *[]models.AttendancePay) error {
	err := s.AttendanceRepository.AttendancePay(listIDs)
	if err != nil {
		return err
	}

	// aqui si me gustaria hacer que si no ha errores crear la expense
	err = s.ExpenseService.ExpenseCreate(listIDs)
	if err != nil {
		err := s.AttendanceRepository.AttendanceRevertPay(listIDs)
		return err
	}

	// aqui si no hubiera error en crear la expenserecien guardar los cambios, o confirmarlos, caso contrario hacer un rollback 

	return nil
}

func (s *AttendanceService) AttendanceDelete(id string) error {
	err := s.AttendanceRepository.AttendanceDelete(id)
	if err != nil {
		return err
	}
	return nil
}