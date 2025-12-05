package schemas

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type MemberResponse struct {
	ID         int64                `json:"id"`
	FirstName  string              `json:"first_name"`
	LastName   string              `json:"last_name"`
	Username   string              `json:"username"`
	Email      string              `json:"email"`
	Address    *string             `json:"address"`
	Phone      *string             `json:"phone"`
	IsAdmin    bool                `json:"is_admin"`
	IsActive   bool                `json:"is_active"`
	Role       RoleResponse        `json:"role"`
	PointSales []PointSaleResponse `json:"point_sales"`
}

type MemberResponseDTO struct {
	ID        int64         `json:"id"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Address   *string      `json:"address"`
	Phone     *string      `json:"phone"`
	IsActive  bool         `json:"is_active"`
	IsAdmin   bool         `json:"is_admin"`
	Role      RoleResponseDTO `json:"role"`
}

type MemberSimpleDTO struct {
	ID        int64    `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Username  string  `json:"username"`
}

type MemberCreate struct {
	FirstName    string  `json:"first_name" validate:"required" example:"John"`
	LastName     string  `json:"last_name" validate:"required" example:"Doe"`
	Username     string  `json:"username" validate:"required" example:"johndoe"`
	Email        string  `json:"email" validate:"email,required" example:"a@b.com"`
	Password     string  `json:"password" validate:"required,password" example:"Password123*"`
	Address      *string `json:"address,omitempty" example:"casita roja|null"`
	Phone        *string `json:"phone,omitempty" example:"123123123|null"`
	RoleID       int64   `json:"role_id" validate:"required" example:"1"`
	PointSaleIDs []int64 `json:"point_sales_ids" validate:"required" example:"1,2,3"`
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password) {
		return false
	}
	return true
}

func (u MemberCreate) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("password", validatePassword)

	err := validate.Struct(u)
	if err == nil {
		return nil
	}

	validatorErr := err.(validator.ValidationErrors)[0]
	field := validatorErr.Field()
	tag := validatorErr.Tag()

	var errorMessage string
	switch tag {
	case "required":
		errorMessage = fmt.Sprintf("el campo %s es obligatorio", field)
	case "email":
		errorMessage = fmt.Sprintf("el campo %s debe ser un email válido", field)
	case "password":
		errorMessage = "el campo password debe tener al menos 8 caracteres, una letra mayúscula, un número y un carácter especial"
	default:
		errorMessage = fmt.Sprintf("el campo %s no cumple la validación %s", field, tag)
	}

	return ErrorResponse(422, fmt.Sprintf("error al validar campo(s): %s", errorMessage), err)
}

type MemberUpdate struct {
	ID           int64    `json:"id" validate:"required" example:"1"`
	FirstName    string  `json:"first_name" validate:"required" example:"John"`
	LastName     string  `json:"last_name" validate:"required" example:"Doe"`
	Username     string  `json:"username" validate:"required" example:"johndoe"`
	Email        string  `json:"email" validate:"email,required" example:"a@b.com"`
	Address      *string `json:"address,omitempty" example:"address|null"`
	Phone        *string `json:"phone,omitempty" example:"phone|null"`
	RoleID       int64    `json:"role_id" validate:"required" example:"1"`
	IsActive     *bool   `json:"is_active" validate:"required" example:"true"`
	PointSaleIDs []int64  `json:"point_sales_ids" validate:"required" example:"1,2,3"`
}

func (u *MemberUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(u)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return ErrorResponse(422, "error desconocido al validar", err)
	}

	var messages []string
	for _, e := range validationErrors {
		msg := fmt.Sprintf("Campo '%s' no cumple la regla '%s'", e.Field(), e.Tag())
		messages = append(messages, msg)
	}

	return ErrorResponse(422, fmt.Sprintf("Error al validar campo(s): %s", strings.Join(messages, ", ")), err)
}

type MemberUpdatePassword struct {
	OldPassword string `json:"old_password" validate:"required,password" example:"Password123*"`
	NewPassword string `json:"new_password" validate:"required,password" example:"Password123*"`
	ConfirmPass string `json:"confirm_pass" validate:"required,password" example:"Password123*"`
}

func (u MemberUpdatePassword) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("password", validatePassword) // registrar antes de Struct()

	// validación de campos con reglas
	err := validate.Struct(u)
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

	// validación manual: confirmación de contraseña
	if u.NewPassword != u.ConfirmPass {
		return ErrorResponse(422, "las contraseñas no coinciden", fmt.Errorf("las contraseñas no coinciden"))
	}

	return nil
}
