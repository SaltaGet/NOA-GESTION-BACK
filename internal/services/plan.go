package services

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

func (s *PlanService) PlanCreate(plan *schemas.PlanCreate) (int64, error) {
	return s.PlanRepository.PlanCreate(plan)
}

func (s *PlanService) PlanGetAll() ([]*schemas.PlanResponseDTO, error) {
	return s.PlanRepository.PlanGetAll()
}

func (s *PlanService) PlanUpdate(plan *schemas.PlanUpdate) (error) {
	return s.PlanRepository.PlanUpdate(plan)
}

func (s *PlanService) PlanGetByID(id int64) (*schemas.PlanResponse, error) {
	return s.PlanRepository.PlanGetByID(id)
}