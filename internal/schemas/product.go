package schemas

import (
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/agnivade/levenshtein"
)

type ProductFullResponse struct {
	ID              int64                     `json:"id"`
	Code            string                    `json:"code"`
	Name            string                    `json:"name"`
	Description     *string                   `json:"description"`
	Category        CategoryResponse          `json:"category"`
	Price           float64                   `json:"price"`
	StockPointSales []*PointSaleStockResponse `json:"stock_point_sales"`
	StockDeposit    float64                   `json:"stock_deposit"`
	Notifier        bool                      `json:"notifier"`
	MinAmount       float64                   `json:"min_amount"`
	PrimaryImage    *string                   `json:"primary_image"`
	SecondaryImage  []string                  `json:"secondary_image"`
	IsVisible       bool                      `json:"is_visible"`
}

type ProductResponse struct {
	ID             int64            `json:"id"`
	Code           string           `json:"code"`
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	Category       CategoryResponse `json:"category"`
	Price          float64          `json:"price"`
	Stock          float64          `json:"stock"`
	Notifier       bool             `json:"notifier"`
	MinAmount      float64          `json:"min_amount"`
	PrimaryImage   *string          `json:"primary_image"`
	SecondaryImage []string         `json:"secondary_image"`
}

type ProductResponseDTO struct {
	ID        int64             `json:"id"`
	Code      string            `json:"code"`
	Name      string            `json:"name"`
	Category  *CategoryResponse `json:"category,omitempty"`
	Price     float64           `json:"price"`
	Stock     float64           `json:"stock"`
	Notifier  bool              `json:"notifier"`
	MinAmount float64           `json:"min_amount"`
}

type ProductSimpleResponse struct {
	ID        int64   `json:"id"`
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Stock     float64 `json:"stock"`
	Notifier  bool    `json:"notifier"`
	MinAmount float64 `json:"min_amount"`
}

type ProductSimpleResponseDTO struct {
	ID    int64   `json:"id"`
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductCreate struct {
	Code        string   `json:"code" validate:"required" example:"ABC123"`
	Name        string   `json:"name" validate:"required" example:"Producto1"`
	Description *string  `json:"description" example:"description|null"`
	CategoryID  int64    `json:"category_id" validate:"required" example:"1"`
	Price       *float64 `json:"price" example:"100.00"`
	Notifier    bool     `json:"notifier" example:"false"`
	MinAmount   float64  `json:"min_amount" example:"10.00"`
}

type ProductUpdate struct {
	ID          int64    `json:"id" validate:"required" example:"1"`
	Code        string   `json:"code" validate:"required" example:"ABC123"`
	Name        string   `json:"name" validate:"required" example:"Producto1"`
	Description *string  `json:"description" example:"description|null"`
	CategoryID  uint     `json:"category_id" validate:"required" example:"1"`
	Price       *float64 `json:"price" example:"100.00"`
	Notifier    bool     `json:"notifier" example:"false"`
	MinAmount   float64  `json:"min_amount" example:"10.00"`
}

type ProductPriceUpdate struct {
	ID    int64   `json:"id" validate:"required" example:"1"`
	Price float64 `json:"price" validate:"required,gte=0" example:"100.00"`
}

type ListPriceUpdate struct {
	ListProductPriceUpdate []ProductPriceUpdate `json:"list" validate:"required,min=1,dive"`
}

type ProductStockWithScore struct {
	Product *ProductStockFullResponse
	Score   float64
	Length  int
}

type ProductValidateImage struct {
	ProductID      int64                  `json:"product_id" validate:"required" example:"1"`
	PrimaryImage   string                 `json:"primary_image" validate:"required,oneof=set keep" example:"set | keep"`
	SecondaryImage ValidateSecondaryImage `json:"secondary_image" validate:"required"`
}

type ValidateSecondaryImage struct {
	Add         *int64   `json:"add" example:"1"`
	KeepUUIDs   []string `json:"keep_uuids" validate:"required" example:"uuid1,uuid2,uuid3"`
	RemoveUUIDs []string `json:"remove_uuids" example:"uuid1,uuid2,uuid3"`
}

type ListVisibilityUpdate struct {
	ListProductVisibilityUpdate []ProductVisibilityUpdate `json:"list" validate:"required,min=1,dive"`
}

type ProductVisibilityUpdate struct {
	ProductID  int64 `json:"product_id" validate:"required" example:"1"`
	Visibility *bool `json:"visibility" validate:"required" example:"true"`
}

type ProductWithScore struct {
	Product *tenant.Product
	Score   float64
	Length  int
}

func CalculateRelevance(search, target string) float64 {
	if search == target {
		return 100.0
	}

	// 2. SEGUNDO: Empieza con el término (score 90)
	if strings.HasPrefix(target, search) {
		return 90.0
	}

	// 3. TERCERO: Contiene el término (score 80)
	if strings.Contains(target, search) {
		return 80.0
	}

	// 4. CUARTO: Similitud por Levenshtein (score 60-79)
	distance := levenshtein.ComputeDistance(search, target)
	maxLen := float64(max(len(search), len(target)))

	if maxLen == 0 {
		return 0
	}

	similarity := (1.0 - float64(distance)/maxLen) * 100

	// Mapear similitud al rango 60-79 para que esté después de "contiene"
	if similarity < 60 {
		return 0
	}

	// Escalar de 60-100 a 60-79
	return 60.0 + (similarity-60.0)*0.475
}

// Función auxiliar para obtener el máximo
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type ProductExcelCreate struct {
	Name        string
	Code        string
	Description *string
	Price       float64
	MinAmount   float64
	Category    string
	Stock       float64
}
