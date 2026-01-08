package schemas

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type AuthLogin struct {
	Username string `json:"username" validate:"required,username"`
	Password string `json:"password" validate:"required"`
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return regexp.MustCompile(`^[^@]+@[^.]+$`).MatchString(username)
}

func (a *AuthLogin) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("username", validateUsername)

	err := validate.Struct(a)
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

	message := fmt.Sprintf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)

	return ErrorResponse(422, message, fmt.Errorf("%s", message))
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
	AcceptedTerms    bool                     `json:"accepted_terms"`
}

type EnvironmentPermissions struct {
	Environment string             `json:"environment"`
	Groups      []GroupPermissions `json:"groups"`
}

type GroupPermissions struct {
	Group       string   `json:"group"`
	Permissions []string `json:"permissions"`
}

type AuthForgotPassword struct {
	Username         string `json:"username" validate:"required"`
	TenantIdentifier string `json:"tenant_identifier" validate:"required"`
}

func (f *AuthForgotPassword) Validate() error {
	validate := validator.New()
	err := validate.Struct(f)
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

type AuthResetPassword struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,password" example:"Password123*"`
	ConfirmPass string `json:"confirm_pass" validate:"required,password" example:"Password123*"`
}

func (f *AuthResetPassword) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("password", validatePassword) // registrar antes de Struct()

	// validación de campos con reglas
	err := validate.Struct(f)
	if err != nil {
		var messages []string
		for _, e := range err.(validator.ValidationErrors) {
			field := e.Field()
			tag := e.Tag()

			var msg string
			switch tag {
			case "required":
				msg = fmt.Sprintf("el campo %s es obligatorio", field)
			case "password":
				msg = fmt.Sprintf("el campo %s debe tener al menos 8 caracteres, una mayúscula, un número y un caracter especial", field)
			default:
				msg = fmt.Sprintf("el campo %s no cumple la validación %s", field, tag)
			}
			messages = append(messages, msg)
		}
		return ErrorResponse(422, strings.Join(messages, ", "), err)
	}

	if f.NewPassword != f.ConfirmPass {
		return ErrorResponse(422, "las contraseñas no coinciden", fmt.Errorf("las contraseñas no coinciden"))
	}

	return nil
}
