package schemas

type DepositResponse struct {
	ID          int64   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description *string `json:"description"`
	Category    CategoryResponse `json:"category"`
	Price       float64 `json:"price"`
	Stock       float64    `json:"stock"`
	PrimaryImage *string `json:"primary_image"`
	SecondaryImage []string `json:"secondary_image"`
}

type DepositResponseStock struct {
	ID          int64   `json:"id"`
	Stock       float64    `json:"stock"`
}

type DepositUpdateStock struct {
	ProductID int64    `json:"product_id" validate:"required" example:"1"`
	Stock     *float64 `json:"stock" validate:"required,gte=0" example:"10"`
	Method    string   `json:"method" validate:"oneof=add subtract set" example:"add|subtract|set"`
}
