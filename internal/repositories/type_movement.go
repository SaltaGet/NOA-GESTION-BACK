package repositories

import (
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"gorm.io/gorm"
)

func (t *TypeMovementRepository) TypeMovementCreate(movementCreate schemas.TypeMovementCreate) error {
	var err error
	switch movementCreate.TypeMovement {
	case "income":
		err = t.DB.Create(&models.TypeIncome{Name: movementCreate.Name}).Error
	case "expense":
		err = t.DB.Create(&models.TypeExpense{Name: movementCreate.Name}).Error
	default:
		return schemas.ErrorResponse(400, "tipo de movimiento no válido", fmt.Errorf("tipo de movimiento no valido: %s", movementCreate.TypeMovement))
	}

	if err != nil {
		if schemas.IsDuplicateError(err) {
			return schemas.ErrorResponse(409, "tipo de movimiento ya existe", err)
		}
		return schemas.ErrorResponse(500, "error al crear tipo de movimiento", err)
	}

	return nil
}

func (t *TypeMovementRepository) TypeMovementGetAll(typeMovement string) ([]*schemas.TypeMovementResponse, error) {
	var typeMovements []*schemas.TypeMovementResponse
	switch typeMovement {
	case "income":
		err := t.DB.Model(&models.TypeIncome{}).Select("id", "name").Scan(&typeMovements).Error
		if err != nil {
			return nil, schemas.ErrorResponse(500, "error al obtner movimientos", err)
		}
	case "expense":
		err := t.DB.Model(&models.TypeExpense{}).Select("id", "name").Scan(&typeMovements).Error
		if err != nil {
			return nil, schemas.ErrorResponse(500, "error al obtner movimientos", err)
		}
	default:
		return nil, schemas.ErrorResponse(400, "tipo de movimiento no válido", fmt.Errorf("tipo de movimiento no valido: %s", typeMovement))
	}
	return typeMovements, nil
}

func (t *TypeMovementRepository) TypeMovementUpdate(movementUpdate schemas.TypeMovementUpdate) error {
	var res *gorm.DB

	switch movementUpdate.TypeMovement {
	case "income":
		res = t.DB.Model(&models.TypeIncome{}).
			Where("id = ?", movementUpdate.ID).
			Update("name", movementUpdate.Name)

	case "expense":
		res = t.DB.Model(&models.TypeExpense{}).
			Where("id = ?", movementUpdate.ID).
			Update("name", movementUpdate.Name)

	default:
		return schemas.ErrorResponse(400,	"tipo de movimiento no válido",	fmt.Errorf("tipo de movimiento no válido: %s", movementUpdate.TypeMovement),
		)
	}

	if res.RowsAffected == 0 {
		return schemas.ErrorResponse(404,	"tipo de movimiento no encontrado",	fmt.Errorf("id %d no encontrado", movementUpdate.ID),
		)
	}

	if res.Error != nil {
		if schemas.IsDuplicateError(res.Error) {
			return schemas.ErrorResponse(409, "tipo de movimiento ya existe", res.Error)
		}

		return schemas.ErrorResponse(500, "error al actualizar tipo de movimiento", res.Error)
	}

	return nil
}

