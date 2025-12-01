package tenant_cache

import (
	"sync"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"gorm.io/gorm"
)

var tenants sync.Map // map[int64]*TenantContainer

// // func GetTenantContainer(db *gorm.DB, tenantID int64) *dependencies.TenantContainer {
// // 	if val, ok := tenants.Load(tenantID); ok {
// // 		tc := val.(*dependencies.TenantContainer)

// // 		sqlDB, err := tc.DB.DB()
// // 		if err == nil && sqlDB.Ping() == nil {
// // 			return tc
// // 		}

// // 		tenants.Delete(tenantID)
// // 	}

// // 	newTC := dependencies.NewTenantContainer(db)
// // 	actual, loaded := tenants.LoadOrStore(tenantID, newTC)

// // 	if loaded {
// // 		return actual.(*dependencies.TenantContainer)
// // 	}

// // 	return newTC
// // }

func GetTenantContainer(db *gorm.DB, tenantID int64) *dependencies.TenantContainer {
	if val, ok := tenants.Load(tenantID); ok {
		tc := val.(*dependencies.TenantContainer)
		return tc
	}

	newTC := dependencies.NewTenantContainer(db)
	actual, loaded := tenants.LoadOrStore(tenantID, newTC)

	if loaded {
		return actual.(*dependencies.TenantContainer)
	}

	return newTC
}


func GetContainerTenantsCache() *sync.Map {
	return &tenants
}