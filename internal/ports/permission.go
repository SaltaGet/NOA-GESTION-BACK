package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"



type PermissionService interface {
	PermissionByRoleID(roleID int64) (permissions *[]string, err error)
	PermissionGetAll() (permissions *[]schemas.PermissionResponse, err error)
	PermissionGetToMe(roleID int64) (permissions *[]schemas.PermissionResponse, err error)
}

type PermissionRepository interface {
	PermissionByRoleID(roleID int64) (permissions *[]string, err error)
	PermissionGetAll() (permissions *[]schemas.PermissionResponse, err error)
	PermissionGetToMe(roleID int64) (permissions *[]schemas.PermissionResponse, err error)
}