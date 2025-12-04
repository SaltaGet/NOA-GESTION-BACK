package schemas

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type PointSaleCreate struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	IsDeposit   *bool   `json:"is_deposit" validate:"required"`
}

func (p *PointSaleCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
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

type PointSaleUpdate struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	IsDeposit   *bool   `json:"is_deposit" validate:"required"`
}

func (p *PointSaleUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
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

type PointSaleResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	IsDeposit   bool    `json:"is_deposit"`
	IsMain      bool    `json:"is_main"`
}

type PointSaleUpdateMain struct {
	ID      int64 `json:"id"`
	NewMain int64 `json:"new_main"`
}

func (p *PointSaleUpdateMain) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
	if err == nil {
		return nil
	}

	if p.ID == p.NewMain {
		return ErrorResponse(422, "El punto de venta no puede ser el mismo", fmt.Errorf("el punto de venta no puede ser el mismo"))
	}


	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	message := fmt.Sprintf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)

	return ErrorResponse(422, message, fmt.Errorf("%s", message))
}
