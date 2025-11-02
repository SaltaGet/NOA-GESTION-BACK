package ports

import "github.com/DanielChachagua/GestionCar/pkg/models"

type ClientService interface {
	ClientGetByID(id string) (client *models.Client, err error)
	ClientGetByName(name string) (clients *[]models.Client, err error)
	ClientGetAll() (clients *[]models.Client, err error)
	ClientCreate(clientCreate *models.ClientCreate) (id string, err error)
	ClientUpdate(clienUpdate *models.ClientUpdate) (err error)
	ClientDelete(id string) (err error)
}

type ClientRepository interface {
	ClientGetByID(id string) (client *models.Client, err error)
	ClientGetByName(name string) (clients *[]models.Client, err error)
	ClientExist(email string, dni string, cuil string) (err error)
	ClientGetAll() (clients *[]models.Client, err error)
	ClientCreate(clientCreate *models.ClientCreate) (id string, err error)
	ClientUpdate(clienUpdate *models.ClientUpdate) (err error)
	ClientDelete(id string) (err error)
}
