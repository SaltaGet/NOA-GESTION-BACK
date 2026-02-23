package tenant_cache

import (
	"fmt"
	"sync"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"gorm.io/gorm"
	"golang.org/x/sync/singleflight"
)

var (
	tenants sync.Map // map[int64]*TenantContainer
	sfGroup singleflight.Group
)

func GetTenantContainer(db *gorm.DB, tenantID int64) *dependencies.TenantContainer {
	// 1. Caso rápido optimista: ya existe
	if val, ok := tenants.Load(tenantID); ok {
		return val.(*dependencies.TenantContainer)
	}

	// 2. Si no existe, usamos singleflight para asegurar que solo una goroutine (request)
	// ejecute NewTenantContainer para este tenant específico a la vez.
	key := fmt.Sprintf("tenant_container_%d", tenantID)
	
	val, err, _ := sfGroup.Do(key, func() (interface{}, error) {
		// Dobre chequeo recomendado (aunque LoadOrStore hace algo similar, esto ahorra inicializar struct basura)
		if v, ok := tenants.Load(tenantID); ok {
			return v, nil
		}
		
		newTC := dependencies.NewTenantContainer(db)
		tenants.Store(tenantID, newTC)
		return newTC, nil
	})

	if err != nil {
		// Si hubiera error, lo manejaríamos, pero NewTenantContainer no devuelve errores actualmente
		return nil
	}

	return val.(*dependencies.TenantContainer)
}

func GetContainerTenantsCache() *sync.Map {
	return &tenants
}


// package tenant_cache

// import (
// 	"sync"

// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
// 	"gorm.io/gorm"
// )

// var tenants sync.Map // map[int64]*TenantContainer

// func GetTenantContainer(db *gorm.DB, tenantID int64) *dependencies.TenantContainer {
// 	if val, ok := tenants.Load(tenantID); ok {
// 		tc := val.(*dependencies.TenantContainer)
// 		return tc
// 	}

// 	newTC := dependencies.NewTenantContainer(db)
// 	actual, loaded := tenants.LoadOrStore(tenantID, newTC)

// 	if loaded {
// 		return actual.(*dependencies.TenantContainer)
// 	}

// 	return newTC
// }


// func GetContainerTenantsCache() *sync.Map {
// 	return &tenants
// }