package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	lru "github.com/hashicorp/golang-lru"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mainDB       *gorm.DB
	tenantDBs    *lru.Cache
	mu           sync.RWMutex
	dbExpiration = 30 * time.Minute
)

type tenantDBEntry struct {
	db       *gorm.DB
	lastUsed time.Time
}

func ConnectDB(dsn string) (*gorm.DB, error) {
	if err := ensureDatabaseExists(dsn); err != nil {
		log.Fatalf("No se pudo crear la base: %v", err)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	setupDBConnection(db, 50, 25)

	if err := db.AutoMigrate(&models.Tenant{}, &models.User{}, &models.UserTenant{}); err != nil {
		log.Fatalf("Error en migración: %v", err)
	}

	return ensureAdmin(db)
}

func ensureAdmin(db *gorm.DB) (*gorm.DB, error) {
	var email string
	db.Model(&models.User{}).Select("email").Where("email = ?", os.Getenv("ADMIN_EMAIL")).Scan(&email)

	if email != "" {
		log.Println("El admin ya existe")
		mainDB = db
		return db, nil
	}

	pass, err := utils.HashPassword(os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		return nil, err
	}

	if err := db.Create(&models.User{
		FirstName: os.Getenv("FIRSTNAME_ADMIN"),
		LastName:  os.Getenv("LASTNAME_ADMIN"),
		Username:  os.Getenv("ADMIN_USERNAME"),
		Email:     os.Getenv("ADMIN_EMAIL"),
		Password:  pass,
		IsAdmin:   true,
	}).Error; err != nil {
		return nil, err
	}

	mainDB = db
	return db, nil
}

func setupDBConnection(db *gorm.DB, maxOpen, maxIdle int) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error al obtener conexión de base: %v", err)
	}
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(3 * time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)
}

func GetTenantDB(uri string) (*gorm.DB, error) {
	key := uri
	if os.Getenv("ENV") != "prod" {
		key = filePathFromURI(uri)
		if _, err := os.Stat(key); os.IsNotExist(err) {
			return nil, fmt.Errorf("la base de datos del tenant no existe: %s", uri)
		}
	}

	mu.RLock()
	if val, ok := tenantDBs.Get(key); ok {
		entry := val.(*tenantDBEntry)
		entry.lastUsed = time.Now()
		mu.RUnlock()
		return entry.db, nil
	}
	mu.RUnlock()

	var db *gorm.DB
	var err error

	db, err = gorm.Open(mysql.Open(uri), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	setupDBConnection(db, 20, 5)

	entry := &tenantDBEntry{db: db, lastUsed: time.Now()}

	mu.Lock()
	tenantDBs.Add(key, entry)
	mu.Unlock()

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

func StartDBJanitor() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		mu.Lock()
		for _, key := range tenantDBs.Keys() {
			if val, ok := tenantDBs.Get(key); ok {
				entry := val.(*tenantDBEntry)
				if time.Since(entry.lastUsed) > dbExpiration {
					if db, err := entry.db.DB(); err == nil {
						db.Close()
					}
					tenantDBs.Remove(key)
				}
			}
		}
		mu.Unlock()
	}
}

func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("No se pudo obtener la conexión de bajo nivel:", err)
	}

	if sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			log.Fatal("Error al cerrar la conexión:", err)
		}
	}
	return nil
}

func CloseAllTenantDBs() error {
	for _, key := range tenantDBs.Keys() {
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

func InitDBCache(maxEntries int) {
	var err error
	tenantDBs, err = lru.New(maxEntries)
	if err != nil {
		log.Println(err)
	}
}
