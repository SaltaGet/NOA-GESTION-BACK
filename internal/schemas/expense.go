package schemas

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type ExpenseBuyResponse struct {
	ID             int64                    `json:"id"`
	Member         MemberSimpleDTO          `json:"member"`
	Supplier       SupplierResponseDTO      `json:"supplier"`
	Description    *string                  `json:"description,omitempty"`
	ExpenseBuyItem []ExpenseBuyItemResponse `json:"expense_buy_items"`
	PayExpenseBuy  []PayExpenseBuyResponse  `json:"pay_expense"`
	Subtotal       float64                  `json:"subtotal"`
	Discount       float64                  `json:"discount"`
	TypeDiscount   string                   `json:"type_discount"`
	Total          float64                  `json:"total"`
	CreatedAt      time.Time                `json:"created_at"`
}

type ExpenseBuyItemResponse struct {
	ID           int64                    `json:"id"`
	Product      ProductSimpleResponseDTO `json:"product"`
	Amount       float64                  `json:"amount"`
	Price        float64                  `json:"price"`
	Discount     float64                  `json:"discount"`
	TypeDiscount string                   `json:"type_discount"`
	Subtotal     float64                  `json:"subtotal"`
	Total        float64                  `json:"total"`
	CreatedAt    time.Time                `json:"created_at"`
}

type PayExpenseBuyResponse struct {
	ID        int64   `json:"id"`
	Total     float64 `json:"total"`
	MethodPay string  `json:"method_pay"`
}

type ExpenseBuyResponseDTO struct {
	ID             int64                    `json:"id"`
	Member         MemberSimpleDTO          `json:"member"`
	Supplier       SupplierResponse         `json:"supplier"`
	RegisterID     *int64                   `json:"register_id,omitempty"`
	Description    *string                  `json:"description,omitempty"`
	ExpenseItemBuy []ExpenseBuyItemResponse `json:"expense_item_buys"`
	PayExpenseBuy  []PayExpenseBuyResponse  `json:"pay_expense"`
	Subtotal       float64                  `json:"subtotal"`
	Discount       float64                  `json:"discount"`
	TypeDiscount   string                   `json:"type_discount"`
	Total          float64                  `json:"total"`
	CreatedAt      time.Time                `json:"created_at"`
}

type ExpenseBuyResponseSimple struct {
	ID            int64                   `json:"id"`
	Supplier      SupplierResponseDTO     `json:"supplier"`
	RegisterID    *int64                  `json:"register_id,omitempty"`
	Description   *string                 `json:"description,omitempty"`
	PayExpenseBuy []PayExpenseBuyResponse `json:"pay_expense"`
	Subtotal      float64                 `json:"subtotal"`
	Discount      float64                 `json:"discount"`
	TypeDiscount  string                  `json:"type_discount"`
	Total         float64                 `json:"total"`
	CreatedAt     time.Time               `json:"created_at"`
}

type ExpenseOtherResponse struct {
	ID          int64               `json:"id"`
	PointSale   *PointSaleResponse  `json:"point_sale,omitempty"`
	Member      *MemberSimpleDTO     `json:"member,omitempty"`
	RegisterID  *int64              `json:"register_id,omitempty"`
	Description *string             `json:"description,omitempty"`
	Total       float64             `json:"total"`
	PayMethod   string              `json:"pay_method"`
	TypeExpense TypeExpenseResponse `json:"type_expense"`
	CreatedAt   time.Time           `json:"created_at"`
}

type TypeExpenseResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ExpenseOtherResponseDTO struct {
	ID          int64               `json:"id"`
	RegisterID  *int64              `json:"register_id,omitempty"`
	Details     *string             `json:"details,omitempty"`
	Total       float64             `json:"total"`
	PayMethod   string              `json:"pay_method"`
	TypeExpense TypeExpenseResponse `json:"type_expense"`
	PointSale   *PointSaleResponse  `json:"point_sale,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
}

type ExpenseBuyCreate struct {
	SupplierID     int64                  `json:"supplier_id" validate:"required"`
	Details        *string                `json:"details"`
	Discount       float64                `json:"discount"`
	TypeDiscount   string                 `json:"type_discount" validate:"oneof=amount percent" example:"amount percent"`
	ExpenseBuyItem []ExpenseBuyItemCreate `json:"expense_item_buys" validate:"required,dive"`
	PayExpenseBuy  []PayExpenseBuyCreate  `json:"pay_expense" validate:"required,max=3,dive"`
	Total          float64                `json:"total" validate:"required"`
}

type ExpenseBuyItemCreate struct {
	ProductID    int64   `json:"product_id" validate:"required"`
	Amount       float64 `json:"amount"`
	Price        float64 `json:"price"`
	Discount     float64 `json:"discount"`
	TypeDiscount string  `json:"type_discount" validate:"oneof=amount percent" example:"amount percent"`
	Total        float64 `json:"total"`
}

type PayExpenseBuyCreate struct {
	Total     float64 `json:"total" validate:"required"`
	MethodPay string  `json:"payment_method" validate:"oneof=cash credit card transfer" example:"cash credit card transfer"`
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

	message := fmt.Sprintf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)

	return ErrorResponse(422, message, fmt.Errorf("%s", message))
}

type ExpenseBuyUpdate struct {
	ID             int64                  `json:"id" validate:"required"`
	SupplierID     int64                  `json:"supplier_id" validate:"required"`
	Details        *string                `json:"details"`
	Discount       float64                `json:"discount"`
	Type           string                 `json:"type_discount" validate:"oneof=amount percent" example:"amount percent"`
	ExpenseBuyItem []ExpenseBuyItemCreate `json:"expense_item_buys" validate:"required,dive"`
	PayExpenseBuy  []PayExpenseBuyCreate  `json:"pay_expense" validate:"required,max=3,dive"`
	Total          float64                `json:"total" validate:"required"`
}

type ExpenseBuyItemUpdate struct {
	Product      ProductSimpleResponseDTO `json:"product"`
	Amount       float64                  `json:"amount"`
	Price        float64                  `json:"price"`
	Discount     float64                  `json:"discount"`
	TypeDiscount string                   `json:"type_discount" validate:"oneof=amount percent" example:"amount percent"`
	SubTotal     float64                  `json:"subtotal"`
	Total        float64                  `json:"total"`
	CreatedAt    time.Time                `json:"created_at"`
}

type PayExpenseBuyUpdate struct {
	Amount    float64 `json:"amount" validate:"required" example:"100.00"`
	MethodPay string  `json:"payment_method" validate:"oneof=cash credit card transfer" example:"cash credit card transfer"`
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

	message := fmt.Sprintf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)

	return ErrorResponse(422, message, fmt.Errorf("%s", message))
}

type ExpenseOtherCreate struct {
	Details       *string `json:"details"`
	Total         float64 `json:"total" validate:"required"`
	PayMethod     string  `json:"payment_method" validate:"oneof=cash credit card transfer" example:"cash credit card transfer"`
	TypeExpenseID int64   `json:"type_expense_id" validate:"required"`
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
	ID        int64   `json:"id" validate:"required"`
	Details   *string `json:"details,omitempty"`
	Total     float64 `json:"total" validate:"required"`
	PayMethod string  `json:"payment_method" validate:"oneof=cash credit card transfer" example:"cash credit card transfer"`
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

	message := fmt.Sprintf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)

	return ErrorResponse(422, message, fmt.Errorf("%s", message))
}
