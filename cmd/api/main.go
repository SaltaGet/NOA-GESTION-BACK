package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"strconv"
	"strings"

	"time"

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

	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"

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

				mpServer := &server.GrpcMPServer{}
				pb.RegisterMPServiceServer(grpcServer, mpServer)

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

// package main

// import (
// 	"bufio"
// 	"bytes"
// 	"crypto/tls"
// 	"encoding/base64"
// 	"encoding/xml"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"os/exec"
// 	"strings"
// 	"time"

// 	"github.com/rs/zerolog/log"
// )





// type FECAERequest struct {
// 	XMLName  xml.Name `xml:"FECAESolicitar"`
// 	Xmlns    string   `xml:"xmlns,attr"`
// 	Auth     Auth     `xml:"Auth"`
// 	FeCAEReq FeCAEReq `xml:"FeCAEReq"`
// }

// type Auth struct {
// 	Token string `xml:"Token"`
// 	Sign  string `xml:"Sign"`
// 	Cuit  int64  `xml:"Cuit"`
// }

// type FeCAEReq struct {
// 	FeCabReq FeCabReq          `xml:"FeCabReq"`
// 	FeDetReq []FECAEDetRequest `xml:"FeDetReq>FECAEDetRequest"`
// }

// type FeCabReq struct {
// 	CantReg  int `xml:"CantReg"`
// 	PtoVta   int `xml:"PtoVta"`
// 	CbteTipo int `xml:"CbteTipo"`
// }

// type FECAEDetRequest struct {
// 	Concepto               int            `xml:"Concepto"`
// 	DocTipo                int            `xml:"DocTipo"`
// 	DocNro                 int64          `xml:"DocNro"`
// 	CbteDesde              int64          `xml:"CbteDesde"`
// 	CbteHasta              int64          `xml:"CbteHasta"`
// 	CbteFch                string         `xml:"CbteFch"`
// 	ImpTotal               float64        `xml:"ImpTotal"`
// 	ImpTotConc             float64        `xml:"ImpTotConc"`
// 	ImpNeto                float64        `xml:"ImpNeto"`
// 	ImpOpEx                float64        `xml:"ImpOpEx"`
// 	ImpTrib                float64        `xml:"ImpTrib"`
// 	ImpIVA                 float64        `xml:"ImpIVA"`
// 	FchServDesde           string         `xml:"FchServDesde,omitempty"`
// 	FchServHasta           string         `xml:"FchServHasta,omitempty"`
// 	FchVtoPago             string         `xml:"FchVtoPago,omitempty"`
// 	MonId                  string         `xml:"MonId"`
// 	MonCotiz               float64        `xml:"MonCotiz"`
// 	CondicionIVAReceptorId int            `xml:"CondicionIVAReceptorId,omitempty"` // ← NOMBRE CORRECTO
// 	Tributos               *TributosArray `xml:"Tributos,omitempty"`
// 	Iva                    *IvaArray      `xml:"Iva,omitempty"`
// }

// type IvaArray struct {
// 	AlicIva []AlicIva `xml:"AlicIva"`
// }

// type AlicIva struct {
// 	Id      int     `xml:"Id"`
// 	BaseImp float64 `xml:"BaseImp"`
// 	Importe float64 `xml:"Importe"`
// }

// type TributosArray struct {
// 	Tributo []Tributo `xml:"Tributo"`
// }

// type Tributo struct {
// 	Id      int     `xml:"Id"`
// 	Desc    string  `xml:"Desc"`
// 	BaseImp float64 `xml:"BaseImp"`
// 	Alic    float64 `xml:"Alic"`
// 	Importe float64 `xml:"Importe"`
// }

// type FECAEResponse struct {
// 	XMLName xml.Name `xml:"Envelope"`
// 	Body    struct {
// 		Response struct {
// 			Results struct {
// 				FeCabResp FeCabResp `xml:"FeCabResp"`
// 				FeDetResp struct {
// 					FECAEDetResponse []FECAEDetResponse `xml:"FECAEDetResponse"`
// 				} `xml:"FeDetResp"`
// 			} `xml:"FECAESolicitarResult"`
// 			Errors struct {
// 				Err []ErrorAfip `xml:"Err"`
// 			} `xml:"Errors"`
// 		} `xml:"FECAESolicitarResponse"`
// 	} `xml:"Body"`
// }

// type FeCabResp struct {
// 	Cuit       int64  `xml:"Cuit"`
// 	PtoVta     int    `xml:"PtoVta"`
// 	CbteTipo   int    `xml:"CbteTipo"`
// 	FchProceso string `xml:"FchProceso"`
// 	CantReg    int    `xml:"CantReg"`
// 	Resultado  string `xml:"Resultado"`
// 	Reproceso  string `xml:"Reproceso"`
// }

// type FECAEDetResponse struct {
// 	Concepto      int    `xml:"Concepto"`
// 	DocTipo       int    `xml:"DocTipo"`
// 	DocNro        int64  `xml:"DocNro"`
// 	CbteDesde     int64  `xml:"CbteDesde"`
// 	CbteHasta     int64  `xml:"CbteHasta"`
// 	CbteFch       string `xml:"CbteFch"`
// 	Resultado     string `xml:"Resultado"`
// 	CAE           string `xml:"CAE"`
// 	CAEFchVto     string `xml:"CAEFchVto"`
// 	Observaciones struct {
// 		Obs []Observacion `xml:"Obs"`
// 	} `xml:"Observaciones"`
// }

// type Observacion struct {
// 	Code int    `xml:"Code"`
// 	Msg  string `xml:"Msg"`
// }

// type ErrorAfip struct {
// 	Code int    `xml:"Code"`
// 	Msg  string `xml:"Msg"`
// }

// type WSFEClient struct {
// 	BaseURL string
// 	Auth    Auth
// 	Client  *http.Client
// }

// // ============================================
// // MODELO DE FACTURA
// // ============================================

// type Factura struct {
// 	PuntoVenta       int
// 	TipoComprobante  int
// 	NumeroDesde      int64
// 	NumeroHasta      int64
// 	TipoDocumento    int
// 	NumeroDocumento  int64
// 	CondicionIVA     int
// 	Concepto         int
// 	Fecha            string
// 	FechaServDesde   string
// 	FechaServHasta   string
// 	FechaVtoPago     string
// 	ImporteNeto      float64
// 	ImporteNoGravado float64
// 	ImporteExento    float64
// 	ImporteIVA       float64
// 	ImporteTributos  float64
// 	ImporteTotal     float64
// 	Alicuotas        []ItemIVA
// 	Tributos         []ItemTributo
// 	MonedaId         string
// 	MonedaCotiz      float64
// }

// type ItemIVA struct {
// 	Codigo        int
// 	BaseImponible float64
// 	Importe       float64
// }

// type ItemTributo struct {
// 	Codigo        int
// 	Descripcion   string
// 	BaseImponible float64
// 	Alicuota      float64
// 	Importe       float64
// }

// func (f *Factura) ToFECAEDetRequest() FECAEDetRequest {
// 	det := FECAEDetRequest{
// 		Concepto:               f.Concepto,
// 		DocTipo:                f.TipoDocumento,
// 		DocNro:                 f.NumeroDocumento,
// 		CbteDesde:              f.NumeroDesde,
// 		CbteHasta:              f.NumeroHasta,
// 		CbteFch:                f.Fecha,
// 		ImpTotal:               f.ImporteTotal,
// 		ImpTotConc:             f.ImporteNoGravado,
// 		ImpNeto:                f.ImporteNeto,
// 		ImpOpEx:                f.ImporteExento,
// 		ImpTrib:                f.ImporteTributos,
// 		ImpIVA:                 f.ImporteIVA,
// 		MonId:                  f.MonedaId,
// 		MonCotiz:               f.MonedaCotiz,
// 		FchServDesde:           f.FechaServDesde,
// 		FchServHasta:           f.FechaServHasta,
// 		FchVtoPago:             f.FechaVtoPago,
// 		CondicionIVAReceptorId: f.CondicionIVA, // ← MAPEO CORRECTO
// 	}

// 	if len(f.Alicuotas) > 0 {
// 		det.Iva = &IvaArray{}
// 		for _, iva := range f.Alicuotas {
// 			det.Iva.AlicIva = append(det.Iva.AlicIva, AlicIva{
// 				Id:      iva.Codigo,
// 				BaseImp: iva.BaseImponible,
// 				Importe: iva.Importe,
// 			})
// 		}
// 	}

// 	if len(f.Tributos) > 0 {
// 		det.Tributos = &TributosArray{}
// 		for _, trib := range f.Tributos {
// 			det.Tributos.Tributo = append(det.Tributos.Tributo, Tributo{
// 				Id:      trib.Codigo,
// 				Desc:    trib.Descripcion,
// 				BaseImp: trib.BaseImponible,
// 				Alic:    trib.Alicuota,
// 				Importe: trib.Importe,
// 			})
// 		}
// 	}

// 	return det
// }

// // ============================================
// // WSAA - AUTENTICACIÓN
// // ============================================

// type LoginCMSResponse struct {
// 	XMLName xml.Name `xml:"Envelope"`
// 	Body    struct {
// 		LoginCmsReturn struct {
// 			Return string `xml:"loginCmsReturn"`
// 		} `xml:"loginCmsResponse"`
// 	} `xml:"Body"`
// }

// type TicketResponse struct {
// 	XMLName xml.Name `xml:"loginTicketResponse"`
// 	Header  struct {
// 		Source      string `xml:"source"`
// 		Destination string `xml:"destination"`
// 		UniqueID    int64  `xml:"uniqueId"`
// 		Generation  string `xml:"generationTime"`
// 		Expiration  string `xml:"expirationTime"`
// 	} `xml:"header"`
// 	Credentials struct {
// 		Token string `xml:"token"`
// 		Sign  string `xml:"sign"`
// 	} `xml:"credentials"`
// }

// type Credentials struct {
// 	Token      string
// 	Sign       string
// 	Expiration time.Time
// 	CUIT       int64
// }

// type WSAAConfig struct {
// 	Homologacion bool
// 	CertFile     string
// 	KeyFile      string
// 	CUIT         int64
// 	Service      string
// }

// type WSAA struct {
// 	Config WSAAConfig
// 	Client *http.Client
// }

// func NewWSAA(config WSAAConfig) (*WSAA, error) {
// 	return &WSAA{
// 		Config: config,
// 		Client: &http.Client{
// 			Timeout: 30 * time.Second,
// 			Transport: &http.Transport{
// 				TLSClientConfig: &tls.Config{
// 					InsecureSkipVerify: config.Homologacion,
// 				},
// 			},
// 		},
// 	}, nil
// }




// func NewWSFEClient(token, sign string, cuit int64, homologacion bool) *WSFEClient {
// 	url := "https://servicios1.afip.gov.ar/wsfev1/service.asmx"
// 	if homologacion {
// 		url = "https://wswhomo.afip.gov.ar/wsfev1/service.asmx"
// 	}

// 	return &WSFEClient{
// 		BaseURL: url,
// 		Auth: Auth{
// 			Token: token,
// 			Sign:  sign,
// 			Cuit:  cuit,
// 		},
// 		Client: &http.Client{Timeout: 30 * time.Second},
// 	}
// }

// func (w *WSFEClient) SolicitarCAE(factura *Factura) (*FECAEDetResponse, error) {
// 	req := FECAERequest{
// 		Xmlns: "http://ar.gov.afip.dif.FEV1/",
// 		Auth:  w.Auth,
// 		FeCAEReq: FeCAEReq{
// 			FeCabReq: FeCabReq{
// 				CantReg:  1,
// 				PtoVta:   factura.PuntoVenta,
// 				CbteTipo: factura.TipoComprobante,
// 			},
// 			FeDetReq: []FECAEDetRequest{factura.ToFECAEDetRequest()},
// 		},
// 	}

// 	soapEnv := w.buildSOAPEnvelope(req)

// 	httpReq, err := http.NewRequest("POST", w.BaseURL, strings.NewReader(soapEnv))
// 	if err != nil {
// 		return nil, err
// 	}

// 	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")
// 	httpReq.Header.Set("SOAPAction", "http://ar.gov.afip.dif.FEV1/FECAESolicitar")

// 	resp, err := w.Client.Do(httpReq)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("WSFE error %d: %s", resp.StatusCode, string(body))
// 	}

// 	var response FECAEResponse
// 	if err := xml.Unmarshal(body, &response); err != nil {
// 		return nil, fmt.Errorf("error parseando respuesta: %v\nBody: %s", err, string(body))
// 	}

// 	if len(response.Body.Response.Errors.Err) > 0 {
// 		errMsg := ""
// 		for _, e := range response.Body.Response.Errors.Err {
// 			errMsg += fmt.Sprintf("Error %d: %s\n", e.Code, e.Msg)
// 		}
// 		return nil, fmt.Errorf("errores de AFIP:\n%s", errMsg)
// 	}

// 	if len(response.Body.Response.Results.FeDetResp.FECAEDetResponse) == 0 {
// 		return nil, fmt.Errorf("sin resultados en la respuesta")
// 	}

// 	detResp := response.Body.Response.Results.FeDetResp.FECAEDetResponse[0]

// 	if detResp.Resultado != "A" {
// 		obsMsg := ""
// 		for _, obs := range detResp.Observaciones.Obs {
// 			obsMsg += fmt.Sprintf("Obs %d: %s\n", obs.Code, obs.Msg)
// 		}
// 		return &detResp, fmt.Errorf("comprobante rechazado:\n%s", obsMsg)
// 	}

// 	return &detResp, nil
// }

// func (w *WSFEClient) buildSOAPEnvelope(req FECAERequest) string {
// 	reqXML, _ := xml.MarshalIndent(req, "    ", "  ")

// 	soap := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
// <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
//   <soap:Body>
//     %s
//   </soap:Body>
// </soap:Envelope>`, string(reqXML))

// 	return soap
// }



// // ============================================
// // FUNCIONES AUXILIARES
// // ============================================

// func saveCredentials(creds *Credentials) error {
// 	data := fmt.Sprintf(`TOKEN=%s
// SIGN=%s
// CUIT=%d
// EXPIRES=%s
// `,
// 		creds.Token,
// 		creds.Sign,
// 		creds.CUIT,
// 		creds.Expiration.Format(time.RFC3339))

// 	return os.WriteFile("credentials.env", []byte(data), 0644)
// }

// func loadCredentials() (*Credentials, error) {
// 	file, err := os.Open("credentials.env")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	creds := &Credentials{}
// 	scanner := bufio.NewScanner(file)

// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		parts := strings.SplitN(line, "=", 2)
// 		if len(parts) != 2 {
// 			continue
// 		}

// 		key := strings.TrimSpace(parts[0])
// 		value := strings.TrimSpace(parts[1])

// 		switch key {
// 		case "TOKEN":
// 			creds.Token = value
// 		case "SIGN":
// 			creds.Sign = value
// 		case "CUIT":
// 			fmt.Sscanf(value, "%d", &creds.CUIT)
// 		case "EXPIRES":
// 			creds.Expiration, _ = time.Parse(time.RFC3339, value)
// 		}
// 	}

// 	return creds, nil
// }

// func (w *WSAA) GetCredentials() (*Credentials, error) {
// 	ticketXML, err := w.createTicketXML()
// 	if err != nil {
// 		return nil, fmt.Errorf("error creando XML ticket: %v", err)
// 	}

//     log.Info().Msg("✅ Ticket XML creado")

// 	signedCMS, err := w.signWithOpenSSL(ticketXML)
// 	if err != nil {
// 		return nil, fmt.Errorf("error firmando CMS: %v", err)
// 	}

//     log.Info().Msg("✅ CMS firmado")

// 	response, err := sendToWSAA(w, signedCMS)
// 	if err != nil {
// 		return nil, fmt.Errorf("error enviando a WSAA: %v", err)
// 	}

//     log.Info().Msg("✅ Respuesta recibida de WSAA")

// 	credentials, err := w.parseResponse(response)
// 	if err != nil {
// 		return nil, fmt.Errorf("error parseando respuesta: %v", err)
// 	}

// 	return credentials, nil
// }

// func sendToWSAA(w *WSAA, cms string) ([]byte, error) {
// 	url := "https://wsaahomo.afip.gov.ar/ws/services/LoginCms"
// 	if !w.Config.Homologacion {
// 		url = "https://wsaa.afip.gov.ar/ws/services/LoginCms"
// 	}

// 	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
// 		<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:wsaa="http://wsaa.view.sua.dvadac.desein.afip.gov">
// 				<soapenv:Header/>
// 				<soapenv:Body>
// 						<wsaa:loginCms>
// 								<wsaa:in0>%s</wsaa:in0>
// 						</wsaa:loginCms>
// 				</soapenv:Body>
// 		</soapenv:Envelope>`, cms)

// 	req, err := http.NewRequest("POST", url, strings.NewReader(soapRequest))
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
// 	req.Header.Set("SOAPAction", "")

// 	resp, err := w.Client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("WSAA error %d: %s", resp.StatusCode, string(body))
// 	}

// 	return body, nil
// }


// func (w *WSAA) parseResponse(data []byte) (*Credentials, error) {
// 	var loginResp LoginCMSResponse
// 	if err := xml.Unmarshal(data, &loginResp); err != nil {
// 		return nil, fmt.Errorf("error parseando SOAP: %v\nRespuesta: %s", err, string(data))
// 	}

// 	ticketXML := loginResp.Body.LoginCmsReturn.Return
// 	if ticketXML == "" {
// 		return nil, fmt.Errorf("respuesta vacía del WSAA")
// 	}

// 	var ticketResp TicketResponse
// 	if err := xml.Unmarshal([]byte(ticketXML), &ticketResp); err != nil {
// 		return nil, fmt.Errorf("error parseando ticket: %v", err)
// 	}

// 	expTime, err := time.Parse("2006-01-02T15:04:05.000-07:00", ticketResp.Header.Expiration)
// 	if err != nil {
// 		expTime, err = time.Parse("2006-01-02T15:04:05", ticketResp.Header.Expiration)
// 		if err != nil {
// 			expTime = time.Now().Add(12 * time.Hour)
// 		}
// 	}

// 	return &Credentials{
// 		Token:      ticketResp.Credentials.Token,
// 		Sign:       ticketResp.Credentials.Sign,
// 		Expiration: expTime,
// 		CUIT:       w.Config.CUIT,
// 	}, nil
// }

// func (w *WSAA) createTicketXML() ([]byte, error) {
// 	now := time.Now()

// 	ticket := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
// <loginTicketRequest version="1.0">
// <header>
//     <uniqueId>%d</uniqueId>
//     <generationTime>%s</generationTime>
//     <expirationTime>%s</expirationTime>
// </header>
// <service>%s</service>
// </loginTicketRequest>`,
// 		now.Unix(),
// 		now.Add(-10*time.Minute).Format("2006-01-02T15:04:05"),
// 		now.Add(12*time.Hour).Format("2006-01-02T15:04:05"),
// 		w.Config.Service,
// 	)

// 	return []byte(ticket), nil
// }

// func (w *WSAA) signWithOpenSSL(data []byte) (string, error) {
// 	tmpXML := "/tmp/ticket.xml"
// 	if err := os.WriteFile(tmpXML, data, 0644); err != nil {
// 		return "", err
// 	}
// 	defer os.Remove(tmpXML)

// 	tmpCMS := "/tmp/ticket.cms"
// 	defer os.Remove(tmpCMS)

// 	keyData, err := os.ReadFile(w.Config.KeyFile)
// 	if err != nil {
// 		return "", fmt.Errorf("no se puede leer la clave privada: %v", err)
// 	}

// 	certData, err := os.ReadFile(w.Config.CertFile)
// 	if err != nil {
// 		return "", fmt.Errorf("no se puede leer el certificado: %v", err)
// 	}

// 	if !bytes.Contains(keyData, []byte("BEGIN")) {
// 		return "", fmt.Errorf("la clave privada no está en formato PEM")
// 	}

// 	if bytes.Contains(keyData, []byte("ENCRYPTED")) {
// 		return "", fmt.Errorf("la clave privada está encriptada. Debe desencriptarla primero con:\nopenssl rsa -in privada.key -out privada_sin_pass.key")
// 	}

// 	log.Printf("🔍 Verificando formato de archivos...")
// 	log.Printf("   Clave: %d bytes, formato PEM: %v", len(keyData), bytes.Contains(keyData, []byte("BEGIN")))
// 	log.Printf("   Cert: %d bytes, formato PEM: %v", len(certData), bytes.Contains(certData, []byte("BEGIN")))

// 	cmd := exec.Command("openssl", "cms", "-sign",
// 		"-in", tmpXML,
// 		"-signer", w.Config.CertFile,
// 		"-inkey", w.Config.KeyFile,
// 		"-outform", "DER",
// 		"-out", tmpCMS,
// 		"-nodetach",
// 		"-binary",
// 		"-md", "sha256",
// 	)

// 	var stderr bytes.Buffer
// 	var stdout bytes.Buffer
// 	cmd.Stderr = &stderr
// 	cmd.Stdout = &stdout

// 	log.Printf("🔧 Ejecutando: openssl cms -sign...")

// 	if err := cmd.Run(); err != nil {
// 		log.Printf("⚠️  Fallo con cms, intentando con smime...")

// 		cmd = exec.Command("openssl", "smime", "-sign",
// 			"-in", tmpXML,
// 			"-signer", w.Config.CertFile,
// 			"-inkey", w.Config.KeyFile,
// 			"-outform", "DER",
// 			"-out", tmpCMS,
// 			"-nodetach",
// 			"-binary",
// 		)

// 		cmd.Stderr = &stderr
// 		cmd.Stdout = &stdout

// 		if err := cmd.Run(); err != nil {
// 			return "", fmt.Errorf("error ejecutando openssl: %v\nStderr: %s\nStdout: %s", err, stderr.String(), stdout.String())
// 		}
// 	}

// 	cmsData, err := os.ReadFile(tmpCMS)
// 	if err != nil {
// 		return "", err
// 	}

//     log.Info().Msgf("✅ CMS generado: %d bytes", len(cmsData))

// 	return base64.StdEncoding.EncodeToString(cmsData), nil
// }


// // ============================================
// // MAIN
// // ============================================

// func main() {
// 	fmt.Println("=== SISTEMA DE FACTURACIÓN ELECTRÓNICA AFIP ===\n")

// 	config := WSAAConfig{
// 		Homologacion: true,
// 		CertFile:     "certificado.crt",
// 		KeyFile:      "privada.key",
// 		CUIT:         20363467076,
// 		Service:      "wsfe",
// 	}

// 	if _, err := os.Stat(config.CertFile); os.IsNotExist(err) {
// 		log.Fatal().Err(err).Msg("Cert file not found")
// 	}
// 	if _, err := os.Stat(config.KeyFile); os.IsNotExist(err) {
// 		log.Fatal().Err(err).Msg("Key file not found")
// 	}

// 	fmt.Println("PASO 1: Obteniendo Token y Sign de AFIP...")
// 	fmt.Println(strings.Repeat("-", 60))

// 	creds, err := loadCredentials()
// 	if err != nil || time.Now().After(creds.Expiration) {
// 		fmt.Println("📡 Solicitando nuevas credenciales al WSAA...")

// 		wsaa, err := NewWSAA(config)
// 		if err != nil {
// 			log.Fatal().Err(err).Msg("Error creando cliente WSAA")
// 		}

// 		creds, err = wsaa.GetCredentials()
// 		if err != nil {
// 			log.Fatal().Err(err).Msg("Error obteniendo credenciales")
// 		}

// 		if err := saveCredentials(creds); err != nil {
// 			log.Printf("⚠️  No se pudieron guardar las credenciales: %v", err)
// 		} else {
// 			fmt.Println("💾 Credenciales guardadas en credentials.env")
// 		}
// 	} else {
// 		fmt.Println("✅ Usando credenciales existentes de credentials.env")
// 	}

// 	fmt.Printf("✅ Token obtenido (válido hasta: %s)\n", creds.Expiration.Format("2006-01-02 15:04:05"))
// 	fmt.Printf("⏱️  Expira en: %.1f horas\n\n", time.Until(creds.Expiration).Hours())

// 	fmt.Println("PASO 2: Creando factura de ejemplo...")
// 	fmt.Println(strings.Repeat("-", 60))

// 	factura := &Factura{
// 		PuntoVenta:       1,
// 		TipoComprobante:  6,
// 		NumeroDesde:      6,
// 		NumeroHasta:      6,
// 		TipoDocumento:    99,
// 		NumeroDocumento:  0,
// 		CondicionIVA:     5, // Consumidor Final
// 		Concepto:         1,
// 		Fecha:            time.Now().Format("20060102"),
// 		ImporteNeto:      1000.00,
// 		ImporteNoGravado: 0.00,
// 		ImporteExento:    0.00,
// 		ImporteIVA:       210.00,
// 		ImporteTributos:  0.00,
// 		ImporteTotal:     1210.00,
// 		Alicuotas: []ItemIVA{
// 			{
// 				Codigo:        5,
// 				BaseImponible: 1000.00,
// 				Importe:       210.00,
// 			},
// 		},
// 		MonedaId:    "PES",
// 		MonedaCotiz: 1,
// 	}

// 	fmt.Printf("📄 Factura B #%d-%08d\n", factura.PuntoVenta, factura.NumeroDesde)
// 	fmt.Printf("   Cliente: Consumidor Final\n")
// 	fmt.Printf("   Neto: $%.2f\n", factura.ImporteNeto)
// 	fmt.Printf("   IVA 21%%: $%.2f\n", factura.ImporteIVA)
// 	fmt.Printf("   Total: $%.2f\n\n", factura.ImporteTotal)

// 	fmt.Println("PASO 3: Solicitando CAE a AFIP...")
// 	fmt.Println(strings.Repeat("-", 60))

// 	wsfe := NewWSFEClient(creds.Token, creds.Sign, config.CUIT, config.Homologacion)

// 	resultado, err := wsfe.SolicitarCAE(factura)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("Error creando cliente WSAA")
// 	}

// 	fmt.Println("\n" + strings.Repeat("=", 60))
// 	fmt.Println("✅ FACTURA AUTORIZADA")
// 	fmt.Println(strings.Repeat("=", 60))
// 	fmt.Printf("\n🎫 CAE: %s\n", resultado.CAE)
// 	fmt.Printf("📅 Vencimiento CAE: %s\n", resultado.CAEFchVto)
// 	fmt.Printf("📋 Comprobante: %04d-%08d-%08d\n",
// 		factura.PuntoVenta,
// 		factura.TipoComprobante,
// 		resultado.CbteDesde)
// 	fmt.Printf("💰 Total: $%.2f\n", factura.ImporteTotal)

// 	if len(resultado.Observaciones.Obs) > 0 {
// 		fmt.Println("\n⚠️  Observaciones:")
// 		for _, obs := range resultado.Observaciones.Obs {
// 			fmt.Printf("   - [%d] %s\n", obs.Code, obs.Msg)
// 		}
// 	}

// 	fmt.Println("\n" + strings.Repeat("=", 60))
// 	fmt.Println("🎉 PROCESO COMPLETADO")
// 	fmt.Println(strings.Repeat("=", 60))

// 	fmt.Println("\n🔍 Consultando factura recién emitida...")
// 	info, err := wsfe.ConsultarFactura(factura.TipoComprobante, factura.PuntoVenta, factura.NumeroDesde)
// 	if err != nil {
// 		log.Printf("❌ No se pudo consultar: %v", err)
// 	} else {
// 		res := info.Body.Response.Result.ResultGet
// 		fmt.Printf("✅ Información recuperada:\n")
// 		fmt.Printf("   CAE: %s\n", res.CodAut)
// 		fmt.Printf("   Fecha: %s\n", res.CbteFch)
// 		fmt.Printf("   Importe: $%.2f\n", res.ImpTotal)
// 		fmt.Printf("   Resultado en AFIP: %s\n", res.Resultado)
// 	}

// 	fmt.Println("🔍 Consultando último número autorizado...")
// 	ultimo, err := wsfe.GetUltimoComprobante(factura.PuntoVenta, factura.TipoComprobante)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("Error")
// 	}

// 	fmt.Printf("✅ Último comprobante autorizado para PtoVta %d, Tipo %d: %d\n", factura.PuntoVenta, factura.TipoComprobante, ultimo)
// }

// type FECompConsultaRequest struct {
// 	XMLName       xml.Name `xml:"FECompConsultar"`
// 	Xmlns         string   `xml:"xmlns,attr"`
// 	Auth          Auth     `xml:"Auth"`
// 	FeCompConsReq struct {
// 		CbteTipo int   `xml:"CbteTipo"`
// 		CbteNro  int64 `xml:"CbteNro"`
// 		PtoVta   int   `xml:"PtoVta"`
// 	} `xml:"FeCompConsReq"`
// }

// type FECompConsultaResponse struct {
// 	XMLName xml.Name `xml:"Envelope"`
// 	Body    struct {
// 		Response struct {
// 			Result struct {
// 				ResultGet struct {
// 					CbteDesde int64   `xml:"CbteDesde"`
// 					CbteHasta int64   `xml:"CbteHasta"`
// 					CbteFch   string  `xml:"CbteFch"`
// 					ImpTotal  float64 `xml:"ImpTotal"`
// 					CodAut    string  `xml:"CodAut"` // Este es el CAE
// 					FchVto    string  `xml:"FchVto"` // Vto del CAE
// 					Resultado string  `xml:"Resultado"`
// 					DocTipo   int     `xml:"DocTipo"`
// 					DocNro    int64   `xml:"DocNro"`
// 					// Puedes agregar más campos según necesites del XSD
// 				} `xml:"ResultGet"`
// 				Errors struct {
// 					Err []ErrorAfip `xml:"Err"`
// 				} `xml:"Errors"`
// 			} `xml:"FECompConsultarResult"`
// 		} `xml:"FECompConsultarResponse"`
// 	} `xml:"Body"`
// }

// func (w *WSFEClient) ConsultarFactura(tipo int, ptoVta int, nro int64) (*FECompConsultaResponse, error) {
// 	req := FECompConsultaRequest{
// 		Xmlns: "http://ar.gov.afip.dif.FEV1/",
// 		Auth:  w.Auth,
// 	}
// 	req.FeCompConsReq.CbteTipo = tipo
// 	req.FeCompConsReq.PtoVta = ptoVta
// 	req.FeCompConsReq.CbteNro = nro

// 	// Usamos la misma lógica de SOAP que ya tienes
// 	reqXML, _ := xml.MarshalIndent(req, "    ", "  ")
// 	soapEnv := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
// <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
//   <soap:Body>
//     %s
//   </soap:Body>
// </soap:Envelope>`, string(reqXML))

// 	httpReq, err := http.NewRequest("POST", w.BaseURL, strings.NewReader(soapEnv))
// 	if err != nil {
// 		return nil, err
// 	}

// 	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")
// 	httpReq.Header.Set("SOAPAction", "http://ar.gov.afip.dif.FEV1/FECompConsultar")

// 	resp, err := w.Client.Do(httpReq)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)

// 	var response FECompConsultaResponse
// 	if err := xml.Unmarshal(body, &response); err != nil {
// 		return nil, fmt.Errorf("error parseando respuesta: %v", err)
// 	}

// 	// Verificar errores de AFIP
// 	if len(response.Body.Response.Result.Errors.Err) > 0 {
// 		return nil, fmt.Errorf("AFIP Error: %s", response.Body.Response.Result.Errors.Err[0].Msg)
// 	}

// 	return &response, nil
// }

// type FECompUltimoAutorizadoRequest struct {
// 	XMLName  xml.Name `xml:"FECompUltimoAutorizado"`
// 	Xmlns    string   `xml:"xmlns,attr"`
// 	Auth     Auth     `xml:"Auth"`
// 	PtoVta   int      `xml:"PtoVta"`
// 	CbteTipo int      `xml:"CbteTipo"`
// }

// type FECompUltimoAutorizadoResponse struct {
// 	XMLName xml.Name `xml:"Envelope"`
// 	Body    struct {
// 		Response struct {
// 			Result struct {
// 				CbteNro  int64 `xml:"CbteNro"`
// 				PtoVta   int   `xml:"PtoVta"`
// 				CbteTipo int   `xml:"CbteTipo"`
// 				Errors   struct {
// 					Err []ErrorAfip `xml:"Err"`
// 				} `xml:"Errors"`
// 			} `xml:"FECompUltimoAutorizadoResult"`
// 		} `xml:"FECompUltimoAutorizadoResponse"`
// 	} `xml:"Body"`
// }

// func (w *WSFEClient) GetUltimoComprobante(puntoVenta, tipoComprobante int) (int64, error) {
// 	req := FECompUltimoAutorizadoRequest{
// 		Xmlns:    "http://ar.gov.afip.dif.FEV1/",
// 		Auth:     w.Auth,
// 		PtoVta:   puntoVenta,
// 		CbteTipo: tipoComprobante,
// 	}

// 	reqXML, _ := xml.Marshal(req)
// 	soapEnv := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
// <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
//   <soap:Body>
//     %s
//   </soap:Body>
// </soap:Envelope>`, string(reqXML))

// 	httpReq, err := http.NewRequest("POST", w.BaseURL, strings.NewReader(soapEnv))
// 	if err != nil {
// 		return 0, err
// 	}

// 	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")
// 	httpReq.Header.Set("SOAPAction", "http://ar.gov.afip.dif.FEV1/FECompUltimoAutorizado")

// 	resp, err := w.Client.Do(httpReq)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)
// 	var response FECompUltimoAutorizadoResponse
// 	if err := xml.Unmarshal(body, &response); err != nil {
// 		return 0, err
// 	}

// 	if len(response.Body.Response.Result.Errors.Err) > 0 {
// 		return 0, fmt.Errorf("AFIP Error: %s", response.Body.Response.Result.Errors.Err[0].Msg)
// 	}

// 	return response.Body.Response.Result.CbteNro, nil
// }
