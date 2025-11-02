package services

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
)

func (s *ServiceService) ServiceGetByID(id string) (*models.Service, error) {
	service, err := s.ServiceRepository.ServiceGetByID(id)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (s *ServiceService) ServiceGetByName(name string) (*[]models.Service, error) {
	services, err := s.ServiceRepository.ServiceGetByName(name)
	if err != nil {
		return nil, err
	}
	return services, nil
}

func (s *ServiceService) ServiceGetAll() (*[]models.Service, error) {
	services, err := s.ServiceRepository.ServiceGetAll()
	if err != nil {
		return nil, err
	}
	return services, nil
}

func (s *ServiceService) ServiceCreate(service *models.ServiceCreate) (string, error) {
	exist, err := s.ServiceRepository.ServiceExistByName(service.Name)
	if err != nil {
		return "", err
	}

	if exist {
		return "", err
	}

	id, err := s.ServiceRepository.ServiceCreate(service)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *ServiceService) ServiceUpdate(service *models.ServiceUpdate) error {
	err := s.ServiceRepository.ServiceUpdate(service)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceService) ServiceDelete(id string) error {
	err := s.ServiceRepository.ServiceDelete(id)
	if err != nil {
		return err
	}
	return nil
}