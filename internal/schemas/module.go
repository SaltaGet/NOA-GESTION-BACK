package schemas

import (
	"time"

)

type ModuleResponse struct {
	ID                     int64   `json:"id"`
	Name                   string  `json:"name"`
	PriceMonthly           float64 `json:"price"`
	PriceYearly            float64 `json:"price_yearly"`
	Description            string  `json:"description"`
	Features               string  `json:"features,omitempty"`
	AmountImagesPerProduct int32   `json:"amount_images_per_product"`
}

type ModuleResponseDTO struct {
	ID                     int64      `json:"id"`
	Name                   string     `json:"name"`
	AmountImagesPerProduct int32      `json:"amount_images_per_product"`
	Expiration             *time.Time `json:"expiration"`
	AcceptTerms            bool       `json:"accept_terms"`
}

type ModuleCreate struct {
	Name                   string  `json:"name" validate:"required" example:"Module1"`
	PriceMonthly           float64 `json:"price" validate:"required,gte=0" example:"100.00"`
	PriceYearly            float64 `json:"price_yearly" validate:"required,gte=0" example:"1000.00"`
	Description            *string `json:"description" example:"description"`
	Features               string  `json:"features,omitempty" example:"features"`
	AmountImagesPerProduct int32   `json:"amount_images_per_product" validate:"required" example:"1"`
}

type ModuleUpdate struct {
	ID                     int64   `json:"id" validate:"required" example:"1"`
	Name                   string  `json:"name" validate:"required" example:"Module1"`
	PriceMonthly           float64 `json:"price" validate:"required,gte=0" example:"100.00"`
	PriceYearly            float64 `json:"price_yearly" validate:"required,gte=0" example:"1000.00"`
	Description            *string `json:"description" example:"description"`
	Features               string  `json:"features,omitempty" example:"features"`
	AmountImagesPerProduct int32   `json:"amount_images_per_product" validate:"required" example:"1"`
}

type ModuleAddTenant struct {
	ModuleID   int64  `json:"module_id" validate:"required" example:"1"`
	TenantID   int64  `json:"tenant_id" validate:"required" example:"1"`
	Expiration string `json:"expiration" validate:"required,datetime=2006-01-02" example:"2023-01-01"`
}
