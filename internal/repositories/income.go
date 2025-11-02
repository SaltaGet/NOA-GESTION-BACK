package repositories

import (
	"errors"
	"time"

	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *IncomeRepository) IncomeGetByID(id string) (*models.IncomeResponse, error) {
	var income models.Income

	err := r.DB.
		Preload("Client").
		Preload("Vehicle").
		Preload("Employee").
		Preload("MovementType").
		Preload("Services").
		First(&income, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrorResponse(404, "Ingreso no encontrado", err)
		}
		return nil, models.ErrorResponse(500, "Error interno al buscar ingreso", err)
	}

	response := models.IncomeResponse{
		ID:      income.ID,
		Ticket:  income.Ticket,
		Details: income.Details,
		Amount:  income.Amount,
		CreatedAt: income.CreatedAt,
		Client: models.ClientResponse{
			ID:        income.Client.ID,
			FirstName: income.Client.FirstName,
			LastName:  income.Client.LastName,
			Cuil:      income.Client.Cuil,
			Dni:       income.Client.Dni,
			Email:     income.Client.Email,
		},
		Vehicle: models.VehicleResponse{
			ID:     income.Vehicle.ID,
			Brand:  income.Vehicle.Brand,
			Model:  income.Vehicle.Model,
			Color:  income.Vehicle.Color,
			Year:   income.Vehicle.Year,
			Domain: income.Vehicle.Domain,
		},
		Employee: models.EmployeeResponse{
			ID:    income.Employee.ID,
			Name:  income.Employee.Name,
			Phone: income.Employee.Phone,
			Email: income.Employee.Email,
		},
		MovementType: models.MovementTypeDTO{
			ID:       income.MovementType.ID,
			Name:     income.MovementType.Name,
			IsIncome: income.MovementType.IsIncome,
		},
	}

	for _, s := range income.Services {
		response.Services = append(response.Services, models.ServiceDTO{
			ID:   s.ID,
			Name: s.Name,
		})
	}
	return &response, nil
}

func (r *IncomeRepository) IncomeGetAll(page, limit int) (*[]models.IncomeDTO, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	var incomes []models.Income

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
		return nil, models.ErrorResponse(500, "Error interno al buscar ingresos", err)
	}

	var incomeDTOs []models.IncomeDTO

	for _, income := range incomes {
		incomeDTO := models.IncomeDTO{
			ID:     income.ID,
			Ticket: income.Ticket,
			Amount: income.Amount,
			CreatedAt: income.CreatedAt,
			Client: models.ClientDTO{
				ID:        income.Client.ID,
				FirstName: income.Client.FirstName,
				LastName:  income.Client.LastName,
			},
			Vehicle: models.VehicleDTO{
				ID:     income.Vehicle.ID,
				Domain: income.Vehicle.Domain,
			},
			Employee: models.EmployeeDTO{
				ID:   income.Employee.ID,
				Name: income.Employee.Name,
			},
			MovementType: models.MovementTypeDTO{
				ID:       income.MovementType.ID,
				Name:     income.MovementType.Name,
				IsIncome: income.MovementType.IsIncome,
			},
		}

		for _, s := range income.Services {
			incomeDTO.Services = append(incomeDTO.Services, models.ServiceDTO{
				ID:   s.ID,
				Name: s.Name,
			})
		}

		incomeDTOs = append(incomeDTOs, incomeDTO)
	}

	return &incomeDTOs, nil
}

// func (r *IncomeRepository) IncomeGetAll() (*[]models.Income, error) {
// 	var incomes []models.Income
// 	if err := r.DB.Limit(100).Order("created_at desc").Find(&incomes).Error; err != nil {
// 		return nil, models.ErrorResponse(500, "Error interno al buscar movimientos", err)
// 	}
// 	return &incomes, nil
// }

func (r *IncomeRepository) IncomeGetToday(page, limit int) (*[]models.IncomeDTO, error) {
	start := time.Now().Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	var incomes []models.Income

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
		return nil, models.ErrorResponse(500, "Error interno al buscar ingresos", err)
	}

	var incomeDTOs []models.IncomeDTO

	for _, income := range incomes {
		incomeDTO := models.IncomeDTO{
			ID:     income.ID,
			Ticket: income.Ticket,
			Amount: income.Amount,
			CreatedAt: income.CreatedAt,
			Client: models.ClientDTO{
				ID:        income.Client.ID,
				FirstName: income.Client.FirstName,
				LastName:  income.Client.LastName,
			},
			Vehicle: models.VehicleDTO{
				ID:     income.Vehicle.ID,
				Domain: income.Vehicle.Domain,
			},
			Employee: models.EmployeeDTO{
				ID:   income.Employee.ID,
				Name: income.Employee.Name,
			},
			MovementType: models.MovementTypeDTO{
				ID:       income.MovementType.ID,
				Name:     income.MovementType.Name,
				IsIncome: income.MovementType.IsIncome,
			},
		}

		for _, s := range income.Services {
			incomeDTO.Services = append(incomeDTO.Services, models.ServiceDTO{
				ID:   s.ID,
				Name: s.Name,
			})
		}

		incomeDTOs = append(incomeDTOs, incomeDTO)
	}

	return &incomeDTOs, nil
}

// func (r *IncomeRepository) IncomeGetToday() (*[]models.Income, error) {
// 	today := time.Now().Format("2006-01-02")
// 	var incomes []models.Income
// 	if err := r.DB.Where("DATE(created_at) = ?", today).Order("created_at desc").Find(&incomes).Error; err != nil {
// 		return nil, models.ErrorResponse(500, "Error interno al buscar movimientos", err)
// 	}
// 	return &incomes, nil
// }

func (r *IncomeRepository) IncomeCreate(income *models.IncomeCreate) (string, error) {
	newID := uuid.NewString()

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var services []models.Service

		if err := tx.Where("id IN ?", income.ServicesID).Find(&services).Error; err != nil {
			return models.ErrorResponse(500, "Error interno al buscar servicios", err)
		}

		if err := tx.Create(&models.Income{
			ID:             newID,
			Ticket:         income.Ticket,
			Details:        income.Details,
			ClientID:       income.ClientID,
			VehicleID:      income.VehicleID,
			EmployeeID:     income.EmployeeID,
			Amount:         income.Amount,
			MovementTypeID: income.MovementTypeID,
			Services:       services,
		}).Error; err != nil {
			return models.ErrorResponse(500, "Error interno al crear movimiento", err)
		}

		return nil
	})

	if err != nil {
		return "", models.ErrorResponse(500, "Error interno al crear movimiento", err)
	}

	return newID, nil
}

func (r *IncomeRepository) IncomeUpdate(incomeUpdate *models.IncomeUpdate) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var income models.Income

		if err := tx.Where("id = ?", incomeUpdate.ID).First(&income).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return models.ErrorResponse(404, "Movimiento no encontrado", err)
			}
			return models.ErrorResponse(500, "Error interno al actualizar movimiento", err)
		}

		income.Ticket = incomeUpdate.Ticket
		income.Details = incomeUpdate.Details
		income.ClientID = incomeUpdate.ClientID
		income.VehicleID = incomeUpdate.VehicleID
		income.EmployeeID = incomeUpdate.EmployeeID
		income.Amount = incomeUpdate.Amount
		income.MovementTypeID = incomeUpdate.MovementTypeID
		income.UpdatedAt = time.Now().UTC()

		var services []models.Service
		if err := tx.Where("id IN ?", incomeUpdate.ServicesID).Find(&services).Error; err != nil {
			return models.ErrorResponse(500, "Error interno al buscar servicios", err)
		}

		if err := tx.Model(&income).Association("Services").Replace(services); err != nil {
			return models.ErrorResponse(500, "Error interno al actualizar servicios", err)
		}

		if err := tx.Save(&income).Error; err != nil {
			return models.ErrorResponse(500, "Error interno al actualizar movimiento", err)
		}

		return nil
	})
}

func (r *IncomeRepository) IncomeDelete(id string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&models.Income{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return models.ErrorResponse(404, "Movimiento no encontrado", err)
			}
			return models.ErrorResponse(500, "Error interno al eliminar movimiento", err)
		}
		return nil
	})
}
