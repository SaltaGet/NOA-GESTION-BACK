package ports

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type PlanRepository interface {
	PlanCreate(plan *schemas.PlanCreate) (int64, error)
	PlanUpdate(plan *schemas.PlanUpdate) error
	PlanGetByID(planID int64) (*schemas.PlanResponse, error)
	PlanGetAll() ([]*schemas.PlanResponseDTO, error)
}

type PlanService interface {
	PlanCreate(plan *schemas.PlanCreate) (int64, error)
	PlanUpdate(plan *schemas.PlanUpdate) error
	PlanGetByID(planID int64) (*schemas.PlanResponse, error)
	PlanGetAll() ([]*schemas.PlanResponseDTO, error)
}
