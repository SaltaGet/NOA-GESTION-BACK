package models

import (
	"time"

	"gorm.io/gorm"
)

type Module struct {
	ID                     int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                   string         `gorm:"not null;uniqueIndex" json:"name"`
	PriceMonthly           float64        `gorm:"type:decimal(10,2)" json:"price"`
	PriceYearly            float64        `gorm:"type:decimal(10,2)" json:"price_yearly"`
	Description            string         `gorm:"type:text" json:"description"`
	Features               string         `gorm:"type:text" json:"features,omitempty"`
	AmountImagesPerProduct int32          `gorm:"not null" json:"amount_images_per_product"`
	CreatedAt              time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Tenants                []TenantModule `gorm:"foreignKey:ModuleID" json:"tenants"`
}

type TenantModule struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ModuleID      int64          `gorm:"not null;uniqueIndex:idx_tenant_module"`
	TenantID      int64          `gorm:"not null;uniqueIndex:idx_tenant_module"`
	Module        Module         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tenant        Tenant         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Expiration    *time.Time     `gorm:"" json:"expiration"`
	AcceptedTerms bool           `gorm:"not null;default:false" json:"accepted_terms"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
