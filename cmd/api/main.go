package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	// "os/signal"
	"strconv"
	"strings"

	// "syscall"
	"time"

	// "github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/initial"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/jobs"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/routes"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/grpc/interceptor"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/grpc/server"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/cache"
	grpc_cache "github.com/SaltaGet/NOA-GESTION-BACK/internal/cache/grpc"
	tenant_cache "github.com/SaltaGet/NOA-GESTION-BACK/internal/cache/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"

	// "github.com/SaltaGet/NOA-GESTION-BACK/internal/services/grpc_serv"
	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// "github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/docs"
	"github.com/rs/zerolog/log"
)

// @title						APP NOA Gestion
// @version					1.0
// @description				This is a api to app noa gestion
// @termsOfService				http://swagger.io/terms/
// @securityDefinitions.apikey	CookieAuth
//
//	@in							cookie
//	@name						access_token
//
// @description				Type "Bearer" followed by a space and the JWT token.
func main() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal().Err(err).Msg("Error cargando .env local")
		}
	}

	logging.InitLogging()

	path := os.Getenv("APP_ROOT")
	if path != "" {
		os.MkdirAll(path+"/backups", os.ModePerm)
	} else {
		log.Fatal().Msg("No está indicada la ruta raiz")
	}

	local := os.Getenv("LOCAL")
	if local == "true" {
		if err := jobs.GenerateSwagger(); err != nil {
			log.Fatal().Err(err).Msg("Error ejecutando swag init")
		}
	}

	// 🔥 Inicializar Redis (opcional, falla gracefully)
	if err := cache.InitRedis(); err != nil {
		log.Warn().Err(err).Msg("⚠️ Advertencia: Redis no disponible")
	}

	// Inicializar cache de tenant DBs
	cacheSize := 1000
	if err := database.InitDBCache(cacheSize); err != nil {
		log.Fatal().Err(err).Msg("Error al inicializar cache de DBs")
	}

	emailCfg := initial.InitEmail()

	db, err := database.ConnectDB(emailCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Error al conectar con la base de datos")
	}

	// Context para shutdown graceful
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tenants := tenant_cache.GetContainerTenantsCache()
	grpcCaches := grpc_cache.GetContainerGrpcCache()

	// Iniciar janitor para limpiar conexiones inactivas
	go database.StartDBJanitor(ctx, tenants, grpcCaches)

	// s := sse.NewServer(&sse.Options{
	// 	Headers: map[string]string{
	// 		"Content-Type":      "text/event-stream",
	// 		"Cache-Control":     "no-cache, no-transform",
	// 		"Connection":        "keep-alive",
	// 		"X-Accel-Buffering": "no",
	// 	},
	// })
	// defer s.Shutdown()

	dep := dependencies.NewApplication(db, emailCfg)

	err = jobs.Migrations(dep)
	if err != nil {
		// log.Fatalf("Error al aplicar migraciones: %v", err)
		log.Err(err).Msg("Error al aplicar migraciones")
	}

	app := fiber.New(fiber.Config{
		IdleTimeout:           30 * time.Second,
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
		ErrorHandler:          customErrorHandler,
		ProxyHeader:           "X-Forwarded-For",
		DisableStartupMessage: false,
		StreamRequestBody:     true,
		// Prefork: true,
	})

	app.Use(middleware.BlockAccess())
	app.Use(middleware.LoggingMiddleware)
	app.Use(middleware.ReadOnlyMiddleware())
	app.Use(middleware.InjectionDepends(dep))

	// Rate limiting global (100 req/min por IP)
	app.Use(middleware.RateLimitMiddleware(100, time.Minute))

	maxAge, err := strconv.Atoi(os.Getenv("MAXAGE"))
	if err != nil {
		maxAge = 300
	}

	credentials, err := strconv.ParseBool(os.Getenv("CREDENTIALS"))
	if err != nil {
		credentials = false
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.ReplaceAll(os.Getenv("ORIGIN"), " ", ""),
		AllowMethods:     os.Getenv("METHODS"),
		AllowHeaders:     os.Getenv("HEADERS"),
		AllowCredentials: credentials,
		MaxAge:           maxAge,
	}))

	// app.Use(limiter.New(limiter.Config{
	// 	Max:        250,
	// 	Expiration: 1 * time.Minute,
	// 	KeyGenerator: func(c *fiber.Ctx) string {
	// 		return c.IP()
	// 	},
	// 	LimitReached: func(c *fiber.Ctx) error {
	// 		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
	// 			"error": "Demasiadas peticiones. Intentá más tarde.",
	// 		})
	// 	},
	// }))

	app.Get("/api/swagger/*", swagger.HandlerDefault)
	app.Get("/api/health", healthCheck)
	// app.Get("/cache/stats", cacheStats)
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	routes.SetupRoutes(app, dep)

	depGrpc := dependencies.NewGrpcApplication(db)

	initBackup := os.Getenv("INIT_BACKUP")
	if initBackup == "true" {
		c := cron.New()
		_, err = c.AddFunc("0 4 * * *", func() {
			log.Info().Msg("⏰ [CRON] Iniciando backup diario...")
			cfg, err := jobs.LoadConfig(dep)
			if err != nil {
				log.Err(err).Msg("❌ [CRON] error leyendo config")
				return
			}
			log.Info().Strs("databases", cfg.Databases).Msg("⏰ Iniciando backup")
			jobs.RunBackup(cfg)
		})
		if err != nil {
			log.Err(err).Msg("❌ Error al crear cron job")
			panic(err)
		}

		c.Start()
		defer c.Stop()
	}

	// Canal para señales del sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Iniciar servidor en goroutine
	go func() {
		port := getEnv("PORT", "3000")
		log.Info().Msgf("🚀 Servidor iniciado en http://localhost:%s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Err(err).Msg("Error al iniciar servidor")
		}
	}()

	var grpcServer *grpc.Server

    // --- SECCIÓN gRPC ACTUALIZADA ---
    go func() {
        portGrpc := getEnv("GRPC_PORT", "50051")
        lis, err := net.Listen("tcp", ":"+portGrpc)
        if err != nil {
            log.Fatal().Msgf("falló al escuchar gRPC: %v", err)
        }

        // Configuramos Keepalive para el servidor
        // Esto ayuda a detectar clientes muertos y mantener conexiones a través de Proxies/Load Balancers
        ka := grpc.KeepaliveParams(keepalive.ServerParameters{
            MaxConnectionIdle:     15 * time.Minute, // Tiempo max que una conexión puede estar ociosa
            MaxConnectionAge:      30 * time.Minute, // Forzar reconexión para balanceo de carga
            MaxConnectionAgeGrace: 5 * time.Minute,  // Tiempo extra para terminar llamadas antes de cerrar
            Time:                  5 * time.Second,  // Ping al cliente cada 5s para ver si sigue vivo
            Timeout:               1 * time.Second,  // Espera 1s por el pong
        })

        grpcServer = grpc.NewServer(
            ka,
            grpc.ChainUnaryInterceptor(
                interceptor.LoggingInterceptor,
                interceptor.AuthInterceptor,
                interceptor.MultiTenantInterceptor(dep),
            ),
        )

        // Registro de servicios
        productServer := &server.GrpcProductServer{}
        pb.RegisterProductServiceServer(grpcServer, productServer)

        tenantServer := &server.GrpcTenantServer{
            GrpcTenantService: depGrpc.TenantGrpcService,
        }
        pb.RegisterTenantServiceServer(grpcServer, tenantServer)

        categoryServer := &server.GrpcCategoryServer{}
        pb.RegisterCategoryServiceServer(grpcServer, categoryServer)

        log.Info().Msgf("🚀 Servidor gRPC iniciado en :%s", portGrpc)
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatal().Msgf("falló al servir gRPC: %v", err)
        }
    }()

	// Esperar señal de terminación
	<-quit
	// log.Println("🛑 Apagando servidor...")

	// Cerrar contexto para detener janitor
	// cancel()

	// Shutdown graceful con timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Err(err).Msg("Error durante shutdown")
	}

	// Cerrar todas las conexiones
	log.Debug().Msg("Cerrando conexiones...")

	if err := database.CloseDB(db); err != nil {
		log.Err(err).Msg("Error al cerrar DB principal")
	}

	if err := database.CloseAllTenantDBs(); err != nil {
		log.Err(err).Msg("Error al cerrar DBs de tenants")
	}

	if err := cache.CloseRedis(); err != nil {
		log.Err(err).Msg("Error al cerrar Redis")
	}

	log.Info().Msg("✅ Servidor apagado correctamente")
}

// healthCheck endpoint de health check
func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
		"redis":  cache.IsAvailable(),
		"time":   time.Now(),
	})
}

// cacheStats endpoint para estadísticas de cache
func cacheStats(c *fiber.Ctx) error {
	if !cache.IsAvailable() {
		return c.Status(503).JSON(fiber.Map{
			"error": "Redis no disponible",
		})
	}

	stats, err := cache.GetCacheStats()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(stats)
}

// customErrorHandler maneja errores de forma centralizada
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
		"code":  code,
	})
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
