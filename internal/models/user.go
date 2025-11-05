package models

import (
	"time"
)

type User struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName  string    `gorm:"not null;size:30" json:"first_name"`
	LastName   string    `gorm:"not null;size:30" json:"last_name"`
	Email      string    `gorm:"unique;not null" json:"email"`
	Username   string    `gorm:"unique;not null" json:"username"`
	Address    *string   `gorm:"size:255;default:null" json:"address"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
