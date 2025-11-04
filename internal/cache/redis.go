package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	ctx    = context.Background()
)

// RedisConfig configuración de Redis
type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	MaxRetries   int
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// InitRedis inicializa la conexión con Redis
func InitRedis() error {
	config := getRedisConfig()

	Client = redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		MaxRetries:   config.MaxRetries,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	// Verificar conexión
	if err := Client.Ping(ctx).Err(); err != nil {
		log.Printf("⚠️  Redis no disponible: %v. La app continuará sin cache.", err)
		Client = nil
		return err
	}

	log.Println("✅ Redis conectado correctamente")
	return nil
}

// getRedisConfig obtiene la configuración desde variables de entorno
func getRedisConfig() RedisConfig {
	return RedisConfig{
		Addr:         getEnv("REDIS_ADDR", "localhost:6379"),
		Password:     getEnv("REDIS_PASSWORD", ""),
		DB:           getEnvInt("REDIS_DB", 0),
		MaxRetries:   getEnvInt("REDIS_MAX_RETRIES", 3),
		PoolSize:     getEnvInt("REDIS_POOL_SIZE", 10),
		MinIdleConns: getEnvInt("REDIS_MIN_IDLE", 5),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

// CloseRedis cierra la conexión
func CloseRedis() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// --- Operaciones Básicas ---

// Set guarda un valor en Redis con TTL
func Set(key string, value interface{}, ttl time.Duration) error {
	if Client == nil {
		return nil // Fallar silenciosamente si Redis no está disponible
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error al serializar: %w", err)
	}

	return Client.Set(ctx, key, data, ttl).Err()
}

// Get obtiene un valor de Redis
func Get(key string, dest interface{}) error {
	if Client == nil {
		return fmt.Errorf("redis no disponible")
	}

	data, err := Client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Delete elimina una o más claves
func Delete(keys ...string) error {
	if Client == nil {
		return nil
	}

	return Client.Del(ctx, keys...).Err()
}

// Exists verifica si una clave existe
func Exists(key string) (bool, error) {
	if Client == nil {
		return false, fmt.Errorf("redis no disponible")
	}

	result, err := Client.Exists(ctx, key).Result()
	return result > 0, err
}

// DeletePattern elimina claves que coincidan con un patrón
func DeletePattern(pattern string) error {
	if Client == nil {
		return nil
	}

	iter := Client.Scan(ctx, 0, pattern, 0).Iterator()
	pipe := Client.Pipeline()

	count := 0
	for iter.Next(ctx) {
		pipe.Del(ctx, iter.Val())
		count++
		
		// Ejecutar en lotes de 100
		if count%100 == 0 {
			if _, err := pipe.Exec(ctx); err != nil {
				return err
			}
			pipe = Client.Pipeline()
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	// Ejecutar restantes
	if count%100 != 0 {
		_, err := pipe.Exec(ctx)
		return err
	}

	return nil
}

// SetNX establece un valor solo si la clave no existe (lock distribuido)
func SetNX(key string, value interface{}, ttl time.Duration) (bool, error) {
	if Client == nil {
		return false, fmt.Errorf("redis no disponible")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	return Client.SetNX(ctx, key, data, ttl).Result()
}

// Increment incrementa un contador
func Increment(key string) (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("redis no disponible")
	}

	return Client.Incr(ctx, key).Result()
}

// IncrementWithExpire incrementa y establece TTL
func IncrementWithExpire(key string, ttl time.Duration) (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("redis no disponible")
	}

	pipe := Client.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, ttl)
	
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, err
	}

	return incr.Val(), nil
}

// --- Helpers ---

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		var result int
		if _, err := fmt.Sscanf(val, "%d", &result); err == nil {
			return result
		}
	}
	return defaultVal
}

// IsAvailable verifica si Redis está disponible
func IsAvailable() bool {
	return Client != nil
}