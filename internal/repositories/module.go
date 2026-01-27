package repositories

import (
	"errors"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *MainRepository) ModuleGet(id int64) (*schemas.ModuleResponse, error) {
	var module models.Module
	if err := r.DB.First(&module, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Modulo no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error al buscar modulo", err)
	}

	var moduleResponse schemas.ModuleResponse
	copier.Copy(&moduleResponse, &module)

	return &moduleResponse, nil
}

func (r *MainRepository) ModuleGetAll() ([]schemas.ModuleResponse, error) {
	var modules []models.Module
	if err := r.DB.Find(&modules).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "Error al obtener modulos", err)
	}

	var modulesResponse []schemas.ModuleResponse
	copier.Copy(&modulesResponse, &modules)

	return modulesResponse, nil
}

func (r *MainRepository) ModuleCreate(moduleCreate *schemas.ModuleCreate) (int64, error) {
	newModule := &models.Module{
		Name:                   moduleCreate.Name,
		AmountImagesPerProduct: moduleCreate.AmountImagesPerProduct,
		PriceMonthly:           moduleCreate.PriceMonthly,
		PriceYearly:            moduleCreate.PriceYearly,
		Description:            *moduleCreate.Description,
		Features:               moduleCreate.Features,
	}

	err := r.DB.Create(&newModule).Error
	if err != nil {
		if schemas.IsDuplicateError(err) {
			return 0, schemas.ErrorResponse(409, "El modulo "+newModule.Name+" ya existe", err)
		}
		return 0, schemas.ErrorResponse(500, "Error al crear modulo", err)
	}

	return newModule.ID, nil
}

func (r *MainRepository) ModuleUpdate(moduleUpdate *schemas.ModuleUpdate) error {
	if err := r.DB.Model(&models.Module{}).Where("id = ?", moduleUpdate.ID).Updates(moduleUpdate).Error; err != nil {
		if schemas.IsDuplicateError(err) {
			return schemas.ErrorResponse(409, "El modulo "+moduleUpdate.Name+" ya existe", err)
		}
		return schemas.ErrorResponse(500, "Error al actualizar modulo", err)
	}

	return nil
}

func (r *MainRepository) ModuleDelete(id int64) error {
	// return r.DB.Delete(&models.Module{}, id).Error
	return nil
}

// func (r *MainRepository) ModuleAddTenant(moduleAddTenant *schemas.ModuleAddTenant) error {
// 	loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
// 	exp, err := time.ParseInLocation("2006-01-02", moduleAddTenant.Expiration, loc)
// 	if err != nil {
// 		return schemas.ErrorResponse(422, "Formato de fecha inválido, debe ser YYYY-MM-DD", err)
// 	}

// 	newModuleAdd := &models.TenantModule{
// 		ModuleID:   moduleAddTenant.ModuleID,
// 		TenantID:   moduleAddTenant.TenantID,
// 	}

// 	if err := r.DB.
// 		Where("tenant_id = ? AND module_id = ?", moduleAddTenant.TenantID, moduleAddTenant.ModuleID).
// 		FirstOrCreate(&newModuleAdd).Error; err != nil {
// 		return schemas.ErrorResponse(500, "Error interno al guardar el modulo para el tenant", errors.New("error interno al guardar el modulo para el tenant"))
// 	}

// 	newModuleAdd.Expiration = &exp

// 	if err := r.DB.Save(&newModuleAdd).Error; err != nil {
// 		return schemas.ErrorResponse(500, "Error interno al guardar el modulo para el tenant", errors.New("error interno al guardar el modulo para el tenant"))
// 	}

//		return nil
//	}
func (r *MainRepository) ModuleAddTenant(moduleAddTenant *schemas.ModuleAddTenant) error {
	if err := r.DB.First(&models.Tenant{}, moduleAddTenant.TenantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "Tenant no encontrado", err)
		}
		return schemas.ErrorResponse(500, "Error al buscar tenant", err)
	}

	if err := r.DB.Where("id = ?", moduleAddTenant.ModuleID).First(&models.Module{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "El Módulo especificado no existe", err)
		}
		return schemas.ErrorResponse(500, "Error interno al validar el Módulo", err)
	}

	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		// Fallback a UTC o manejar error si la locación falla
		loc = time.UTC
	}

	exp, err := time.ParseInLocation("2006-01-02", moduleAddTenant.Expiration, loc)
	if err != nil {
		return schemas.ErrorResponse(422, "Formato de fecha inválido, debe ser YYYY-MM-DD", err)
	}

	newModuleAdd := models.TenantModule{
		ModuleID:   moduleAddTenant.ModuleID,
		TenantID:   moduleAddTenant.TenantID,
		Expiration: &exp,
	}

	// Usamos Upsert: Si coinciden TenantID y ModuleID, actualiza la Expiration
	err = r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tenant_id"}, {Name: "module_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"expiration"}),
	}).Create(&newModuleAdd).Error

	if err != nil {
		return schemas.ErrorResponse(500, "Error al procesar el módulo", err)
	}

	return nil
}

func (r *MainRepository) ModuleGetByTenantID(tenantID int64) ([]schemas.ModuleResponseDTO, error) {
	var modules []models.Module

	err := r.DB.
		Joins("JOIN tenant_modules tm ON tm.module_id = modules.id").
		Where("tm.tenant_id = ?", tenantID).
		Where("tm.deleted_at IS NULL").
		Preload("Tenants", "tenant_id = ?", tenantID).
		Find(&modules).Error

	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al procesar modulo", err)
	}

	modulesResponse := make([]schemas.ModuleResponseDTO, 0, len(modules))
	for _, module := range modules {
		if len(module.Tenants) == 0 {
			continue
		}

		modulesResponse = append(modulesResponse, schemas.ModuleResponseDTO{
			ID:                     module.ID,
			Name:                   module.Name,
			AmountImagesPerProduct: module.AmountImagesPerProduct,
			Expiration:             module.Tenants[0].Expiration,
			AcceptTerms: module.Tenants[0].AcceptedTerms,
		})
	}

	return modulesResponse, nil
}
