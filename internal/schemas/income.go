package schemas

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type SaleIncomeCreate struct {
	Items    []ItemSaleIncomeCreate `json:"items" validate:"required,dive"`
	Pay      []PayCreate            `json:"pay" validate:"required,dive,max=3"`
	Discount float32                `json:"discount"`
	Type     string                 `json:"type_discount" validate:"oneof=amount percent"`
	Total    float32                `json:"total"`
	IsBudget bool                   `json:"is_budget"`
}

type ItemSaleIncomeCreate struct {
	ProductID    int64   `json:"product_id" validate:"required"`
	Amount       float32 `json:"amount" validate:"required"`
	Discount     float32 `json:"discount"`
	TypeDiscount string  `json:"type_discount" validate:"oneof=amount percent"`
}

type PayCreate struct {
	Amount    float32 `json:"amount" validate:"required"`
	MethodPay string  `json:"method_pay" validate:"oneof=cash credit card transfer"`
}

func (i *SaleIncomeCreate) Validate() error {
	validate := validator.New()

	if err := validate.Struct(i); err != nil {
		validationErr := err.(validator.ValidationErrors)[0]
		field := validationErr.Field()
		tag := validationErr.Tag()
		param := validationErr.Param()
		return ErrorResponse(422, fmt.Sprintf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param), err)
	}

	var sumPay float32
	for _, p := range i.Pay {
		sumPay += p.Amount
	}

	if sumPay > i.Total {
		message := fmt.Sprintf("la suma de los pagos (%.2f) no puede superar el total de la venta (%.2f)", sumPay, i.Total)
		return ErrorResponse(422, message, fmt.Errorf(message))
	}

	return nil
}

type SaleIncomeUpdate struct {
	ID       int64                  `json:"id" validate:"required"`
	Items    []ItemSaleIncomeUpdate `json:"items" validate:"required,dive"`
	Pay      []PayUpdate            `json:"pay" validate:"required,dive,max=3"`
	Discount float32                `json:"discount"`
	Type     string                 `json:"type_discount" validate:"oneof=amount percent"`
	Total    float32                `json:"total"`
	IsBudget bool                   `json:"is_budget"`
}

type ItemSaleIncomeUpdate struct {
	ProductID    int64   `json:"product_id" validate:"required"`
	Amount       float32 `json:"amount" validate:"required"`
	Discount     float32 `json:"discount"`
	TypeDiscount string  `json:"type_discount" validate:"oneof=amount percent"`
}

type PayUpdate struct {
	Amount    float32 `json:"amount" validate:"required"`
	MethodPay string  `json:"method_pay" validate:"oneof=cash credit card transfer"`
}

func (i *SaleIncomeUpdate) Validate() error {
	validate := validator.New()

	if err := validate.Struct(i); err != nil {
		validationErr := err.(validator.ValidationErrors)[0]
		field := validationErr.Field()
		tag := validationErr.Tag()
		param := validationErr.Param()
		return ErrorResponse(422, fmt.Sprintf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param), err)
	}

	var sumPay float32
	for _, p := range i.Pay {
		sumPay += p.Amount
	}

	if sumPay > i.Total {
		message := fmt.Sprintf("la suma de los pagos (%.2f) no puede superar el total de la venta (%.2f)", sumPay, i.Total)
		return ErrorResponse(422, message, fmt.Errorf(message))
	}

	return nil
}

type SaleIncomeResponse struct {
	ID        int64                    `json:"id"`
	Member    MemberSimpleDTO          `json:"member"`
	Client    ClientResponseDTO        `json:"client"`
	Items     []SaleIncomeItemResponse `json:"items"`
	Pay       []PayResponse            `json:"pay"`
	Total     float64                  `json:"total"`
	IsBudget  bool                     `json:"is_budget"`
	CreatedAt time.Time                `json:"created_at"`
}

type SaleIncomeItemResponse struct {
	ID        int64                    `json:"id"`
	Product   ProductSimpleResponseDTO `json:"product"`
	Amount    float64                  `json:"quantity"`
	Price     float64                  `json:"price"`
	Discount  float64                  `json:"discount"`
	Type      string                   `json:"type_discount"`
	SubTotal  float64                  `json:"subtotal"`
	Total     float64                  `json:"total"`
	CreatedAt time.Time                `json:"created_at"`
}

type PayResponse struct {
	ID        int64   `json:"id"`
	Amount    float64 `json:"amount"`
	MethodPay string  `json:"method_pay"`
}


type SaleIncomeResponseDTO struct {
	ID            int64           `json:"id"`
	Member        MemberSimpleDTO `json:"member"`
	Client        ClientSimpleDTO `json:"client"`
	Pay           []PayResponse   `json:"pay"`
	Total         float64         `json:"total"`
	PaymentMethod string          `json:"payment_method"`
	CreatedAt     time.Time       `json:"created_at"`
}

type SaleIncomeSimpleResponse struct {
	ID        int64                    `json:"id"`
	Items     []ProductSimpleResponseDTO `json:"items"`
	Pay       []PayResponse            `json:"pay"`
	Total     float64                  `json:"total"`
	IsBudget  bool                     `json:"is_budget"`
	CreatedAt time.Time                `json:"created_at"`
}

type SaleIncomeItemResponseDTO struct {
	ID        int64                    `json:"id"`
	Product   ProductSimpleResponseDTO `json:"product"`
	Amount    float64                  `json:"quantity"`
	Total     float64                  `json:"total"`
	CreatedAt time.Time                `json:"created_at"`
}