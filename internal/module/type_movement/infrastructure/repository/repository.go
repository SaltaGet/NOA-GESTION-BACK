package repository


import (
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"gorm.io/gorm"
)

func (t *TypeMovementRepository) TypeMovementGetAll(typeMovement string) ([]*schemas.TypeMovementResponse, error) {
	var typeMovements []*schemas.TypeMovementResponse
	switch typeMovement {
	case "income":
		err := t.DB.Model(&models.TypeIncome{}).Select("id", "name").Scan(&typeMovements).Error
		if err != nil {
			return nil, schemas.HandlerErrorGorm(err, "Tipo de ingreso", schemas.Read)
		}
	case "expense":
		err := t.DB.Model(&models.TypeExpense{}).Select("id", "name").Scan(&typeMovements).Error
		if err != nil {
			return nil, schemas.HandlerErrorGorm(err, "Tipo de egreso", schemas.Read)
		}
	default:
		return nil, schemas.ErrorResponse(400, "tipo de movimiento no válido", fmt.Errorf("tipo de movimiento no valido: %s", typeMovement))
	}
	return typeMovements, nil
}

func (t *TypeMovementRepository) TypeMovementCreate(memberID int64, movementCreate schemas.TypeMovementCreate) error {
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var err error
		switch movementCreate.TypeMovement {
		case "income":
			typeIncome := models.TypeIncome{Name: movementCreate.Name}
			err = tx.Create(&typeIncome).Error
		case "expense":
			typeExpense := models.TypeExpense{Name: movementCreate.Name}
			err = tx.Create(&typeExpense).Error
		default:
			return schemas.ErrorResponse(400, "tipo de movimiento no válido", fmt.Errorf("tipo de movimiento no valido: %s", movementCreate.TypeMovement))
		}

		if err != nil {
			return schemas.HandlerErrorGorm(err, "Tipo de movimiento", schemas.Create)
		}

		return nil
	})

	return err
}

// TypeMovementUpdate actualiza un tipo de movimiento con auditoría
func (t *TypeMovementRepository) TypeMovementUpdate(memberID int64, movementUpdate schemas.TypeMovementUpdate) error {
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var res *gorm.DB
		switch movementUpdate.TypeMovement {
		case "income":
			// Obtener estado anterior
			var oldIncome models.TypeIncome
			if err := tx.First(&oldIncome, movementUpdate.ID).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Tipo de ingreso", schemas.Read)
			}
			// Actualizar
			res = tx.Model(&models.TypeIncome{}).
				Where("id = ?", movementUpdate.ID).
				Update("name", movementUpdate.Name)			
			
			var newIncome models.TypeIncome
			tx.First(&newIncome, movementUpdate.ID)
		case "expense":
			// Obtener estado anterior
			var oldExpense models.TypeExpense
			if err := tx.First(&oldExpense, movementUpdate.ID).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Tipo de egreso", schemas.Read)
			}

			// Actualizar
			res = tx.Model(&models.TypeExpense{}).
				Where("id = ?", movementUpdate.ID).
				Update("name", movementUpdate.Name)

			// Obtener estado nuevo
			var newExpense models.TypeExpense
			tx.First(&newExpense, movementUpdate.ID)
		default:
			return schemas.ErrorResponse(400, "tipo de movimiento no válido", fmt.Errorf("tipo de movimiento no válido: %s", movementUpdate.TypeMovement))
		}

		if res.RowsAffected == 0 {
			return schemas.ErrorResponse(404, "tipo de movimiento no encontrado", fmt.Errorf("id %d no encontrado", movementUpdate.ID))
		}

		if res.Error != nil {
			return schemas.HandlerErrorGorm(res.Error, "Tipo de movimiento", schemas.Update)
		}

		return nil
	})

	return err
}
