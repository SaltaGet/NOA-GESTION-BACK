package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

type ModuleRepository interface {
	// ModuleGetAll() ([]*schemas.Module, error)
	// ModuleGetById(id int64) (*schemas.Module, error)
	// ModuleCreate(module *schemas.Module) error
	// ModuleUpdate(module *schemas.Module) error
	// ModuleDelete(id int64) error
	ModuleGetByTenantID(tenantID int64) ([]schemas.ModuleResponseDTO, error)
}

type ModuleService interface {
	// ModuleGetAll() ([]*schemas.Module, error)
	// ModuleGetById(id int64) (*schemas.Module, error)
	// ModuleCreate(module *schemas.Module) error
	// ModuleUpdate(module *schemas.Module) error
	// ModuleDelete(id int64) error
}