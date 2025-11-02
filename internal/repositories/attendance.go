package repositories

import (
	"errors"
	"fmt"

	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *AttendanceRepository) AttendanceGetByID(id string) (*models.AttendanceDTO, error) {
	var attendanceDTO models.AttendanceDTO
	if err := r.DB.Model(&models.Attendance{}).
		Select("id, employee_id, attendance, hours, date, amount, is_holiday, is_paid, created_at, updated_at").
		Where("id = ?", id).
		First(&attendanceDTO).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrorResponse(404, "Asistencia no encontrada", err)
		}
		return nil, models.ErrorResponse(500, "Error interno al buscar la asistencia", err)
	}
	return &attendanceDTO, nil
}

func (r *AttendanceRepository) AttendanceGetAll() (*[]models.AttendanceDTO, error) {
	var attendancesDTO []models.AttendanceDTO
	if err := r.DB.Model(&models.Attendance{}).
		Select("id, employee_id, attendance, hours, date, amount, is_holiday, is_paid, created_at, updated_at").
		Find(&attendancesDTO).
		Error; err != nil {
		return nil, models.ErrorResponse(500, "Error al buscar las asistencias", err)
	}
	return &attendancesDTO, nil

}

func (r *AttendanceRepository) AttendanceGetByEmployeeID(userID string) (*[]models.AttendanceDTO, error) {
	var attendances []models.AttendanceDTO
	if err := r.DB.Model(&models.Attendance{}).
		Select("id, employee_id, attendance, hours, date, amount, is_holiday, is_paid, created_at, updated_at").
		Where("employee_id = ?", userID).
		Find(&attendances).Error; err != nil {
		return nil, models.ErrorResponse(500, "Error interno al buscar asistencias", err)
	}
	return &attendances, nil
}

func (r *AttendanceRepository) AttendanceCreate(attendance *models.AttendanceCreate) (string, error) {
	newId := uuid.NewString()
	if err := r.DB.Create(&models.Attendance{
		ID:         newId,
		EmployeeID: attendance.EmployeeID,
		Attendance: attendance.Attendance,
		Hours:      attendance.Hours,
		Date:       attendance.Date,
		Amount:     attendance.Amount,
		IsHoliday:  attendance.IsHoliday,
	}).Error; err != nil {
		return "", models.ErrorResponse(500, "Error interno al crear la asistencia", err)
	}

	return newId, nil
}

func (r *AttendanceRepository) AttendanceUpdate(attendance *models.AttendanceUpdate) error {
	if err := r.DB.Where("id = ?", attendance.ID).Updates(&models.Attendance{
		Attendance: attendance.Attendance,
		Hours:      attendance.Hours,
		Date:       attendance.Date,
		Amount:     attendance.Amount,
		IsHoliday:  attendance.IsHoliday,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrorResponse(404, "Asistencia no encontrada", err)
		}
		return models.ErrorResponse(500, "Error interno al actualizar la asistencia", err)
	}
	return nil
}

func (r *AttendanceRepository) AttendanceDelete(id string) error {
	var attendance models.Attendance
	if err := r.DB.Where("id = ?", id).Delete(&attendance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrorResponse(404, "Asistencia no encontrada", err)
		}
		return models.ErrorResponse(500, "Error interno al eliminar la asistencia", err)
	}
	return nil
}

func (r *AttendanceRepository) AttendanceGetByDate(date_start string, date_end string) (*[]models.AttendanceDTO, error) {
	var attendances []models.AttendanceDTO
	if err := r.DB.Model(&models.Attendance{}).
		Select("id, employee_id, attendance, hours, date, amount, is_holiday, is_paid, created_at, updated_at").
		Where("DATE(date) >= ? AND DATE(date) <= ?", date_start, date_end).
		Find(&attendances).Error; err != nil {
		return nil, models.ErrorResponse(500, "Error interno al buscar las asistencias", err)
	}
	return &attendances, nil
}

func (r *AttendanceRepository) AttendanceUpdatePay(listIDs []string) error {
	if err := r.DB.
		Model(&models.Attendance{}).
		Where("id IN (?)", listIDs).
		Update("is_paid", gorm.Expr("NOT is_paid")).
		Error; err != nil {
		return models.ErrorResponse(500, "Error interno al marcar las asistencias como pagadas", err)
	}

	return nil
}

func (r *AttendanceRepository) AttendancePay(listIDs *[]models.AttendancePay) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var values [][]any
	
		for _, attendance := range *listIDs {
			values = append(values, []any{attendance.AttendanceID, attendance.EmployeeID})
		}

		var count int64
		if err := tx.
			Model(&models.Attendance{}).
			Where("(id, employee_id) IN ?", values).
			Where("is_paid = ?", true).
			Count(&count).
			Error; err != nil {
			return models.ErrorResponse(500, "Error interno al verificar las asistencias", err)
		}
	
		if count > 0 {
			return models.ErrorResponse(400, "existen asistencias que ya fueron pagadas", fmt.Errorf("existen asistencias que ya fueron pagadas"))
		}
	
		if err := tx.
			Model(&models.Attendance{}).
			Where("(id, employee_id) IN ?", values).
			Update("is_paid = ?", true).
			Error; err != nil {
			return models.ErrorResponse(500, "Error interno al marcar las asistencias como pagadas", err)
		}
	
		return nil
	})
}
