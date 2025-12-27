package services

import (
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
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

func (t *TenantService) TenantGetAll() (*[]schemas.TenantResponse, error) {
	tenants, err := t.TenantRepository.TenantGetAll()
	if err != nil {
		return nil, err
	}
	return tenants, nil
}

func (t *TenantService) TenantGetConnectionByIdentifier(tenantIdentifier string) (*models.Tenant, error) {
	return t.TenantRepository.TenantGetConnectionByIdentifier(tenantIdentifier)
}

func (t *TenantService) TenantGetConections() ([]*models.Tenant, error) {
	conections, err := t.TenantRepository.TenantGetConections()
	if err != nil {
		return nil, err
	}
	return conections, nil
}

func (t *TenantService) TenantCreateByUserID(adminID int64, tenantCreate *schemas.TenantCreate, userID int64) (int64, error) {
	id, err := t.TenantRepository.TenantCreateByUserID(adminID, tenantCreate, userID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (t *TenantService) TenantUserCreate(adminID int64, tenantUserCreate *schemas.TenantUserCreate) (int64, error) {
	oldPassword := tenantUserCreate.UserCreate.Password
	id, err := t.TenantRepository.TenantUserCreate(adminID, tenantUserCreate)
	if err != nil {
		return 0, err
	}

	username := fmt.Sprintf("%s@%s", tenantUserCreate.UserCreate.Username, tenantUserCreate.TenantCreate.Identifier)
	go t.EmailService.SendEmail(tenantUserCreate.UserCreate.Email, "Bienvenido a NOA-GESTION", utils.WelcomeUser(username, oldPassword))


	return id, nil
}

func (t *TenantService) TenantUpdate(adminID int64, userID int64, tenant *schemas.TenantUpdate) error {
	return t.TenantRepository.TenantUpdate(adminID, userID, tenant)
}

func (t *TenantService) TenantUpdateExpiration(adminID int64, tenantUpdateExpiration *schemas.TenantUpdateExpiration) error {
	return t.TenantRepository.TenantUpdateExpiration(adminID, tenantUpdateExpiration)
}