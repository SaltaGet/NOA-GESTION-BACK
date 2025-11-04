package services

import (
"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
)

func (u *UserService) UserGetByListID(ids []string) (*[]schemas.UserDTO, error) {
	users, err := u.UserRepository.UserGetByListID(ids)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserService) UserCreate(userCreate *schemas.UserCreate) (string, error) {
	existingUser, err := u.UserRepository.UserGetExistByUsernameEmail(userCreate.Username, userCreate.Email)
	if err != nil {
		return "",err
	}
	if existingUser {
		return "", err
	}

	
	userCreate.Password, err = utils.HashPassword(userCreate.Password)
	if err != nil {
		return "", err
	}

	id, err := u.UserRepository.UserCreate(userCreate)
	if err != nil {
		return "", err
	}

	return id, nil
}