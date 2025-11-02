package repositories

import (
	"errors"
	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *ServiceRepository) ServiceGetByID(id string) (*models.Service, error) {
	var service models.Service
	if err := r.DB.Where("id = ?", id).First(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrorResponse(404, "Servicio no encontrado", err)
		}
		return nil, models.ErrorResponse(500, "Error interno al buscar servicio", err)
	}
	return &service, nil
}

func (r *ServiceRepository) ServiceExistByName(name string) (bool, error) {
	var service models.Service
	if err := r.DB.Where("name = ?", name).First(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *ServiceRepository) ServiceGetByName(name string) (*[]models.Service, error) {
	var services []models.Service
	if err := r.DB.Limit(5).Where("name LIKE ?", "%"+name+"%").Find(&services).Error; err != nil {
		return nil, models.ErrorResponse(500, "Error interno al buscar servicio", err)
	}

	return &services, nil
}

func (r *ServiceRepository) ServiceGetAll() (*[]models.Service, error) {
	var services []models.Service
	if err := r.DB.Find(&services).Error; err != nil {
		return nil, models.ErrorResponse(500, "Error interno al buscar servicios", err)
	}
	return &services, nil
}

func (r *ServiceRepository) ServiceCreate(service *models.ServiceCreate) (string, error) {
	newID := uuid.NewString()
	if err := r.DB.Create(&models.Service{
		ID:   newID,
		Name: service.Name,
	}).Error; err != nil {
		return "", models.ErrorResponse(500, "Error interno al crear servicio", err)
	}
	return newID, nil
}

func (r *ServiceRepository) ServiceUpdate(service *models.ServiceUpdate) error {
	if err := r.DB.Where("id = ?", service.ID).First(&models.Service{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrorResponse(404, "Servicio no encontrado", err)
		}
		return models.ErrorResponse(500, "Error interno al buscar servicio", err)
	}
	s := models.Service{
		ID:   service.ID,
		Name: service.Name,
	}
	if err := r.DB.Save(&s).Error; err != nil {
		return models.ErrorResponse(500, "Error interno al actualizar servicio", err)
	}
	return nil
}

func (r *ServiceRepository) ServiceDelete(id string) error {
	if err := r.DB.Where("id = ?", id).Delete(&models.Service{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrorResponse(404, "Servicio no encontrado", err)
		}
		return models.ErrorResponse(500, "Error interno al eliminar servicio", err)
	}
	return nil
}
