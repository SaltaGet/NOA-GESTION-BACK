package repositories

import (
	"errors"

	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *VehicleRepository) VehicleGetByID(id string) (*models.Vehicle, error) {
	var vehicle models.Vehicle
	if err := r.DB.Where("id = ?", id).First(&vehicle).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrorResponse(404, "Vehiculo no encontrado", err)
		}
		return nil, models.ErrorResponse(500, "Error interno al buscar vehiculo", err)
	}
	return &vehicle, nil
}

func (r *VehicleRepository) VehicleGetByDomain(domain string) (*[]models.Vehicle, error) {
	var vehicles []models.Vehicle
	if err := r.DB.Preload("Client").Where("domain LIKE ?", "%"+domain+"%").Find(&vehicles).Error; err != nil {
		return nil, models.ErrorResponse(500, "Error interno al buscar vehiculo", err)
	}
	return &vehicles, nil
}

func (r *VehicleRepository) VehicleExistByDomain(domain string) (bool, error) {
	var vehicle models.Vehicle
	if err := r.DB.Preload("Client").Where("domain = ?", domain).First(&vehicle).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, models.ErrorResponse(500, "Error interno al buscar vehiculo", err)
	}
	return true, nil
}

func (r *VehicleRepository) VehicleGetAll() (*[]models.Vehicle, error) {
	var vehicles []models.Vehicle
	if err := r.DB.Find(&vehicles).Error; err != nil {
		return nil, models.ErrorResponse(500, "Error interno al buscar vehiculo", err)
	}
	return &vehicles, nil
}

func (r *VehicleRepository) VehicleGetByClientID(clientID string) (*[]models.Vehicle, error) {
	var vehicles []models.Vehicle
	if err := r.DB.Where("client_id = ?", clientID).Find(&vehicles).Error; err != nil {
		return nil, models.ErrorResponse(500, "Error interno al buscar vehiculo", err)
	}
	return &vehicles, nil
}

func (r *VehicleRepository) VehicleCreate(vehicle *models.VehicleCreate) (string, error) {
	newVehicle := models.Vehicle{
		ID: uuid.NewString(),
		Brand: vehicle.Brand,
		Model: vehicle.Model,
		Color: vehicle.Color,
		Year: vehicle.Year,
		Domain: vehicle.Domain,
		ClientID: vehicle.ClientID,
	}
	if err := r.DB.Create(&newVehicle).Error; err != nil {
		return "", models.ErrorResponse(500, "Error interno al crear vehiculo", err)
	}
	return newVehicle.ID, nil
}

func (r *VehicleRepository) VehicleUpdate(vehicle *models.VehicleUpdate) error {
	if err := r.DB.Where("id = ?", vehicle.ID).Updates(&vehicle).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrorResponse(404, "Vehiculo no encontrado", err)
		}
		return models.ErrorResponse(500, "Error interno al actualizar vehiculo", err) 
	}

	return nil
}

func (r *VehicleRepository) VehicleDelete(id string) error {
	var vehicle models.Vehicle
	if err := r.DB.Where("id = ?", id).Delete(&vehicle).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrorResponse(404, "Vehiculo no encontrado", err)
		}
		return models.ErrorResponse(500, "Error interno al eliminar vehiculo", err)
	}
	return nil
}


// func (r *Repository) GetVehicleByClientIDAndDomain(clientID, domain string) (*models.Vehicle, error) {
// 	var vehicle models.Vehicle
// 	if err := r.DB.Where("client_id = ? AND domain = ?", clientID, domain).First(&vehicle).Error; err != nil {
// 		return nil, err
// 	}
// 	return &vehicle, nil
// }

// func (r *Repository) GetVehicleByClientIDAndDomainLike(clientID, domain string) ([]models.Vehicle, error) {
// 	var vehicles []models.Vehicle
// 	if err := r.DB.Where("client_id = ? AND domain LIKE ?", clientID, "%"+domain+"%").Find(&vehicles).Error; err != nil {
// 		return nil, err
// 	}
// 	return vehicles, nil
// }
