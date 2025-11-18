package schemas

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type UserDTO struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

type UserCreate struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

func (u *UserCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(u)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	message := fmt.Sprintf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)

	return ErrorResponse(422, message, fmt.Errorf("%s", message))
}

type UserUpdate struct {
	ID        int64  `json:"id" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

func (u *UserUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(u)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	message := fmt.Sprintf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)

	return ErrorResponse(422, message, fmt.Errorf("%s", message))
}
