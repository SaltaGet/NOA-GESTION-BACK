package repository

import (
	"fmt"
	"math"
	"time"
)

type IncomeOtherCreate struct {
	Total        float64 `json:"total" validate:"required"`
	TypeIncomeID int64   `json:"type_income_id" validate:"required"`
	Details      *string `json:"details"`
	MethodIncome string  `json:"method_income" validate:"oneof=cash credit card transfer" example:"cash credit card transfer"`
}

type IncomeOtherUpdate struct {
	ID           int64   `json:"id" validate:"required"`
	Total        float64 `json:"total" validate:"required"`
	TypeIncomeID int64   `json:"type_income_id" validate:"required"`
	Details      *string `json:"details"`
	MethodIncome string  `json:"method_income" validate:"oneof=cash credit card transfer" example:"cash credit card transfer"`
}

type IncomeOtherResponse struct {
	ID             int64              `json:"id"`
	Member         *MemberSimpleDTO   `json:"member,omitempty"`
	Total          float64            `json:"total"`
	TypeIncome     TypeIncomeResponse `json:"type_income"`
	Details        *string            `json:"details,omitempty"`
	MethodIncome   string             `json:"method_income"`
	PointSale      *PointSaleResponse `json:"point_sale,omitempty"`
	CashRegisterID *int64             `json:"cash_register_id,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
}

type TypeIncomeResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
