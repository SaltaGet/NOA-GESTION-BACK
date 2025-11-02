package ports

import "github.com/DanielChachagua/GestionCar/pkg/models"

type UserRepository interface {
	// Update(user *user.UserUpdate) (err error)
	UserGetByID(id string) (user *models.User, err error) 
	UserGetByListID(ids []string) (users *[]models.UserDTO, err error)
	UserGetByUsername(username string) (user *models.User, err error) 
	UserGetExistByUsernameEmail(username string, email string) (exist bool, err error) 
	// UserGetByIdentifier(identifier string) (user *models.User, err error) 
	// UserGetByEmail(email string) (user *models.User, err error) 
	// UserExist(identifier string, email string) (exist bool, err error) 
	UserCreate(user *models.UserCreate) (id string, err error)
	UserTenantAdd(userID, tenantID string) (err error)
}

type UserService interface {
	UserGetByListID(ids []string) (users *[]models.UserDTO, err error)
	UserCreate(user *models.UserCreate) (id string, err error)
	// Update(user *models.UserUpdate) (err error)
}