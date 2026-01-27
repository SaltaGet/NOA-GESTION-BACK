package schemas

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type EcommerceResponse struct {
	ID                int64                   `json:"id"`
	PaymentID         string                  `json:"payment_id"`
	ExternalReference string                  `json:"external_reference"`
	Status            string                  `json:"status"`
	Total             float64                 `json:"total"`
	DeliveryStatus    string                  `json:"delivery_status"`
	DeliveryID        *string                 `json:"delivery_id"`
	DateCreated       string                  `json:"date_created"`
	DateApproved      string                  `json:"date_approved"`
	TransactionAmount float64                 `json:"transaction_amount"`
	NetReceivedAmount float64                 `json:"net_received_amount"`
	PayerFirstName    string                  `json:"payer_first_name"`
	PayerLastName     string                  `json:"payer_last_name"`
	PayerEmail        string                  `json:"payer_email"`
	PayMethod         string                  `json:"pay_method"`
	OperationType     string                  `json:"operation_type"`
	Message           string                  `json:"message,omitempty"`
	Items             []EcommerceItemResponse `json:"items"`
}

type EcommerceItemResponse struct {
	Product    ProductSimpleResponseDTO   `json:"product"`
	Amount       float64 `json:"amount"`
	Price_Cost   float64 `json:"price_cost"`
	Price        float64 `json:"price"`
	Discount     float64 `json:"discount"`
	TypeDiscount string  `json:"type_discount"`
	Subtotal     float64 `json:"subtotal"`
	Total        float64 `json:"total"`
}

type EcommerceResponseDTO struct {
	ID                int64   `json:"id"`
	ExternalReference string  `json:"external_reference"`
	Status            string  `json:"status"`
	Total             float64 `json:"total"`
	DateCreated       string  `json:"date_created"`
	PayerEmail        string  `json:"payer_email"`
}

type EcommerceStatusUpdate struct {
	ID                int64  `json:"id" validate:"required"`
	NewStatus         string `json:"new_status" validate:"required"`
}

func (e *EcommerceStatusUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(e)
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