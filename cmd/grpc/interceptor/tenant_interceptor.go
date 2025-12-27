package interceptor

import (
	"context"
	"strings"

	grpc_cache "github.com/SaltaGet/NOA-GESTION-BACK/internal/cache/grpc"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
    "github.com/rs/zerolog/log"
)

func MultiTenantInterceptor(
    // resolver *grpc_cache.DBResolver, 
    deps *dependencies.MainContainer) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
        if isExcludedFromMultiTenant(info.FullMethod) {
			log.Debug().Str("method", info.FullMethod).Msg("⏭️ Método excluido de multitenant")
			return handler(ctx, req)
		}
        
        md, _ := metadata.FromIncomingContext(ctx)
        tenantIdentifier := md.Get("x-tenant-identifier")[0]

        // 1. Obtener la conexión específica
        // db, tenant, err := deps.GetConnection(tenantIdentifier, deps)
        tenant, err := deps.TenantController.TenantService.TenantGetConnectionByIdentifier(tenantIdentifier)
        if err != nil {
            return nil, status.Error(codes.Unauthenticated, "Tenant inválido")
        }

        db, err := database.GetTenantDB(tenant.Connection, tenant.ID)
        if err != nil {
            return nil, status.Error(codes.Internal, "Error al obtener la conexión de la base de datos del tenant")
        }

        depsGrpc := grpc_cache.GetGRPCContainer(db, tenant.ID)

        // 2. Inyectar la conexión en el contexto de ESTA petición
        // newCtx := context.WithValue(ctx, "db_conn", db)
        newCtx := context.WithValue(ctx, "deps_grpc", depsGrpc)

        // 3. Continuar la ejecución
        return handler(newCtx, req)
    }
}

func isExcludedFromMultiTenant(fullMethod string) bool {
	// Lista de servicios completos que NO requieren multitenant
	excludedServices := []string{
		"/tenant.TenantService/",  // Todo el servicio de gestión de tenants
		"/auth.AuthService/",      // Servicio de autenticación (si lo tienes)
		"/admin.AdminService/",    // Servicio de administración (si lo tienes)
	}

	for _, excluded := range excludedServices {
		if strings.HasPrefix(fullMethod, excluded) {
			return true
		}
	}

	// Lista de métodos específicos excluidos (opcional)
	excludedMethods := []string{
		// Ejemplo: "/product.ProductService/GetPublicProducts",
	}

	for _, excluded := range excludedMethods {
		if fullMethod == excluded {
			return true
		}
	}

	return false
}