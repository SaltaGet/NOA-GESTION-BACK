package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"


type ClientService interface {
	ClientGetByID(id int64) (client *schemas.ClientResponse, err error)
	ClientGetByFilter(search string) (clients *[]schemas.ClientResponseDTO, err error)
	ClientGetAll(limit, page int64, search *map[string]string) (*[]schemas.ClientResponseDTO, int64, error)
	ClientCreate(memberID int64, clientCreate *schemas.ClientCreate) (id int64, err error)
	ClientUpdate(clienUpdate *schemas.ClientUpdate) (err error)
	ClientUpdateCredit(pointSaleID int64, clienUpdateCredit *schemas.ClientUpdateCredit) (err error)
	ClientDelete(id int64) (err error)
}

type ClientRepository interface {
	ClientGetByID(id int64) (client *schemas.ClientResponse, err error)
	ClientGetByFilter(search string) (clients *[]schemas.ClientResponseDTO, err error)
	ClientGetAll(limit, page int64, search *map[string]string) (*[]schemas.ClientResponseDTO, int64, error)
	ClientCreate(memberID int64, clientCreate *schemas.ClientCreate) (id int64, err error)
	ClientUpdate(clienUpdate *schemas.ClientUpdate) (err error)
	ClientUpdateCredit(pointSaleID int64, clienUpdateCredit *schemas.ClientUpdateCredit) (err error)
	ClientDelete(id int64) (err error)
}
