package schemas

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type ModuleResponse struct {
	ID                     int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                   string  `gorm:"not null;uniqueIndex" json:"name"`
	PriceMonthly           float64 `gorm:"type:decimal(10,2)" json:"price"`
	PriceYearly            float64 `gorm:"type:decimal(10,2)" json:"price_yearly"`
	Description            string  `gorm:"type:text" json:"description"`
	Features               string  `gorm:"type:text" json:"features,omitempty"`
	AmountImagesPerProduct int32   `gorm:"not null" json:"amount_images_per_product"`
}

type ModuleResponseDTO struct {
	ID                     int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                   string     `gorm:"not null;uniqueIndex" json:"name"`
	AmountImagesPerProduct int32      `gorm:"not null" json:"amount_images_per_product"`
	Expiration             *time.Time `gorm:"" json:"expiration"`
}

type ModuleCreate struct {
	Name                   string  `json:"name" validate:"required" example:"Module1"`
	PriceMonthly           float64 `json:"price" validate:"required,gte=0" example:"100.00"`
	PriceYearly            float64 `json:"price_yearly" validate:"required,gte=0" example:"1000.00"`
	Description            *string  `json:"description" example:"description"`
	Features               string  `json:"features,omitempty" example:"features"`
	AmountImagesPerProduct int32   `json:"amount_images_per_product" validate:"required" example:"1"`
}

func (m *ModuleCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(m)
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

type ModuleUpdate struct {
	ID                     int64   `json:"id" validate:"required" example:"1"`
	Name                   string  `json:"name" validate:"required" example:"Module1"`
	PriceMonthly           float64 `json:"price" validate:"required,gte=0" example:"100.00"`
	PriceYearly            float64 `json:"price_yearly" validate:"required,gte=0" example:"1000.00"`
	Description            *string  `json:"description" example:"description"`
	Features               string  `json:"features,omitempty" example:"features"`
	AmountImagesPerProduct int32   `json:"amount_images_per_product" validate:"required" example:"1"`
}

func (m *ModuleUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(m)
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

type ModuleAddTenant struct {
	ModuleID int64 `json:"module_id" validate:"required" example:"1"`
	TenantID int64 `json:"tenant_id" validate:"required" example:"1"`
	Expiration string `json:"expiration" validate:"required,datetime=2006-01-02" example:"2023-01-01"`
}

func (m *ModuleAddTenant) Validate() error {
	validate := validator.New()
	err := validate.Struct(m)
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