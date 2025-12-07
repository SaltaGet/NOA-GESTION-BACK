package schemas

import (
	"fmt"
	"math"
	"time"

	"github.com/go-playground/validator/v10"
)

type IncomeSaleCreate struct {
	Items    []ItemIncomeSaleCreate `json:"items" validate:"required,dive"`
	Pay      []PayCreate            `json:"pay" validate:"required,max=3,dive"`
	ClientID int64                  `json:"client_id" validate:"required"`
	Discount float64                `json:"discount" validate:"min=0"`
	Type     string                 `json:"type_discount" validate:"oneof=amount percent" example:"amount percent"`
	Total    float64                `json:"total"`
	IsBudget bool                   `json:"is_budget"`
}

type ItemIncomeSaleCreate struct {
	ProductID    int64   `json:"product_id" validate:"required"`
	Amount       float64 `json:"amount" validate:"required"`
	Discount     float64 `json:"discount" validate:"min=0"`
	TypeDiscount string  `json:"type_discount" validate:"oneof=amount percent"`
}

type PayCreate struct {
	Total     float64 `json:"total" validate:"required" example:"1000"`
	MethodPay string  `json:"method_pay" validate:"oneof=cash credit card transfer" example:"cash credit card transfer"`
}

func (i *IncomeSaleCreate) Validate() error {
	validate := validator.New()

	if err := validate.Struct(i); err != nil {
		validationErr := err.(validator.ValidationErrors)[0]
		field := validationErr.Field()
		tag := validationErr.Tag()
		param := validationErr.Param()
		return ErrorResponse(422, fmt.Sprintf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param), err)
	}

	var sumPay float64
	for _, p := range i.Pay {
		sumPay += p.Total
	}

	if math.Abs(sumPay-i.Total) > 1 {
		message := fmt.Sprintf("la diferencia entre la suma de pagos (%.2f) y el total (%.2f) no puede ser mayor que 1",
			sumPay, i.Total)

		return ErrorResponse(422, message, fmt.Errorf("%s", message))
	}

	return nil
}

type IncomeSaleUpdate struct {
	ID       int64                  `json:"id" validate:"required"`
	Items    []ItemIncomeSaleUpdate `json:"items" validate:"required,dive"`
	Pay      []PayUpdate            `json:"pay" validate:"required,max=3,dive"`
	ClientID int64                  `json:"client_id" validate:"required"`
	Discount float64                `json:"discount"`
	Type     string                 `json:"type_discount" validate:"oneof=amount percent" example:"amount percent"`
	Total    float64                `json:"total"`
	IsBudget bool                   `json:"is_budget"`
}

type ItemIncomeSaleUpdate struct {
	ProductID    int64   `json:"product_id" validate:"required"`
	Amount       float64 `json:"amount" validate:"required"`
	Discount     float64 `json:"discount"`
	TypeDiscount string  `json:"type_discount" validate:"oneof=amount percent"`
}

type PayUpdate struct {
	Total     float64 `json:"total" validate:"required"`
	MethodPay string  `json:"method_pay" validate:"oneof=cash credit card transfer"`
}

func (i *IncomeSaleUpdate) Validate() error {
	validate := validator.New()

	if err := validate.Struct(i); err != nil {
		validationErr := err.(validator.ValidationErrors)[0]
		field := validationErr.Field()
		tag := validationErr.Tag()
		param := validationErr.Param()
		return ErrorResponse(422, fmt.Sprintf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param), err)
	}

	var sumPay float64
	for _, p := range i.Pay {
		sumPay += p.Total
	}

	if math.Abs(sumPay-i.Total) > 1 {
		message := fmt.Sprintf("la diferencia entre la suma de pagos (%.2f) y el total (%.2f) no puede ser mayor que 1",
			sumPay, i.Total)

		return ErrorResponse(422, message, fmt.Errorf("%s", message))
	}

	return nil
}

type IncomeOtherCreate struct {
	Total        float64 `json:"total" validate:"required"`
	TypeIncomeID int64   `json:"type_income_id" validate:"required"`
	Details      *string `json:"details"`
	MethodIncome string  `json:"method_income" validate:"oneof=cash credit card transfer" example:"cash credit card transfer"`
}

func (i *IncomeOtherCreate) Validate() error {
	validate := validator.New()

	if err := validate.Struct(i); err != nil {
		validationErr := err.(validator.ValidationErrors)[0]
		field := validationErr.Field()
		tag := validationErr.Tag()
		param := validationErr.Param()
		return ErrorResponse(422, fmt.Sprintf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param), err)
	}

	return nil
}

type IncomeOtherUpdate struct {
	ID           int64   `json:"id" validate:"required"`
	Total        float64 `json:"total" validate:"required"`
	TypeIncomeID int64   `json:"type_income_id" validate:"required"`
	Details      *string `json:"details"`
	MethodIncome string  `json:"method_income" validate:"oneof=cash credit card transfer" example:"cash credit card transfer"`
}

func (i *IncomeOtherUpdate) Validate() error {
	validate := validator.New()

	if err := validate.Struct(i); err != nil {
		validationErr := err.(validator.ValidationErrors)[0]
		field := validationErr.Field()
		tag := validationErr.Tag()
		param := validationErr.Param()
		return ErrorResponse(422, fmt.Sprintf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param), err)
	}

	return nil
}

type IncomeSaleResponse struct {
	ID        int64                    `json:"id"`
	Member    MemberSimpleDTO          `json:"member"`
	Client    ClientResponseDTO        `json:"client"`
	Items     []IncomeSaleItemResponse `json:"items"`
	Pay       []PayResponse            `json:"pay"`
	SubTotal  float64                  `json:"subtotal"`
	Discount  float64                  `json:"discount"`
	Type      string                   `json:"type_discount"`
	Total     float64                  `json:"total"`
	IsBudget  bool                     `json:"is_budget"`
	CreatedAt time.Time                `json:"created_at"`
}

type IncomeSaleItemResponse struct {
	ID           int64                    `json:"id"`
	Product      ProductSimpleResponseDTO `json:"product"`
	Amount       float64                  `json:"quantity"`
	Price        float64                  `json:"price"`
	Discount     float64                  `json:"discount"`
	TypeDiscount string                   `json:"type_discount"`
	SubTotal     float64                  `json:"subtotal"`
	Total        float64                  `json:"total"`
	CreatedAt    time.Time                `json:"created_at"`
}

type IncomeSaleResponseDTO struct {
	ID        int64           `json:"id"`
	Member    MemberSimpleDTO `json:"member"`
	Client    ClientSimpleDTO `json:"client"`
	Pay       []PayResponse   `json:"pay"`
	Total     float64         `json:"total"`
	CreatedAt time.Time       `json:"created_at"`
}

type IncomeSaleSimpleResponse struct {
	ID        int64                      `json:"id"`
	Items     []IncomeSaleItemResponseDTO `json:"items"`
	Pay       []PayResponse              `json:"pay"`
	Total     float64                    `json:"total"`
	IsBudget  bool                       `json:"is_budget"`
	CreatedAt time.Time                  `json:"created_at"`
}

type IncomeSaleItemResponseDTO struct {
	ID        int64                    `json:"id"`
	Product   ProductSimpleResponseDTO `json:"product"`
	Amount    float64                  `json:"quantity"`
	Total     float64                  `json:"total"`
	CreatedAt time.Time                `json:"created_at"`
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
