package ports

import "github.com/DanielChachagua/GestionCar/pkg/models"

type MovementTypeService interface {
	MovementTypeGetByID(id string) (movementType *models.MovementType, err error)
	MovementTypeGetAll(isIncome bool) (movementTypes *[]models.MovementType, err error)
	MovementTypeCreate(movementType *models.MovementTypeCreate) (id string, err error)
	MovementTypeUpdate(movementTypeUpdate *models.MovementTypeUpdate) (err error)
	MovementTypeDelete(id string) (err error)
}

type MovementTypeRepository interface {
	MovementTypeGetByID(id string) (movementType *models.MovementType, err error)
	MovementTypeGetAll(isIncome bool) (movementTypes *[]models.MovementType, err error)
	MovementTypeCreate(movementType *models.MovementTypeCreate) (id string, err error)
	MovementTypeUpdate(movementTypeUpdate *models.MovementTypeUpdate) (err error)
	MovementTypeDelete(id string) (err error)
}
