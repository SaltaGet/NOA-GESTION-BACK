package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/master"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/cache"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/utils"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/aarondl/sqlboiler/v4/boil"
	lru "github.com/hashicorp/golang-lru"
	_ "github.com/jackc/pgx/v5/stdlib" // Driver para sql.Open
	"github.com/rs/zerolog/log"
)

//go:embed schemas_db/main.sql
var mainSchemaSQL string

var (
	mainDB            *sql.DB
	tenantDBs         *lru.Cache
	mu                sync.RWMutex
	dbExpiration      = 30 * time.Minute
	tenantConnections sync.Map // Cache de connection strings desencriptadas
	tenantLocks       sync.Map // Locks por tenant para evitar conexiones duplicadas
)

type tenantDBEntry struct {
	db       *sql.DB
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
		MaxOpenConns:    getEnvInt("TENANT_DB_MAX_OPEN", 5),
		MaxIdleConns:    getEnvInt("TENANT_DB_MAX_IDLE", 2),
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

func ConnectDB(cfg *schemas.EmailConfig) (*sql.DB, error) {
	dsn := os.Getenv("URI_DB")
	if dsn == "" {
		return nil, fmt.Errorf("la variable de entorno URI_DB no esta definida")
	}

	if err := EnsureDatabaseExists(dsn); err != nil {
		log.Fatal().Err(err).Msg("No se pudo crear la base")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	setupDBConnection(db, mainDBConfig)

	if _, err := db.ExecContext(context.Background(), mainSchemaSQL); err != nil {
		log.Fatal().Err(err).Msg("Error en al crear tablas")
	}

	err = ensurePlans(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Error en migración de planes")
	}

	mainDB = db

	if err := ApplyAuditAdminTriggers(db); err != nil {
		log.Error().Err(err).Msg("Error al aplicar triggers de auditoria")
	}

	return ensureAdmin(db, cfg)
}

func ensureAdmin(db *sql.DB, cfg *schemas.EmailConfig) (*sql.DB, error) {
	parentCtx := context.Background()

	raw := os.Getenv("ADMIN_EMAIL")
	if raw == "" {
		log.Warn().Msg("No se definió ADMIN_EMAIL")
		return db, nil
	}

	emailList := strings.Split(strings.ReplaceAll(raw, " ", ""), ",")
	for _, email := range emailList {
		ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
		defer cancel()
		// 1. Verificar si el admin ya existe por su email
		exists, err := master.Admins(master.AdminWhere.Email.EQ(email)).Exists(ctx, db)
		if err != nil {
			return nil, fmt.Errorf("error al verificar existencia de admin %s: %w", email, err)
		}

		// 2. Si ya existe, pasamos al siguiente email de la lista
		if exists {
			log.Info().Msgf("Admin %s ya existe. Omitiendo creación y envío de correo.", email)
			continue
		}

		// 3. Preparar datos para el NUEVO admin
		var pass string
		if os.Getenv("ENV") == "dev" {
			pass = "123456"
		} else {
			// Nota: Aquí deberías guardar el HASH de la pass en la DB, no la pass plana
			pass, _ = utils.GenerateRandomString(6)
		}

		userName := strings.Split(email, "@")[0]
		admin := master.Admin{
			Email:        email,
			Password:     pass, // TODO: HashPass(pass)
			Username:     userName,
			IsSuperAdmin: true,
			FirstName:    "Admin",
			LastName:     "Admin",
		}

		// 4. Insertar solo si no existía
		if err := admin.Insert(ctx, db, boil.Infer()); err != nil {
			return nil, fmt.Errorf("error al insertar admin %s: %w", email, err)
		}

		// 5. ENVIAR EMAIL: Solo llegamos aquí si el admin es realmente nuevo
		log.Info().Msgf("Nuevo admin creado: %s. Enviando credenciales...", email)

		utils.SendEmail(
			admin.Email,
			"Bienvenido a NOA-GESTION",
			utils.WelcomeAdmin(admin.Email, admin.Username, pass),
			cfg,
		)
	}

	return db, nil
}

func ensurePlans(db *sql.DB) error {
	ctx := context.Background()

	_, err := db.Exec("SELECT 'public.plans'::regclass")
	if err != nil {
		log.Warn().Msg("Esperando a que la tabla plans sea visible...")
		time.Sleep(500 * time.Millisecond) // Un respiro para el catálogo de PG
	}

	plan := master.Plan{
		Name:            "Básico",
		Description:     "Plan básico",
		Features:        "Nada es básico, así que no esperes mucho",
		AmountPointSale: 1,
		AmountMember:    5,
		AmountProduct:   2000,
	}

	query := `
		INSERT INTO public.plans (name, price_mounthly, price_yearly, description, features, amount_point_sale, amount_member, amount_product, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		ON CONFLICT (name) DO NOTHING;
	`
	_, err = db.ExecContext(ctx, query,
		plan.Name,
		25.00,
		250.00,
		plan.Description,
		plan.Features,
		plan.AmountPointSale,
		plan.AmountMember,
		plan.AmountProduct,
	)

	if err != nil {
		return fmt.Errorf("error al asegurar plan: %w", err)
	}

	return nil
}

func setupDBConnection(db *sql.DB, config DBConfig) {
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
}

// GetTenantDB obtiene o crea una conexión de tenant con patrón double-check locking
func GetTenantDB(encryptedConn string, tenantID int64) (*sql.DB, error) {
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
		if err := entry.db.Ping(); err == nil {
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
		if db := entry.db; db != nil {
			db.Close()
		}
		tenantDBs.Remove(tenantID)
	}
	mu.Unlock()
}

func openTenantDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("error al abrir DB de tenant: %w", err)
	}

	setupDBConnection(db, tenantDBConfig)
	return db, nil
}

func EnsureDatabaseExists(dsn string) error {
	u, err := url.Parse(dsn)
	if err != nil {
		return fmt.Errorf("error al parsear DSN: %w", err)
	}

	// Extract database name from path (e.g., /ospam -> ospam)
	dbName := strings.TrimPrefix(u.Path, "/")
	if dbName == "" {
		// If DSN doesn't have a path, maybe it is just the host?
		// Postgres DSNs usually have the DB name in the path.
		return fmt.Errorf("nombre de base de datos no encontrado en DSN")
	}

	// Connect to 'postgres' database to check existence/create new DB
	// We modify the path to point to 'postgres'
	u.Path = "/postgres"
	dsnWithoutDB := u.String()

	db, err := sql.Open("pgx", dsnWithoutDB)
	if err != nil {
		return fmt.Errorf("error al conectar a postgres default db: %w", err)
	}
	defer db.Close()

	// Check if DB exists
	var exists bool
	query := `SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)`
	err = db.QueryRow(query, dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error verificando existencia de DB: %w", err)
	}

	if !exists {
		log.Info().Msgf("La base de datos %s no existe, creándola...", dbName)
		_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName))
		if err != nil {
			return fmt.Errorf("error creando base de datos: %w", err)
		}
		log.Info().Msg("Base de datos creada exitosamente")
	}

	return nil
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
			if db := entry.db; db != nil {
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

func CloseDB(db *sql.DB) error {
	sqlDB := db
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
			if db := entry.db; db != nil {
				db.Close()
			}
			tenantDBs.Remove(key)
		}
	}
	return nil
}

func GetMainDB() *sql.DB {
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

var TriggerAuditAdmin = `
	CREATE OR REPLACE FUNCTION audit_trigger_function_admin()
	RETURNS TRIGGER AS $$
	DECLARE
			current_member TEXT;
			current_tx_id BIGINT;
	BEGIN
			current_member := current_setting('app.current_member_id', true);
			current_tx_id := txid_current();

			IF current_member IS NULL OR current_member = '' OR current_member = '0' THEN
					RETURN NEW;
			END IF;

			INSERT INTO audit_log_admins (
					transaction_id,
					admin_id,
					method,
					path,
					old_value,
					new_value,
					created_at
			)
			VALUES (
					current_tx_id,
					current_member::BIGINT,
					LOWER(TG_OP),
					TG_TABLE_NAME,
					CASE WHEN TG_OP = 'INSERT' THEN NULL ELSE to_jsonb(OLD) END,
					CASE WHEN TG_OP = 'DELETE' THEN NULL ELSE to_jsonb(NEW) END,
					NOW()
			);

			RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;
`

func ApplyAuditAdminTriggers(db *sql.DB) error {
	// 1. Crear la función del trigger
	if _, err := db.Exec(TriggerAuditAdmin); err != nil {
		return err
	}

	// 2. Obtener TODAS las tablas del esquema público
	var tables []string
	queryTables := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_type = 'BASE TABLE'
		AND table_name != 'audit_log_admins' 
		AND table_name != 'migrations';
	`
	rows, err := db.Query(queryTables)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return err
		}
		tables = append(tables, table)
	}

	// 3. Aplicar el trigger a cada una
	for _, table := range tables {
		db.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS tr_audit_%s ON %s", table, table))

		createQuery := fmt.Sprintf(`
      CREATE TRIGGER tr_audit_%s
      AFTER INSERT OR UPDATE OR DELETE ON "%s"
      FOR EACH ROW EXECUTE FUNCTION audit_trigger_function_admin();`,
			table, table)

		if _, err := db.Exec(createQuery); err != nil {
			return fmt.Errorf("error al crear trigger en tabla %s: %w", table, err)
		}
	}
	return nil
}
