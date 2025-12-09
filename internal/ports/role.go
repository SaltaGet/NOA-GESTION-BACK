package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"


type RoleService interface {
	RoleGetByID(id int64) (role *schemas.RoleResponse, err error)
	RoleGetAll() (roles *[]schemas.RoleResponse, err error)
	RoleCreate(roleCrate *schemas.RoleCreate) (id int64, err error)
	RoleUpdate(roleUpdate *schemas.RoleUpdate) (err error)
}

type RoleRepository interface {
	RoleGetByID(id int64) (role *schemas.RoleResponse, err error)
	RoleGetAll() (roles *[]schemas.RoleResponse, err error)
	RoleCreate(roleCrate *schemas.RoleCreate) (id int64, err error)
	RoleUpdate(roleUpdate *schemas.RoleUpdate) (err error)
}