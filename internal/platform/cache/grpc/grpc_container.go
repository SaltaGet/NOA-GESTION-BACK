package grpc_cache

import (
	// "context"
	// "context"
	"context"
	"sync"

	// "github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	// "github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"gorm.io/gorm"
)

var grpcContainers sync.Map // map[int64]*TenantContainer

func GetGRPCContainer(db *gorm.DB, tenantID int64) *dependencies.GrpcContainer {
	if val, ok := grpcContainers.Load(tenantID); ok {
		tc := val.(*dependencies.GrpcContainer)
		return tc
	}

	newTC := dependencies.NewGrpcContainer(db)
	actual, loaded := grpcContainers.LoadOrStore(tenantID, newTC)

	if loaded {
		return actual.(*dependencies.GrpcContainer)
	}

	return newTC
}


func GetContainerGrpcCache() *sync.Map {
	return &grpcContainers
}



// type DBResolver struct {
//     connections map[string]*gorm.DB
//     mu          sync.RWMutex
// }

// func (r *DBResolver) GetConnection(tenantIdentifier string, deps *dependencies.MainContainer) (*gorm.DB, *models.Tenant, error) {
//     r.mu.RLock()
//     db, exists := r.connections[tenantIdentifier]
//     r.mu.RUnlock()

//     if exists {
//         return db, nil, nil
//     }

// 		tenant, err := deps.TenantController.TenantService.TenantGetConnectionByIdentifier(tenantIdentifier)
// 		if err != nil {
// 			return nil, nil, err
// 		}

// 		db, err = database.GetTenantDB(tenant.Connection, tenant.ID)
// 		if err != nil {
// 			return nil, nil, err
// 		}

// 		r.mu.Lock()
// 		r.connections[tenantIdentifier] = db
// 		r.mu.Unlock()

//     // Aquí podrías inicializar una nueva conexión si no existe (con un Lock de escritura)
//     return db, tenant, nil
// }

// type contextKey string


// // Helper para extraer el contenedor de forma segura
func GetGrpcContainerFromContext(ctx context.Context) *dependencies.GrpcContainer {
    val := ctx.Value("deps_grpc")
    if val == nil {
        return nil 
    }

    container, ok := val.(*dependencies.GrpcContainer)
    if !ok {
        return nil
    }
    return container
}