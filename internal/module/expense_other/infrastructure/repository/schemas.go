package repository

import (
	"time"
)

type ExpenseOtherResponse struct {
	ID          int64               `json:"id"`
	PointSale   *PointSaleResponse  `json:"point_sale,omitempty"`
	Member      *MemberSimpleDTO    `json:"member,omitempty"`
	CashRegisterID  *int64              `json:"cash_register_id,omitempty"`
	Details     *string             `json:"details,omitempty"`
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

type ExpenseOtherCreate struct {
	Details       *string `json:"details"`
	Total         float64 `json:"total" validate:"required"`
	PayMethod     string  `json:"payment_method" validate:"oneof=cash credit card transfer" example:"cash credit card transfer"`
	TypeExpenseID int64   `json:"type_expense_id" validate:"required"`
}

type ExpenseOtherUpdate struct {
	ID            int64   `json:"id" validate:"required"`
	Details       *string `json:"details,omitempty"`
	Total         float64 `json:"total" validate:"required"`
	PayMethod     string  `json:"payment_method" validate:"oneof=cash credit card transfer" example:"cash credit card transfer"`
	TypeExpenseID int64   `json:"type_expense_id" validate:"required" example:"1"`
}
