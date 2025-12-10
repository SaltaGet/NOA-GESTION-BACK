package models

import "time"

type Plan struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string    `gorm:"type:varchar(255);uniqueIndex" json:"name"`
	PriceMounthly   float64   `gorm:"type:decimal(10,2)" json:"price"`
	PriceYearly     float64   `gorm:"type:decimal(10,2)" json:"price_yearly"`
	Description     string    `gorm:"type:text" json:"description"`
	Features        string    `gorm:"type:text" json:"features,omitempty"`
	AmountPointSale int64     `gorm:"not null" json:"amount_point_sale"`
	AmountMember    int64     `gorm:"not null" json:"amount_member"`
	AmountProduct   int64     `gorm:"not null" json:"amount_product"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Tenants         []Tenant  `gorm:"foreignKey:PlanID" json:"tenants,omitempty"`
}
