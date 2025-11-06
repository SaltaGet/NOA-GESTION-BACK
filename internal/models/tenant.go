package models

import (
	"time"
)

type Tenant struct {
	ID          int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string       `gorm:"not null" json:"name"`
	Identifier  string       `gorm:"not null;unique" json:"identifier"`
	Address     string       `gorm:"not null" json:"address"`
	Phone       string       `gorm:"not null" json:"phone"`
	Email       string       `gorm:"Index;not null" json:"email"`
	CuitPdv     string       `gorm:"size:50;uniqueIndex;not null" json:"cuit_pdv"`
	IsActive    bool         `gorm:"not null;default:true" json:"is_active"`
	PlanID      int64        `gorm:"not null" json:"plan_id"`
	Connection  string       `gorm:"not null" json:"connection"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	UserTenants []UserTenant `gorm:"foreignKey:TenantID" json:"user_tenants"`
	Plan        Plan         `gorm:"foreignKey:PlanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"plan"`
}
