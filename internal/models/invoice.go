package models

import (
	"time"

	"gorm.io/datatypes"
)

type Invoice struct {
	ID          int64       `gorm:"primaryKey" json:"id"`
	InvoiceData datatypes.JSON `gorm:"type:jsonb" json:"invoice_data"`
	CreatedAt   time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

