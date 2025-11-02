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

type AuthResult struct {
	ID         string
	FirstName  string
	LastName   string
	Username   string
	IsAdmin    bool
	Tenant     *Tenant
	Role       *Role
	Permissions *[]string
}

type AuthenticatedUser struct {
	ID            string   `json:"id"`
	FirstName     string   `json:"first_name"`
	LastName      string   `json:"last_name"`
	Username      string   `json:"username"`
	IsAdmin       bool     `json:"is_admin"`
	IsAdminTenant bool     `json:"is_admin_tenant"`
	RoleID        *string  `json:"role_id"`
	RoleName      *string  `json:"role_name"`
	Permissions   []string `json:"permissions"`
	TenantID      *string  `json:"tenant_id"`
	TenantName    *string  `json:"tenant_name"`
	Identifier    *string  `json:"identifier"`
}
