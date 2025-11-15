package ports

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)


type TypeMovementRepository interface {
	TypeMovementGetAll(typeMovement string) ([]*schemas.TypeMovementResponse, error)
	TypeMovementCreate(movementCreate schemas.TypeMovementCreate) (error)
	TypeMovementUpdate(movementUpdate schemas.TypeMovementUpdate) (error)
}

type TypeMovementService interface {
	TypeMovementGetAll(typeMovement string) ([]*schemas.TypeMovementResponse, error)
	TypeMovementCreate(movementCreate schemas.TypeMovementCreate) (error)
	TypeMovementUpdate(movementUpdate schemas.TypeMovementUpdate) (error)
}



