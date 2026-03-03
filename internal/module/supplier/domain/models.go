package domain


import (
	"time"

	"gorm.io/gorm"
)

// Proveedor
type Supplier struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	CompanyName string         `gorm:"type:varchar(100);not null" json:"company_name"`
	Identifier  *string        `gorm:"type:varchar(20);unique" json:"cuit,omitempty"`
	Address     *string        `gorm:"type:varchar(255)" json:"address"`
	DebtLimit   *float64       `json:"debt_limit,omitempty"`
	Email       *string        `gorm:"type:varchar(100);unique" json:"email" validate:"email,omitempty"`
	Phone       *string        `gorm:"type:varchar(20)" json:"phone,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
