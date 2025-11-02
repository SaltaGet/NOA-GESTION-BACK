package ports

import "github.com/DanielChachagua/GestionCar/pkg/models"

type TenantRepository interface {
	TenantGetByID(id string) (tenant *models.Tenant, err error) 
	TenantGetByIdentifier(identifier string) (tenant *models.Tenant, err error) 
	TenantGetAll(userID string) (tenants *[]models.TenantResponse, err error)
	TenantGetConections() (conections *[]string, err error)
	TenantCreateByUserID(tenantCreate *models.TenantCreate, userID string) (id string, err error)
	TenantUserCreate(tenantUserCreate *models.TenantUserCreate) (id string, err error)
	TenantUpdate(userID string, tenant *models.TenantUpdate) (err error)
}

type TenantService interface {
	TenantGetByID(tenantID string) (tenant *models.Tenant, err error) 
	TenantGetByIdentifier(identifier string) (tenant *models.Tenant, err error) 
	TenantGetAll(userID string) (tenants *[]models.TenantResponse, err error)
	TenantGetConections() (conections *[]string, err error)
	TenantCreateByUserID(tenantCreate *models.TenantCreate, userID string) (id string, err error)
	TenantUserCreate(tenantUserCreate *models.TenantUserCreate) (id string, err error)
	TenantUpdate(userID string, tenant *models.TenantUpdate) (err error)
}