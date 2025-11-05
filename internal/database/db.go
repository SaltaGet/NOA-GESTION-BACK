package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/cache"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	lru "github.com/hashicorp/golang-lru"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mainDB            *gorm.DB
	tenantDBs         *lru.Cache
	mu                sync.RWMutex
	dbExpiration      = 30 * time.Minute
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

func ConnectDB() (*gorm.DB, error) {
	dsn := os.Getenv("URI_DB")
	if dsn == "" {
		return nil, fmt.Errorf("la variable de entorno URI_DB no esta definida")
	}

	if err := ensureDatabaseExists(dsn); err != nil {
		log.Fatalf("No se pudo crear la base: %v", err)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	setupDBConnection(db, mainDBConfig)

	if err := db.AutoMigrate(&models.Tenant{}, &models.User{}, &models.UserTenant{}, &models.Plan{}, &models.Admin{}); err != nil {
		log.Fatalf("Error en migración: %v", err)
	}

	ensurePlans(db)

	return ensureAdmin(db)
}

func ensureAdmin(db *gorm.DB) (*gorm.DB, error) {
	var email string
	db.Model(&models.Admin{}).Select("email").Where("email = ?", os.Getenv("ADMIN_EMAIL")).Scan(&email)

	if email != "" {
		log.Println("El admin ya existe")
		mainDB = db
		return db, nil
	}

	if err := db.Create(&models.Admin{
		FirstName:     os.Getenv("FIRSTNAME_ADMIN"),
		LastName:      os.Getenv("LASTNAME_ADMIN"),
		Username:      os.Getenv("ADMIN_USERNAME"),
		Email:         os.Getenv("ADMIN_EMAIL"),
		Password:      os.Getenv("ADMIN_PASSWORD"),
		IsSuperAdmin: true,
	}).Error; err != nil {
		return nil, err
	}

	mainDB = db
	return db, nil
}

func ensurePlans(db *gorm.DB) error {
	plans := []models.Plan{
		{Name: "Free"}, {Name: "Basic"}, {Name: "Premium"},
	}

	err := db.Create(&plans)
	if err != nil {
		return err.Error
	}
	return nil
}

func setupDBConnection(db *gorm.DB, config DBConfig) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error al obtener conexión de base: %v", err)
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
		entry.lastUsed = time.Now()
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
		entry.lastUsed = time.Now()
		mu.RUnlock()
		return entry.db, nil
	}
	mu.RUnlock()

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
	// Validar archivo SQLite en modo dev
	if os.Getenv("ENV") != "prod" {
		key := filePathFromURI(connStr)
		if _, err := os.Stat(key); os.IsNotExist(err) {
			return nil, fmt.Errorf("la base de datos del tenant no existe: %s", connStr)
		}
	}

	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error al abrir DB de tenant: %w", err)
	}

	setupDBConnection(db, tenantDBConfig)
	return db, nil
}

func ensureDatabaseExists(dsn string) error {
	if os.Getenv("ENV") != "prod" {
		return nil
	}

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

func filePathFromURI(uri string) string {
	uri = strings.TrimPrefix(uri, "file:")
	if idx := strings.Index(uri, "?"); idx != -1 {
		uri = uri[:idx]
	}
	return uri
}

// StartDBJanitor limpia conexiones inactivas periódicamente
func StartDBJanitor(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cleanupInactiveConnections()
		}
	}
}

func cleanupInactiveConnections() {
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
			log.Printf("Conexión de tenant %v cerrada por inactividad", key)
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

func GetMainDB() *gorm.DB {
	return mainDB
}

func InitDBCache(maxEntries int) error {
	var err error
	tenantDBs, err = lru.New(maxEntries)
	if err != nil {
		return fmt.Errorf("error al inicializar cache LRU: %w", err)
	}
	return nil
}


// package database

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"
// 	"sync"
// 	"time"

// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
// 	lru "github.com/hashicorp/golang-lru"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// var (
// 	mainDB            *gorm.DB
// 	tenantDBs         *lru.Cache
// 	mu                sync.RWMutex
// 	dbExpiration      = 30 * time.Minute
// 	tenantConnections sync.Map // Cache de connection strings desencriptadas
// 	tenantLocks       sync.Map // Locks por tenant para evitar conexiones duplicadas
// )

// type tenantDBEntry struct {
// 	db       *gorm.DB
// 	lastUsed time.Time
// }

// type DBConfig struct {
// 	MaxOpenConns    int
// 	MaxIdleConns    int
// 	ConnMaxLifetime time.Duration
// 	ConnMaxIdleTime time.Duration
// }

// var (
// 	mainDBConfig = DBConfig{
// 		MaxOpenConns:    getEnvInt("MAIN_DB_MAX_OPEN", 50),
// 		MaxIdleConns:    getEnvInt("MAIN_DB_MAX_IDLE", 25),
// 		ConnMaxLifetime: 3 * time.Hour,
// 		ConnMaxIdleTime: 30 * time.Minute,
// 	}
// 	tenantDBConfig = DBConfig{
// 		MaxOpenConns:    getEnvInt("TENANT_DB_MAX_OPEN", 20),
// 		MaxIdleConns:    getEnvInt("TENANT_DB_MAX_IDLE", 5),
// 		ConnMaxLifetime: 3 * time.Hour,
// 		ConnMaxIdleTime: 30 * time.Minute,
// 	}
// )

// func getEnvInt(key string, defaultVal int) int {
// 	if val := os.Getenv(key); val != "" {
// 		var result int
// 		if _, err := fmt.Sscanf(val, "%d", &result); err == nil {
// 			return result
// 		}
// 	}
// 	return defaultVal
// }

// func ConnectDB() (*gorm.DB, error) {
// 	dsn := os.Getenv("URI_DB")
// 	if dsn == "" {
// 		return nil, fmt.Errorf("la variable de entorno URI_DB no esta definida")
// 	}

// 	if err := ensureDatabaseExists(dsn); err != nil {
// 		log.Fatalf("No se pudo crear la base: %v", err)
// 	}

// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	setupDBConnection(db, mainDBConfig)

// 	if err := db.AutoMigrate(&models.Tenant{}, &models.User{}, &models.UserTenant{}, &models.Plan{}); err != nil {
// 		log.Fatalf("Error en migración: %v", err)
// 	}

// 	return ensureAdmin(db)
// }

// func ensureAdmin(db *gorm.DB) (*gorm.DB, error) {
// 	var email string
// 	db.Model(&models.User{}).Select("email").Where("email = ?", os.Getenv("ADMIN_EMAIL")).Scan(&email)

// 	if email != "" {
// 		log.Println("El admin ya existe")
// 		mainDB = db
// 		return db, nil
// 	}

// 	if err := db.Create(&models.User{
// 		FirstName:     os.Getenv("FIRSTNAME_ADMIN"),
// 		LastName:      os.Getenv("LASTNAME_ADMIN"),
// 		Username:      os.Getenv("ADMIN_USERNAME"),
// 		Email:         os.Getenv("ADMIN_EMAIL"),
// 		Password:      os.Getenv("ADMIN_PASSWORD"),
// 		IsAdminTenant: true,
// 		IsAdmin:       true,
// 	}).Error; err != nil {
// 		return nil, err
// 	}

// 	mainDB = db
// 	return db, nil
// }

// func setupDBConnection(db *gorm.DB, config DBConfig) {
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		log.Fatalf("Error al obtener conexión de base: %v", err)
// 	}
// 	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
// 	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
// 	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
// 	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)
// }

// // GetTenantDB obtiene o crea una conexión de tenant con patrón double-check locking
// func GetTenantDB(encryptedConn string, tenantID int64) (*gorm.DB, error) {
// 	// 1️⃣ Verificación rápida si ya existe la conexión GORM
// 	mu.RLock()
// 	if val, ok := tenantDBs.Get(tenantID); ok {
// 		entry := val.(*tenantDBEntry)
// 		entry.lastUsed = time.Now()
// 		mu.RUnlock()
// 		return entry.db, nil
// 	}
// 	mu.RUnlock()

// 	// 2️⃣ Obtener lock específico del tenant para evitar conexiones duplicadas
// 	lockInterface, _ := tenantLocks.LoadOrStore(tenantID, &sync.Mutex{})
// 	tenantLock := lockInterface.(*sync.Mutex)
// 	tenantLock.Lock()
// 	defer tenantLock.Unlock()

// 	// 3️⃣ Double-check: verificar nuevamente si otro goroutine ya creó la conexión
// 	mu.RLock()
// 	if val, ok := tenantDBs.Get(tenantID); ok {
// 		entry := val.(*tenantDBEntry)
// 		entry.lastUsed = time.Now()
// 		mu.RUnlock()
// 		return entry.db, nil
// 	}
// 	mu.RUnlock()

// 	// 4️⃣ Obtener connection string desencriptada
// 	connStr, err := getDecryptedConnection(encryptedConn, tenantID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 5️⃣ Abrir nueva conexión
// 	db, err := openTenantDB(connStr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 6️⃣ Guardar en cache
// 	entry := &tenantDBEntry{db: db, lastUsed: time.Now()}
// 	mu.Lock()
// 	tenantDBs.Add(tenantID, entry)
// 	mu.Unlock()

// 	return db, nil
// }

// // getDecryptedConnection obtiene la connection string desencriptada con cache
// func getDecryptedConnection(encryptedConn string, tenantID int64) (string, error) {
// 	// Verificar cache
// 	if val, ok := tenantConnections.Load(tenantID); ok {
// 		return val.(string), nil
// 	}

// 	// Desencriptar
// 	if encryptedConn == "" {
// 		return "", fmt.Errorf("connection string vacía para tenant %d", tenantID)
// 	}

// 	decrypted, err := utils.Decrypt(encryptedConn)
// 	if err != nil {
// 		return "", fmt.Errorf("error al desencriptar connection: %w", err)
// 	}

// 	// Guardar en cache
// 	tenantConnections.Store(tenantID, decrypted)
// 	return decrypted, nil
// }

// // InvalidateTenantConnection invalida el cache de una connection string
// func InvalidateTenantConnection(tenantID int64) {
// 	tenantConnections.Delete(tenantID)
	
// 	mu.Lock()
// 	if val, ok := tenantDBs.Get(tenantID); ok {
// 		entry := val.(*tenantDBEntry)
// 		if db, err := entry.db.DB(); err == nil {
// 			db.Close()
// 		}
// 		tenantDBs.Remove(tenantID)
// 	}
// 	mu.Unlock()
// }

// func openTenantDB(connStr string) (*gorm.DB, error) {
// 	// Validar archivo SQLite en modo dev
// 	if os.Getenv("ENV") != "prod" {
// 		key := filePathFromURI(connStr)
// 		if _, err := os.Stat(key); os.IsNotExist(err) {
// 			return nil, fmt.Errorf("la base de datos del tenant no existe: %s", connStr)
// 		}
// 	}

// 	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
// 	if err != nil {
// 		return nil, fmt.Errorf("error al abrir DB de tenant: %w", err)
// 	}

// 	setupDBConnection(db, tenantDBConfig)
// 	return db, nil
// }

// func ensureDatabaseExists(dsn string) error {
// 	if os.Getenv("ENV") != "prod" {
// 		return nil
// 	}

// 	parts := strings.Split(dsn, "/")
// 	if len(parts) < 2 {
// 		return fmt.Errorf("DSN inválido: %s", dsn)
// 	}
// 	dbNameAndParams := parts[1]
// 	dbName := strings.Split(dbNameAndParams, "?")[0]

// 	dsnWithoutDB := strings.Split(dsn, "/")[0] + "/?charset=utf8mb4&parseTime=True&loc=Local"

// 	db, err := sql.Open("mysql", dsnWithoutDB)
// 	if err != nil {
// 		return fmt.Errorf("error al conectar sin base: %w", err)
// 	}
// 	defer db.Close()

// 	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName))
// 	return err
// }

// func filePathFromURI(uri string) string {
// 	uri = strings.TrimPrefix(uri, "file:")
// 	if idx := strings.Index(uri, "?"); idx != -1 {
// 		uri = uri[:idx]
// 	}
// 	return uri
// }

// // StartDBJanitor limpia conexiones inactivas periódicamente
// func StartDBJanitor(ctx context.Context) {
// 	ticker := time.NewTicker(1 * time.Hour)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case <-ticker.C:
// 			cleanupInactiveConnections()
// 		}
// 	}
// }

// func cleanupInactiveConnections() {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	keys := tenantDBs.Keys()
// 	for _, key := range keys {
// 		val, ok := tenantDBs.Peek(key) // Peek no actualiza LRU
// 		if !ok {
// 			continue
// 		}

// 		entry := val.(*tenantDBEntry)
// 		if time.Since(entry.lastUsed) > dbExpiration {
// 			if db, err := entry.db.DB(); err == nil {
// 				db.Close()
// 			}
// 			tenantDBs.Remove(key)
// 			log.Printf("Conexión de tenant %v cerrada por inactividad", key)
// 		}
// 	}
// }

// func CloseDB(db *gorm.DB) error {
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		return fmt.Errorf("no se pudo obtener la conexión de bajo nivel: %w", err)
// 	}

// 	if sqlDB != nil {
// 		if err := sqlDB.Close(); err != nil {
// 			return fmt.Errorf("error al cerrar la conexión: %w", err)
// 		}
// 	}
// 	return nil
// }

// func CloseAllTenantDBs() error {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	keys := tenantDBs.Keys()
// 	for _, key := range keys {
// 		if val, ok := tenantDBs.Get(key); ok {
// 			entry := val.(*tenantDBEntry)
// 			if db, err := entry.db.DB(); err == nil {
// 				db.Close()
// 			}
// 			tenantDBs.Remove(key)
// 		}
// 	}
// 	return nil
// }

// func GetMainDB() *gorm.DB {
// 	return mainDB
// }

// func InitDBCache(maxEntries int) error {
// 	var err error
// 	tenantDBs, err = lru.New(maxEntries)
// 	if err != nil {
// 		return fmt.Errorf("error al inicializar cache LRU: %w", err)
// 	}
// 	return nil
// }

// // package database

// // import (
// // 	"database/sql"
// // 	"fmt"
// // 	"log"
// // 	"os"
// // 	"strings"
// // 	"sync"
// // 	"time"

// // 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
// // 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
// // 	lru "github.com/hashicorp/golang-lru"
// // 	"gorm.io/driver/mysql"
// // 	"gorm.io/gorm"
// // )

// // var (
// // 	mainDB       *gorm.DB
// // 	tenantDBs    *lru.Cache
// // 	mu           sync.RWMutex
// // 	dbExpiration = 30 * time.Minute
// // 	tenantConnections sync.Map
// // )

// // type tenantDBEntry struct {
// // 	db       *gorm.DB
// // 	lastUsed time.Time
// // }

// // func ConnectDB() (*gorm.DB, error) {
// // 	dsn := os.Getenv("URI_DB")
// // 	if dsn == "" {
// // 		return nil, fmt.Errorf("la variable de entorno URI_DB no esta definida")
// // 	}

// // 	if err := ensureDatabaseExists(dsn); err != nil {
// // 		log.Fatalf("No se pudo crear la base: %v", err)
// // 	}

// // 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	setupDBConnection(db, 50, 25)

// // 	if err := db.AutoMigrate(&models.Tenant{}, &models.User{}, &models.UserTenant{}, &models.Plan{}); err != nil {
// // 		log.Fatalf("Error en migración: %v", err)
// // 	}

// // 	return ensureAdmin(db)
// // }

// // func ensureAdmin(db *gorm.DB) (*gorm.DB, error) {
// // 	var email string
// // 	db.Model(&models.User{}).Select("email").Where("email = ?", os.Getenv("ADMIN_EMAIL")).Scan(&email)

// // 	if email != "" {
// // 		log.Println("El admin ya existe")
// // 		mainDB = db
// // 		return db, nil
// // 	}

// // 	if err := db.Create(&models.User{
// // 		FirstName: os.Getenv("FIRSTNAME_ADMIN"),
// // 		LastName:  os.Getenv("LASTNAME_ADMIN"),
// // 		Username:  os.Getenv("ADMIN_USERNAME"),
// // 		Email:     os.Getenv("ADMIN_EMAIL"),
// // 		Password:  os.Getenv("ADMIN_PASSWORD"),
// // 		IsAdminTenant: true,
// // 		IsAdmin:   true,
// // 	}).Error; err != nil {
// // 		return nil, err
// // 	}

// // 	mainDB = db
// // 	return db, nil
// // }

// // func setupDBConnection(db *gorm.DB, maxOpen, maxIdle int) {
// // 	sqlDB, err := db.DB()
// // 	if err != nil {
// // 		log.Fatalf("Error al obtener conexión de base: %v", err)
// // 	}
// // 	sqlDB.SetMaxOpenConns(maxOpen)
// // 	sqlDB.SetMaxIdleConns(maxIdle)
// // 	sqlDB.SetConnMaxLifetime(3 * time.Hour)
// // 	sqlDB.SetConnMaxIdleTime(30 * time.Minute)
// // }

// // func GetTenantDB(encryptedConn string, tenantID int64) (*gorm.DB, error) {
// // 	// 1️⃣ Intentar obtener la conexión desencriptada en cache
// // 	connStr, ok := getCachedConnection(tenantID)
// // 	if !ok {
// // 		// Si no existe, desencriptar y guardar
// // 		decrypted, err := utils.Decrypt(encryptedConn)
// // 		if err != nil {
// // 			return nil, err
// // 		}

// // 		connStr = decrypted
// // 		cacheConnection(tenantID, connStr)
// // 	}

// // 	// 2️⃣ Verificar si ya existe la conexión GORM activa (en caché LRU)
// // 	mu.RLock()
// // 	if val, ok := tenantDBs.Get(tenantID); ok {
// // 		entry := val.(*tenantDBEntry)
// // 		entry.lastUsed = time.Now()
// // 		mu.RUnlock()
// // 		return entry.db, nil
// // 	}
// // 	mu.RUnlock()

// // 	// 3️⃣ Abrir una nueva conexión
// // 	db, err := openTenantDB(connStr)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	// 4️⃣ Guardar en cache
// // 	entry := &tenantDBEntry{db: db, lastUsed: time.Now()}
// // 	mu.Lock()
// // 	tenantDBs.Add(tenantID, entry)
// // 	mu.Unlock()

// // 	return db, nil
// // }

// // func getCachedConnection(tenantID int64) (string, bool) {
// // 	if val, ok := tenantConnections.Load(tenantID); ok {
// // 		return val.(string), true
// // 	}
// // 	return "", false
// // }

// // func cacheConnection(tenantID int64, conn string) {
// // 	tenantConnections.Store(tenantID, conn)
// // }

// // func openTenantDB(connStr string) (*gorm.DB, error) {
// // 	// Si estás en modo dev, validar archivo SQLite
// // 	if os.Getenv("ENV") != "prod" {
// // 		key := filePathFromURI(connStr)
// // 		if _, err := os.Stat(key); os.IsNotExist(err) {
// // 			return nil, fmt.Errorf("la base de datos del tenant no existe: %s", connStr)
// // 		}
// // 	}

// // 	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	setupDBConnection(db, 20, 5)
// // 	return db, nil
// // }


// // // func GetTenantDB(uri string) (*gorm.DB, error) {
// // // 	key := uri
// // // 	if os.Getenv("ENV") != "prod" {
// // // 		key = filePathFromURI(uri)
// // // 		if _, err := os.Stat(key); os.IsNotExist(err) {
// // // 			return nil, fmt.Errorf("la base de datos del tenant no existe: %s", uri)
// // // 		}
// // // 	}

// // // 	mu.RLock()
// // // 	if val, ok := tenantDBs.Get(key); ok {
// // // 		entry := val.(*tenantDBEntry)
// // // 		entry.lastUsed = time.Now()
// // // 		mu.RUnlock()
// // // 		return entry.db, nil
// // // 	}
// // // 	mu.RUnlock()

// // // 	var db *gorm.DB
// // // 	var err error

// // // 	db, err = gorm.Open(mysql.Open(uri), &gorm.Config{})
// // // 	if err != nil {
// // // 		return nil, err
// // // 	}

// // // 	setupDBConnection(db, 20, 5)

// // // 	entry := &tenantDBEntry{db: db, lastUsed: time.Now()}

// // // 	mu.Lock()
// // // 	tenantDBs.Add(key, entry)
// // // 	mu.Unlock()

// // // 	return db, nil
// // // }



// // func ensureDatabaseExists(dsn string) error {
// // 	if os.Getenv("ENV") != "prod" {
// // 		return nil
// // 	}

// // 	parts := strings.Split(dsn, "/")
// // 	if len(parts) < 2 {
// // 		return fmt.Errorf("DSN inválido: %s", dsn)
// // 	}
// // 	dbNameAndParams := parts[1]
// // 	dbName := strings.Split(dbNameAndParams, "?")[0]

// // 	dsnWithoutDB := strings.Split(dsn, "/")[0] + "/?charset=utf8mb4&parseTime=True&loc=Local"

// // 	db, err := sql.Open("mysql", dsnWithoutDB)
// // 	if err != nil {
// // 		return fmt.Errorf("error al conectar sin base: %w", err)
// // 	}
// // 	defer db.Close()

// // 	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName))
// // 	return err
// // }

// // func filePathFromURI(uri string) string {
// // 	uri = strings.TrimPrefix(uri, "file:")
// // 	if idx := strings.Index(uri, "?"); idx != -1 {
// // 		uri = uri[:idx]
// // 	}
// // 	return uri
// // }

// // func StartDBJanitor() {
// // 	ticker := time.NewTicker(1 * time.Hour)
// // 	defer ticker.Stop()

// // 	for range ticker.C {
// // 		mu.Lock()
// // 		for _, key := range tenantDBs.Keys() {
// // 			if val, ok := tenantDBs.Get(key); ok {
// // 				entry := val.(*tenantDBEntry)
// // 				if time.Since(entry.lastUsed) > dbExpiration {
// // 					if db, err := entry.db.DB(); err == nil {
// // 						db.Close()
// // 					}
// // 					tenantDBs.Remove(key)
// // 				}
// // 			}
// // 		}
// // 		mu.Unlock()
// // 	}
// // }

// // func CloseDB(db *gorm.DB) error {
// // 	sqlDB, err := db.DB()
// // 	if err != nil {
// // 		log.Fatal("No se pudo obtener la conexión de bajo nivel:", err)
// // 	}

// // 	if sqlDB != nil {
// // 		if err := sqlDB.Close(); err != nil {
// // 			log.Fatal("Error al cerrar la conexión:", err)
// // 		}
// // 	}
// // 	return nil
// // }

// // func CloseAllTenantDBs() error {
// // 	for _, key := range tenantDBs.Keys() {
// // 		if val, ok := tenantDBs.Get(key); ok {
// // 			entry := val.(*tenantDBEntry)
// // 			if db, err := entry.db.DB(); err == nil {
// // 				db.Close()
// // 			}
// // 			tenantDBs.Remove(key)
// // 		}
// // 	}
// // 	return nil
// // }

// // func GetMainDB() *gorm.DB {
// // 	return mainDB
// // }

// // func InitDBCache(maxEntries int) {
// // 	var err error
// // 	tenantDBs, err = lru.New(maxEntries)
// // 	if err != nil {
// // 		log.Println(err)
// // 	}
// // }
