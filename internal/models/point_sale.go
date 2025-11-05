package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type PointSale struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Description *string   `gorm:"size:200" json:"description"`
	IsDeposit   bool      `gorm:"not null;default:false" json:"is_deposit"`
	CreatedAt   time.Time  `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`
	Members       []Member `gorm:"many2many:member_point_sales;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"members"`
}

func (p *PointSale) BeforeCreate(tx *gorm.DB) (err error) {
	p.Name = strings.ToLower(p.Name)
	return
}