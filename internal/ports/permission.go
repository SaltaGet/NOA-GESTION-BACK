package ports

import "github.com/DanielChachagua/GestionCar/pkg/models"


type PermissionService interface {
	PermissionByRoleID(roleID string) (permissions *[]string, err error)
	PermissionGetAll() (permissions *[]models.PermissionResponse, err error)
	PermissionGetToMe(roleID string) (permissions *[]models.PermissionResponse, err error)
}

type PermissionRepository interface {
	PermissionByRoleID(roleID string) (permissions *[]string, err error)
	PermissionGetAll() (permissions *[]models.PermissionResponse, err error)
	PermissionGetToMe(roleID string) (permissions *[]models.PermissionResponse, err error)
}