package schemas

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type TypeMovementCreate struct {
	Name         string `json:"name" validate:"required"`
	TypeMovement string `json:"type_movement" validate:"required,oneof=income expense"`
}

func (t *TypeMovementCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(t)
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

type TypeMovementUpdate struct {
	ID           int64  `json:"id" validate:"required"`
	Name         string `json:"name" validate:"required"`
	TypeMovement string `json:"type_movement" validate:"required,oneof=income expense"`
}

func (t *TypeMovementUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(t)
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

type TypeMovementResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
