package repositories

import (
	"errors"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *ExpenseRepository) ExpenseGetByID(id string) (*schemas.ExpenseResponse, error) {
	var expense schemas.Expense

	err := r.DB.
		Preload("PurchaseOrder").
		Preload("PurchaseOrder.PurchaseProducts").
		Preload("PurchaseOrder.PurchaseProducts.Product").
		Preload("PurchaseOrder.Supplier").
		Preload("MovementType").
		First(&expense, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Ingreso no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al buscar ingreso", err)
	}

	response := schemas.ExpenseResponse{
		ID:        expense.ID,
		Details:   expense.Details,
		Amount:    expense.Amount,
		CreatedAt: expense.CreatedAt,
		MovementType: schemas.MovementTypeDTO{
			ID:       expense.MovementType.ID,
			Name:     expense.MovementType.Name,
			IsIncome: expense.MovementType.IsIncome,
		},
	}

	return &response, nil
}

func (r *ExpenseRepository) ExpenseGetAll(page, limit int) (*[]schemas.ExpenseDTO, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	var expenses []schemas.Expense

	err := r.DB.
		Preload("PurchaseOrder").
		Preload("PurchaseOrder.PurchaseProducts").
		Preload("PurchaseOrder.PurchaseProducts.Product").
		Preload("PurchaseOrder.Supplier").
		Preload("MovementType").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&expenses).Error

	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al buscar ingreso", err)
	}

	var expenseDTOs []schemas.ExpenseDTO

	return &expenseDTOs, nil
}

func (r *ExpenseRepository) ExpenseGetToday(page, limit int) (*[]schemas.ExpenseDTO, error) {
	start := time.Now().Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	var expenses []schemas.Expense

	err := r.DB.
		Preload("PurchaseOrder").
		Preload("PurchaseOrder.PurchaseProducts").
		Preload("PurchaseOrder.PurchaseProducts.Product").
		Preload("PurchaseOrder.Supplier").
		Preload("MovementType").
		Limit(limit).
		Where("created_at >= ? AND created_at < ?", start, end).
		Offset(offset).
		Order("created_at DESC").
		Find(&expenses).Error

	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al buscar ingreso", err)
	}

	var expenseDTOs []schemas.ExpenseDTO

	for _, expense := range expenses {
		response := schemas.ExpenseDTO{
			ID:        expense.ID,
			Amount:    expense.Amount,
			CreatedAt: expense.CreatedAt,
			MovementType: schemas.MovementTypeDTO{
				ID:       expense.MovementType.ID,
				Name:     expense.MovementType.Name,
				IsIncome: expense.MovementType.IsIncome,
			},
		}

		expenseDTOs = append(expenseDTOs, response)
	}

	return &expenseDTOs, nil
}

// func (r *ExpenseRepository) ExpenseGetByID(id string) (*schemas.ExpenseResponse, error) {
// 	var expense schemas.Expense
// 	if err := r.DB.Where("id = ?", id).First(&expense).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, schemas.ErrorResponse(404, "Movimiento no encontrado", err)
// 		}
// 		return nil, schemas.ErrorResponse(500, "Error interno al buscar movimiento", err)
// 	}
// 	return &expense, nil
// }
// func (r *ExpenseRepository) ExpenseGetAll() (*[]schemas.Expense, error) {
// 	var expenses []schemas.Expense
// 	if err := r.DB.Limit(100).Order("created_at desc").Find(&expenses).Error; err != nil {
// 		return nil, schemas.ErrorResponse(500, "Error interno al buscar movimientos", err)
// 	}
// 	return &expenses, nil
// }

// func (r *ExpenseRepository) ExpenseGetToday() (*[]schemas.Expense, error) {
// 	today := time.Now().Format("2006-01-02")
// 	var expenses []schemas.Expense
// 	if err := r.DB.Where("DATE(created_at) = ?", today).Order("created_at desc").Find(&expenses).Error; err != nil {
// 		return nil, schemas.ErrorResponse(500, "Error interno al buscar movimientos", err)
// 	}
// 	return &expenses, nil
// }

func (r *ExpenseRepository) ExpenseCreate(expense *schemas.ExpenseCreate) (string, error) {
	newID := uuid.NewString()
	if err := r.DB.Create(&schemas.Expense{
		ID:             newID,
		Details:        expense.Details,
		PurchaseOrderID:     expense.PurchaseOrderID,
		MovementTypeID: expense.MovementTypeID,
		Amount:         expense.Amount,
	}).Error; err != nil {
		return "", schemas.ErrorResponse(500, "Error interno al crear movimiento", err)
	}
	return newID, nil
}

func (r *ExpenseRepository) ExpenseUpdate(expense *schemas.ExpenseUpdate) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", expense.ID).
			Updates(&schemas.Expense{
				Details: expense.Details,
				PurchaseOrderID:     expense.PurchaseOrderID,
				MovementTypeID: expense.MovementTypeID,
				Amount:         expense.Amount,
			}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Movimiento no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error interno al actualizar movimiento", err)
		}
		return nil
	})
}

func (r *ExpenseRepository) ExpenseDelete(id string) error {
	if err := r.DB.Where("id = ?", id).Delete(&schemas.Expense{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "Movimiento no encontrado", err)
		}
		return schemas.ErrorResponse(500, "Error interno al eliminar movimiento", err)
	}
	return nil
}
