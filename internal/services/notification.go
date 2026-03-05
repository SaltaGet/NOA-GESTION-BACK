package services

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (s *NotificationService) NotificationStock(tenantID int64) ([]*schemas.ProductSimpleResponse, error) {
	return s.NotificationRepository.NotificationStock(tenantID)
}
