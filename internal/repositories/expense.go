package repositories

import (
	"errors"
	"time"

	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *ExpenseRepository) ExpenseGetByID(id string) (*models.ExpenseResponse, error) {
	var expense models.Expense

	err := r.DB.
		Preload("PurchaseOrder").
		Preload("PurchaseOrder.PurchaseProducts").
		Preload("PurchaseOrder.PurchaseProducts.Product").
		Preload("PurchaseOrder.Supplier").
		Preload("MovementType").
		First(&expense, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrorResponse(404, "Ingreso no encontrado", err)
		}
		return nil, models.ErrorResponse(500, "Error interno al buscar ingreso", err)
	}

	response := models.ExpenseResponse{
		ID:        expense.ID,
		Details:   expense.Details,
		Amount:    expense.Amount,
		CreatedAt: expense.CreatedAt,
		MovementType: models.MovementTypeDTO{
			ID:       expense.MovementType.ID,
			Name:     expense.MovementType.Name,
			IsIncome: expense.MovementType.IsIncome,
		},
	}

	if expense.PurchaseOrder != nil {
		var purchaseProducts []models.PurchaseProductResponse

		for _, pp := range expense.PurchaseOrder.PurchaseProducts {
			purchaseProducts = append(purchaseProducts, models.PurchaseProductResponse{
				ID:              pp.ID,
				ProductID:       pp.ProductID,
				PurchaseOrderID: pp.PurchaseOrderID,
				ExpiredAt:       pp.ExpiredAt,
				UnitPrice:       pp.UnitPrice,
				Quantity:        pp.Quantity,
				TotalPrice:      pp.TotalPrice,
				CreatedAt:       pp.CreatedAt,
				Product: models.ProductDTO{
					ID:         pp.Product.ID,
					Identifier: pp.Product.Identifier,
					Name:       pp.Product.Name,
					CreatedAt:  pp.Product.CreatedAt,
				},
			})
		}

		response.PurchaseOrder = &models.PurchaseOrderResponse{
			ID:          expense.PurchaseOrder.ID,
			OrderNumber: expense.PurchaseOrder.OrderNumber,
			OrderDate:   expense.PurchaseOrder.OrderDate,
			Amount:      expense.PurchaseOrder.Amount,
			CreatedAt:   expense.PurchaseOrder.CreatedAt,
			UpdatedAt:   expense.PurchaseOrder.UpdatedAt,
			Supplier: models.SupplierResponse{
				ID:        expense.PurchaseOrder.Supplier.ID,
				Name:      expense.PurchaseOrder.Supplier.Name,
				Phone:     expense.PurchaseOrder.Supplier.Phone,
				Email:     expense.PurchaseOrder.Supplier.Email,
				CreatedAt: expense.PurchaseOrder.Supplier.CreatedAt,
			},
			PurchaseProducts: purchaseProducts,
		}
	}

	return &response, nil
}

func (r *ExpenseRepository) ExpenseGetAll(page, limit int) (*[]models.ExpenseDTO, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	var expenses []models.Expense

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
		return nil, models.ErrorResponse(500, "Error interno al buscar ingreso", err)
	}

	var expenseDTOs []models.ExpenseDTO

	for _, expense := range expenses {
		response := models.ExpenseDTO{
			ID:        expense.ID,
			Amount:    expense.Amount,
			CreatedAt: expense.CreatedAt,
			MovementType: models.MovementTypeDTO{
				ID:       expense.MovementType.ID,
				Name:     expense.MovementType.Name,
				IsIncome: expense.MovementType.IsIncome,
			},
		}

		if expense.PurchaseOrder != nil {
			response.PurchaseOrder = &models.PurchaseOrderDTO{
				ID:          expense.PurchaseOrder.ID,
				OrderNumber: expense.PurchaseOrder.OrderNumber,
				OrderDate:   expense.PurchaseOrder.OrderDate,
				Amount:      expense.PurchaseOrder.Amount,
				CreatedAt:   expense.PurchaseOrder.CreatedAt,
			}
		}

		expenseDTOs = append(expenseDTOs, response)
	}

	return &expenseDTOs, nil
}

func (r *ExpenseRepository) ExpenseGetToday(page, limit int) (*[]models.ExpenseDTO, error) {
	start := time.Now().Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	var expenses []models.Expense

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
		return nil, models.ErrorResponse(500, "Error interno al buscar ingreso", err)
	}

	var expenseDTOs []models.ExpenseDTO

	for _, expense := range expenses {
		response := models.ExpenseDTO{
			ID:        expense.ID,
			Amount:    expense.Amount,
			CreatedAt: expense.CreatedAt,
			MovementType: models.MovementTypeDTO{
				ID:       expense.MovementType.ID,
				Name:     expense.MovementType.Name,
				IsIncome: expense.MovementType.IsIncome,
			},
		}

		if expense.PurchaseOrder != nil {
			response.PurchaseOrder = &models.PurchaseOrderDTO{
				ID:          expense.PurchaseOrder.ID,
				OrderNumber: expense.PurchaseOrder.OrderNumber,
				OrderDate:   expense.PurchaseOrder.OrderDate,
				Amount:      expense.PurchaseOrder.Amount,
				CreatedAt:   expense.PurchaseOrder.CreatedAt,
			}
		}

		expenseDTOs = append(expenseDTOs, response)
	}

	return &expenseDTOs, nil
}

// func (r *ExpenseRepository) ExpenseGetByID(id string) (*models.ExpenseResponse, error) {
// 	var expense models.Expense
// 	if err := r.DB.Where("id = ?", id).First(&expense).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, models.ErrorResponse(404, "Movimiento no encontrado", err)
// 		}
// 		return nil, models.ErrorResponse(500, "Error interno al buscar movimiento", err)
// 	}
// 	return &expense, nil
// }
// func (r *ExpenseRepository) ExpenseGetAll() (*[]models.Expense, error) {
// 	var expenses []models.Expense
// 	if err := r.DB.Limit(100).Order("created_at desc").Find(&expenses).Error; err != nil {
// 		return nil, models.ErrorResponse(500, "Error interno al buscar movimientos", err)
// 	}
// 	return &expenses, nil
// }

// func (r *ExpenseRepository) ExpenseGetToday() (*[]models.Expense, error) {
// 	today := time.Now().Format("2006-01-02")
// 	var expenses []models.Expense
// 	if err := r.DB.Where("DATE(created_at) = ?", today).Order("created_at desc").Find(&expenses).Error; err != nil {
// 		return nil, models.ErrorResponse(500, "Error interno al buscar movimientos", err)
// 	}
// 	return &expenses, nil
// }

func (r *ExpenseRepository) ExpenseCreate(expense *models.ExpenseCreate) (string, error) {
	newID := uuid.NewString()
	if err := r.DB.Create(&models.Expense{
		ID:             newID,
		Details:        expense.Details,
		PurchaseOrderID:     expense.PurchaseOrderID,
		MovementTypeID: expense.MovementTypeID,
		Amount:         expense.Amount,
	}).Error; err != nil {
		return "", models.ErrorResponse(500, "Error interno al crear movimiento", err)
	}
	return newID, nil
}

func (r *ExpenseRepository) ExpenseUpdate(expense *models.ExpenseUpdate) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", expense.ID).
			Updates(&models.Expense{
				Details: expense.Details,
				PurchaseOrderID:     expense.PurchaseOrderID,
				MovementTypeID: expense.MovementTypeID,
				Amount:         expense.Amount,
			}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return models.ErrorResponse(404, "Movimiento no encontrado", err)
			}
			return models.ErrorResponse(500, "Error interno al actualizar movimiento", err)
		}
		return nil
	})
}

func (r *ExpenseRepository) ExpenseDelete(id string) error {
	if err := r.DB.Where("id = ?", id).Delete(&models.Expense{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrorResponse(404, "Movimiento no encontrado", err)
		}
		return models.ErrorResponse(500, "Error interno al eliminar movimiento", err)
	}
	return nil
}
