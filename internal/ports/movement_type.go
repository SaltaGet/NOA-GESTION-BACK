package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"


type MovementTypeService interface {
	MovementTypeGetByID(id string) (movementType *schemas.MovementType, err error)
	MovementTypeGetAll(isIncome bool) (movementTypes *[]schemas.MovementType, err error)
	MovementTypeCreate(movementType *schemas.MovementTypeCreate) (id string, err error)
	MovementTypeUpdate(movementTypeUpdate *schemas.MovementTypeUpdate) (err error)
	MovementTypeDelete(id string) (err error)
}

type MovementTypeRepository interface {
	MovementTypeGetByID(id string) (movementType *schemas.MovementType, err error)
	MovementTypeGetAll(isIncome bool) (movementTypes *[]schemas.MovementType, err error)
	MovementTypeCreate(movementType *schemas.MovementTypeCreate) (id string, err error)
	MovementTypeUpdate(movementTypeUpdate *schemas.MovementTypeUpdate) (err error)
	MovementTypeDelete(id string) (err error)
}
