package services

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

func (m *ModuleService) ModuleGet(id int64) (*schemas.ModuleResponse, error) {
	return m.ModuleRepository.ModuleGet(id)
}

func (m *ModuleService) ModuleGetAll() ([]schemas.ModuleResponse, error) {
	return m.ModuleRepository.ModuleGetAll()
}

func (m *ModuleService) ModuleCreate(module *schemas.ModuleCreate) (int64, error) {
	return m.ModuleRepository.ModuleCreate(module)
}

func (m *ModuleService) ModuleUpdate(module *schemas.ModuleUpdate) error {
	return m.ModuleRepository.ModuleUpdate(module)
}

func (m *ModuleService) ModuleDelete(id int64) error {
	return m.ModuleRepository.ModuleDelete(id)
}

func (m *ModuleService) ModuleAddTenant(module *schemas.ModuleAddTenant) error {
	return m.ModuleRepository.ModuleAddTenant(module)
}