package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type PointSale struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Description *string        `gorm:"type:varchar(200)" json:"description"`
	Number      int64          `gorm:"not null" json:"number"`
	IsDeposit   bool           `gorm:"not null;default:false" json:"is_deposit"`
	IsMain      bool           `gorm:"not null;default:false" json:"is_main"`
	CreatedAt   time.Time      `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeleteAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Members     []Member       `gorm:"many2many:member_point_sales;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"members"`
}

func (p *PointSale) BeforeCreate(tx *gorm.DB) (err error) {
	p.Name = strings.ToLower(p.Name)
	return
}
