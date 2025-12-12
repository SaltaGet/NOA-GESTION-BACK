package services

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (p *PermissionService) PermissionByRoleID(roleID int64) (*[]string, error) {
	return p.PermissionRepository.PermissionByRoleID(roleID)
}

func (p *PermissionService) PermissionGetAll() (*[]schemas.PermissionResponse, error) {
	return p.PermissionRepository.PermissionGetAll()
}

func (p *PermissionService) PermissionGetToMe(roleID int64) (*[]schemas.PermissionResponse, error) {
	return p.PermissionRepository.PermissionGetToMe(roleID)
}

func (p *PermissionService) PermissionUpdateAll() (error) {
	return p.PermissionRepository.PermissionUpdateAll()
}

