package models

import (
	"time"
)

type Product struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Identifier string    `gorm:"not null;unique" json:"identifier"`
	Name       string    `gorm:"not null" json:"name"`
	Stock      float32   `gorm:"not null;min:0;default:0" json:"stock"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
