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

func (r *RoleService) RoleCreate(memberID int64, roleCreate *schemas.RoleCreate) (int64, error) {
	return r.RoleRepository.RoleCreate(memberID, roleCreate)
}

func (r *RoleService) RoleUpdate(memberID int64, roleUpdate *schemas.RoleUpdate) error {
	return r.RoleRepository.RoleUpdate(memberID, roleUpdate)
}