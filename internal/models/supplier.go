package models

import (
	"time"

	"gorm.io/gorm"
)

// Proveedor
type Supplier struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	CompanyName string         `gorm:"not null;size:100" json:"company_name"`
	Identifier  *string        `gorm:"unique;size:20" json:"cuit,omitempty"`
	Address     *string        `gorm:"size:150" json:"address"`
	DebtLimit   *float64       `json:"debt_limit,omitempty"`
	Email       *string        `gorm:"unique" json:"email" validate:"email,omitempty"`
	Phone       *string        `gorm:"size:20" json:"phone,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
