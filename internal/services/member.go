package services

import (
	// "github.com/SaltaGet/NOA-GESTION-BACK/internal/cache"
	// "github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	// "github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (m *MemberService) MemberGetByID(id int64) (*schemas.MemberResponse, error) {
	return m.MemberRepository.MemberGetByID(id)
}

func (m *MemberService) MemberGetPermissionByUserID(userID int64) (*schemas.MemberResponse, error) {
	return m.MemberRepository.MemberGetPermissionByUserID(userID)
}

func (m *MemberService) MemberGetAll(limit, page int, search *map[string]string) ([]*schemas.MemberResponseDTO, int64, error) {
	return m.MemberRepository.MemberGetAll(limit, page, search)
}

func (m *MemberService) MemberCreate(memberID int64, memeberCreate *schemas.MemberCreate, plan *schemas.PlanResponseDTO) (int64, error) {
	amountMember, err := m.MemberRepository.MemberCount()
	if err != nil {
		return 0, err
	}

	if amountMember >= plan.AmountMember {
		return 0, schemas.ErrorResponse(400, "El plan no permite agregar mas miembros", fmt.Errorf("el plan no permite agregar mas miembros"))
	}

	return m.MemberRepository.MemberCreate(memberID, memeberCreate)
}

func (m *MemberService) MemberUpdate(memberID int64, memeberUpdate *schemas.MemberUpdate) error {
	return m.MemberRepository.MemberUpdate(memberID, memeberUpdate)
}

func (m *MemberService) MemberUpdatePassword(memberID int64, memberUpdatePassword *schemas.MemberUpdatePassword) error {
	return m.MemberRepository.MemberUpdatePassword(memberID, memberUpdatePassword)
}

func (m *MemberService) MemberDelete(memberID int64, id int64) (error) {
	return m.MemberRepository.MemberDelete(memberID, id)
}

// ********************************************************************************************************************************

// EXAMPLE USE CACHE

// UpdateMemberRole actualiza el rol de un miembro e invalida el cache
// func (s *MemberService) UpdateMemberRole(memberID, roleID int64, tenantID int64, connection string) error {
// 	db, err := database.GetTenantDB(connection, tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	// Obtener el miembro actual para saber qué usuario invalidar
// 	var member models.Member
// 	if err := db.Where("id = ?", memberID).First(&member).Error; err != nil {
// 		return err
// 	}

// 	// Actualizar el rol
// 	if err := db.Model(&member).Update("role_id", roleID).Error; err != nil {
// 		return err
// 	}

// 	// 🔥 Invalidar cache del usuario para que obtenga los nuevos permisos
// 	// if cache.IsAvailable() {
// 	// 	_ = cache.InvalidateAllUserVersions(member.UserID)
// 	// }

// 	return nil
// }

// // DeactivateMember desactiva un miembro e invalida su cache inmediatamente
// func (s *MemberService) DeactivateMember(memberID int64, tenantID int64, connection string) error {
// 	db, err := database.GetTenantDB(connection, tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	var member models.Member
// 	if err := db.Where("id = ?", memberID).First(&member).Error; err != nil {
// 		return err
// 	}

// 	// Desactivar
// 	if err := db.Model(&member).Update("is_active", false).Error; err != nil {
// 		return err
// 	}

// 	// 🔥 Invalidar cache para que pierda acceso inmediatamente
// 	// if cache.IsAvailable() {
// 	// 	_ = cache.InvalidateAllUserVersions(member.UserID)
// 	// }

// 	return nil
// }

// // --- Ejemplo de TenantService ---

// // UpdateTenantConnection actualiza la connection string de un tenant
// func (s *TenantService) UpdateTenantConnection(tenantID int64, newConnection string) error {
// 	db := database.GetMainDB()

// 	// Actualizar en DB
// 	if err := db.Model(&models.Tenant{}).
// 		Where("id = ?", tenantID).
// 		Update("connection", newConnection).Error; err != nil {
// 		return err
// 	}

// 	// 🔥 Invalidar cache de la connection antigua
// 	database.InvalidateTenantConnection(tenantID)

// 	return nil
// }

// // DeactivateTenant desactiva un tenant y todos sus usuarios
// func (s *TenantService) DeactivateTenant(tenantID int64) error {
// 	db := database.GetMainDB()

// 	// Desactivar tenant
// 	if err := db.Model(&models.Tenant{}).
// 		Where("id = ?", tenantID).
// 		Update("is_active", false).Error; err != nil {
// 		return err
// 	}

// 	// Obtener todos los usuarios del tenant
// 	// var userTenants []models.UserTenant
// 	// db.Where("tenant_id = ?", tenantID).Find(&userTenants)

// 	// 🔥 Invalidar cache de todos los usuarios
// 	// if cache.IsAvailable() {
// 	// 	for _, ut := range userTenants {
// 	// 		_ = cache.InvalidateAllUserVersions(ut.UserID)
// 	// 	}
// 	// }

// 	// Invalidar connection del tenant
// 	database.InvalidateTenantConnection(tenantID)

// 	return nil
// }

// // --- Ejemplo de RoleService ---

// // UpdateRolePermissions actualiza permisos de un rol
// func (s *RoleService) UpdateRolePermissions(roleID int64, permissionIDs []int64, tenantID int64, connection string) error {
// 	db, err := database.GetTenantDB(connection, tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	var role models.Role
// 	if err := db.Preload("Permissions").Where("id = ?", roleID).First(&role).Error; err != nil {
// 		return err
// 	}

// 	// Actualizar permisos
// 	var permissions []models.Permission
// 	db.Where("id IN ?", permissionIDs).Find(&permissions)

// 	if err := db.Model(&role).Association("Permissions").Replace(permissions); err != nil {
// 		return err
// 	}

// 	// 🔥 Invalidar cache de permisos del rol
// 	if cache.IsAvailable() {
// 		_ = cache.InvalidateRolePermissions(roleID)
// 	}

// 	// Invalidar cache de todos los miembros con este rol
// 	var members []models.Member
// 	db.Where("role_id = ?", roleID).Find(&members)

// 	// if cache.IsAvailable() {
// 	// 	for _, member := range members {
// 	// 		_ = cache.InvalidateAllUserVersions(member.UserID)
// 	// 	}
// 	// }

// 	return nil
// }

// // --- Ejemplo de AuthService ---

// // Login genera token con versión para cache
// func (s *AuthService) Login(username, password string, tenantIdentifier *string) (string, error) {
// 	// ... validar credenciales ...

// 	// authUser, err := s.AuthCurrentUser(user.ID, tenantID, memberID, -1)
// 	// if err != nil {
// 	// 	return "", err
// 	// }

// 	// // Generar token con versión (timestamp actual)
// 	// token, err := utils.GenerateTokenWithVersion(authUser, time.Now().Unix())
// 	// if err != nil {
// 	// 	return "", err
// 	// }

// 	// // Opcional: guardar en cache
// 	// if cache.IsAvailable() {
// 	// 	_ = cache.SetAuthUser(authUser, time.Now().Unix())
// 	// }

// 	// return token, nil
// 	return "", nil
// }

// // Logout invalida el refresh token
// func (s *AuthService) Logout(userID int64) error {
// 	if cache.IsAvailable() {
// 		_ = cache.RevokeRefreshToken(userID)
// 		_ = cache.InvalidateAllUserVersions(userID)
// 	}
// 	return nil
// }

// // RefreshToken renueva el access token
// func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
// 	// Verificar refresh token
// 	// userID, err := utils.VerifyRefreshToken(refreshToken)
// 	// if err != nil {
// 	// 	return "", err
// 	// }

// 	// // Validar que el refresh token coincida con el guardado
// 	// if cache.IsAvailable() {
// 	// 	storedToken, err := cache.GetRefreshToken(userID)
// 	// 	if err != nil || storedToken != refreshToken {
// 	// 		return "", fmt.Errorf("refresh token inválido")
// 	// 	}
// 	// }

// 	// // Generar nuevo access token
// 	// authUser, err := s.AuthCurrentUser(userID, -1, -1, -1)
// 	// if err != nil {
// 	// 	return "", err
// 	// }

// 	// return utils.GenerateTokenWithVersion(authUser, time.Now().Unix())
// 	return "", nil
// }
