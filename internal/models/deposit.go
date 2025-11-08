package models

import "time"

type Deposit struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID int64      `gorm:"not null" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID;references:ID" json:"-"`
	Stock     float64   `gorm:"not null;default:0" json:"stock"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`
}