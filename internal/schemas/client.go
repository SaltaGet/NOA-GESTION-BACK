package schemas

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ClientCreate struct {
	FirstName   string  `json:"first_name" validate:"required" example:"Jorge"`
	LastName    string  `json:"last_name" validate:"required" example:"Lopez"`
	CompanyName *string `json:"company_name" example:"John Company | null"`
	Identifier  *string `json:"identifier" example:"30000000 | null"`
	Email       *string `json:"email" validate:"email" example:"john@example.com | null"`
	Phone       *string `json:"phone" example:"1111111111 | null"`
	Address     *string `json:"address" example:" Calle 123 | null"`
}

func (c *ClientCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(c)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type ClientUpdate struct {
	ID          string  `json:"id" validate:"required"`
	FirstName   string  `json:"first_name" validate:"required" example:"Jorge"`
	LastName    string  `json:"last_name" validate:"required" example:"Lopez"`
	CompanyName *string `json:"company_name" example:"John Company | null"`
	Identifier  *string `json:"identifier" example:"30000000 | null"`
	Email       *string `json:"email" validate:"email" example:"john@example.com | null"`
	Phone       *string `json:"phone" example:"1111111111 | null"`
	Address     *string `json:"address" example:" Calle 123 | null"`
}

func (c *ClientUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(c)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type ClientResponseDTO struct {
	ID          string  `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	CompanyName *string `json:"company_name,omitempty"`
	Identifier  *string `json:"identifier,omitempty"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
}

type ClientResponse struct {
	ID           string    `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	CompanyName  *string   `json:"company_name,omitempty"`
	Identifier   *string   `json:"identifier,omitempty"`
	Email        *string   `json:"email,omitempty"`
	Phone        *string   `json:"phone,omitempty"`
	Address      *string   `json:"address,omitempty"`
	MemberCreate MemberResponseDTO `json:"member_create"`
}

type ClientSimpleDTO struct {
	ID           string    `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	CompanyName  *string   `json:"company_name,omitempty"`
}


