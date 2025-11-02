package schemas

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type Expense struct {
	ID              string         `gorm:"primaryKey;size:36" json:"id"`
	Details         string         `json:"details"`
	PurchaseOrderID *string        `gorm:"size:36" json:"purchase_order_id"`
	MovementTypeID  string         `gorm:"not null;size:36" json:"movement_type_id"`
	Amount          float32        `gorm:"not null" json:"amount"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	PurchaseOrder   *PurchaseOrder `gorm:"foreignKey:PurchaseOrderID" json:"purchase_order"`
	MovementType    MovementType   `gorm:"foreignKey:MovementTypeID;references:ID" json:"movement_type"`
}

type ExpenseCreate struct {
	Details        string  `json:"details" validate:"required"`
	PurchaseOrderID *string `json:"purchase_order_id"`
	MovementTypeID string  `json:"movement_type_id" validate:"required"`
	Amount         float32 `json:"amount" validate:"required"`
}

func (e *ExpenseCreate) Validate() error {
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

type ExpenseUpdate struct {
	ID      string `json:"id"`
	Details string `json:"details" validate:"required"`
	PurchaseOrderID *string `json:"purchase_order_id"`
	MovementTypeID string  `json:"movement_type_id" validate:"required"`
	Amount         float32 `json:"amount" validate:"required"`
}

func (e *ExpenseUpdate) Validate() error {
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

type ExpenseResponse struct {
	ID            string          `json:"id"`
	Details       string          `json:"details"`
	Amount        float32         `json:"amount"`
	CreatedAt     time.Time       `json:"created_at"`
	PurchaseOrder *PurchaseOrderResponse  `json:"purchase_order"`
	MovementType  MovementTypeDTO `json:"movement_type"`
}

type ExpenseDTO struct {
	ID            string            `json:"id"`
	Amount        float32           `json:"amount"`
	CreatedAt     time.Time         `json:"created_at"`
	PurchaseOrder *PurchaseOrderDTO `json:"purchase_order"`
	MovementType  MovementTypeDTO   `json:"movement_type"`
}
