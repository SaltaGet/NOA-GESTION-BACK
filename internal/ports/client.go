package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"


type ClientService interface {
	ClientGetByID(id string) (client *schemas.Client, err error)
	ClientGetByName(name string) (clients *[]schemas.Client, err error)
	ClientGetAll() (clients *[]schemas.Client, err error)
	ClientCreate(clientCreate *schemas.ClientCreate) (id string, err error)
	ClientUpdate(clienUpdate *schemas.ClientUpdate) (err error)
	ClientDelete(id string) (err error)
}

type ClientRepository interface {
	ClientGetByID(id string) (client *schemas.Client, err error)
	ClientGetByName(name string) (clients *[]schemas.Client, err error)
	ClientExist(email string, dni string, cuil string) (err error)
	ClientGetAll() (clients *[]schemas.Client, err error)
	ClientCreate(clientCreate *schemas.ClientCreate) (id string, err error)
	ClientUpdate(clienUpdate *schemas.ClientUpdate) (err error)
	ClientDelete(id string) (err error)
}
