package services

import (
"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (m *MovementTypeService) MovementTypeCreate(movementType *schemas.MovementTypeCreate) (string, error) {
	id, err := m.MovementTypeRepository.MovementTypeCreate(movementType)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (m *MovementTypeService) MovementTypeUpdate(movementType *schemas.MovementTypeUpdate) error {
	err := m.MovementTypeRepository.MovementTypeUpdate(movementType)
	if err != nil {
		return err
	}

	return nil
}

func (m *MovementTypeService) MovementTypeDelete(id string) error {
	err := m.MovementTypeRepository.MovementTypeDelete(id)
	if err != nil {
		return err
	}
	return nil
}

func (m *MovementTypeService) MovementTypeGetByID(id string) (*schemas.MovementType, error) {
	movementType, err := m.MovementTypeRepository.MovementTypeGetByID(id)
	if err != nil {
		return nil, err
	}
	return movementType, nil
}

func (m *MovementTypeService) MovementTypeGetAll(isIncome bool) (*[]schemas.MovementType, error) {
	movementTypes, err := m.MovementTypeRepository.MovementTypeGetAll(isIncome)
	if err != nil {
		return nil, err
	}
	return movementTypes, nil
}