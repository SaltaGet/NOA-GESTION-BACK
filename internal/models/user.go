package models

import (
	"strings"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	"gorm.io/gorm"
)

type User struct {
	ID            int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName     string       `gorm:"not null;size:30" json:"first_name"`
	LastName      string       `gorm:"not null;size:30" json:"last_name"`
	Username      string       `gorm:"size:30;not null;uniqueIndex" json:"username"`
	Email         string       `gorm:"unique;not null" json:"email"`
	Password      string       `gorm:"not null" json:"password"`
	Address       *string      `gorm:"size:255;default:null" json:"address"`
	IsAdmin       bool         `gorm:"not null;default:false" json:"is_admin"`
	IsAdminTenant bool         `gorm:"not null;default:false" json:"is_admin_tenant"`
	CreatedAt     time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	UserTenants   []UserTenant `gorm:"foreignKey:UserID" json:"user_tenants"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword

	var address string
	if address == "" {
		address = "Sin dirección"
	} else {
		address = strings.ToLower(*u.Address)
	}

	u.Email = strings.ToLower(u.Email)
	u.FirstName = strings.ToLower(u.FirstName)
	u.LastName = strings.ToLower(u.LastName)
	u.Address = &address

	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if tx.Statement.Changed("Password") {
		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}
	return
}
