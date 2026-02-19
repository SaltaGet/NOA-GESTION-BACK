package repositories

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (r *MainRepository) NotificationStock(tenantID int64) ([]*models.Product, error) {
	var tenant models.Tenant
	if err := r.DB.Select("connection").Where("id = ?", tenantID).First(tenant).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Tenant", schemas.Read)
	}

	dbTenant, err := database.GetTenantDB(tenant.Connection, tenantID)
	if err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Tenant", schemas.Read)
	}

	var products []*models.Product
	err = dbTenant.Model(&models.Product{}).
		Joins("JOIN stock_deposits sd ON sd.product_id = products.id").
		Where("sd.stock <= products.min_amount").
		Where("products.notifier = ?", true).
		Preload("StockDeposit").
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}
