package repositories

import (
	"errors"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *IncomeRepository) IncomeGetByID(id string) (*schemas.IncomeResponse, error) {
	var income schemas.IncomeResponse

	err := r.DB.
		Preload("Client").
		Preload("Vehicle").
		Preload("Employee").
		Preload("MovementType").
		Preload("Services").
		First(&income, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Ingreso no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al buscar ingreso", err)
	}

	response := schemas.IncomeResponse{
		ID:      income.ID,
		Ticket:  income.Ticket,
		Details: income.Details,
		Amount:  income.Amount,
		CreatedAt: income.CreatedAt,
		Client: schemas.ClientResponse{
			ID:        income.Client.ID,
			FirstName: income.Client.FirstName,
			LastName:  income.Client.LastName,
			Email:     income.Client.Email,
		},
		MovementType: schemas.MovementTypeDTO{
			ID:       income.MovementType.ID,
			Name:     income.MovementType.Name,
			IsIncome: income.MovementType.IsIncome,
		},
	}

	return &response, nil
}

func (r *IncomeRepository) IncomeGetAll(page, limit int) (*[]schemas.IncomeDTO, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	var incomes []schemas.IncomeResponse

	err := r.DB.Preload("Client").
		Preload("Vehicle").
		Preload("Employee").
		Preload("MovementType").
		Preload("Services").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&incomes).Error

	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al buscar ingresos", err)
	}

	var incomeDTOs []schemas.IncomeDTO

	for _, income := range incomes {
		incomeDTO := schemas.IncomeDTO{
			ID:     income.ID,
			Ticket: income.Ticket,
			Amount: income.Amount,
			CreatedAt: income.CreatedAt,
			Client: schemas.ClientResponseDTO{
				ID:        income.Client.ID,
				FirstName: income.Client.FirstName,
				LastName:  income.Client.LastName,
			},
			MovementType: schemas.MovementTypeDTO{
				ID:       income.MovementType.ID,
				Name:     income.MovementType.Name,
				IsIncome: income.MovementType.IsIncome,
			},
		}

		incomeDTOs = append(incomeDTOs, incomeDTO)
	}

	return &incomeDTOs, nil
}

// func (r *IncomeRepository) IncomeGetAll() (*[]schemas.Income, error) {
// 	var incomes []schemas.Income
// 	if err := r.DB.Limit(100).Order("created_at desc").Find(&incomes).Error; err != nil {
// 		return nil, schemas.ErrorResponse(500, "Error interno al buscar movimientos", err)
// 	}
// 	return &incomes, nil
// }

func (r *IncomeRepository) IncomeGetToday(page, limit int) (*[]schemas.IncomeDTO, error) {
	start := time.Now().Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	var incomes []schemas.IncomeResponse

	err := r.DB.Preload("Client").
		Preload("Vehicle").
		Preload("Employee").
		Preload("MovementType").
		Preload("Services").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Where("created_at >= ? AND created_at < ?", start, end).
		Find(&incomes).Error

	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al buscar ingresos", err)
	}

	var incomeDTOs []schemas.IncomeDTO

	for _, income := range incomes {
		incomeDTO := schemas.IncomeDTO{
			ID:     income.ID,
			Ticket: income.Ticket,
			Amount: income.Amount,
			CreatedAt: income.CreatedAt,
			Client: schemas.ClientResponseDTO{
				ID:        income.Client.ID,
				FirstName: income.Client.FirstName,
				LastName:  income.Client.LastName,
			},
			MovementType: schemas.MovementTypeDTO{
				ID:       income.MovementType.ID,
				Name:     income.MovementType.Name,
				IsIncome: income.MovementType.IsIncome,
			},
		}

		incomeDTOs = append(incomeDTOs, incomeDTO)
	}

	return &incomeDTOs, nil
}

// func (r *IncomeRepository) IncomeGetToday() (*[]schemas.Income, error) {
// 	today := time.Now().Format("2006-01-02")
// 	var incomes []schemas.Income
// 	if err := r.DB.Where("DATE(created_at) = ?", today).Order("created_at desc").Find(&incomes).Error; err != nil {
// 		return nil, schemas.ErrorResponse(500, "Error interno al buscar movimientos", err)
// 	}
// 	return &incomes, nil
// }

func (r *IncomeRepository) IncomeCreate(income *schemas.IncomeCreate) (string, error) {
	newID := uuid.NewString()

	err := r.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&schemas.IncomeResponse{
			ID:             newID,
			Ticket:         income.Ticket,
			Details:        income.Details,
			Amount:         income.Amount,
		}).Error; err != nil {
			return schemas.ErrorResponse(500, "Error interno al crear movimiento", err)
		}

		return nil
	})

	if err != nil {
		return "", schemas.ErrorResponse(500, "Error interno al crear movimiento", err)
	}

	return newID, nil
}

func (r *IncomeRepository) IncomeUpdate(incomeUpdate *schemas.IncomeUpdate) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var income schemas.IncomeResponse

		if err := tx.Where("id = ?", incomeUpdate.ID).First(&income).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Movimiento no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error interno al actualizar movimiento", err)
		}

		income.Ticket = incomeUpdate.Ticket
		income.Details = incomeUpdate.Details
		income.Amount = incomeUpdate.Amount

		if err := tx.Save(&income).Error; err != nil {
			return schemas.ErrorResponse(500, "Error interno al actualizar movimiento", err)
		}

		return nil
	})
}

func (r *IncomeRepository) IncomeDelete(id string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&models.IncomeSale{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Movimiento no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error interno al eliminar movimiento", err)
		}
		return nil
	})
}
