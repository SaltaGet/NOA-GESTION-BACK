package schemas

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ProductFullResponse struct {
	ID              int64                     `json:"id"`
	Code            string                    `json:"code"`
	Name            string                    `json:"name"`
	Description     *string                    `json:"description"`
	Category        CategoryResponse          `json:"category"`
	Price           float64                   `json:"price"`
	StockPointSales []*PointSaleStockResponse `json:"stock_point_sales"`
	StockDeposit    float64          `json:"stock_deposit"`
	Notifier        bool                      `json:"notifier"`
	MinAmount       float64                   `json:"min_amount"`
}

type ProductResponse struct {
	ID          int64            `json:"id"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Category    CategoryResponse `json:"category"`
	Price       float64          `json:"price"`
	Stock       float64          `json:"stock"`
	Notifier    bool             `json:"notifier"`
	MinAmount   float64          `json:"min_amount"`
}

type ProductResponseDTO struct {
	ID        int64             `json:"id"`
	Code      string            `json:"code"`
	Name      string            `json:"name"`
	Category  *CategoryResponse `json:"category,omitempty"`
	Price     float64           `json:"price"`
	Stock     float64           `json:"stock"`
	Notifier  bool              `json:"notifier"`
	MinAmount float64           `json:"min_amount"`
}

type ProductSimpleResponse struct {
	ID        int64   `json:"id"`
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Stock     float64 `json:"stock"`
	Notifier  bool    `json:"notifier"`
	MinAmount float64 `json:"min_amount"`
}

type ProductSimpleResponseDTO struct {
	ID    int64   `json:"id"`
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductCreate struct {
	Code        string   `json:"code" validate:"required" example:"ABC123"`
	Name        string   `json:"name" validate:"required" example:"Producto1"`
	Description *string  `json:"description" example:"description|null"`
	CategoryID  int64    `json:"category_id" validate:"required" example:"1"`
	Price       *float64 `json:"price" example:"100.00"`
	Notifier    bool     `json:"notifier" example:"false"`
	MinAmount   float64  `json:"min_amount" example:"10.00"`
}

func (p *ProductCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
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

type ProductUpdate struct {
	ID          int64    `json:"id" validate:"required" example:"1"`
	Code        string   `json:"code" validate:"required" example:"ABC123"`
	Name        string   `json:"name" validate:"required" example:"Producto1"`
	Description *string  `json:"description" example:"description|null"`
	CategoryID  uint     `json:"category_id" validate:"required" example:"1"`
	Price       *float64 `json:"price" example:"100.00"`
	Notifier    bool     `json:"notifier" example:"false"`
	MinAmount   float64  `json:"min_amount" example:"10.00"`
}

func (p *ProductUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
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

type ProductPriceUpdate struct {
	ID    int64   `json:"id" validate:"required" example:"1"`
	Price float64 `json:"price" validate:"required,gte=0" example:"100.00"`
}

type ListPriceUpdate struct {
	ListProductPriceUpdate []ProductPriceUpdate `json:"list" validate:"required,min=1,dive"`
}

func (p *ListPriceUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
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

type ProductStockWithScore struct {
    Product *ProductStockFullResponse
    Score   float64
    Length  int
}

type ProductValidateImage struct {
	ProductID      int64                  `json:"product_id" validate:"required" example:"1"`
	PrimaryImage   string                 `json:"primary_image" validate:"required,oneof=add update keep" example:"add | update | keep"`
	SecondaryImage ValidateSecondaryImage `json:"secondary_image" validate:"required"`
}

type ValidateSecondaryImage struct {
	Keep   *int64 `json:"keep" validate:"required" example:"1"`
	Add    *int64 `json:"add" validate:"required" example:"1"`
	Remove *int64 `json:"remove" validate:"required" example:"1"`
}

func (p *ProductValidateImage) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
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