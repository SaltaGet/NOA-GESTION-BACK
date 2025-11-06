package models

import (
	"time"

	"gorm.io/gorm"
)

type Client struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName      string    `gorm:"not null;size:30" json:"first_name"`
	LastName       string    `gorm:"not null;size:30" json:"last_name"`
	CompanyName    string    `gorm:"not null;size:100" json:"company_name"`
	Identifier     *string   `gorm:"unique;size:20" json:"cuit,omitempty"`
	Email          *string   `gorm:"unique" json:"email" validate:"email"`
	Phone          *string   `json:"phone"`
	Address        *string   `json:"address"`
	MemberCreateID int64     `gorm:"not null" json:"member_create_id"`
	MemberCreate   Member    `gorm:"foreignKey:MemberCreateID" json:"member_create"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
