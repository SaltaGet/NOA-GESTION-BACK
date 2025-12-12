package repositories

import (
	"errors"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func (r *MainRepository) PlanCreate(adminID int64, planCreate *schemas.PlanCreate) (int64, error) {
	plan := &models.Plan{
		Name:            planCreate.Name,
		PriceMounthly:   planCreate.PriceMounthly,
		PriceYearly:     planCreate.PriceYearly,
		Description:     planCreate.Description,
		Features:        planCreate.Features,
		AmountPointSale: planCreate.AmountPointSale,
		AmountMember:    planCreate.AmountMember,
	}

	if err := r.DB.Create(plan).Error; err != nil {
		if schemas.IsDuplicateError(err) {
			return 0, schemas.ErrorResponse(409, "El plan "+plan.Name+" ya existe", err)
		}
		return 0, schemas.ErrorResponse(500, "Error al crear el plan", err)
	}

	go database.SaveAuditAdminAsync(r.DB, models.AuditLogAdmin{
		AdminID: adminID,
		Method:  "create",
		Path:    "plan",
	}, nil, plan)

	return plan.ID, nil
}

func (r *MainRepository) PlanUpdate(adminID int64, planUpdate *schemas.PlanUpdate) error {
	var plan, oldPlan models.Plan
	if err := r.DB.First(&plan, planUpdate.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "El plan no encopntrado", err)
		}
		return schemas.ErrorResponse(500, "Error al buscar el plan", err)
	}

	oldPlan = plan

	plan.Name = planUpdate.Name
	plan.PriceMounthly = planUpdate.PriceMounthly
	plan.PriceYearly = planUpdate.PriceYearly
	plan.Description = planUpdate.Description
	plan.Features = planUpdate.Features
	plan.AmountPointSale = planUpdate.AmountPointSale
	plan.AmountMember = planUpdate.AmountMember

	if err := r.DB.Save(&plan).Error; err != nil {
		if schemas.IsDuplicateError(err) {
			return schemas.ErrorResponse(409, "El plan "+plan.Name+" ya existe", err)
		}
		return schemas.ErrorResponse(500, "Error al actualizar el plan", err)
	}

	go database.SaveAuditAdminAsync(r.DB, models.AuditLogAdmin{
		AdminID: adminID,
		Method:  "update",
		Path:    "plan",
	}, oldPlan, plan)

	return nil
}

func (r *MainRepository) PlanGetAll() ([]*schemas.PlanResponseDTO, error) {
	var plans []models.Plan
	if err := r.DB.Find(&plans).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "Error al obtener los planes", err)
	}
	var plansResponse []*schemas.PlanResponseDTO
	copier.Copy(&plansResponse, &plans)

	return plansResponse, nil
}

func (r *MainRepository) PlanGetByID(id int64) (*schemas.PlanResponse, error) {
	var plan models.Plan
	if err := r.DB.Preload("Tenants").First(&plan, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "El plan no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error al obtener el plan", err)
	}

	var planResponse schemas.PlanResponse
	copier.Copy(&planResponse, &plan)

	return &planResponse, nil
}
