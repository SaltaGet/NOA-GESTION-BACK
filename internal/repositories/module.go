package repositories

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (r *MainRepository) ModuleGetByTenantID(tenantID int64) ([]schemas.ModuleResponseDTO, error) {
	var modules []models.Module

	err := r.DB.
		Joins("JOIN tenant_modules tm ON tm.module_id = modules.id").
		Where("tm.tenant_id = ?", tenantID).
		Where("tm.deleted_at IS NULL").
		Preload("Tenants", "tenant_id = ?", tenantID).
		Find(&modules).Error

	if err != nil {
		return nil, err
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
		})
	}

	return modulesResponse, nil
}
