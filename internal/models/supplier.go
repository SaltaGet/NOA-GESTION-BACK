package models

import (
	"time"
)

// Proveedor
type Supplier struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Address   string    `gorm:"not null" json:"address"`
	Phone     string    `gorm:"not null" json:"phone"`
	Email     string    `gorm:"not null" json:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
