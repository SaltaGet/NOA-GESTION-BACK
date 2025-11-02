package services

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
)

func (r *RoleService) RoleGetAll() (*[]models.RoleResponse, error) {
	roles, err := r.RoleRepository.RoleGetAll()
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleService) RoleCreate(roleCrate *models.RoleCreate) (id string, err error) {
	id, err = r.RoleRepository.RoleCreate(roleCrate)
	if err != nil {
		return "", err
	}
	return id, nil
}