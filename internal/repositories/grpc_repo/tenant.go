package grpc_repo

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
)

func (r *GrpcMainRepository) ListTenants() ([]models.Tenant, error) {
	var tenants []models.Tenant
	if err := r.DB.Find(&tenants).Error; err != nil {
		return nil, err
	}

	return tenants, nil
}