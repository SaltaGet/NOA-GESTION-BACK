package models

import "time"

type PayTenant struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID    int64     `gorm:"not null" json:"tenant_id"`
	AdminID     int64     `gorm:"not null" json:"admin_id"`
	AmountMonth int64     `gorm:"not null" json:"amount"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	Tenant    Tenant    `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE;" json:"-"`
	Admin     Admin     `gorm:"foreignKey:AdminID;references:ID" json:"admin"`
	PayDetail PayDetail `gorm:"foreignKey:PayTenantID" json:"detail"`
}

type PayDetail struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PayTenantID int64     `gorm:"not null" json:"pay_tenant_id"`
	PayID       *string   `json:"pay_id"`
	Amount      float64   `gorm:"not null" json:"amount"`
	MethodPay   string    `gorm:"not null;default:cash" json:"method_pay"`
	StatePay    string    `gorm:"not null;default:pending" json:"state_pay"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
