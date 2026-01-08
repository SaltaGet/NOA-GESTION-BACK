package models

import "time"

type SettingTenant struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID       int64     `gorm:"uniqueIndex;not null" json:"tenant_id"`
	Logo           *string   `gorm:"type:varchar(255)" json:"logo"`
	FrontPage      *string   `gorm:"type:varchar(255)" json:"front_page"`
	Title          *string   `gorm:"type:varchar(255)" json:"title"`
	Slogan         *string   `gorm:"type:text" json:"slogan"`
	PrimaryColor   *string   `gorm:"type:varchar(255)" json:"primary_color"`
	SecondaryColor *string   `gorm:"type:varchar(255)" json:"secondary_color"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
