package models

import "time"

type Product struct {
	ID          int64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string   `gorm:"size:50;not null;uniqueIndex" json:"code"`
	Name        string   `gorm:"size:100;not null" json:"name"`
	Description *string  `gorm:"size:200" json:"description"`
	Price       float64  `gorm:"not null,default:0" json:"price"`
	CategoryID  int64    `gorm:"not null" json:"category_id"`
	Category    Category `gorm:"foreignKey:CategoryID;references:ID" json:"category"`

	Notifier  bool    `gorm:"not null;default:false" json:"notifier"`
	MinAmount float64 `gorm:"not null;default:0" json:"min_amount"`

	CreatedAt       time.Time         `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt       time.Time         `gorm:"autoUpdateTime:milli" json:"updated_at"`
	StockDeposit    *Deposit          `gorm:"foreignKey:ProductID" json:"stock_deposit"`
	StockPointSales []*StockPointSale `gorm:"foreignKey:ProductID" json:"stock_point_sales"`
}
