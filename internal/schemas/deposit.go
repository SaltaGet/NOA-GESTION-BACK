package schemas

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type DepositResponse struct {
	ID          int64   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description *string `json:"description"`
	Category    CategoryResponse `json:"category"`
	Price       float64 `json:"price"`
	Stock       float64    `json:"stock"`
}

type DepositUpdateStock struct {
	ProductID int64    `json:"product_id" validate:"required" example:"1"`
	Stock     *float64 `json:"stock" validate:"required,gte=0" example:"10"`
	Method    string   `json:"method" validate:"oneof=add subtract set" example:"add|subtract|set"`
}


func (d *DepositUpdateStock) Validate() error {
	validate := validator.New()
	err := validate.Struct(d)
	if err == nil {
		return nil
	}

	validatorErr := err.(validator.ValidationErrors)[0]
	field := validatorErr.Field()
	tag := validatorErr.Tag()
	params := validatorErr.Param()

	errorMessage := field + " " + tag + " " + params
	return ErrorResponse(422, fmt.Sprintf("error al validar campo(s): %s", errorMessage), err)
}