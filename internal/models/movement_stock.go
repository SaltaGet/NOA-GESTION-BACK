package models

import "time"

type MovementStock struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	MemberID    uint      `gorm:"not null" json:"member_id"`
	Member      Member    `gorm:"foreignKey:MemberID;references:ID" json:"member"`
	ProductID   uint      `gorm:"not null" json:"product_id"`
	Product     Product   `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	Amount      float64   `gorm:"not null" json:"amount"`
	FromID      uint      `gorm:"not null" json:"from_id"`
	FromType    string    `gorm:"not null" json:"from_type" validate:"oneof=deposit point_sale"`
	ToID        uint      `gorm:"not null" json:"to_id"`
	ToType      string    `gorm:"not null" json:"to_type" validate:"oneof=deposit point_sale"`
	IgnoreStock bool      `gorm:"not null;default:false" json:"ignore_stock"`
	CreatedAt   time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
}
