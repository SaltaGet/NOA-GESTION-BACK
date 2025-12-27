package ports

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)


type TenantRepository interface {
	TenantGetByID(id int64) (tenant *models.Tenant, err error) 
	TenantGetByIdentifier(identifier string) (tenant *models.Tenant, err error) 
	TenantGetAll() (tenants *[]schemas.TenantResponse, err error)
	TenantGetConnectionByIdentifier(tenantIdentifier string) (connection *models.Tenant, err error)
	TenantGetConections() (conections []*models.Tenant, err error)
	TenantCreateByUserID(adminID int64, tenantCreate *schemas.TenantCreate, userID int64) (id int64, err error)
	TenantUserCreate(adminID int64, tenantUserCreate *schemas.TenantUserCreate) (id int64, err error)
	TenantUpdate(adminID int64, userID int64, tenant *schemas.TenantUpdate) (err error)
	TenantUpdateExpiration(adminID int64, tenantUpdateExpiration *schemas.TenantUpdateExpiration) (err error)
}

type TenantService interface {
	TenantGetByID(tenantID int64) (tenant *models.Tenant, err error) 
	TenantGetByIdentifier(identifier string) (tenant *models.Tenant, err error) 
	TenantGetAll() (tenants *[]schemas.TenantResponse, err error)
	TenantGetConnectionByIdentifier(tenantIdentifier string) (connection *models.Tenant, err error)
	TenantGetConections() (conections []*models.Tenant, err error)
	TenantCreateByUserID(adminID int64, tenantCreate *schemas.TenantCreate, userID int64) (id int64, err error)
	TenantUserCreate(adminID int64, tenantUserCreate *schemas.TenantUserCreate) (id int64, err error)
	TenantUpdate(adminID int64, userID int64, tenant *schemas.TenantUpdate) (err error)
	TenantUpdateExpiration(adminID int64, tenantUpdateExpiration *schemas.TenantUpdateExpiration) (err error)
}