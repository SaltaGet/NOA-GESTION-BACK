package repositories

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/DanielChachagua/GestionCar/pkg/database"
	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/DanielChachagua/GestionCar/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// func (r *MainRepository) TenantGetByID(tenantID string) (*models.Tenant, error) {
// 	var tenant models.Tenant
// 	err := r.DB.
// 		Where("id = ?", tenantID).
// 		First(&tenant).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, models.ErrorResponse(404, "Tenant not found", err)
// 		}
// 		return nil, models.ErrorResponse(500, "Error retrieving tenant", err)
// 	}

// 	if !tenant.IsActive {
// 		return nil, models.ErrorResponse(403, "Tenant is inactive", nil)
// 	}

// 	return &tenant, nil
// }
func (r *MainRepository) TenantGetByID(tenantID string) (*models.Tenant, error) {
    var tenant models.Tenant
    err := r.DB.
        Where("id = ?", tenantID).
        First(&tenant).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, models.ErrorResponse(404, "Tenant not found", err)
        }
        return nil, models.ErrorResponse(500, "Error interno retrieving tenant", err)
    }

    if !tenant.IsActive {
        return nil, models.ErrorResponse(403, "Tenant is inactive", fmt.Errorf("tenant is inactive"))
    }

    return &tenant, nil
}


func (r *MainRepository) TenantGetByIdentifier(identifier string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.DB.
		Where("identifier = ?", identifier).
		First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return nil, models.ErrorResponse(500, "Error interno al obtener los tenants", err)
	}

	if !tenant.IsActive {
		return nil, models.ErrorResponse(403, "Tenant is inactive", fmt.Errorf("tenant is inactive"))
	}

	return &tenant, nil
}

func (r *MainRepository) TenantGetAll(userID string) (*[]models.TenantResponse, error) {
	var tenants []models.TenantResponse

	err := r.DB.
		Model(&models.Tenant{}).
		Select(`tenants.id, tenants.name, tenants.address, tenants.phone,
                tenants.email, tenants.is_active,
                user_tenants.is_active AS user_is_active,
                tenants.created_at, tenants.updated_at`).
		Joins("JOIN user_tenants ON tenants.id = user_tenants.tenant_id").
		Where("user_tenants.user_id = ?", userID).
		Scan(&tenants).Error

	if err != nil {
		return nil, models.ErrorResponse(500, "Error interno retrieving tenants", err)
	}

	return &tenants, nil
}

func (r *MainRepository) TenantGetConections() (*[]string, error) {
	var connections []string
	if err := r.DB.Debug().
	Model(&models.Tenant{}).
	Select("connection").
	Scan(&connections).Error; err != nil {
		return nil, models.ErrorResponse(500, "Error interno al obtener las connections", err)
	}

	var connectionsDecrypted []string

	for _, connection := range connections {
		connectionDecrypted, err := utils.Decrypt(connection)
		if err != nil {
			return nil, models.ErrorResponse(500, "Error interno al obtener las connections", err)
		}
		connectionsDecrypted = append(connectionsDecrypted, connectionDecrypted)
	}

	return &connectionsDecrypted, nil
}

func (r *MainRepository) TenantCreateByUserID(tenantCreate *models.TenantCreate, userID string) (string, error) {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	tenantName := strings.ReplaceAll(tenantCreate.Name, " ", "_")
	identifier := strings.ReplaceAll(tenantCreate.Identifier, " ", "_")
	uri := fmt.Sprintf("%s%s_%s.db%s", os.Getenv("URI_PATH"), tenantName, identifier, os.Getenv("URI_CONFIG"))
	connection, err := utils.Encrypt(uri)
	if err != nil {
		return "", models.ErrorResponse(500, "Error interno al obtener connection", err)
	}

	tenant := &models.Tenant{
		ID:         uuid.NewString(),
		Name:       tenantCreate.Name,
		Address:    tenantCreate.Address,
		Phone:      tenantCreate.Phone,
		Email:      tenantCreate.Email,
		CuitPdv:    tenantCreate.CuitPdv,
		Connection: connection,
		Identifier: identifier,
	}

	if err := tx.Create(tenant).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrInvalidData){
			return "", models.ErrorResponse(400, "Los campos email, cuit_pdv y identifier deben ser únicos, algun campo ya existe", err)
		}
		return "", models.ErrorResponse(500, "Error interno al crear tenant", err)
	}

	var user models.User
	if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", models.ErrorResponse(404, "User no encontrado", err)
		}
		return "", models.ErrorResponse(500, "Error interno al obtener user", err)
	}

	if err := tx.Create(&models.UserTenant{
		UserID:    user.ID,
		TenantID:  tenant.ID,
		IsAdmin:   true,
	}).Error; err != nil {
		tx.Rollback()
		return "", models.ErrorResponse(500, "Error interno al crear user-tenant", err)
	}

	err = database.PrepareDB(uri)
	if err != nil {
		tx.Rollback()
		return "", models.ErrorResponse(500, "Error interno al crear la base de datos del tenant", err)
	}

	if err := tx.Commit().Error; err != nil {
		return "", models.ErrorResponse(500, "Error interno al crear user-tenant", err)
	}

	return tenant.ID, nil
}

func (r *MainRepository) TenantUserCreate(tenantUserCreate *models.TenantUserCreate) (string, error) {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	tenantName := strings.ReplaceAll(tenantUserCreate.TenantCreate.Name, " ", "_")
	identifier := strings.ReplaceAll(tenantUserCreate.TenantCreate.Identifier, " ", "_")
	uri := fmt.Sprintf("%s%s_%s.db%s", os.Getenv("URI_PATH"), tenantName, identifier, os.Getenv("URI_CONFIG"))
	connection, err := utils.Encrypt(uri)
	if err != nil {
		return "", err
	}

	tenant := &models.Tenant{
		ID:         uuid.NewString(),
		Name:       tenantUserCreate.TenantCreate.Name,
		Address:    tenantUserCreate.TenantCreate.Address,
		Phone:      tenantUserCreate.TenantCreate.Phone,
		Email:      tenantUserCreate.TenantCreate.Email,
		CuitPdv:    tenantUserCreate.TenantCreate.CuitPdv,
		Connection: connection,
		Identifier: identifier,
	}

	if err := tx.Create(tenant).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrInvalidData){
			return "", models.ErrorResponse(400, "Los campos email, cuit_pdv y identifier deben ser únicos, algun campo ya existe", err)
		}
		return "", models.ErrorResponse(500, "Error interno creating tenant", err)
	}

	pass, err :=utils.HashPassword(tenantUserCreate.UserCreate.Password)
	if err != nil {
		return "", models.ErrorResponse(500, "Error interno hashed password", err)
	}

	user := &models.User{
		ID:       uuid.NewString(),
		FirstName: tenantUserCreate.UserCreate.FirstName,
		LastName:  tenantUserCreate.UserCreate.LastName,
		Username: tenantUserCreate.UserCreate.Username,
		Email:    tenantUserCreate.UserCreate.Email,
		Password: pass,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return "", models.ErrorResponse(500, "Error interno al crear usuario", err)
	}

	if err := tx.Create(&models.UserTenant{
		UserID:    user.ID,
		TenantID:  tenant.ID,
		IsAdmin:   true,
	}).Error; err != nil {
		tx.Rollback()
		return "", models.ErrorResponse(500, "Error interno al crear tenant", err)
	}

	err = database.PrepareDB(uri)
	if err != nil {
		tx.Rollback()
		return "", models.ErrorResponse(500, "Error interno al crear la base de datos del tenant", err)
	}

	if err := tx.Commit().Error; err != nil {
		return "", models.ErrorResponse(500, "Error interno al crear user-tenant", err)
	}

	return tenant.ID, nil
}

func (r *MainRepository) TenantUpdate(userID string, tenant *models.TenantUpdate) error {
	var userTenant models.UserTenant

	err := r.DB.First(&userTenant, "user_id = ? AND tenant_id = ?", userID, tenant.ID).Error
	if err != nil {
		return models.ErrorResponse(404, "User-tenant relationship not found", err)
	}

	if !userTenant.IsAdmin {
		return models.ErrorResponse(403, "No tienes permisos para actualizar el tenant", fmt.Errorf("no tienes permisos para actualizar el tenant"))
	}

	if err := r.DB.Model(&models.Tenant{}).Updates(tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrorResponse(404, "Tenant not found", err)
		}
		return models.ErrorResponse(500, "Error interno updating tenant", err)
	}

	return nil
}
