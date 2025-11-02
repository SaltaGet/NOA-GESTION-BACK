package services

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
)

// v VehicleService

func (v *VehicleService) VehicleCreate(vehicleCreate *models.VehicleCreate) (string , error) {
	exist, err := v.VehicleRepository.VehicleExistByDomain(vehicleCreate.Domain)
	if err != nil {
		return "", err
	}

	if exist {
		return "", err
	}

	vehicle, err := v.VehicleRepository.VehicleCreate(vehicleCreate)

	if err != nil {
		return "", err
	}

	return vehicle, nil
}

func (v *VehicleService) VehicleGetAll() (*[]models.Vehicle, error) {
	vehicles, err := v.VehicleRepository.VehicleGetAll()
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (v *VehicleService) VehicleGetByID(id string) (*models.Vehicle, error) {
	vehicle, err := v.VehicleRepository.VehicleGetByID(id)
	if err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (v *VehicleService) VehicleGetByDomain(domain string) (*[]models.Vehicle, error) {
	vehicle, err := v.VehicleRepository.VehicleGetByDomain(domain)
	if err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (v *VehicleService) VehicleGetByClientID(clientID string) (*[]models.Vehicle, error) {
	vehicles, err := v.VehicleRepository.VehicleGetByClientID(clientID)
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (v *VehicleService) VehicleUpdate(vehicleUpdate *models.VehicleUpdate) error {
	err := v.VehicleRepository.VehicleUpdate(vehicleUpdate)

	if err != nil {
		return err
	}
	return nil
}

func (v *VehicleService) VehicleDelete(id string) (error) {
	err := v.VehicleRepository.VehicleDelete(id)
	if err != nil {
		return err
	}
	return nil
}
