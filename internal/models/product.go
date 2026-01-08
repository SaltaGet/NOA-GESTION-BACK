package models

import (
	"strings"
	"time"

	"github.com/agnivade/levenshtein"
)

type Product struct {
	ID              int64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code            string   `gorm:"size:50;not null;uniqueIndex" json:"code"`
	Name            string   `gorm:"size:100;not null" json:"name"`
	Description     *string  `gorm:"type:text" json:"description"`
	Price           float64  `gorm:"not null,default:0" json:"price"`
	CategoryID      int64    `gorm:"not null" json:"category_id"`
	Category        Category `gorm:"foreignKey:CategoryID;references:ID" json:"category"`
	PrimaryImage    *string  `gorm:"size:255;default:null" json:"primary_image"`
	SecondaryImages *string  `gorm:"type:text;default:null" json:"secondary_images"`
	IsVisible       bool     `gorm:"not null;default:false" json:"is_visible"`

	Notifier  bool    `gorm:"not null;default:false" json:"notifier"`
	MinAmount float64 `gorm:"not null;default:0" json:"min_amount"`

	CreatedAt       time.Time         `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt       time.Time         `gorm:"autoUpdateTime:milli" json:"updated_at"`
	StockDeposit    *Deposit          `gorm:"foreignKey:ProductID" json:"stock_deposit"`
	StockPointSales []*StockPointSale `gorm:"foreignKey:ProductID" json:"stock_point_sales"`
}

type ProductWithScore struct {
	Product *Product
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
