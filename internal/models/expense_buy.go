package models

import "time"

type Expense struct {
	ID           int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	Details      string       `json:"details"`
	Amount       float32      `gorm:"not null" json:"amount"`
	CreatedAt    time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	MovementType MovementType `gorm:"foreignKey:MovementTypeID;references:ID" json:"movement_type"`
}
