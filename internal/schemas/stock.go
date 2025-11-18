package schemas

type ProductStockFullResponse struct {
	ID          int64                 `json:"id"`
	Code        string                `json:"code"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Category    CategoryResponseStock `json:"category"`
	Price       float64               `json:"price"`
	Stock       float64               `json:"stock"`
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
}
