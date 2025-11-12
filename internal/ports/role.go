package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"


type RoleService interface {
	RoleGetAll() (roles *[]schemas.RoleResponse, err error)
	RoleCreate(roleCrate *schemas.RoleCreate) (id int64, err error)
}

type RoleRepository interface {
	RoleGetAll() (roles *[]schemas.RoleResponse, err error)
	RoleCreate(roleCrate *schemas.RoleCreate) (id int64, err error)
}