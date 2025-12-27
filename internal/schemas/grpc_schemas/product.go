package grpc_schemas

type CategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ProductResponse struct {
	ID          int64            `json:"id"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	Category    CategoryResponse `json:"category"`
	Price       float64          `json:"price"`
	Stock       float64          `json:"stock"`
	UrlImage    []string         `json:"url_image"`
}

type ProductResponseDTO struct {
	ID          int64            `json:"id"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	Category    CategoryResponse `json:"category"`
	Price       float64          `json:"price"`
	Stock       float64          `json:"stock"`
	UrlImage    *string          `json:"url_image"`
}