package models

import "time"

type Invoice struct {
	ID string `gorm:"primaryKey" json:"id"`
    ClientID *int64 `gorm:"not null;index" json:"client_id"`
    ClientEcommerceID *int64 `gorm:"not null;index" json:"client_ecommerce_id"`
	TypeVoucher int64 `gorm:"not null" json:"type_voucher"`
	PointSale int64 `gorm:"not null" json:"point_sale"`
    NumberVoucher int64 `gorm:"not null" json:"number_voucher"`
    DateVoucher string `gorm:"not null" json:"date_voucher"`
    Cae string `gorm:"not null" json:"cae"`
    CaeMaturity string `gorm:"not null" json:"cae_maturity"`
    Status string `gorm:"not null" json:"status"`
    IsTest bool `gorm:"not null;default:false" json:"is_test"`
    Subtotal float64 `gorm:"not null" json:"subtotal"`
    Description *string `gorm:"type:text" json:"description"`
    Discount float64 `gorm:"not null;default:0" json:"discount"`
    TypeDiscount string `gorm:"not null;default:percent" json:"type_discount" validate:"oneof=amount percent"`
    Total float64 `gorm:"not null" json:"total"`
    IsBudget bool `gorm:"not null;default:false" json:"is_budget"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
