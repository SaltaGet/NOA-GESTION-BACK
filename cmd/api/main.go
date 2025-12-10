package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	// "os/signal"
	"strconv"
	"strings"

	// "syscall"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/initial"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/jobs"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/routes"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/cache"
	tenant_cache "github.com/SaltaGet/NOA-GESTION-BACK/internal/cache/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"

	// "github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	// "github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	// "github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/docs"
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
			log.Fatalf("Error cargando .env local: %v", err)
		}
	}

	path := os.Getenv("APP_ROOT")
	if path != "" {
		os.MkdirAll(path+"/backups", os.ModePerm)
	} else {
		log.Fatal("No está indicada la ruta raiz")
	}

	local := os.Getenv("LOCAL")
	if local == "true" {
		if err := jobs.GenerateSwagger(); err != nil {
			log.Fatalf("Error ejecutando swag init: %v", err)
		}
	}

	// 🔥 Inicializar Redis (opcional, falla gracefully)
	if err := cache.InitRedis(); err != nil {
		log.Printf("⚠️  Advertencia: Redis no disponible - %v", err)
	}

	// Inicializar cache de tenant DBs
	cacheSize := 100
	if err := database.InitDBCache(cacheSize); err != nil {
		log.Fatalf("Error al inicializar cache de DBs: %v", err)
	}

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	// Context para shutdown graceful
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tenants := tenant_cache.GetContainerTenantsCache()

	// Iniciar janitor para limpiar conexiones inactivas
	go database.StartDBJanitor(ctx, tenants)

	// s := sse.NewServer(&sse.Options{
	// 	Headers: map[string]string{
	// 		"Content-Type":      "text/event-stream",
	// 		"Cache-Control":     "no-cache, no-transform",
	// 		"Connection":        "keep-alive",
	// 		"X-Accel-Buffering": "no",
	// 	},
	// })
	// defer s.Shutdown()

	emailCfg := initial.InitEmail()

	dep := dependencies.NewApplication(db, emailCfg)

	err = jobs.Migrations(dep)
	if err != nil {
		// log.Fatalf("Error al aplicar migraciones: %v", err)
		logging.ERROR("Error al aplicar migraciones: %v", err)
	}

	app := fiber.New(fiber.Config{
		IdleTimeout:           30 * time.Second,
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
		ErrorHandler:          customErrorHandler,
		ProxyHeader:           "X-Forwarded-For",
		DisableStartupMessage: false,
		StreamRequestBody:     true,
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
	// 	Max:        120,
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
	// app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	routes.SetupRoutes(app, dep)

	initBackup := os.Getenv("INIT_BACKUP")
	if initBackup == "true" {
		c := cron.New()
		env := os.Getenv("ENV")
		if env == "prod" {
			// _, err = c.AddFunc("0 4 * * *", func() {
			_, err = c.AddFunc("@every 1m", func() {
				logging.INFO("⏰ [CRON] Iniciando backup diario...")
				cfg, err := jobs.LoadConfig(dep)
				if err != nil {
					logging.ERROR("❌ [CRON] error leyendo config: %s", err.Error())
					return
				}
				// if err := jobs.ExampleRestore(cfg, "string_tenant2"); err != nil {
				// 	// log.Fatal(err)
				// 	logging.ERROR("❌ [CRON] error restaurando DB: %s", err.Error())
				// }
				fmt.Println("⏰ Iniciando backup:", cfg.Databases)
				jobs.RunBackup(cfg)
			})
			if err != nil {
				logging.ERROR("❌ Error al crear cron job: %s", err.Error())
				panic(err)
			}
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
		log.Printf("🚀 Servidor iniciado en http://localhost:%s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Printf("Error al iniciar servidor: %v", err)
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
		log.Printf("Error durante shutdown: %v", err)
	}

	// Cerrar todas las conexiones
	log.Println("Cerrando conexiones...")

	if err := database.CloseDB(db); err != nil {
		log.Printf("Error al cerrar DB principal: %v", err)
	}

	if err := database.CloseAllTenantDBs(); err != nil {
		log.Printf("Error al cerrar DBs de tenants: %v", err)
	}

	if err := cache.CloseRedis(); err != nil {
		log.Printf("Error al cerrar Redis: %v", err)
	}

	log.Println("✅ Servidor apagado correctamente")
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
