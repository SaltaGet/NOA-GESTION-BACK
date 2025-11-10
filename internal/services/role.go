package services

import (
"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (r *RoleService) RoleGetAll() (*[]schemas.RoleResponse, error) {
	return r.RoleRepository.RoleGetAll()
}
