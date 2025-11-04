package services

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (t *TenantService) TenantGetByID(id int64) (*models.Tenant, error) {
	tenant, err := t.TenantRepository.TenantGetByID(id)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

func (t *TenantService) TenantGetByIdentifier(identifier string) (*models.Tenant, error) {
	tenant, err := t.TenantRepository.TenantGetByIdentifier(identifier)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

func (t *TenantService) TenantGetAll(userID int64) (*[]schemas.TenantResponse, error) {
	tenants, err := t.TenantRepository.TenantGetAll(userID)
	if err != nil {
		return nil, err
	}
	return tenants, nil
}

func (t *TenantService) TenantGetConections() (*[]string, error) {
	conections, err := t.TenantRepository.TenantGetConections()
	if err != nil {
		return nil, err
	}
	return conections, nil
}

func (t *TenantService) TenantCreateByUserID(tenantCreate *schemas.TenantCreate, userID int64) (int64, error) {
	id, err := t.TenantRepository.TenantCreateByUserID(tenantCreate, userID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (t *TenantService) TenantUserCreate(tenantUserCreate *schemas.TenantUserCreate) (int64, error) {
	id, err := t.TenantRepository.TenantUserCreate(tenantUserCreate)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (t *TenantService) TenantUpdate(userID int64, tenant *schemas.TenantUpdate) error {
	return t.TenantRepository.TenantUpdate(userID, tenant)
}
