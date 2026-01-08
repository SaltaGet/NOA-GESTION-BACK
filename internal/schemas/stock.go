package schemas

type ProductStockFullResponse struct {
	ID          int64                 `json:"id"`
	Code        string                `json:"code"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Category    CategoryResponseStock `json:"category"`
	Price       float64               `json:"price"`
	Stock       float64               `json:"stock"`
	PrimaryImage *string `json:"primary_image"`
	SecondaryImage []string `json:"secondary_image"`
}

type ProductStockFullResponseCategory struct {
	ID           int64   `json:"id"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	CategoryID   int64   `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Price        float64 `json:"price"`
	Stock        float64 `json:"stock"`
	PrimaryImage *string `json:"primary_image"`
	SecondaryImages *string `json:"secondary_images"`
}
