package repositories

import (
	"errors"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *MovementTypeRepository) MovementTypeGetByID(id string) (*schemas.MovementType, error) {
		var movementType schemas.MovementType
		if err := r.DB.Where("id = ?", id).First(&movementType).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, schemas.ErrorResponse(404, "Movimiento no encontrado", err)
			}
			return nil, schemas.ErrorResponse(500, "Error interno al buscar movimiento", err)
		}
		return &movementType, nil
}

func (r *MovementTypeRepository) MovementTypeGetAll(isIncome bool) (*[]schemas.MovementType, error) {
		var movementTypes []schemas.MovementType
		if err := r.DB.Where("is_income = ?", isIncome).Find(&movementTypes).Error; err != nil {
			return nil, schemas.ErrorResponse(500, "Error interno al buscar movimientos", err)
		}
		return &movementTypes, nil
}

func (r *MovementTypeRepository) MovementTypeCreate(movementType *schemas.MovementTypeCreate) (string, error) {
	newID := uuid.NewString()
			if err := r.DB.Create(&schemas.MovementType{
				ID: newID,
				Name: movementType.Name,
				IsIncome: movementType.IsIncome,
			}).Error; err != nil {
				return "", schemas.ErrorResponse(500, "Error interno al crear movimiento", err)
			}
			return newID, nil
}

func (r *MovementTypeRepository) MovementTypeUpdate(movementTypeUpdate *schemas.MovementTypeUpdate) error {
			if err := r.DB.Model(&schemas.MovementType{}).Where("id = ?", movementTypeUpdate.ID).Updates(&schemas.MovementType{
				Name: movementTypeUpdate.Name,
				IsIncome: movementTypeUpdate.IsIncome,
			}).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return schemas.ErrorResponse(404, "Movimiento no encontrado", err)
				}
				return schemas.ErrorResponse(500, "Error interno al actualizar movimiento", err)
			}
			return nil
}

func (r *MovementTypeRepository) MovementTypeDelete(id string) error {
		var movementType schemas.MovementType
		if err := r.DB.Where("id = ?", id).Delete(&movementType).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Movimiento no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error interno al eliminar movimiento", err)
		}
		return nil
}