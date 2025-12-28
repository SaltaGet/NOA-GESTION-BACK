package schemas

import "time"

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
