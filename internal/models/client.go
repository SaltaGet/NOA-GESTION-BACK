package models

import (
	"time"
)

type Client struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string    `gorm:"not null;size:30" json:"first_name"`
	LastName  string    `gorm:"not null;size:30" json:"last_name"`
	Cuil      string    `gorm:"not null;unique;size:30" json:"cuil"`
	Dni       string    `gorm:"not null;unique;size:30" json:"dni"`
	Email     string    `gorm:"not null;unique" json:"email" validate:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
