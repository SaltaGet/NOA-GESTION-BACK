package repositories

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	"gorm.io/gorm"
)

// func (r *MainRepository) TenantGetByID(tenantID string) (*schemas.Tenant, error) {
// 	var tenant schemas.Tenant
// 	err := r.DB.
// 		Where("id = ?", tenantID).
// 		First(&tenant).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, schemas.ErrorResponse(404, "Tenant not found", err)
// 		}
// 		return nil, schemas.ErrorResponse(500, "Error retrieving tenant", err)
// 	}

// 	if !tenant.IsActive {
// 		return nil, schemas.ErrorResponse(403, "Tenant is inactive", nil)
// 	}

// 	return &tenant, nil
// }
func (r *MainRepository) TenantGetByID(tenantID int64) (*models.Tenant, error) {
    var tenant models.Tenant
    err := r.DB.
        Where("id = ?", tenantID).
        First(&tenant).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, schemas.ErrorResponse(404, "Tenant not found", err)
        }
        return nil, schemas.ErrorResponse(500, "Error interno retrieving tenant", err)
    }

    if !tenant.IsActive {
        return nil, schemas.ErrorResponse(403, "Tenant is inactive", fmt.Errorf("tenant is inactive"))
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
			return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al obtener los tenants", err)
	}

	if !tenant.IsActive {
		return nil, schemas.ErrorResponse(403, "Tenant es inactivo", fmt.Errorf("tenant es inactivo"))
	}

	return &tenant, nil
}

func (r *MainRepository) TenantGetAll() (*[]schemas.TenantResponse, error) {
	var tenants []schemas.TenantResponse

	err := r.DB.
		Model(&models.Tenant{}).
		Select(`tenants.id, tenants.name, tenants.address, tenants.phone,
                tenants.email, tenants.is_active,
                tenants.created_at, tenants.updated_at`).
		Joins("JOIN user_tenants ON tenants.id = user_tenants.tenant_id").
		Scan(&tenants).Error

	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno retrieving tenants", err)
	}

	return &tenants, nil
}

func (r *MainRepository) TenantGetConections() (*[]string, error) {
	var connections []string
	if err := r.DB.Debug().
	Model(&models.Tenant{}).
	Select("connection").
	Scan(&connections).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al obtener las connections", err)
	}

	var connectionsDecrypted []string

	for _, connection := range connections {
		connectionDecrypted, err := utils.Decrypt(connection)
		if err != nil {
			return nil, schemas.ErrorResponse(500, "Error interno al obtener las connections", err)
		}
		connectionsDecrypted = append(connectionsDecrypted, connectionDecrypted)
	}

	return &connectionsDecrypted, nil
}

func (r *MainRepository) TenantCreateByUserID(tenantCreate *schemas.TenantCreate, userID int64) (int64, error) {
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
		return 0, schemas.ErrorResponse(500, "Error interno al obtener connection", err)
	}

	tenant := &models.Tenant{
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
			return 0, schemas.ErrorResponse(400, "Los campos email, cuit_pdv y identifier deben ser únicos, algun campo ya existe", err)
		}
		return 0, schemas.ErrorResponse(500, "Error interno al crear tenant", err)
	}

	var user models.User
	if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, schemas.ErrorResponse(404, "User no encontrado", err)
		}
		return 0, schemas.ErrorResponse(500, "Error interno al obtener user", err)
	}

	if err := tx.Create(&models.UserTenant{
		UserID:    user.ID,
		TenantID:  tenant.ID,
	}).Error; err != nil {
		tx.Rollback()
		return 0, schemas.ErrorResponse(500, "Error interno al crear user-tenant", err)
	}

	// generar pass generic ***

	memberAdmin := &models.Member{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username: user.Username,
		Email:    user.Email,
		Password: "1",
		IsAdmin:  true,
		Address: user.Address,
		RoleID: 1,
	}

	// enviar email ***

	err = database.PrepareDB(uri, *memberAdmin)
	if err != nil {
		tx.Rollback()
		return 0, schemas.ErrorResponse(500, "Error interno al crear la base de datos del tenant", err)
	}

	if err := tx.Commit().Error; err != nil {
		return 0, schemas.ErrorResponse(500, "Error interno al crear user-tenant", err)
	}

	return tenant.ID, nil
}

func (r *MainRepository) TenantUserCreate(tenantUserCreate *schemas.TenantUserCreate) (int64, error) {
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
		return 0, err
	}

	tenant := &models.Tenant{
		Name:       tenantUserCreate.TenantCreate.Name,
		Address:    tenantUserCreate.TenantCreate.Address,
		Phone:      tenantUserCreate.TenantCreate.Phone,
		Email:      tenantUserCreate.TenantCreate.Email,
		CuitPdv:    tenantUserCreate.TenantCreate.CuitPdv,
		Connection: connection,
		Identifier: identifier,
		PlanID:     tenantUserCreate.TenantCreate.PlanID,
	}

	if err := tx.Create(tenant).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrInvalidData){
			return 0, schemas.ErrorResponse(400, "Los campos email, cuit_pdv y identifier deben ser únicos, algun campo ya existe", err)
		}
		if schemas.IsDuplicateError(err) {
			if strings.Contains(err.Error(), "email") {
				return 0, schemas.ErrorResponse(409, "El email del tenant ya existe", err)
			} else if strings.Contains(err.Error(), "identifier") {
				return 0, schemas.ErrorResponse(409, "El identificador del tenant ya existe", err)
			}	else if strings.Contains(err.Error(), "cuit_pdv") {
				return 0, schemas.ErrorResponse(409, "El cuit del tenant ya existe", err)
			}
		}
		return 0, schemas.ErrorResponse(500, "Error interno creating tenant", err)
	}

	user := &models.User{
		FirstName: tenantUserCreate.UserCreate.FirstName,
		LastName:  tenantUserCreate.UserCreate.LastName,
		Email:    tenantUserCreate.UserCreate.Email,
		Address: &tenantUserCreate.TenantCreate.Address,
		Username: tenantUserCreate.UserCreate.Username,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrInvalidData){
			if schemas.IsDuplicateError(err) {
				if strings.Contains(err.Error(), "email") {
					return 0, schemas.ErrorResponse(409, "El email del usuario ya existe", err)
				} else if strings.Contains(err.Error(), "username") {
					return 0, schemas.ErrorResponse(409, "El username del usuario ya existe", err)
				}
			}
			return 0, schemas.ErrorResponse(400, "Los campos email e identifier deben ser únicos, algun campo ya existe", err)
		}
		return 0, schemas.ErrorResponse(500, "Error interno creating tenant", err)
	}

	if err := tx.Create(&models.UserTenant{
		UserID:    user.ID,
		TenantID:  tenant.ID,
	}).Error; err != nil {
		tx.Rollback()
		return 0, schemas.ErrorResponse(500, "Error interno al crear tenant", err)
	}

	memberAdmin := &models.Member{
		FirstName: tenantUserCreate.UserCreate.FirstName,
		LastName:  tenantUserCreate.UserCreate.LastName,
		Username: tenantUserCreate.UserCreate.Username,
		Email:    tenantUserCreate.UserCreate.Email,
		Password: tenantUserCreate.UserCreate.Password,
		IsAdmin:  true,
		Address: &tenantUserCreate.TenantCreate.Address,
		RoleID: 1,
	}

	err = database.PrepareDB(uri, *memberAdmin)
	if err != nil {
		tx.Rollback()
		return 0, schemas.ErrorResponse(500, "Error interno al crear la base de datos del tenant", err)
	}

	if err := tx.Commit().Error; err != nil {
		return 0, schemas.ErrorResponse(500, "Error interno al crear user-tenant", err)
	}

	return tenant.ID, nil
}

func (r *MainRepository) TenantUpdate(userID int64, tenant *schemas.TenantUpdate) error {
	// var userTenant models.UserTenant

	// err := r.DB.First(&userTenant, "user_id = ? AND tenant_id = ?", userID, tenant.ID).Error
	// if err != nil {
	// 	return schemas.ErrorResponse(404, "User-tenant relationship not found", err)
	// }

	// // if !userTenant.IsAdmin {
	// // 	return schemas.ErrorResponse(403, "No tienes permisos para actualizar el tenant", fmt.Errorf("no tienes permisos para actualizar el tenant"))
	// // }

	// if err := r.DB.Model(&models.Tenant{}).Updates(tenant).Error; err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		return schemas.ErrorResponse(404, "Tenant not found", err)
	// 	}
	// 	return schemas.ErrorResponse(500, "Error interno updating tenant", err)
	// }

	return nil
}
