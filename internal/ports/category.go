package ports

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)


type CategoryRepository interface {
	CategoryGetByID(id int64) (*models.Category, error)
	CategoryGetAll() ([]*models.Category, error)
	CategoryCreate(categoryCreate *schemas.CategoryCreate) (int64, error)
	CategoryUpdate(categoryUpdate *schemas.CategoryUpdate) error
	CategoryDelete(id int64) error
}

type CategoryService interface {
	CategoryGetByID(id int64) (*schemas.CategoryResponse, error)
	CategoryGetAll() ([]*schemas.CategoryResponse, error)
	CategoryCreate(categoryCreate *schemas.CategoryCreate) (int64, error)
	CategoryUpdate(categoryUpdate *schemas.CategoryUpdate) error
	CategoryDelete(id int64) error
}