package models

import "time"

type CashRegister struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PointSaleID  int64     `gorm:"not null" json:"point_sale_id"`
	PointSale    PointSale `gorm:"foreignKey:PointSaleID;references:ID" json:"point_sale"`
	MemberOpenID int64     `gorm:"not null" json:"user_open_id"`
	MemberOpen   Member    `gorm:"foreignKey:MemberOpenID;references:ID" json:"member_open"`
	OpenAmount   float64   `gorm:"" json:"open_amount"`
	HourOpen     time.Time `gorm:"" json:"hour_open"`

	MemberCloseID *int64     `json:"user_close_id"`
	MemberClose   *Member    `gorm:"foreignKey:MemberCloseID;references:ID" json:"member_close"`
	CloseAmount   *float64   `gorm:"" json:"close_amount"`
	HourClose     *time.Time `gorm:"" json:"hour_close"`

	IsClose   bool      `gorm:"not null,default:false" json:"is_close"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`
}
