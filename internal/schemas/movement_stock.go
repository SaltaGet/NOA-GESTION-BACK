package schemas

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type MovementStockResponse struct {
	ID          int64            `json:"id"`
	Member      MemberSimpleDTO `json:"member"`
	Product     ProductResponse `json:"product"`
	Amount      float64         `json:"amount"`
	FromID      int64            `json:"from_id"`
	FromType    string          `json:"from_type"`
	ToID        int64            `json:"to_id"`
	ToType      string          `json:"to_type"`
	IgnoreStock bool            `json:"ignore_stock"`
	CreatedAt   time.Time       `json:"created_at"`
}

type MovementStockResponseDTO struct {
	ID          int64                  `json:"id"`
	Member      MemberSimpleDTO       `json:"member"`
	Product     ProductSimpleResponse `json:"product"`
	Amount      float64               `json:"amount"`
	FromID      int64                  `json:"from_id"`
	FromType    string                `json:"from_type"`
	ToID        int64                  `json:"to_id"`
	ToType      string                `json:"to_type"`
	IgnoreStock bool                  `json:"ignore_stock"`
	CreatedAt   time.Time             `json:"created_at"`
}

type MovementStock struct {
	ProductID int64    `json:"product_id" validate:"required" example:"1"`
	Amount    float64 `json:"amount" validate:"required" example:"10"`

	FromType string `json:"from_type" validate:"oneof=deposit point_sale" example:"deposit|point_sale"`
	FromID   int64   `json:"from_id" validate:"required" example:"1"`

	ToType string `json:"to_type" validate:"oneof=deposit point_sale" example:"deposit|point_sale"`
	ToID   int64   `json:"to_id" validate:"required" example:"1"`

	IgnoreStock *bool `json:"ignore_stock" validate:"required" example:"false"`
}

// func (m *MovementStock) Validate() error {
// 	validate := validator.New()
// 	err := validate.Struct(m)
// 	if err == nil {
// 		return nil
// 	}

// 	validatorErr := err.(validator.ValidationErrors)[0]
// 	field := validatorErr.Field()
// 	tag := validatorErr.Tag()
// 	params := validatorErr.Param()

//		errorMessage := field + " " + tag + " " + params
//		return ErrorResponse(422, fmt.Sprintf("error al validar campo(s): %s", errorMessage), err)
//	}
func (m *MovementStock) Validate() error {
	validate := validator.New()
	err := validate.Struct(m)
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

type MovementStockList struct {
	ProductID         int64                `json:"product_id" validate:"required" example:"1"`
	MovementStockItem []MovementStockItem `json:"movement_stock_item" validate:"required,dive"`
}

type MovementStockItem struct {
	Amount float64 `json:"amount" validate:"required" example:"10"`

	FromType string `json:"from_type" validate:"oneof=deposit point_sale" example:"deposit|point_sale"`
	FromID   int64   `json:"from_id" validate:"required" example:"1"`

	ToType      string `json:"to_type" validate:"oneof=deposit point_sale" example:"deposit|point_sale"`
	ToID        int64   `json:"to_id" validate:"required" example:"1"`
	IgnoreStock *bool  `json:"ignore_stock" validate:"required" example:"false"`
}

// func (m *MovementStock) Validate() error {
// 	validate := validator.New()
// 	err := validate.Struct(m)
// 	if err == nil {
// 		return nil
// 	}

// 	validatorErr := err.(validator.ValidationErrors)[0]
// 	field := validatorErr.Field()
// 	tag := validatorErr.Tag()
// 	params := validatorErr.Param()

//		errorMessage := field + " " + tag + " " + params
//		return ErrorResponse(422, fmt.Sprintf("error al validar campo(s): %s", errorMessage), err)
//	}
func (m *MovementStockList) Validate() error {
	validate := validator.New()
	err := validate.Struct(m)
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
