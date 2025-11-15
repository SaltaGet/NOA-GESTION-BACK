package services

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

func (t *TypeMovementService) TypeMovementCreate(movementCreate schemas.TypeMovementCreate) (error) {
	return t.TypeMovementRepository.TypeMovementCreate(movementCreate)
}

func (t *TypeMovementService) TypeMovementGetAll(typeMovement string) ([]*schemas.TypeMovementResponse, error) {
	return t.TypeMovementRepository.TypeMovementGetAll(typeMovement)
}

func (t *TypeMovementService) TypeMovementUpdate(movementUpdate schemas.TypeMovementUpdate) (error) {
	return t.TypeMovementRepository.TypeMovementUpdate(movementUpdate)
}