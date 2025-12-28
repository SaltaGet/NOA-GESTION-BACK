package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/cache"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	lru "github.com/hashicorp/golang-lru"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mainDB            *gorm.DB
	tenantDBs         *lru.Cache
	mu                sync.RWMutex
	dbExpiration      = 2 * time.Minute
	tenantConnections sync.Map // Cache de connection strings desencriptadas
	tenantLocks       sync.Map // Locks por tenant para evitar conexiones duplicadas
)

type tenantDBEntry struct {
	db       *gorm.DB
	lastUsed time.Time
}

type DBConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

var (
	mainDBConfig = DBConfig{
		MaxOpenConns:    getEnvInt("MAIN_DB_MAX_OPEN", 50),
		MaxIdleConns:    getEnvInt("MAIN_DB_MAX_IDLE", 25),
		ConnMaxLifetime: 3 * time.Hour,
		ConnMaxIdleTime: 30 * time.Minute,
	}
	tenantDBConfig = DBConfig{
		MaxOpenConns:    getEnvInt("TENANT_DB_MAX_OPEN", 20),
		MaxIdleConns:    getEnvInt("TENANT_DB_MAX_IDLE", 5),
		ConnMaxLifetime: 3 * time.Hour,
		ConnMaxIdleTime: 30 * time.Minute,
	}
)

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		var result int
		if _, err := fmt.Sscanf(val, "%d", &result); err == nil {
			return result
		}
	}
	return defaultVal
}

func ConnectDB(cfg *schemas.EmailConfig) (*gorm.DB, error) {
	dsn := os.Getenv("URI_DB")
	if dsn == "" {
		return nil, fmt.Errorf("la variable de entorno URI_DB no esta definida")
	}

	if err := EnsureDatabaseExists(dsn); err != nil {
		log.Fatal().Err(err).Msg("No se pudo crear la base")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	setupDBConnection(db, mainDBConfig)

	if err := db.AutoMigrate(
		&models.Tenant{}, 
		&models.User{}, 
		&models.UserTenant{}, 
		&models.Plan{}, 
		&models.Admin{}, 
		&models.PayTenant{},
		&models.PayDetail{},
		&models.Module{},
		&models.TenantModule{},
		&models.Feedback{},
		&models.News{},
		&models.AuditLogAdmin{}); err != nil {
		log.Fatal().Err(err).Msg("Error en migración")
	}

	err = ensurePlans(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Error en migración de planes")
	}

	return ensureAdmin(db, cfg)
}

func ensureAdmin(db *gorm.DB, cfg *schemas.EmailConfig) (*gorm.DB, error) {
	raw := os.Getenv("ADMIN_EMAIL")
	if raw == "" {
		log.Warn().Msg("No se definió ADMIN_EMAIL")
		return db, nil
	}

	emailList := strings.Split(strings.ReplaceAll(raw, " ", ""), ",")

	var existing []string
	if err := db.Model(&models.Admin{}).
		Where("email IN (?)", emailList).
		Pluck("email", &existing).Error; err != nil {
		return nil, err
	}

	exists := map[string]bool{}
	for _, e := range existing {
		exists[e] = true
	}

	var adminsToCreate []models.Admin
	var passwords = map[string]string{} // guardar contraseñas generadas

	for _, email := range emailList {
		if exists[email] {
			continue
		}

		pass, err := utils.GenerateRandomString(6)
		if err != nil {
			return nil, err
		}

		if os.Getenv("ENV") == "dev" {
			pass = "123456"
		}

		userName := strings.Split(email, "@")[0]

		passwords[email] = pass

		adminsToCreate = append(adminsToCreate, models.Admin{
			Email:        email,
			Password:     pass,
			Username:     userName,
			IsSuperAdmin: true,
			FirstName:    "Admin",
			LastName:     "Admin",
		})
	}

	// Crear en DB
	if len(adminsToCreate) > 0 {
		if err := db.Create(&adminsToCreate).Error; err != nil {
			return nil, err
		}

		// Ahora que existen en la DB, enviar los emails
		for _, adm := range adminsToCreate {
			pass := passwords[adm.Email]

			utils.SendEmail(
				adm.Email,
				"Bienvenido a NOA-GESTION",
				utils.WelcomeAdmin(adm.Email, adm.Username, pass),
				cfg,
			)
		}
	}

	return db, nil
}



func ensurePlans(db *gorm.DB) error {
	plan := models.Plan{
		Name:            "Básico",
		PriceMounthly:   25,
		PriceYearly:     250,
		Description:     "Plan básico",
		Features:        "emmmm, nada es básico, asi que no esperes mucho",
		AmountPointSale: 1,
		AmountMember:    5,
		AmountProduct: 1000,
	}

	err := db.Create(&plan).Error
	if err != nil {
		if schemas.IsDuplicateError(err) {
			log.Warn().Msg("El plan básico ya existe")
			return nil
		}
		return err
	}
	return nil
}

func setupDBConnection(db *gorm.DB, config DBConfig) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Error al obtener conexión de base")
	}
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)
}

// GetTenantDB obtiene o crea una conexión de tenant con patrón double-check locking
func GetTenantDB(encryptedConn string, tenantID int64) (*gorm.DB, error) {
	// 1️⃣ Verificación rápida si ya existe la conexión GORM
	mu.RLock()
	if val, ok := tenantDBs.Get(tenantID); ok {
		entry := val.(*tenantDBEntry)
		entry.lastUsed = time.Now() // Actualizamos el uso para el LRU
		mu.RUnlock()
		return entry.db, nil
	}
	mu.RUnlock()

	// 2️⃣ Obtener lock específico del tenant para evitar conexiones duplicadas
	lockInterface, _ := tenantLocks.LoadOrStore(tenantID, &sync.Mutex{})
	tenantLock := lockInterface.(*sync.Mutex)
	tenantLock.Lock()
	defer tenantLock.Unlock()

	// 3️⃣ Double-check: verificar nuevamente si otro goroutine ya creó la conexión
	mu.RLock()
	if val, ok := tenantDBs.Get(tenantID); ok {
		entry := val.(*tenantDBEntry)
		sqlDB, err := entry.db.DB()
		if err == nil && sqlDB.Ping() == nil {
			entry.lastUsed = time.Now()
			mu.RUnlock()
			return entry.db, nil
		}
		// La conexión en cache está cerrada, la removemos para crear una nueva.
		entry.lastUsed = time.Now()
		mu.RUnlock()
		mu.Lock()
		tenantDBs.Remove(tenantID)
		mu.Unlock()
		return entry.db, nil
	} else {
		mu.RUnlock()
	}

	// 4️⃣ Obtener connection string desencriptada
	connStr, err := getDecryptedConnection(encryptedConn, tenantID)
	if err != nil {
		return nil, err
	}

	// 5️⃣ Abrir nueva conexión
	db, err := openTenantDB(connStr)
	if err != nil {
		return nil, err
	}

	// 6️⃣ Guardar en cache
	entry := &tenantDBEntry{db: db, lastUsed: time.Now()}
	mu.Lock()
	tenantDBs.Add(tenantID, entry)
	mu.Unlock()

	return db, nil
}

// getDecryptedConnection obtiene la connection string desencriptada con cache
func getDecryptedConnection(encryptedConn string, tenantID int64) (string, error) {
	// 1️⃣ Verificar cache de Redis primero
	if cache.IsAvailable() {
		if conn, err := cache.GetTenantConnection(tenantID); err == nil {
			return conn, nil
		}
	}

	// 2️⃣ Verificar cache en memoria (sync.Map)
	if val, ok := tenantConnections.Load(tenantID); ok {
		return val.(string), nil
	}

	// 3️⃣ Desencriptar
	if encryptedConn == "" {
		return "", fmt.Errorf("connection string vacía para tenant %d", tenantID)
	}

	decrypted, err := utils.Decrypt(encryptedConn)
	if err != nil {
		return "", fmt.Errorf("error al desencriptar connection: %w", err)
	}

	// 4️⃣ Guardar en ambos caches
	tenantConnections.Store(tenantID, decrypted)
	if cache.IsAvailable() {
		_ = cache.SetTenantConnection(tenantID, decrypted)
	}

	return decrypted, nil
}

// InvalidateTenantConnection invalida el cache de una connection string
func InvalidateTenantConnection(tenantID int64) {
	tenantConnections.Delete(tenantID)

	mu.Lock()
	if val, ok := tenantDBs.Get(tenantID); ok {
		entry := val.(*tenantDBEntry)
		if db, err := entry.db.DB(); err == nil {
			db.Close()
		}
		tenantDBs.Remove(tenantID)
	}
	mu.Unlock()
}

func openTenantDB(connStr string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error al abrir DB de tenant: %w", err)
	}

	setupDBConnection(db, tenantDBConfig)
	return db, nil
}

func EnsureDatabaseExists(dsn string) error {
	parts := strings.Split(dsn, "/")
	if len(parts) < 2 {
		return fmt.Errorf("DSN inválido: %s", dsn)
	}
	dbNameAndParams := parts[1]
	dbName := strings.Split(dbNameAndParams, "?")[0]

	dsnWithoutDB := strings.Split(dsn, "/")[0] + "/?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return fmt.Errorf("error al conectar sin base: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName))
	return err
}

func FilePathFromURI(uri string) string {
	uri = strings.TrimPrefix(uri, "file:")
	if idx := strings.Index(uri, "?"); idx != -1 {
		uri = uri[:idx]
	}
	return uri
}

// StartDBJanitor limpia conexiones inactivas periódicamente
func StartDBJanitor(ctx context.Context, tenants, gprcCache *sync.Map) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cleanupInactiveConnections(tenants, gprcCache)
		}
	}
}

func cleanupInactiveConnections(tenants *sync.Map, gprcCache *sync.Map) {
	mu.Lock()
	defer mu.Unlock()

	keys := tenantDBs.Keys()
	for _, key := range keys {
		val, ok := tenantDBs.Peek(key) // Peek no actualiza LRU
		if !ok {
			continue
		}

		entry := val.(*tenantDBEntry)
		if time.Since(entry.lastUsed) > dbExpiration {
			if db, err := entry.db.DB(); err == nil {
				db.Close()
			}
			tenantDBs.Remove(key)
			log.Info().Msgf("Conexión de tenant %v cerrada por inactividad", key)
			tenants.Delete(key.(int64))
			log.Info().Msgf("Conexión de tenant cache %v cerrada por inactividad", key)
			gprcCache.Delete(key.(int64))
			log.Info().Msgf("Conexión de gRPC cache %v cerrada por inactividad", key)
		}
	}
}

func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("no se pudo obtener la conexión de bajo nivel: %w", err)
	}

	if sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("error al cerrar la conexión: %w", err)
		}
	}
	return nil
}

func CloseAllTenantDBs() error {
	mu.Lock()
	defer mu.Unlock()

	keys := tenantDBs.Keys()
	for _, key := range keys {
		if val, ok := tenantDBs.Get(key); ok {
			entry := val.(*tenantDBEntry)
			if db, err := entry.db.DB(); err == nil {
				db.Close()
			}
			tenantDBs.Remove(key)
		}
	}
	return nil
}

// func GetMainDB() *gorm.DB {
// 	return mainDB
// }

func InitDBCache(maxEntries int) error {
	var err error
	tenantDBs, err = lru.New(maxEntries)
	if err != nil {
		return fmt.Errorf("error al inicializar cache LRU: %w", err)
	}
	return nil
}
