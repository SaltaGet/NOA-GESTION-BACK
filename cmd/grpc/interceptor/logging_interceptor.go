// ============================================
// internal/grpc/interceptor/logging.go
// ============================================
package interceptor

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor registra todas las peticiones gRPC con métricas de tiempo
func LoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	// Obtener IP del cliente
	clientIP := getClientIP(ctx)

	// Obtener tenant si existe en metadatos
	tenantID := getTenantFromMetadata(ctx)

	// Ejecutar el handler
	resp, err := handler(ctx, req)

	// Calcular duración
	duration := time.Since(start)
	durationMs := float64(duration.Nanoseconds()) / 1_000_000.0

	// Obtener código de estado
	statusCode := codes.OK
	if err != nil {
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}
	}

	// Crear log con campos estructurados
	logEvent := log.Info()

	// Si hay error, cambiar a nivel error
	if err != nil {
		logEvent = log.Error().Err(err)
	}

	logEvent.
		Str("ip", clientIP).
		Str("method", info.FullMethod).
		Str("code", statusCode.String()).
		Int("code_number", int(statusCode)).
		Float64("duration_ms", durationMs).
		Str("tenant", tenantID)

	// Agregar información adicional dependiendo del tiempo de respuesta
	if durationMs > 1000 {
		logEvent.Msg("⚠️ gRPC request SLOW")
	} else if err != nil {
		logEvent.Msg("❌ gRPC request failed")
	} else {
		logEvent.Msg("✅ gRPC request completed")
	}

	return resp, err
}

// getClientIP extrae la IP del cliente del contexto
func getClientIP(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "unknown"
	}
	return p.Addr.String()
}

// getTenantFromMetadata extrae el tenant_id de los metadatos
func getTenantFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "-"
	}

	tenantHeaders := md.Get("x-tenant-identifier")
	if len(tenantHeaders) == 0 {
		return "-"
	}

	return tenantHeaders[0]
}

// ============================================
// Versión con métricas adicionales (opcional)
// ============================================

// LoggingInterceptorDetailed incluye más información de debugging
func LoggingInterceptorDetailed(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	clientIP := getClientIP(ctx)
	tenantID := getTenantFromMetadata(ctx)
	userAgent := getUserAgent(ctx)

	// Log de inicio (opcional, útil para debugging)
	log.Debug().
		Str("ip", clientIP).
		Str("method", info.FullMethod).
		Str("tenant", tenantID).
		Msg("→ gRPC request started")

	resp, err := handler(ctx, req)

	duration := time.Since(start)
	durationMs := float64(duration.Nanoseconds()) / 1_000_000.0

	statusCode := codes.OK
	errorMsg := ""
	if err != nil {
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
			errorMsg = st.Message()
		}
	}

	logEvent := log.Info()
	if err != nil {
		logEvent = log.Error()
	}

	logEvent.
		Str("ip", clientIP).
		Str("method", info.FullMethod).
		Str("code", statusCode.String()).
		Int("code_number", int(statusCode)).
		Float64("duration_ms", durationMs).
		Str("tenant", tenantID).
		Str("user_agent", userAgent)

	if errorMsg != "" {
		logEvent.Str("error_msg", errorMsg)
	}

	// Categorizar por performance
	switch {
	case durationMs > 5000:
		logEvent.Msg("🔴 gRPC request CRITICAL SLOW")
	case durationMs > 1000:
		logEvent.Msg("🟡 gRPC request SLOW")
	case err != nil:
		logEvent.Msg("❌ gRPC request failed")
	default:
		logEvent.Msg("✅ gRPC request completed")
	}

	return resp, err
}

// getUserAgent extrae el user-agent de los metadatos
func getUserAgent(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "-"
	}

	userAgentHeaders := md.Get("user-agent")
	if len(userAgentHeaders) == 0 {
		return "-"
	}

	return userAgentHeaders[0]
}

// ============================================
// main.go - Configuración con logging
// ============================================

// go func() {
// 	log.Info().Msg("🚀 Iniciando servidor gRPC en :50051...")

// 	lis, err := net.Listen("tcp", ":50051")
// 	if err != nil {
// 		log.Fatal().Msgf("❌ Falló al escuchar gRPC: %v", err)
// 	}

// 	// Configurar interceptores en orden
// 	grpcServer := grpc.NewServer(
// 		grpc.ChainUnaryInterceptor(
// 			interceptor.LoggingInterceptor,          // 1º: Logging (al inicio para medir todo)
// 			interceptor.AuthInterceptor,             // 2º: Autenticación
// 			interceptor.MultiTenantInterceptor(dep), // 3º: MultiTenant
// 		),
// 	)

// 	// Registrar servicios
// 	productServer := &server.GrpcProductServer{}
// 	pb.RegisterProductServiceServer(grpcServer, productServer)

// 	tenantServer := &server.GrpcTenantServer{
// 		MainContainer: dep,
// 	}
// 	pb.RegisterTenantServiceServer(grpcServer, tenantServer)

// 	log.Info().Msg("✅ Servidor gRPC listo en :50051")

// 	if err := grpcServer.Serve(lis); err != nil {
// 		log.Fatal().Msgf("❌ Falló al servir gRPC: %v", err)
// 	}
// }()

// ============================================
// Ejemplo de logs generados
// ============================================

/*
VERSIÓN SIMPLE (LoggingInterceptor):

{"level":"info","ip":"127.0.0.1:54321","method":"/product.ProductService/ListProducts","code":"OK","code_number":0,"duration_ms":45.23,"tenant":"company-abc","message":"✅ gRPC request completed"}

{"level":"error","error":"Tenant inválido","ip":"127.0.0.1:54322","method":"/product.ProductService/GetProduct","code":"Unauthenticated","code_number":16,"duration_ms":12.45,"tenant":"invalid-tenant","message":"❌ gRPC request failed"}

{"level":"info","ip":"127.0.0.1:54323","method":"/product.ProductService/ListProducts","code":"OK","code_number":0,"duration_ms":1250.67,"tenant":"company-xyz","message":"⚠️ gRPC request SLOW"}


VERSIÓN DETALLADA (LoggingInterceptorDetailed):

{"level":"debug","ip":"127.0.0.1:54321","method":"/product.ProductService/ListProducts","tenant":"company-abc","message":"→ gRPC request started"}

{"level":"info","ip":"127.0.0.1:54321","method":"/product.ProductService/ListProducts","code":"OK","code_number":0,"duration_ms":45.23,"tenant":"company-abc","user_agent":"grpc-go/1.60.0","message":"✅ gRPC request completed"}
*/