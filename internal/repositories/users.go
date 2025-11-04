package repositories

import (
	"errors"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *MainRepository) UserGetByID(id string) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Usuario no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al buscar usuario", err)
	}

	return &user, nil
}

func (r *MainRepository) UserGetByListID(ids []string) (*[]schemas.UserDTO, error) {
	var users []schemas.UserDTO

	err := r.DB.
		Model(&models.User{}).
		Select("id, first_name, last_name, username, email").
		Where("id IN ?", ids).
		Scan(&users).Error
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *MainRepository) UserGetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, schemas.ErrorResponse(500, "Error interno al buscar usuario", err)
	}

	return &user, nil
}

func (r *MainRepository) UserGetExistByUsernameEmail(username string, email string) (bool, error) {
	err := r.DB.Where("email = ? OR username = ?", email, username).First(&models.User{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, schemas.ErrorResponse(500, "Error interno al buscar usuario", err)
	}

	return true, nil
}

func (r *MainRepository) UserCreate(user *schemas.UserCreate) (string, error) {
	newID := uuid.NewString()
	err := r.DB.Create(&models.User{Username: user.Username, Email: user.Email, Password: user.Password}).Error
	if err != nil {
		return "" , schemas.ErrorResponse(500, "Error interno al crear usuario", err)
	}
	return newID, nil
}


func (m *MainRepository) UserTenantAdd(userID, tenantID int64) (err error) {
	if err := m.DB.Create(&models.UserTenant{
		UserID: userID,
		TenantID: tenantID,
	}).Error; err != nil {
		return schemas.ErrorResponse(500, "Error interno al agregar usuario a tenant", err)
	}

	return nil
}
