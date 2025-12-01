package repositories

import (
	"errors"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"gorm.io/gorm"
)

func (r *MainRepository) NotificationStock(tenantID int64) ([]*models.Product, error) {
	var tenant models.Tenant
	if err := r.DB.Select("connection").Where("id = ?", tenantID).First(tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "tenant no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "error al obtener el tenant", err)
	}

	dbTenant, err := database.GetTenantDB(tenant.Connection, tenantID)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "error al conectar con la base de datos del tenant", err)
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
