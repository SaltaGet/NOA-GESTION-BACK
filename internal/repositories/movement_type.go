package repositories

import (
	"errors"

	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *MovementTypeRepository) MovementTypeGetByID(id string) (*models.MovementType, error) {
		var movementType models.MovementType
		if err := r.DB.Where("id = ?", id).First(&movementType).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, models.ErrorResponse(404, "Movimiento no encontrado", err)
			}
			return nil, models.ErrorResponse(500, "Error interno al buscar movimiento", err)
		}
		return &movementType, nil
}

func (r *MovementTypeRepository) MovementTypeGetAll(isIncome bool) (*[]models.MovementType, error) {
		var movementTypes []models.MovementType
		if err := r.DB.Where("is_income = ?", isIncome).Find(&movementTypes).Error; err != nil {
			return nil, models.ErrorResponse(500, "Error interno al buscar movimientos", err)
		}
		return &movementTypes, nil
}

func (r *MovementTypeRepository) MovementTypeCreate(movementType *models.MovementTypeCreate) (string, error) {
	newID := uuid.NewString()
			if err := r.DB.Create(&models.MovementType{
				ID: newID,
				Name: movementType.Name,
				IsIncome: movementType.IsIncome,
			}).Error; err != nil {
				return "", models.ErrorResponse(500, "Error interno al crear movimiento", err)
			}
			return newID, nil
}

func (r *MovementTypeRepository) MovementTypeUpdate(movementTypeUpdate *models.MovementTypeUpdate) error {
			if err := r.DB.Model(&models.MovementType{}).Where("id = ?", movementTypeUpdate.ID).Updates(&models.MovementType{
				Name: movementTypeUpdate.Name,
				IsIncome: movementTypeUpdate.IsIncome,
			}).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return models.ErrorResponse(404, "Movimiento no encontrado", err)
				}
				return models.ErrorResponse(500, "Error interno al actualizar movimiento", err)
			}
			return nil
}

func (r *MovementTypeRepository) MovementTypeDelete(id string) error {
		var movementType models.MovementType
		if err := r.DB.Where("id = ?", id).Delete(&movementType).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return models.ErrorResponse(404, "Movimiento no encontrado", err)
			}
			return models.ErrorResponse(500, "Error interno al eliminar movimiento", err)
		}
		return nil
}