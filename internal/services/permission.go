package services

import "github.com/DanielChachagua/GestionCar/pkg/models"

func (p *PermissionService) PermissionByRoleID(roleID string) (*[]string, error) {
	permissions, err := p.PermissionRepository.PermissionByRoleID(roleID)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (p *PermissionService) PermissionGetAll() (*[]models.PermissionResponse, error) {
	permissions, err := p.PermissionRepository.PermissionGetAll()
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (p *PermissionService) PermissionGetToMe(roleID string) (*[]models.PermissionResponse, error) {
	permissions, err := p.PermissionRepository.PermissionGetToMe(roleID)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}