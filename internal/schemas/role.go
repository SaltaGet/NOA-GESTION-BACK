package schemas

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type RoleResponseDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type RoleResponse struct {
	ID          int64                `json:"id"`
	Name        string               `json:"name"`
	Permissions []PermissionResponseDTO `json:"permissions"`
}

type RolePermissionRow struct {
	RoleID      int64
	RoleName    string
	PermID      int64
	PermCode    string
	PermGroup   string
	Environment string
}

type RoleCreate struct {
	Name string `json:"name" validate:"required"`
	PermissionsID []int64 `json:"permissions_id" validate:"required,dive"`
}

func (r *RoleCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
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