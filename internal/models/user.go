package models

import "time"
 
type User struct {
	ID        string    `gorm:"primaryKey;size:36" json:"id"`
	FirstName string    `gorm:"not null;size:30" json:"first_name"`
	LastName  string    `gorm:"not null;size:30" json:"last_name"`
	Username  string    `gorm:"unique;size:30;not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email" validate:"email"`
	Password  string    `gorm:"not null" json:"password"`
	IsAdmin      bool    `gorm:"not null;default:false" json:"is_admin"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	UserTenants []UserTenant `gorm:"foreignKey:UserID" json:"user_tenants"`
}