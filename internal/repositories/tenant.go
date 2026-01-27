package repositories

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

//		return &tenant, nil
//	}
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
                tenants.email, tenants.is_active, tenants.expiration,
                tenants.created_at, tenants.updated_at`).
		Joins("JOIN user_tenants ON tenants.id = user_tenants.tenant_id").
		Scan(&tenants).Error

	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno retrieving tenants", err)
	}

	return &tenants, nil
}

func (r *MainRepository) TenantGetConnectionByIdentifier(tenantIdentifier string) (*models.Tenant, error) {
	var tenant *models.Tenant
	err := r.DB.Select("id", "connection").Where("identifier = ?", tenantIdentifier).First(&tenant).Error
	if err != nil {
		return nil, err
	}

	return tenant, nil
}

func (r *MainRepository) TenantGetConections() ([]*models.Tenant, error) {
	var tenants []*models.Tenant
	if err := r.DB.Debug().
		Model(&models.Tenant{}).
		Select("id", "name", "connection", "identifier").
		Find(&tenants).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al obtener las connections", err)
	}

	for _, tenant := range tenants {
		connectionDecrypted, err := utils.Decrypt(tenant.Connection)
		if err != nil {
			return nil, schemas.ErrorResponse(500, "Error interno al obtener las connections", err)
		}
		tenant.Connection = connectionDecrypted
	}

	return tenants, nil
}

// func (r *MainRepository) TenantCreateByUserID(tenantCreate *schemas.TenantCreate, userID int64) (int64, error) {
// 	tx := r.DB.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	tenantName := strings.ReplaceAll(tenantCreate.Name, " ", "_")
// 	identifier := strings.ReplaceAll(tenantCreate.Identifier, " ", "_")
// 	uri := fmt.Sprintf("%s%s_%s%s", os.Getenv("URI_PATH"), tenantName, identifier, os.Getenv("URI_CONFIG"))
// 	connection, err := utils.Encrypt(uri)
// 	if err != nil {
// 		return 0, schemas.ErrorResponse(500, "Error interno al obtener connection", err)
// 	}

// 	tenant := &models.Tenant{
// 		Name:       tenantCreate.Name,
// 		Address:    tenantCreate.Address,
// 		Phone:      tenantCreate.Phone,
// 		Email:      tenantCreate.Email,
// 		CuitPdv:    tenantCreate.CuitPdv,
// 		Connection: connection,
// 		Identifier: identifier,
// 	}

// 	if err := tx.Create(tenant).Error; err != nil {
// 		tx.Rollback()
// 		if errors.Is(err, gorm.ErrInvalidData) {
// 			return 0, schemas.ErrorResponse(400, "Los campos email, cuit_pdv y identifier deben ser únicos, algun campo ya existe", err)
// 		}
// 		return 0, schemas.ErrorResponse(500, "Error interno al crear tenant", err)
// 	}

// 	var user models.User
// 	if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
// 		tx.Rollback()
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return 0, schemas.ErrorResponse(404, "User no encontrado", err)
// 		}
// 		return 0, schemas.ErrorResponse(500, "Error interno al obtener user", err)
// 	}

// 	if err := tx.Create(&models.UserTenant{
// 		UserID:   user.ID,
// 		TenantID: tenant.ID,
// 	}).Error; err != nil {
// 		tx.Rollback()
// 		return 0, schemas.ErrorResponse(500, "Error interno al crear user-tenant", err)
// 	}

// 	// generar pass generic ***

// 	memberAdmin := &models.Member{
// 		FirstName: user.FirstName,
// 		LastName:  user.LastName,
// 		Username:  user.Username,
// 		Email:     user.Email,
// 		Password:  "1",
// 		IsAdmin:   true,
// 		Address:   user.Address,
// 		RoleID:    1,
// 	}

// 	// enviar email ***

// 	err = database.PrepareDB(uri, *memberAdmin)
// 	if err != nil {
// 		tx.Rollback()
// 		return 0, schemas.ErrorResponse(500, "Error interno al crear la base de datos del tenant", err)
// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		return 0, schemas.ErrorResponse(500, "Error interno al crear user-tenant", err)
// 	}

// 	return tenant.ID, nil
// }

func (r *MainRepository) TenantCreateByUserID(adminID int64, tenantCreate *schemas.TenantCreate, userID int64) (int64, error) {
	var tenantID int64

	err := r.DB.Transaction(func(tx *gorm.DB) error {

		// Normalizar strings
		tenantName := strings.ReplaceAll(tenantCreate.Name, " ", "_")
		identifier := strings.ReplaceAll(tenantCreate.Identifier, " ", "_")

		uri := fmt.Sprintf("%s%s_%s%s",
			os.Getenv("URI_PATH"),
			tenantName,
			identifier,
			os.Getenv("URI_CONFIG"),
		)

		connection, err := utils.Encrypt(uri)
		if err != nil {
			return schemas.ErrorResponse(500, "Error interno al obtener connection", err)
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

		// Crear tenant
		if err := tx.Create(tenant).Error; err != nil {
			if errors.Is(err, gorm.ErrInvalidData) {
				return schemas.ErrorResponse(400,
					"Los campos email, cuit_pdv y identifier deben ser únicos, algun campo ya existe",
					err)
			}
			return schemas.ErrorResponse(500, "Error interno al crear tenant", err)
		}

		tenantID = tenant.ID

		// Buscar usuario
		var user models.User
		if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "User no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error interno al obtener user", err)
		}

		// Crear relación user-tenant
		if err := tx.Create(&models.UserTenant{
			UserID:   user.ID,
			TenantID: tenant.ID,
		}).Error; err != nil {
			return schemas.ErrorResponse(500, "Error interno al crear user-tenant", err)
		}

		// Crear miembro administrador genérico
		memberAdmin := &models.Member{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			Email:     user.Email,
			Password:  "1",
			IsAdmin:   true,
			Address:   user.Address,
			RoleID:    1,
		}

		// Crear DB del tenant
		err = database.PrepareDB(uri, *memberAdmin)
		if err != nil {
			return schemas.ErrorResponse(500, "Error interno al crear la base de datos del tenant", err)
		}

		return nil // todo ok, commit automático
	})

	if err != nil {
		return 0, err // devuelve error del Transaction
	}

	go database.SaveAuditAdminAsync(r.DB, models.AuditLogAdmin{
		AdminID: adminID,
		Method:  "create",
		Path:    "plan",
	}, nil, tenantCreate)

	return tenantID, nil
}

// func (r *MainRepository) TenantUserCreate(tenantUserCreate *schemas.TenantUserCreate) (int64, error) {
// 	tx := r.DB.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	tenantName := strings.ReplaceAll(tenantUserCreate.TenantCreate.Name, " ", "_")
// 	identifier := strings.ReplaceAll(tenantUserCreate.TenantCreate.Identifier, " ", "_")
// 	uri := fmt.Sprintf("%s%s_%s%s", os.Getenv("URI_PATH"), tenantName, identifier, os.Getenv("URI_CONFIG"))
// 	connection, err := utils.Encrypt(uri)
// 	if err != nil {
// 		return 0, err
// 	}

// 	tenant := &models.Tenant{
// 		Name:       tenantUserCreate.TenantCreate.Name,
// 		Address:    tenantUserCreate.TenantCreate.Address,
// 		Phone:      tenantUserCreate.TenantCreate.Phone,
// 		Email:      tenantUserCreate.TenantCreate.Email,
// 		CuitPdv:    tenantUserCreate.TenantCreate.CuitPdv,
// 		Connection: connection,
// 		Identifier: identifier,
// 		PlanID:     tenantUserCreate.TenantCreate.PlanID,
// 	}

// 	if err := tx.Create(tenant).Error; err != nil {
// 		tx.Rollback()
// 		if errors.Is(err, gorm.ErrInvalidData) {
// 			return 0, schemas.ErrorResponse(400, "Los campos email, cuit_pdv y identifier deben ser únicos, algun campo ya existe", err)
// 		}
// 		if schemas.IsDuplicateError(err) {
// 			if strings.Contains(err.Error(), "email") {
// 				return 0, schemas.ErrorResponse(409, "El email del tenant ya existe", err)
// 			} else if strings.Contains(err.Error(), "identifier") {
// 				return 0, schemas.ErrorResponse(409, "El identificador del tenant ya existe", err)
// 			} else if strings.Contains(err.Error(), "cuit_pdv") {
// 				return 0, schemas.ErrorResponse(409, "El cuit del tenant ya existe", err)
// 			}
// 		}
// 		return 0, schemas.ErrorResponse(500, "Error interno creating tenant", err)
// 	}

// 	user := &models.User{
// 		FirstName: tenantUserCreate.UserCreate.FirstName,
// 		LastName:  tenantUserCreate.UserCreate.LastName,
// 		Email:     tenantUserCreate.UserCreate.Email,
// 		Address:   &tenantUserCreate.TenantCreate.Address,
// 		Username:  tenantUserCreate.UserCreate.Username,
// 	}

// 	if err := tx.Create(user).Error; err != nil {
// 		tx.Rollback()
// 		if errors.Is(err, gorm.ErrInvalidData) {
// 			if schemas.IsDuplicateError(err) {
// 				if strings.Contains(err.Error(), "email") {
// 					return 0, schemas.ErrorResponse(409, "El email del usuario ya existe", err)
// 				} else if strings.Contains(err.Error(), "username") {
// 					return 0, schemas.ErrorResponse(409, "El username del usuario ya existe", err)
// 				}
// 			}
// 			return 0, schemas.ErrorResponse(400, "Los campos email e identifier deben ser únicos, algun campo ya existe", err)
// 		}
// 		return 0, schemas.ErrorResponse(500, "Error interno creating tenant", err)
// 	}

// 	if err := tx.Create(&models.UserTenant{
// 		UserID:   user.ID,
// 		TenantID: tenant.ID,
// 	}).Error; err != nil {
// 		tx.Rollback()
// 		return 0, schemas.ErrorResponse(500, "Error interno al crear tenant", err)
// 	}

// 	memberAdmin := &models.Member{
// 		FirstName: tenantUserCreate.UserCreate.FirstName,
// 		LastName:  tenantUserCreate.UserCreate.LastName,
// 		Username:  tenantUserCreate.UserCreate.Username,
// 		Email:     tenantUserCreate.UserCreate.Email,
// 		Password:  tenantUserCreate.UserCreate.Password,
// 		IsAdmin:   true,
// 		Address:   &tenantUserCreate.TenantCreate.Address,
// 		RoleID:    1,
// 	}

// 	err = database.PrepareDB(uri, *memberAdmin)
// 	if err != nil {
// 		tx.Rollback()
// 		return 0, schemas.ErrorResponse(500, "Error interno al crear la base de datos del tenant", err)
// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		return 0, schemas.ErrorResponse(500, "Error interno al crear user-tenant", err)
// 	}

//		return tenant.ID, nil
//	}
func (r *MainRepository) TenantUserCreate(adminID int64, tenantUserCreate *schemas.TenantUserCreate) (int64, error) {
	var tenantID int64
	err := r.DB.Transaction(func(tx *gorm.DB) error {

		tenantName := strings.ReplaceAll(tenantUserCreate.TenantCreate.Name, " ", "_")
		identifier := strings.ReplaceAll(tenantUserCreate.TenantCreate.Identifier, " ", "_")

		uri := fmt.Sprintf(
			"%s%s_%s%s",
			os.Getenv("URI_PATH"),
			tenantName,
			identifier,
			os.Getenv("URI_CONFIG"),
		)

		connection, err := utils.Encrypt(uri)
		if err != nil {
			return err
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

		// CREAR TENANT
		if err := tx.Create(tenant).Error; err != nil {
			if errors.Is(err, gorm.ErrInvalidData) {
				return schemas.ErrorResponse(400,
					"Los campos email, cuit_pdv y identifier deben ser únicos. Algún campo ya existe", err)
			}
			if schemas.IsDuplicateError(err) {
				switch {
				case strings.Contains(err.Error(), "email"):
					return schemas.ErrorResponse(409, "El email del tenant ya existe", err)
				case strings.Contains(err.Error(), "identifier"):
					return schemas.ErrorResponse(409, "El identificador del tenant ya existe", err)
				case strings.Contains(err.Error(), "cuit_pdv"):
					return schemas.ErrorResponse(409, "El cuit del tenant ya existe", err)
				}
			}
			if strings.Contains(err.Error(), "identifier invalid") {
				return schemas.ErrorResponse(409, "el identificador solo puede contener letras minúsculas, números y guiones, no puede contener espacios", err)
			}
			return schemas.ErrorResponse(500, "Error interno creando tenant", err)
		}

		tenantID = tenant.ID

		// CREAR USER
		user := &models.User{
			FirstName: tenantUserCreate.UserCreate.FirstName,
			LastName:  tenantUserCreate.UserCreate.LastName,
			Email:     tenantUserCreate.UserCreate.Email,
			Address:   &tenantUserCreate.TenantCreate.Address,
			Username:  tenantUserCreate.UserCreate.Username,
		}

		if err := tx.Create(user).Error; err != nil {
			if schemas.IsDuplicateError(err) {
				switch {
				case strings.Contains(err.Error(), "email"):
					return schemas.ErrorResponse(409, "El email del usuario ya existe", err)
				case strings.Contains(err.Error(), "username"):
					return schemas.ErrorResponse(409, "El username del usuario ya existe", err)
				}
			}
			return schemas.ErrorResponse(500, "Error interno creando usuario", err)
		}

		// CREAR RELACIÓN USER-TENANT
		if err := tx.Create(&models.UserTenant{
			UserID:   user.ID,
			TenantID: tenant.ID,
		}).Error; err != nil {
			return schemas.ErrorResponse(500, "Error interno al crear tenant", err)
		}

		// CREAR ADMIN DEL TENANT
		memberAdmin := &models.Member{
			FirstName: tenantUserCreate.UserCreate.FirstName,
			LastName:  tenantUserCreate.UserCreate.LastName,
			Username:  tenantUserCreate.UserCreate.Username,
			Email:     tenantUserCreate.UserCreate.Email,
			Password:  tenantUserCreate.UserCreate.Password,
			IsAdmin:   true,
			Address:   &tenantUserCreate.TenantCreate.Address,
			RoleID:    1,
		}

		if err := database.PrepareDB(uri, *memberAdmin); err != nil {
			return schemas.ErrorResponse(500, "Error interno al crear la base de datos del tenant", err)
		}

		return nil // commit automático
	})

	if err != nil {
		return 0, err
	}

	tenantUserCreate.UserCreate.Password = ""
	go database.SaveAuditAdminAsync(r.DB, models.AuditLogAdmin{
		AdminID: adminID,
		Method:  "create",
		Path:    "plan",
	}, nil, tenantUserCreate)

	return tenantID, nil
}

func (r *MainRepository) TenantUpdate(adminID, userID int64, tenant *schemas.TenantUpdate) error {
	var userTenant models.UserTenant

	err := r.DB.First(&userTenant, "user_id = ? AND tenant_id = ?", userID, tenant.ID).Error
	if err != nil {
		return schemas.ErrorResponse(404, "User-tenant relationship not found", err)
	}

	// if !userTenant.IsAdmin {
	// 	return schemas.ErrorResponse(403, "No tienes permisos para actualizar el tenant", fmt.Errorf("no tienes permisos para actualizar el tenant"))
	// }

	var tenantOld models.Tenant
	var tenantUpdates models.Tenant
	err = r.DB.First(&tenantOld, tenant.ID).Error
	if err != nil {
		return schemas.ErrorResponse(404, "Tenant not found", err)
	}

	if err := r.DB.Model(&models.Tenant{}).Updates(tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "Tenant not found", err)
		}
		return schemas.ErrorResponse(500, "Error interno updating tenant", err)
	}

	r.DB.First(&tenantUpdates, tenant.ID)

	go database.SaveAuditAdminAsync(r.DB, models.AuditLogAdmin{
		AdminID: adminID,
		Method:  "create",
		Path:    "plan",
	}, tenantOld, tenantUpdates)

	return nil
}

func (r *MainRepository) TenantUpdateExpiration(adminID int64, tenantUpdateExpiration *schemas.TenantUpdateExpiration) error {
	loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
	exp, err := time.ParseInLocation("2006-01-02", tenantUpdateExpiration.Expiration, loc)
	if err != nil {
		return schemas.ErrorResponse(422, "Formato de fecha inválido, debe ser YYYY-MM-DD", err)
	}

	var tenantExist models.Tenant
	var tenantSave models.Tenant
	err = r.DB.First(&tenantExist, tenantUpdateExpiration.ID).Error
	if err != nil {
		return schemas.ErrorResponse(404, "Tenant not found", err)
	}

	tenantSave = tenantExist
	tenantExist.Expiration = &exp

	if err := r.DB.Save(&tenantExist).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "Tenant not found", err)
		}
		return schemas.ErrorResponse(500, "Error interno updating tenant", err)
	}

	go database.SaveAuditAdminAsync(r.DB, models.AuditLogAdmin{
		AdminID: adminID,
		Method:  "create",
		Path:    "plan",
	}, tenantSave, tenantExist)

	return nil
}

func (r *MainRepository) TenantUpdateTerms(tenantID int64, tenantUpdateTerms *schemas.TenantUpdateTerms) error {
	// Opción recomendada: Usar Select para forzar la actualización de estos campos específicos
	err := r.DB.Model(&models.Tenant{}).
		Where("id = ?", tenantID).
		Select("AcceptedTerms", "IP", "DateAccepted"). // Asegúrate de que los nombres coincidan con el modelo Tenant
		Updates(tenantUpdateTerms).Error
	if err != nil {
		return schemas.ErrorResponse(500, "Error interno al actualizar terminos tenant", err)
	}

	return nil
}

func (r *MainRepository) TenantGetSettings(tenantID int64) (*schemas.TenantSettingsResponse, error) {
	var settings models.SettingTenant
	if err := r.DB.Where("tenant_id = ?", tenantID).First(&settings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Tenant not found", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno retrieving tenant", err)
	}

	response := schemas.TenantSettingsResponse{
		Logo:           settings.Logo,
		FrontPage:      settings.FrontPage,
		Title:          settings.Title,
		Slogan:         settings.Slogan,
		PrimaryColor:   settings.PrimaryColor,
		SecondaryColor: settings.SecondaryColor,
		Phone:          settings.Phone,
	}

	return &response, nil
}

func (r *MainRepository) TenantUpdateSettings(tenantID int64, tenantUpdateSettings *schemas.TenantUpdateSettings) error {
	// 1. Mapeamos al modelo real de la base de datos
	settings := models.SettingTenant{
		TenantID:       tenantID,
		Title:          tenantUpdateSettings.Title,
		Slogan:         tenantUpdateSettings.Slogan,
		PrimaryColor:   tenantUpdateSettings.PrimaryColor,
		SecondaryColor: tenantUpdateSettings.SecondaryColor,
		Phone: tenantUpdateSettings.Phone,
	}

	// 2. Usamos Clauses para definir qué pasa si hay un conflicto en el tenant_id
	err := r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tenant_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "slogan", "primary_color", "secondary_color", "phone", "updated_at"}),
	}).Create(&settings).Error

	if err != nil {
		return schemas.ErrorResponse(500, "Error al crear/actualizar configuraciones", err)
	}

	return nil
}
