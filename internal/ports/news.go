package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

type NewsRepository interface {
	NewsGetAll() ([]schemas.NewsResponseDTO, error)
	NewsGetByID(id int64) (*schemas.NewsResponse, error)
	NewsCreate(adminID int64, news *schemas.NewsCreate) (int64, error)
	NewsUpdate(adminID int64, news *schemas.NewsUpdate) (error)
	NewsDelete(adminID int64, id int64) (error)
}

type NewsServices interface {
	NewsGetAll() ([]schemas.NewsResponseDTO, error)
	NewsGetByID(id int64) (*schemas.NewsResponse, error)
	NewsCreate(adminID int64, news *schemas.NewsCreate) (int64, error)
	NewsUpdate(adminID int64, news *schemas.NewsUpdate) (error)
	NewsDelete(adminID int64, id int64) (error)
}