package models

import "time"

type StockPointSale struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID   uint      `gorm:"not null" json:"product_id"`
	Product     Product   `gorm:"foreignKey:ProductID;references:ID" json:"-"`
	PointSaleID uint      `gorm:"not null" json:"point_sale_id"`
	PointSale   PointSale `gorm:"foreignKey:PointSaleID;references:ID" json:"point_sale"`
	Stock       float64   `gorm:"not null;default:0" json:"stock"`
	CreatedAt   time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`
}
