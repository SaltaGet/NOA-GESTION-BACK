package schemas

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type TenantCreate struct {
	Name       string `json:"name" validate:"required"`
	Identifier string `json:"identifier" validate:"required"`
	Address    string `json:"address" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	CuitPdv    string `json:"cuit_pdv" validate:"required"`
	PlanID     int64  `json:"plan_id" validate:"required"`
}

func (t *TenantCreate) Validate() error {
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

type TenantUpdate struct {
	ID      string `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

func (t *TenantUpdate) Validate() error {
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

type TenantResponse struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	IsActive   bool      `json:"is_active"`
	Expiration time.Time `json:"expiration"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type TenantUserCreate struct {
	TenantCreate TenantCreate `json:"tenant_create" validate:"required"`
	UserCreate   UserCreate   `json:"user_create" validate:"required"`
}

func (t *TenantUserCreate) Validate() error {
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

type TenantUpdateExpiration struct {
	ID         int64  `json:"id" validate:"required"`
	Expiration string `json:"expiration" validate:"required,datetime=2006-01-02" example:"2023-01-01"`
}

func (t *TenantUpdateExpiration) Validate() error {
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

type TenantUpdateTerms struct {
	AcceptedTerms bool      `json:"accepted_terms" validate:"required"`
	IP            string    `json:"ip" validate:"required"`
	DateAccepted  time.Time `json:"date_aceept" validate:"required"`
}

func (t *TenantUpdateTerms) Validate() error {
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

type TenantUpdateSettings struct {
	Title          *string `json:"title,omitempty" example:"Mi tienda"`
	Slogan         *string `json:"slogan,omitempty" example:"Mi tienda"`
	PrimaryColor   *string `json:"primary_color,omitempty" example:"#FF0000"`
	SecondaryColor *string `json:"secondary_color,omitempty" example:"#FF0000"`
}

func (t *TenantUpdateSettings) Validate() error {
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

type TenantUpdateSettingsWithTenant struct {
	Title          *string `json:"title,omitempty" example:"Mi tienda"`
	Slogan         *string `json:"slogan,omitempty" example:"Mi tienda"`
	PrimaryColor   *string `json:"primary_color,omitempty" example:"#FF0000"`
	SecondaryColor *string `json:"secondary_color,omitempty" example:"#FF0000"`
	TenantID       int64   `json:"tenant_id" validate:"required"`
}
