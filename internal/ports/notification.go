package ports

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type NotificationRepository interface {
	NotificationStock(tenantID int64) ([]*models.Product, error)
}

type NotificationService interface {
	NotificationStock(tenantID int64) ([]*schemas.ProductSimpleResponse, error)
}
