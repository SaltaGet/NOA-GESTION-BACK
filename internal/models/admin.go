package models

import (
	"strings"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	"gorm.io/gorm"
)

type Admin struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName    string         `gorm:"not null;size:30" json:"first_name"`
	LastName     string         `gorm:"not null;size:30" json:"last_name"`
	Username     string         `gorm:"size:30;not null;uniqueIndex" json:"username"`
	Email        string         `gorm:"unique;not null" json:"email"`
	Password     string         `gorm:"not null" json:"password"`
	IsSuperAdmin bool           `gorm:"not null;default:false" json:"is_admin"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (u *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword

	u.Email = strings.ToLower(u.Email)
	u.FirstName = strings.ToLower(u.FirstName)
	u.LastName = strings.ToLower(u.LastName)

	return
}

func (u *Admin) BeforeUpdate(tx *gorm.DB) (err error) {
	if tx.Statement.Changed("Password") {
		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}
	return
}
