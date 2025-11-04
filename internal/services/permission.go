package services

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (p *PermissionService) PermissionByRoleID(roleID int64) (*[]string, error) {
	permissions, err := p.PermissionRepository.PermissionByRoleID(roleID)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (p *PermissionService) PermissionGetAll() (*[]schemas.PermissionResponse, error) {
	permissions, err := p.PermissionRepository.PermissionGetAll()
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (p *PermissionService) PermissionGetToMe(roleID int64) (*[]schemas.PermissionResponse, error) {
	permissions, err := p.PermissionRepository.PermissionGetToMe(roleID)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}