package schemas

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type SupplierCreate struct {
	Name        string   `json:"name" validate:"required"`
	CompanyName string   `json:"company_name" validate:"required"`
	Identifier  *string  `json:"identifier"`
	Address     *string  `json:"address"`
	DebtLimit   *float64 `json:"debt_limit"`
	Email       *string  `json:"email"`
	Phone       *string  `json:"phone"`
}

func (s *SupplierCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(s)
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

type SupplierUpdate struct {
	ID          int64    `json:"id" validate:"required"`
	Name        string   `json:"name" validate:"required"`
	CompanyName string   `json:"company_name" validate:"required"`
	Identifier  *string  `json:"identifier"`
	Address     *string  `json:"address"`
	DebtLimit   *float64 `json:"debt_limit"`
	Email       *string  `json:"email" validate:"omitempty,email"`
	Phone       *string  `json:"phone"`
}

func (s *SupplierUpdate) Validate() error {
	if s.ID == 1 {
		return ErrorResponse(400, "no se puede editar el proveedor por defecto", fmt.Errorf("no se puede editar el proveedor por defecto"))
	}

	email := s.Email
	if email != nil && *email == "" {
		s.Email = nil
	}

	validate := validator.New()
	err := validate.Struct(s)
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

type SupplierResponse struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	CompanyName string   `json:"company_name"`
	Identifier  *string  `json:"identifier"`
	Address     *string  `json:"address"`
	DebtLimit   *float64 `json:"debt_limit"`
	Email       *string  `json:"email"`
	Phone       *string  `json:"phone"`
	CreatedAt   time.Time `json:"created_at"`
}

type SupplierResponseDTO struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	CompanyName string   `json:"company_name"`
}
