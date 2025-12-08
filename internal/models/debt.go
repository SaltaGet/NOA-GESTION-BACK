package models

import "time"

type Debt struct {
	ID        int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID  int64   `gorm:"not null" json:"client_id"`
	Client    Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`
	IncomeSaleID int64   `gorm:"not null" json:"income_sale_id"`
	IncomeSale IncomeSale `gorm:"foreignKey:IncomeSaleID;references:ID" json:"income_sale"`
	IncomeOtherID *int64 `json:"income_other_id"`
	IncomeOther *IncomeOther `gorm:"foreignKey:IncomeOtherID;references:ID" json:"income_other"`
	Total     float64 `gorm:"not null" json:"total"`
	IsCancel  bool    `gorm:"not null;default:false" json:"is_cancel"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`
}