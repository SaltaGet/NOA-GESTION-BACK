package domain

import (
	"time"
)

type User struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName  string    `gorm:"type:varchar(30);not null" json:"first_name"`
	LastName   string    `gorm:"type:varchar(30);not null" json:"last_name"`
	Email      string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Username   string    `gorm:"type:varchar(50);unique;not null" json:"username"`
	Address    *string   `gorm:"type:varchar(255);default:null" json:"address"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
