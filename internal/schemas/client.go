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
	Email       *string `json:"email" validate:"omitempty,email" example:"john@example.com | null"`
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

	message := fmt.Sprintf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)

	return ErrorResponse(422, message, fmt.Errorf("%s", message))
}

type ClientUpdate struct {
	ID          int64  `json:"id" validate:"required"`
	FirstName   string  `json:"first_name" validate:"required" example:"Jorge"`
	LastName    string  `json:"last_name" validate:"required" example:"Lopez"`
	CompanyName *string `json:"company_name" example:"John Company | null"`
	Identifier  *string `json:"identifier" example:"30000000 | null"`
	Email       *string `json:"email" validate:"email" example:"john@example.com | null"`
	Phone       *string `json:"phone" example:"1111111111 | null"`
	Address     *string `json:"address" example:" Calle 123 | null"`
}

func (c *ClientUpdate) Validate() error {
	if c.ID == 1 {
		return ErrorResponse(400, "no se puede eliminar el cliente Consumidor Final", fmt.Errorf("no se puede eliminar el cliente Consumidor Final"))
	}

	validate := validator.New()
	err := validate.Struct(c)
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

type ClientResponseDTO struct {
	ID          int64  `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	CompanyName *string `json:"company_name,omitempty"`
	Identifier  *string `json:"identifier,omitempty"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
}

type ClientResponse struct {
	ID           int64            `json:"id"`
	FirstName    string            `json:"first_name"`
	LastName     string            `json:"last_name"`
	CompanyName  *string           `json:"company_name,omitempty"`
	Identifier   *string           `json:"identifier,omitempty"`
	Email        *string           `json:"email,omitempty"`
	Phone        *string           `json:"phone,omitempty"`
	Address      *string           `json:"address,omitempty"`
	MemberCreate *MemberSimpleDTO `json:"member,omitempty"`
	Pay          []PayDebtResponse     `json:"pay"`
}

type ClientSimpleDTO struct {
	ID          int64  `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	CompanyName *string `json:"company_name,omitempty"`
}
