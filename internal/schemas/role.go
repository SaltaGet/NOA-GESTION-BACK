package schemas

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Role struct {
	ID          string       `gorm:"primaryKey;size:36" json:"id"`
	Name        string       `gorm:"not null;unique" json:"name"`
	Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"permissions"`
}

type RoleCreate struct {
	Name string `json:"name" validate:"required"`
	PermissionsID []string `json:"permissions_id" validate:"required,dive,uuid"`
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

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type RoleUpdate struct {
	ID   string `json:"id" validate:"required,uuid"`
	Name string `json:"name" validate:"required"`
	PermissionsID []string `json:"permissions_id" validate:"required,dive,uuid"`
}

func (r *RoleUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type RoleDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoleResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Permissions []PermissionResponse `json:"permissions"`
}

type RolePermissionRow struct {
    RoleID      string
    RoleName    string
    PermID      string
    PermCode    string
    PermDetails string
    PermGroup   string
}