package schemas

import (
	"time"
)

type MovementStockResponse struct {
	ID          int64                    `json:"id"`
	Member      MemberSimpleDTO          `json:"member"`
	Product     ProductSimpleResponseDTO `json:"product"`
	Amount      float64                  `json:"amount"`
	FromID      int64                    `json:"from_id"`
	FromType    string                   `json:"from_type"`
	ToID        int64                    `json:"to_id"`
	ToType      string                   `json:"to_type"`
	IgnoreStock bool                     `json:"ignore_stock"`
	CreatedAt   time.Time                `json:"created_at"`
}

type MovementStockResponseDTO struct {
	ID          int64                    `json:"id"`
	Member      MemberSimpleDTO          `json:"member"`
	Product     ProductSimpleResponseDTO `json:"product"`
	Amount      float64                  `json:"amount"`
	FromID      int64                    `json:"from_id"`
	FromType    string                   `json:"from_type"`
	ToID        int64                    `json:"to_id"`
	ToType      string                   `json:"to_type"`
	IgnoreStock bool                     `json:"ignore_stock"`
	CreatedAt   time.Time                `json:"created_at"`
}

type MovementStock struct {
	ProductID int64   `json:"product_id" validate:"required" example:"1"`
	Amount    float64 `json:"amount" validate:"required" example:"10"`

	FromType string `json:"from_type" validate:"oneof=deposit point_sale" example:"deposit|point_sale"`
	FromID   int64  `json:"from_id" validate:"required" example:"1"`

	ToType string `json:"to_type" validate:"oneof=deposit point_sale" example:"deposit|point_sale"`
	ToID   int64  `json:"to_id" validate:"required" example:"1"`

	IgnoreStock *bool `json:"ignore_stock" validate:"required" example:"false"`
}

type MovementStockList struct {
	ProductID         int64               `json:"product_id" validate:"required" example:"1"`
	MovementStockItem []MovementStockItem `json:"movement_stock_item" validate:"required,dive"`
}

type MovementStockItem struct {
	Amount float64 `json:"amount" validate:"required" example:"10"`

	FromType string `json:"from_type" validate:"oneof=deposit point_sale" example:"deposit|point_sale"`
	FromID   int64  `json:"from_id" validate:"required" example:"1"`

	ToType      string `json:"to_type" validate:"oneof=deposit point_sale" example:"deposit|point_sale"`
	ToID        int64  `json:"to_id" validate:"required" example:"1"`
	IgnoreStock *bool  `json:"ignore_stock" validate:"required" example:"false"`
}

