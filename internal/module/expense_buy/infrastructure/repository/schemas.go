package repository

import (
	"fmt"
	"math"
	"time"
)

type ExpenseBuyResponse struct {
	ID             int64                    `json:"id"`
	Member         MemberSimpleDTO          `json:"member"`
	Supplier       SupplierResponseDTO      `json:"supplier"`
	Details        *string                  `json:"details,omitempty"`
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

type TypeExpenseResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
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

func (e *ExpenseBuyCreate) ValidateIntegrity() error {
    var sumPay float64
    for _, p := range e.PayExpenseBuy {
        sumPay += p.Total
    }

    // Usamos math.Abs para evitar problemas con decimales (punto flotante)
    if math.Abs(sumPay-e.Total) > 0.01 { // Usamos 0.01 si manejas centavos
        return fmt.Errorf("la suma de los pagos (%.2f) no coincide con el total (%.2f)", sumPay, e.Total)
    }
    return nil
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

func (e *ExpenseBuyUpdate) ValidateIntegrity() error {
    var sumPay float64
    for _, p := range e.PayExpenseBuy {
        sumPay += p.Total
    }

    // Usamos math.Abs para evitar problemas con decimales (punto flotante)
    if math.Abs(sumPay-e.Total) > 0.01 { // Usamos 0.01 si manejas centavos
        return fmt.Errorf("la suma de los pagos (%.2f) no coincide con el total (%.2f)", sumPay, e.Total)
    }
    return nil
}
