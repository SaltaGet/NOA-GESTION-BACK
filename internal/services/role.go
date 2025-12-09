package services

import (
"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (r *RoleService) RoleGetByID(id int64) (*schemas.RoleResponse, error) {
	return r.RoleRepository.RoleGetByID(id)
}

func (r *RoleService) RoleGetAll() (*[]schemas.RoleResponse, error) {
	return r.RoleRepository.RoleGetAll()
}

func (r *RoleService) RoleCreate(roleCreate *schemas.RoleCreate) (int64, error) {
	return r.RoleRepository.RoleCreate(roleCreate)
}

func (r *RoleService) RoleUpdate(roleUpdate *schemas.RoleUpdate) error {
	return r.RoleRepository.RoleUpdate(roleUpdate)
}