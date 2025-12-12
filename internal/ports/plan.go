package ports

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type PlanRepository interface {
	PlanCreate(adminID int64, plan *schemas.PlanCreate) (int64, error)
	PlanUpdate(adminID int64, plan *schemas.PlanUpdate) error
	PlanGetByID(planID int64) (*schemas.PlanResponse, error)
	PlanGetAll() ([]*schemas.PlanResponseDTO, error)
}

type PlanService interface {
	PlanCreate(adminID int64, plan *schemas.PlanCreate) (int64, error)
	PlanUpdate(adminID int64, plan *schemas.PlanUpdate) error
	PlanGetByID(planID int64) (*schemas.PlanResponse, error)
	PlanGetAll() ([]*schemas.PlanResponseDTO, error)
}
