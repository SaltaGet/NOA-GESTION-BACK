package schemas

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type ExpenseBuyResponse struct {
	ID             int64                    `json:"id"`
	PointSale      *PointSaleResponse       `json:"point_sale,omitempty"`
	Member         MemberSimpleDTO          `json:"member"`
	Supplier       SupplierResponse         `json:"supplier"`
	RegisterID     *int64                   `json:"register_id,omitempty"`
	Description    *string                  `json:"description,omitempty"`
	ExpenseItemBuy []ExpenseBuyItemResponse `json:"expense_item_buys"`
	PayExpenseBuy  []PayExpenseBuyResponse  `json:"pay_expense"`
	Subtotal       float64                  `json:"subtotal"`
	Discount       float64                  `json:"discount"`
	Type           string                   `json:"type_discount"`
	Total          float64                  `json:"total"`
	CreatedAt      time.Time                `json:"created_at"`
}

type ExpenseBuyItemResponse struct {
	ID           int64                    `json:"id"`
	Product      ProductSimpleResponseDTO `json:"product"`
	Amount       float64                  `gorm:"not null" json:"amount"`
	Price        float64                  `gorm:"not null" json:"price"`
	Discount     float64                  `gorm:"not null;default:0" json:"discount"`
	TypeDiscount string                   `gorm:"not null;default:percent" json:"type_discount" validate:"oneof=amount percent"`
	SubTotal     float64                  `gorm:"not null" json:"subtotal"`
	Total        float64                  `gorm:"not null" json:"total"`
	CreatedAt    time.Time                `gorm:"autoCreateTime:milli" json:"created_at"`
}

type PayExpenseBuyResponse struct {
	ID        int64   `json:"id"`
	Amount    float64 `json:"amount"`
	MethodPay string  `json:"payment_method"`
}

type ExpenseBuyResponseDTO struct {
	ID             int64                    `json:"id"`
	PointSale      *PointSaleResponse       `json:"point_sale,omitempty"`
	Member         MemberSimpleDTO          `json:"member"`
	Supplier       SupplierResponse         `json:"supplier"`
	RegisterID     *int64                   `json:"register_id,omitempty"`
	Description    *string                  `json:"description,omitempty"`
	ExpenseItemBuy []ExpenseBuyItemResponse `json:"expense_item_buys"`
	PayExpenseBuy  []PayExpenseBuyResponse  `json:"pay_expense"`
	Subtotal       float64                  `json:"subtotal"`
	Discount       float64                  `json:"discount"`
	Type           string                   `json:"type_discount"`
	Total          float64                  `json:"total"`
	CreatedAt      time.Time                `json:"created_at"`
}

type ExpenseBuyResponseSimple struct {
	ID             int64                    `json:"id"`
	Supplier       SupplierResponseDTO         `json:"supplier"`
	RegisterID     *int64                   `json:"register_id,omitempty"`
	Description    *string                  `json:"description,omitempty"`
	PayExpenseBuy  []PayExpenseBuyResponse  `json:"pay_expense"`
	Subtotal       float64                  `json:"subtotal"`
	Discount       float64                  `json:"discount"`
	Type           string                   `json:"type_discount"`
	Total          float64                  `json:"total"`
	CreatedAt      time.Time                `json:"created_at"`
}

type ExpenseOtherResponse struct {
	ID            int64              `json:"id"`
	PointSale     *PointSaleResponse `json:"point_sale,omitempty"`
	Member        MemberSimpleDTO    `json:"member"`
	RegisterID    *int64             `json:"register_id,omitempty"`
	Description   *string            `json:"description,omitempty"`
	Total         float64            `json:"total"`
	PayMethod string             `json:"pay_method"`
	CreatedAt     time.Time          `json:"created_at"`
}

type ExpenseOtherResponseDTO struct {
	ID            int64              `json:"id"`
	RegisterID    *int64             `json:"register_id,omitempty"`
	Description   *string            `json:"description,omitempty"`
	Total         float64            `json:"total"`
	PayMethod string             `json:"pay_method"`
	CreatedAt     time.Time          `json:"created_at"`
}

type ExpenseBuyCreate struct {
	SupplierID     int64                  `json:"supplier_id" validate:"required"`
	Details        *string                `json:"details"`
	Discount       float64                `json:"discount"`
	Type           string                 `json:"type_discount" validate:"oneof=amount percent"`
	ExpenseBuyItem []ExpenseBuyItemCreate `json:"expense_item_buys" validate:"required,dive"`
	PayExpenseBuy  []PayExpenseBuyCreate  `json:"pay_expense" validate:"required,dive,max=3"`
	Total          float64                `json:"total" validate:"required"`
}

type ExpenseBuyItemCreate struct {
	ID           int64                    `json:"id"`
	Product      ProductSimpleResponseDTO `json:"product"`
	Amount       float64                  `json:"amount"`
	Price        float64                  `json:"price"`
	Discount     float64                  `json:"discount"`
	TypeDiscount string                   `json:"type_discount" validate:"oneof=amount percent"`
	SubTotal     float64                  `json:"subtotal"`
	Total        float64                  `json:"total"`
	CreatedAt    time.Time                `json:"created_at"`
}

type PayExpenseBuyCreate struct {
	Amount    float64 `json:"amount" validate:"required"`
	MethodPay string  `json:"payment_method" validate:"oneof=cash credit card transfer"`
}

func (e *ExpenseBuyCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(e)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type ExpenseBuyUpdate struct {
	ID             int64                  `json:"id" validate:"required"`
	SupplierID     int64                  `json:"supplier_id" validate:"required"`
	Details        *string                `json:"details"`
	Discount       float64                `json:"discount"`
	Type           string                 `json:"type_discount" validate:"oneof=amount percent"`
	ExpenseBuyItem []ExpenseBuyItemCreate `json:"expense_item_buys" validate:"required,dive"`
	PayExpenseBuy  []PayExpenseBuyCreate  `json:"pay_expense" validate:"required,dive,max=3"`
	Total          float64                `json:"total" validate:"required"`
}

type ExpenseBuyItemUpdate struct {
	ID           int64                    `json:"id"`
	Product      ProductSimpleResponseDTO `json:"product"`
	Amount       float64                  `json:"amount"`
	Price        float64                  `json:"price"`
	Discount     float64                  `json:"discount"`
	TypeDiscount string                   `json:"type_discount" validate:"oneof=amount percent"`
	SubTotal     float64                  `json:"subtotal"`
	Total        float64                  `json:"total"`
	CreatedAt    time.Time                `json:"created_at"`
}

type PayExpenseBuyUpdate struct {
	Amount    float64 `json:"amount" validate:"required"`
	MethodPay string  `json:"payment_method" validate:"oneof=cash credit card transfer"`
}

func (e *ExpenseBuyUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(e)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type ExpenseOtherCreate struct {
	Description *string `json:"description"`
	Total       float64 `json:"total" validate:"required"`
	PayMethod   string  `json:"payment_method" validate:"oneof=cash credit card transfer"`
}

func (e *ExpenseOtherCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(e)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type ExpenseOtherUpdate struct {
	ID          int64   `json:"id" validate:"required"`
	Description *string `json:"description"`
	Total       float64 `json:"total" validate:"required"`
	PayMethod   string  `json:"payment_method" validate:"oneof=cash credit card transfer"`
}

func (e *ExpenseOtherUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(e)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}