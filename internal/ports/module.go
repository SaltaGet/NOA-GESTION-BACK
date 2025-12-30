package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

type ModuleRepository interface {
	// ModuleGetAll() ([]*schemas.Module, error)
	// ModuleGetById(id int64) (*schemas.Module, error)
	// ModuleCreate(module *schemas.Module) error
	// ModuleUpdate(module *schemas.Module) error
	// ModuleDelete(id int64) error
	ModuleGet(id int64) (*schemas.ModuleResponse, error)
	ModuleGetAll() ([]schemas.ModuleResponse, error)
	ModuleCreate(module *schemas.ModuleCreate) (int64, error)
	ModuleUpdate(module *schemas.ModuleUpdate) error
	ModuleDelete(id int64) error
	ModuleAddTenant(module *schemas.ModuleAddTenant) error
	ModuleGetByTenantID(tenantID int64) ([]schemas.ModuleResponseDTO, error)
}

type ModuleService interface {
	ModuleGet(id int64) (*schemas.ModuleResponse, error)
	ModuleGetAll() ([]schemas.ModuleResponse, error)
	ModuleCreate(module *schemas.ModuleCreate) (int64, error)
	ModuleUpdate(module *schemas.ModuleUpdate) error
	ModuleDelete(id int64) error
	ModuleAddTenant(module *schemas.ModuleAddTenant) error
}
