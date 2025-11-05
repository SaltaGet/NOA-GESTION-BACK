package ports

import (
	// "github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	// "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)


type UserRepository interface {
	// Update(user *user.UserUpdate) (err error)
	// UserGetByID(id string) (user *models.User, err error) 
	// UserGetByListID(ids []string) (users *[]schemas.UserDTO, err error)
	// UserGetByUsername(username string) (user *models.User, err error) 
	// UserGetExistByUsernameEmail(username string, email string) (exist bool, err error) 
	// // UserGetByIdentifier(identifier string) (user *schemas.User, err error) 
	// // UserGetByEmail(email string) (user *schemas.User, err error) 
	// // UserExist(identifier string, email string) (exist bool, err error) 
	// UserCreate(user *schemas.UserCreate) (id string, err error)
	// UserTenantAdd(userID, tenantID int64) (err error)
}

type UserService interface {
	// UserGetByListID(ids []string) (users *[]schemas.UserDTO, err error)
	// UserCreate(user *schemas.UserCreate) (id string, err error)
	// Update(user *schemas.UserUpdate) (err error)
}