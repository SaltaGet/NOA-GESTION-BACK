package services

import (
"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (r *RoleService) RoleGetAll() (*[]schemas.RoleResponse, error) {
	return r.RoleRepository.RoleGetAll()
}

func (r *RoleService) RoleCreate(roleCreate *schemas.RoleCreate) (int64, error) {
	return r.RoleRepository.RoleCreate(roleCreate)
}