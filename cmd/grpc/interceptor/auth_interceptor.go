package interceptor

import (
	"context"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no se encontraron metadatos")
	}

	values := md["x-internal-secret"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "api key no proporcionada")
	}

	apiKey := values[0]
	expectedKey := os.Getenv("INTERNAL_SERVICE_KEY")

	// 3. Validar
	if apiKey != expectedKey {
		return nil, status.Errorf(codes.Unauthenticated, "credenciales inválidas")
	}

	// 4. Continuar con la ejecución (como c.Next() en Fiber)
	return handler(ctx, req)
}