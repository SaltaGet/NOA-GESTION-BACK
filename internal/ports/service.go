package ports

import "github.com/DanielChachagua/GestionCar/pkg/models"

type ServiceRepository interface {
	ServiceGetByID(id string) (service *models.Service, err error)
	ServiceExistByName(name string) (exist bool, err error)
	ServiceGetByName(name string) (services *[]models.Service, err error)
	ServiceGetAll() (services *[]models.Service, err error)
	ServiceCreate(service *models.ServiceCreate) (id string, err error)
	ServiceUpdate(service *models.ServiceUpdate) (err error)
	ServiceDelete(id string) error
}

type ServiceService interface {
	ServiceGetByID(id string) (service *models.Service, err error)
	ServiceGetByName(name string) (services *[]models.Service, err error)
	ServiceGetAll() (services *[]models.Service, err error)
	ServiceCreate(service *models.ServiceCreate) (id string, err error)
	ServiceUpdate(service *models.ServiceUpdate) (err error)
	ServiceDelete(id string) error
}

