package services

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
)

func (t *TenantService) TenantGetByID(id string) (*models.Tenant, error) {
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

func (t *TenantService) TenantGetAll(userID string) (*[]models.TenantResponse, error) {
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

func (t *TenantService) TenantCreateByUserID(tenantCreate *models.TenantCreate, userID string) (string, error) {
	id, err := t.TenantRepository.TenantCreateByUserID(tenantCreate, userID)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (t *TenantService) TenantUserCreate(tenantUserCreate *models.TenantUserCreate) (string, error) {
	id, err := t.TenantRepository.TenantUserCreate(tenantUserCreate)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (t *TenantService) TenantUpdate(userID string, tenant *models.TenantUpdate) error {
	return t.TenantRepository.TenantUpdate(userID, tenant)
}
