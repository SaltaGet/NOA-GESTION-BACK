package repositories

import (
	"errors"
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	"gorm.io/gorm"
)

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

// func (t *TypeMovementRepository) TypeMovementCreate(movementCreate schemas.TypeMovementCreate) error {
// 	var err error
// 	switch movementCreate.TypeMovement {
// 	case "income":
// 		err = t.DB.Create(&models.TypeIncome{Name: movementCreate.Name}).Error
// 	case "expense":
// 		err = t.DB.Create(&models.TypeExpense{Name: movementCreate.Name}).Error
// 	default:
// 		return schemas.ErrorResponse(400, "tipo de movimiento no válido", fmt.Errorf("tipo de movimiento no valido: %s", movementCreate.TypeMovement))
// 	}

// 	if err != nil {
// 		if schemas.IsDuplicateError(err) {
// 			return schemas.ErrorResponse(409, "tipo de movimiento ya existe", err)
// 		}
// 		return schemas.ErrorResponse(500, "error al crear tipo de movimiento", err)
// 	}

// 	return nil
// }

// func (t *TypeMovementRepository) TypeMovementUpdate(movementUpdate schemas.TypeMovementUpdate) error {
// 	var res *gorm.DB

// 	switch movementUpdate.TypeMovement {
// 	case "income":
// 		res = t.DB.Model(&models.TypeIncome{}).
// 			Where("id = ?", movementUpdate.ID).
// 			Update("name", movementUpdate.Name)

// 	case "expense":
// 		res = t.DB.Model(&models.TypeExpense{}).
// 			Where("id = ?", movementUpdate.ID).
// 			Update("name", movementUpdate.Name)

// 	default:
// 		return schemas.ErrorResponse(400,	"tipo de movimiento no válido",	fmt.Errorf("tipo de movimiento no válido: %s", movementUpdate.TypeMovement),
// 		)
// 	}

// 	if res.RowsAffected == 0 {
// 		return schemas.ErrorResponse(404,	"tipo de movimiento no encontrado",	fmt.Errorf("id %d no encontrado", movementUpdate.ID),
// 		)
// 	}

// 	if res.Error != nil {
// 		if schemas.IsDuplicateError(res.Error) {
// 			return schemas.ErrorResponse(409, "tipo de movimiento ya existe", res.Error)
// 		}

// 		return schemas.ErrorResponse(500, "error al actualizar tipo de movimiento", res.Error)
// 	}

// 	return nil
// }

// TypeMovementCreate crea un nuevo tipo de movimiento con auditoría
func (t *TypeMovementRepository) TypeMovementCreate(memberID int64, movementCreate schemas.TypeMovementCreate) error {
	var movementSave any
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		switch movementCreate.TypeMovement {
		case "income":
			typeIncome := models.TypeIncome{Name: movementCreate.Name}
			err = tx.Create(&typeIncome).Error
			movementSave = typeIncome

		case "expense":
			typeExpense := models.TypeExpense{Name: movementCreate.Name}
			err = tx.Create(&typeExpense).Error
			movementSave = typeExpense

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
	})

	if err == nil {
		// Guardar auditoría
		go database.SaveAuditAsync(t.DB, models.AuditLog{
			MemberID: memberID,
			Method:   "create",
			Path:     utils.Ternary(movementCreate.TypeMovement == "income", "type-income", "type-expense"),
		}, nil, movementSave)
	}

	return err
}

// TypeMovementUpdate actualiza un tipo de movimiento con auditoría
func (t *TypeMovementRepository) TypeMovementUpdate(memberID int64, movementUpdate schemas.TypeMovementUpdate) error {
	var oldTypeMovement, newTypeMovement any
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		var res *gorm.DB
		switch movementUpdate.TypeMovement {
		case "income":
			// Obtener estado anterior
			var oldIncome models.TypeIncome
			if err := tx.First(&oldIncome, movementUpdate.ID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return schemas.ErrorResponse(404, "tipo de movimiento no encontrado", err)
				}
				return schemas.ErrorResponse(500, "error al obtener tipo de movimiento", err)
			}
			oldTypeMovement = oldIncome

			// Actualizar
			res = tx.Model(&models.TypeIncome{}).
				Where("id = ?", movementUpdate.ID).
				Update("name", movementUpdate.Name)			
			
			var newIncome models.TypeIncome
			tx.First(&newIncome, movementUpdate.ID)
			newTypeMovement = newIncome
		case "expense":
			// Obtener estado anterior
			var oldExpense models.TypeExpense
			if err := tx.First(&oldExpense, movementUpdate.ID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return schemas.ErrorResponse(404, "tipo de movimiento no encontrado", err)
				}
				return schemas.ErrorResponse(500, "error al obtener tipo de movimiento", err)
			}
			oldTypeMovement = oldExpense

			// Actualizar
			res = tx.Model(&models.TypeExpense{}).
				Where("id = ?", movementUpdate.ID).
				Update("name", movementUpdate.Name)

			// Obtener estado nuevo
			var newExpense models.TypeExpense
			tx.First(&newExpense, movementUpdate.ID)
			newTypeMovement = newExpense

		default:
			return schemas.ErrorResponse(400, "tipo de movimiento no válido", fmt.Errorf("tipo de movimiento no válido: %s", movementUpdate.TypeMovement))
		}

		if res.RowsAffected == 0 {
			return schemas.ErrorResponse(404, "tipo de movimiento no encontrado", fmt.Errorf("id %d no encontrado", movementUpdate.ID))
		}

		if res.Error != nil {
			if schemas.IsDuplicateError(res.Error) {
				return schemas.ErrorResponse(409, "tipo de movimiento ya existe", res.Error)
			}
			return schemas.ErrorResponse(500, "error al actualizar tipo de movimiento", res.Error)
		}

		return nil
	})

	if err == nil {
		// Guardar auditoría
		go database.SaveAuditAsync(t.DB, models.AuditLog{
			MemberID: memberID,
			Method:   "update",
			Path:     "type-movement",
		}, oldTypeMovement, newTypeMovement)
	}

	return err
}
