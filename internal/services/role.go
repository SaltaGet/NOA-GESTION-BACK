package services

import (
"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (r *RoleService) RoleGetAll() (*[]schemas.RoleResponse, error) {
	roles, err := r.RoleRepository.RoleGetAll()
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleService) RoleCreate(roleCrate *schemas.RoleCreate) (id string, err error) {
	id, err = r.RoleRepository.RoleCreate(roleCrate)
	if err != nil {
		return "", err
	}
	return id, nil
}