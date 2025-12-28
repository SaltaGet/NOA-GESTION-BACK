package schemas

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type PlanResponse struct {
	ID              int64            `json:"id"`
	Name            string           `json:"name"`
	PriceMounthly   float64          `json:"price"`
	PriceYearly     float64          `json:"price_yearly"`
	Description     string           `json:"description"`
	Features        string           `json:"features,omitempty"`
	AmountPointSale int64            `json:"amount_point_sale"`
	AmountMember    int64            `json:"amount_member"`
	AmountProduct   int64            `json:"amount_product"`
	Tenants         []TenantResponse `json:"tenants,omitempty"`
}

type PlanResponseDTO struct {
	ID              int64               `json:"id"`
	Name            string              `json:"name"`
	PriceMounthly   float64             `json:"price"`
	PriceYearly     float64             `json:"price_yearly"`
	Description     string              `json:"description"`
	Features        string              `json:"features,omitempty"`
	AmountPointSale int64               `json:"amount_point_sale"`
	AmountMember    int64               `json:"amount_member"`
	AmountProduct   int64               `json:"amount_product"`
	Modules         []ModuleResponseDTO `json:"modules,omitempty"`
}

type PlanCreate struct {
	Name            string  `json:"name" validate:"required" example:"Plan1"`
	PriceMounthly   float64 `json:"price" validate:"required,gte=0" example:"100.00"`
	PriceYearly     float64 `json:"price_yearly" validate:"required,gte=0" example:"1000.00"`
	Description     string  `json:"description" validate:"required" example:"description"`
	Features        string  `json:"features" example:"features"`
	AmountPointSale int64   `json:"amount_point_sale" validate:"required" example:"1"`
	AmountMember    int64   `json:"amount_member" validate:"required" example:"5"`
	AmountProduct   int64   `json:"amount_product"`
}

func (p *PlanCreate) Validate() error {
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

type PlanUpdate struct {
	ID              int64   `json:"id" validate:"required" example:"1"`
	Name            string  `json:"name" validate:"required" example:"Plan1"`
	PriceMounthly   float64 `json:"price" validate:"required,gte=0" example:"100.00"`
	PriceYearly     float64 `json:"price_yearly" validate:"required,gte=0" example:"1000.00"`
	Description     string  `json:"description" validate:"required" example:"description"`
	Features        string  `json:"features" example:"features"`
	AmountPointSale int64   `json:"amount_point_sale" validate:"required" example:"1"`
	AmountMember    int64   `json:"amount_member" validate:"required" example:"5"`
	AmountProduct   int64   `json:"amount_product"`
}

func (p *PlanUpdate) Validate() error {
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
