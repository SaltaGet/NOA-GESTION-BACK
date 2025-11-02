package ports

import "github.com/DanielChachagua/GestionCar/pkg/models"

type VehicleService interface {
	VehicleGetByID(id string) (vehicle *models.Vehicle, err error)
	VehicleGetAll() (vehicles *[]models.Vehicle, err error)
	VehicleGetByDomain(domain string) (vehicles *[]models.Vehicle, err error)
	VehicleGetByClientID(clientID string) (vehicles *[]models.Vehicle, err error)
	VehicleCreate(vehicleCreate *models.VehicleCreate) (id string, err error)
	VehicleUpdate(vehicleUpdate *models.VehicleUpdate) (err error)
	VehicleDelete(id string) (err error)
}

type VehicleRepository interface {
	VehicleGetByID(id string) (vehicle *models.Vehicle, err error)
	VehicleGetAll() (vehicles *[]models.Vehicle, err error)
	VehicleGetByDomain(domain string) (vehicles *[]models.Vehicle, err error)
	VehicleExistByDomain(domain string) (exist bool, err error)
	VehicleGetByClientID(clientID string) (vehicles *[]models.Vehicle, err error)
	VehicleCreate(vehicleCreate *models.VehicleCreate) (id string, err error)
	VehicleUpdate(vehicleUpdate *models.VehicleUpdate) (err error)
	VehicleDelete(id string) (err error)
}
