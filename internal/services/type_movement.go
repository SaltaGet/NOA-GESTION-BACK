package services

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

func (t *TypeMovementService) TypeMovementCreate(memberID int64, movementCreate schemas.TypeMovementCreate) (error) {
	return t.TypeMovementRepository.TypeMovementCreate(memberID, movementCreate)
}

func (t *TypeMovementService) TypeMovementGetAll(typeMovement string) ([]*schemas.TypeMovementResponse, error) {
	return t.TypeMovementRepository.TypeMovementGetAll(typeMovement)
}

func (t *TypeMovementService) TypeMovementUpdate(memberID int64, movementUpdate schemas.TypeMovementUpdate) (error) {
	return t.TypeMovementRepository.TypeMovementUpdate(memberID, movementUpdate)
}