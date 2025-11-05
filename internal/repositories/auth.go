package repositories

import (
	"errors"
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"gorm.io/gorm"
)

// func (r *MainRepository) AuthUserGetByID(userID int64) (*models.User, error) {
// 	var user models.User
// 	if err := r.DB.Where("id = ?", userID).First(&user).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
// 		}
// 		return nil, schemas.ErrorResponse(500, "Error al intentar loguearse", err)
// 	}
// 	return &user, nil
// }

// func (r *MainRepository) AuthUserGetByUsername(username string) (*models.User, error) {
// 	var user models.User
// 	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
// 		}
// 		return nil, schemas.ErrorResponse(500, "Error al intentar loguearse", err)
// 	}
// 	return &user, nil
// }

func (r *MainRepository) AuthTenantGetByID(tenantID int64) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.DB.
		Where("id = ?", tenantID).
		First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al obtener los tenants", err)
	}

	if !tenant.IsActive {
		return nil, schemas.ErrorResponse(403, "Tenant esta inactivo", fmt.Errorf("credenciales incorrectas"))
	}

	// if len(tenant.UserTenants) == 0 {
	// 	return nil, schemas.ErrorResponse(403, "No tienes permiso para acceder al tenant", fmt.Errorf("sin permiso para acceder al tenant"))
	// }

	// if !tenant.UserTenants[0].IsActive {
	// 	return nil, schemas.ErrorResponse(403, "No tienes permiso para acceder al tenant", fmt.Errorf("sin permiso para acceder al tenant"))
	// }

	return &tenant, nil
}

func (r *MainRepository) AuthTenantGetByIdentifier(identifier string) (*models.Tenant, error) {
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
		return nil, schemas.ErrorResponse(403, "Tenant esta inactivo", fmt.Errorf("credenciales incorrectas"))
	}

	// if len(tenant.UserTenants) == 0 {
	// 	return nil, schemas.ErrorResponse(403, "No tienes permiso para acceder al tenant", fmt.Errorf("sin permiso para acceder al tenant"))
	// }

	// if !tenant.UserTenants[0].IsActive {
	// 	return nil, schemas.ErrorResponse(403, "No tienes permiso para acceder al tenant", fmt.Errorf("sin permiso para acceder al tenant"))
	// }

	return &tenant, nil
}

func (r *MainRepository) AuthMemberGetByUserID(userID int64, connection string, tenantID int64) (*models.Member, error) {
	db, err := database.GetTenantDB(connection, tenantID)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al recibir la conexión", err)
	}

	var member models.Member
	if err := db.Where("user_id = ?", userID).First(&member).Error; err != nil {
		return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
	}

	if !member.IsActive {
		return nil, schemas.ErrorResponse(403, "Miembro inactivo", fmt.Errorf("miembro inactivo"))
	}

	return &member, nil
}

func (r *MainRepository) AuthMemberGetByID(id int64, connection string, tenantID int64) (*models.Member, *[]string, error) {
	// db, err := database.GetTenantDB(connection)
	db, err := database.GetTenantDB(connection, tenantID)
	if err != nil {
		return nil, nil, schemas.ErrorResponse(500, "Error al recibir la conexión", err)
	}

	var member models.Member
	if err := db.Where("id = ?", id).First(&member).Error; err != nil {
		return nil, nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
	}

	if !member.IsActive {
		return nil, nil, schemas.ErrorResponse(403, "Miembro inactivo", fmt.Errorf("miembro inactivo"))
	}

	perm := make([]string, len(member.Role.Permissions))
	for i, p := range member.Role.Permissions {
		perm[i] = p.Code
	}

	return &member, &perm, nil
}

func (r *MainRepository) AuthMemberGetByUsername(username string, connection string, tenantID int64) (*models.Member, error) {
	// db, err := database.GetTenantDB(connection)
	db, err := database.GetTenantDB(connection, tenantID)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al recibir la conexión", err)
	}

	var member models.Member
	if err := db.Where("username = ?", username).First(&member).Error; err != nil {
		return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
	}

	if !member.IsActive {
		return nil, schemas.ErrorResponse(403, "Miembro inactivo", fmt.Errorf("miembro inactivo"))
	}

	return &member, nil
}

func (r *MainRepository) AuthAdminGetByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	if err := r.DB.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
	}

	return &admin, nil
}

func (r *MainRepository) AuthAdminGetByID(id int64) (*models.Admin, error) {
	var admin models.Admin
	if err := r.DB.Where("id = ?", id).First(&admin).Error; err != nil {
		return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
	}

	return &admin, nil
}

func (r *MainRepository) AuthPointSale(pointSaleID int64, connection string, tenantID, memberID int64) (*models.PointSale, error) {
	db, err := database.GetTenantDB(connection, tenantID)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al recibir la conexión", err)
	}

	var pointSale models.PointSale
	err = db.Model(&models.PointSale{}).
		Joins("JOIN member_point_sales mp ON mp.point_sale_id = point_sales.id").
		Where("mp.member_id = ?", memberID).
		Where("id = ?", pointSaleID).
		First(&pointSale).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(401, "Credenciales incorrectas", err)
		}
		return nil, schemas.ErrorResponse(500, "Error al obtener tenant", err)
	}

	return &pointSale, nil
}
