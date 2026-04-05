package schemas

import "github.com/volatiletech/null/v8"

type ProductStockFullResponse struct {
	ID             int64                 `json:"id"`
	Code           string                `json:"code"`
	Name           string                `json:"name"`
	Description    string                `json:"description"`
	Category       CategoryResponseStock `json:"category"`
	Price          float64               `json:"price"`
	Stock          float64               `json:"stock"`
	PrimaryImage   *string               `json:"primary_image"`
	SecondaryImage []string              `json:"secondary_image"`
}

type ProductStockFullResponseCategory struct {
	ID              int64       `json:"id" boil:"id"`
	Code            string      `json:"code" boil:"code"`
	Name            string      `json:"name" boil:"name"`
	Description     null.String `json:"description" boil:"description"`
	CategoryID      int64       `json:"category_id" boil:"category_id"`
	CategoryName    string      `json:"category_name" boil:"category_name"`
	Price           float64     `json:"price" boil:"price"`
	Stock           float64     `json:"stock" boil:"stock"`
	PrimaryImage    null.String `json:"primary_image" boil:"primary_image"`
	SecondaryImages null.String `json:"secondary_images" boil:"secondary_images"`
}
