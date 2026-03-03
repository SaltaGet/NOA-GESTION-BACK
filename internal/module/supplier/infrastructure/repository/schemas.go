package repository


import (
	"time"

)

type SupplierCreate struct {
	Name        string   `json:"name" validate:"required"`
	CompanyName string   `json:"company_name" validate:"required"`
	Identifier  *string  `json:"identifier"`
	Address     *string  `json:"address"`
	DebtLimit   *float64 `json:"debt_limit"`
	Email       *string  `json:"email"`
	Phone       *string  `json:"phone"`
}

type SupplierUpdate struct {
	ID          int64    `json:"id" validate:"required"`
	Name        string   `json:"name" validate:"required"`
	CompanyName string   `json:"company_name" validate:"required"`
	Identifier  *string  `json:"identifier"`
	Address     *string  `json:"address"`
	DebtLimit   *float64 `json:"debt_limit"`
	Email       *string  `json:"email" validate:"omitempty,email"`
	Phone       *string  `json:"phone"`
}

type SupplierResponse struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	CompanyName string   `json:"company_name"`
	Identifier  *string  `json:"identifier"`
	Address     *string  `json:"address"`
	DebtLimit   *float64 `json:"debt_limit"`
	Email       *string  `json:"email"`
	Phone       *string  `json:"phone"`
	CreatedAt   time.Time `json:"created_at"`
}

type SupplierResponseDTO struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	CompanyName string   `json:"company_name"`
}
