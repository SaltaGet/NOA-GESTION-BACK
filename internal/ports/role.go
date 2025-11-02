package ports

import "github.com/DanielChachagua/GestionCar/pkg/models"

type RoleService interface {
	RoleGetAll() (roles *[]models.RoleResponse, err error)
	RoleCreate(roleCrate *models.RoleCreate) (id string, err error)
}

type RoleRepository interface {
	RoleGetAll() (roles *[]models.RoleResponse, err error)
	RoleCreate(roleCrate *models.RoleCreate) (id string, err error)
}