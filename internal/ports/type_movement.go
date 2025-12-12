package ports

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)


type TypeMovementRepository interface {
	TypeMovementGetAll(typeMovement string) ([]*schemas.TypeMovementResponse, error)
	TypeMovementCreate(memberID int64, movementCreate schemas.TypeMovementCreate) (error)
	TypeMovementUpdate(memberID int64, movementUpdate schemas.TypeMovementUpdate) (error)
}

type TypeMovementService interface {
	TypeMovementGetAll(typeMovement string) ([]*schemas.TypeMovementResponse, error)
	TypeMovementCreate(memberID int64, movementCreate schemas.TypeMovementCreate) (error)
	TypeMovementUpdate(memberID int64, movementUpdate schemas.TypeMovementUpdate) (error)
}



