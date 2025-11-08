package schemas

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type AuthLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (a *AuthLogin) Validate() error {
	validate := validator.New()
	err := validate.Struct(a)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type AuthLoginAdmin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (a *AuthLoginAdmin) Validate() error {
	validate := validator.New()
	err := validate.Struct(a)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type AuthenticatedUser struct {
	ID               int64    `json:"id"`
	FirstName        string   `json:"first_name"`
	LastName         string   `json:"last_name"`
	Username         string   `json:"username"`
	IsAdmin          bool     `json:"is_admin"`
	RoleID           int64    `json:"role_id,omitempty"`
	RoleName         string   `json:"role_name,omitempty"`
	Permissions      []string `json:"permissions,omitempty"`
	TenantID         int64    `json:"tenant_id,omitempty"`
	TenantName       string   `json:"tenant_name,omitempty"`
	TenantIdentifier string   `json:"tenant_identifier,omitempty"`
}

type AuthPointSaleContext struct {
	ID   int64   `json:"id"`
	Name string `json:"name"`
}