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
	ID               int64                    `json:"id"`
	FirstName        string                   `json:"first_name"`
	LastName         string                   `json:"last_name"`
	Username         string                   `json:"username"`
	IsAdmin          bool                     `json:"is_admin"`
	RoleID           int64                    `json:"role_id"`
	RoleName         string                   `json:"role_name"`
	Permissions      []EnvironmentPermissions `json:"permissions"`
	TenantID         int64                    `json:"tenant_id"`
	TenantName       string                   `json:"tenant_name"`
	TenantIdentifier string                   `json:"tenant_identifier"`
	ListPermissions  []string                 `json:"list_permissions"`
}

type EnvironmentPermissions struct {
	Environment string             `json:"environment"`
	Groups      []GroupPermissions `json:"groups"`
}

type GroupPermissions struct {
	Group       string   `json:"group"`
	Permissions []string `json:"permissions"`
}
