package models

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Identifier  string         `gorm:"not null;unique" json:"identifier"`
	Address     string         `gorm:"not null" json:"address"`
	Phone       string         `gorm:"not null" json:"phone"`
	Email       string         `gorm:"uniqueIndex;not null" json:"email"`
	CuitPdv     string         `gorm:"size:50;uniqueIndex;not null" json:"cuit_pdv"`
	IsActive    bool           `gorm:"not null;default:true" json:"is_active"`
	PlanID      int64          `gorm:"not null" json:"plan_id"`
	Connection  string         `gorm:"not null" json:"connection"`
	Expiration  *time.Time     `gorm:"" json:"expiration"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	PayTenant   []PayTenant    `gorm:"foreignKey:TenantID" json:"pay_tenants"`
	UserTenants []UserTenant   `gorm:"foreignKey:TenantID" json:"user_tenants"`
	Plan        Plan           `gorm:"foreignKey:PlanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"plan"`
	Modules     []TenantModule `gorm:"foreignKey:TenantID" json:"modules"`
}

func (t *Tenant) BeforeCreate(tx *gorm.DB) (err error) {
	t.Identifier = strings.ToLower(strings.TrimSpace(t.Identifier))

	var validSubdomain = regexp.MustCompile(`^[a-z0-9-]+$`)

	if !validSubdomain.MatchString(t.Identifier) {
		return errors.New("identifier invalid - only lowercase letters, numbers, and hyphens are allowed")
	}

	return nil
}
