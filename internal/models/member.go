package models

import (
	"gorm.io/gorm"
	"time"
)

type Member struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string `gorm:"not null;size:30" json:"first_name"`
	LastName  string `gorm:"not null;size:30" json:"last_name"`
	Username  string `gorm:"unique;size:30;not null" json:"username"`
	Email     string `gorm:"unique;not null" json:"email" validate:"email"`
	Password  string `gorm:"not null" json:"password"`
	IsActive  bool   `gorm:"not null;default:true" json:"is_active"`
	RoleID    int64  `gorm:"not null;size:36" json:"role_id"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Role      Role           `gorm:"foreignKey:RoleID;references:ID" json:"role"`
}
