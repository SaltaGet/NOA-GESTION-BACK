package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

const (
	// TTLs
	AuthUserTTL    = 5 * time.Minute
	TenantInfoTTL  = 10 * time.Minute
	PermissionsTTL = 15 * time.Minute

	// Prefijos de claves
	authUserPrefix    = "auth:user:"
	tenantPrefix      = "tenant:"
	permissionsPrefix = "perms:role:"
)

// --- Cache de Usuario Autenticado ---

// GetAuthUser obtiene usuario autenticado del cache
func GetAuthUser(userID int64, tenantID int64) (*schemas.AuthenticatedUser, error) {
	if Client == nil {
		return nil, fmt.Errorf("redis no disponible")
	}

	key := authUserKey(userID, tenantID)
	var authUser schemas.AuthenticatedUser

	if err := Get(key, &authUser); err != nil {
		return nil, err
	}

	return &authUser, nil
}

// SetAuthUser guarda usuario autenticado en cache
func SetAuthUser(authUser *schemas.AuthenticatedUser) error {
	if Client == nil {
		return nil // Fallar silenciosamente
	}

	key := authUserKey(authUser.ID, authUser.TenantID)
	return Set(key, authUser, AuthUserTTL)
}

// InvalidateAuthUser invalida cache de un usuario específico
func InvalidateAuthUser(memberID, tenantID int64) error {
	if Client == nil {
		return nil
	}

	key := authUserKey(memberID, tenantID)
	return Delete(key)
}

// InvalidateAllUserVersions invalida todas las versiones de un usuario
func InvalidateAllUserVersions(memberID int64) error {
	if Client == nil {
		return nil
	}

	pattern := fmt.Sprintf("%s%d:v*", authUserPrefix, memberID)
	return DeletePattern(pattern)
}

func authUserKey(memberID, tenantID int64) string {
	return fmt.Sprintf("%s%d:%d", authUserPrefix, memberID, tenantID)
}

// --- Cache de Tenant ---

// GetTenantConnection obtiene la connection string desencriptada del cache
func GetTenantConnection(tenantID int64) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("redis no disponible")
	}

	key := tenantConnectionKey(tenantID)
	var connection string

	if err := Get(key, &connection); err != nil {
		return "", err
	}

	return connection, nil
}

// SetTenantConnection guarda la connection string en cache
func SetTenantConnection(tenantID int64, connection string) error {
	if Client == nil {
		return nil
	}

	key := tenantConnectionKey(tenantID)
	return Set(key, connection, TenantInfoTTL)
}

// InvalidateTenantConnection invalida el cache de un tenant
func InvalidateTenantConnection(tenantID int64) error {
	if Client == nil {
		return nil
	}

	// Invalidar connection
	connKey := tenantConnectionKey(tenantID)
	if err := Delete(connKey); err != nil {
		return err
	}

	// Invalidar todos los usuarios de ese tenant
	pattern := fmt.Sprintf("%s*", authUserPrefix)
	return DeletePattern(pattern)
}

func tenantConnectionKey(tenantID int64) string {
	return fmt.Sprintf("%s%d:conn", tenantPrefix, tenantID)
}

// --- Cache de Permisos ---

// GetRolePermissions obtiene permisos de un rol del cache
func GetRolePermissions(roleID int64) ([]string, error) {
	if Client == nil {
		return nil, fmt.Errorf("redis no disponible")
	}

	key := rolePermissionsKey(roleID)
	var permissions []string

	if err := Get(key, &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

// SetRolePermissions guarda permisos de un rol en cache
func SetRolePermissions(roleID int64, permissions []string) error {
	if Client == nil {
		return nil
	}

	key := rolePermissionsKey(roleID)
	return Set(key, permissions, PermissionsTTL)
}

// InvalidateRolePermissions invalida permisos de un rol
func InvalidateRolePermissions(roleID int64) error {
	if Client == nil {
		return nil
	}

	key := rolePermissionsKey(roleID)
	return Delete(key)
}

func rolePermissionsKey(roleID int64) string {
	return fmt.Sprintf("%s%d", permissionsPrefix, roleID)
}

// --- Rate Limiting ---

// CheckRateLimit verifica límite de peticiones
func CheckRateLimit(identifier string, maxRequests int, window time.Duration) (bool, error) {
	if Client == nil {
		return true, nil // Sin Redis, permitir
	}

	key := fmt.Sprintf("ratelimit:%s", identifier)
	
	count, err := IncrementWithExpire(key, window)
	if err != nil {
		return true, err // En caso de error, permitir
	}

	return count <= int64(maxRequests), nil
}

// --- Session Management ---

// StoreRefreshToken guarda un refresh token
func StoreRefreshToken(userID int64, token string, ttl time.Duration) error {
	if Client == nil {
		return nil
	}

	key := fmt.Sprintf("refresh:%d", userID)
	return Set(key, token, ttl)
}

// GetRefreshToken obtiene un refresh token
func GetRefreshToken(userID int64) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("redis no disponible")
	}

	key := fmt.Sprintf("refresh:%d", userID)
	var token string

	if err := Get(key, &token); err != nil {
		return "", err
	}

	return token, nil
}

// RevokeRefreshToken revoca un refresh token
func RevokeRefreshToken(userID int64) error {
	if Client == nil {
		return nil
	}

	key := fmt.Sprintf("refresh:%d", userID)
	return Delete(key)
}

// --- Distributed Lock ---

// AcquireLock adquiere un lock distribuido
func AcquireLock(lockKey string, ttl time.Duration) (bool, error) {
	if Client == nil {
		return true, nil // Sin Redis, permitir
	}

	key := fmt.Sprintf("lock:%s", lockKey)
	return SetNX(key, time.Now().Unix(), ttl)
}

// ReleaseLock libera un lock distribuido
func ReleaseLock(lockKey string) error {
	if Client == nil {
		return nil
	}

	key := fmt.Sprintf("lock:%s", lockKey)
	return Delete(key)
}

// --- Cache Stats ---

// GetCacheStats obtiene estadísticas del cache
func GetCacheStats() (map[string]interface{}, error) {
	if Client == nil {
		return nil, fmt.Errorf("redis no disponible")
	}

	info, err := Client.Info(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"connected": true,
		"info":      info,
	}

	return stats, nil
}

// WarmupCache precarga datos frecuentes en el cache
func WarmupCache(data map[string]interface{}) error {
	if Client == nil {
		return nil
	}

	pipe := Client.Pipeline()
	for key, value := range data {
		jsonData, err := json.Marshal(value)
		if err != nil {
			continue
		}
		pipe.Set(context.Background(), key, jsonData, AuthUserTTL)
	}

	_, err := pipe.Exec(context.Background())
	return err
}



// --- Blacklist para revocación inmediata ---

// BlacklistUser agrega un usuario a la blacklist temporal
func BlacklistUser(userID int64, duration time.Duration) error {
	if Client == nil {
		return nil
	}

	key := blacklistKey(userID)
	return Set(key, true, duration)
}

// IsBlacklisted verifica si un usuario está en la blacklist
func IsBlacklisted(userID int64) bool {
	if Client == nil {
		return false
	}

	key := blacklistKey(userID)
	exists, err := Exists(key)
	return err == nil && exists
}

// RemoveFromBlacklist elimina un usuario de la blacklist
func RemoveFromBlacklist(userID int64) error {
	if Client == nil {
		return nil
	}

	key := blacklistKey(userID)
	return Delete(key)
}

func blacklistKey(userID int64) string {
	return fmt.Sprintf("blacklist:user:%d", userID)
}

// BlacklistTenant agrega todos los usuarios de un tenant a la blacklist
func BlacklistTenant(tenantID int64, duration time.Duration) error {
	if Client == nil {
		return nil
	}

	key := blacklistTenantKey(tenantID)
	return Set(key, true, duration)
}

// IsTenantBlacklisted verifica si un tenant está en la blacklist
func IsTenantBlacklisted(tenantID int64) bool {
	if Client == nil {
		return false
	}

	key := blacklistTenantKey(tenantID)
	exists, err := Exists(key)
	return err == nil && exists
}

func blacklistTenantKey(tenantID int64) string {
	return fmt.Sprintf("blacklist:tenant:%d", tenantID)
}