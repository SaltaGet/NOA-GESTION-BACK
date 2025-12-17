package models

import (
	"strings"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	"gorm.io/gorm"
)

type Member struct {
	ID        int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string  `gorm:"not null;size:30" json:"first_name"`
	LastName  string  `gorm:"not null;size:30" json:"last_name"`
	Username  string  `gorm:"unique;size:30;not null" json:"username"`
	Email     string  `gorm:"unique;not null" json:"email" validate:"email"`
	Password  string  `gorm:"not null" json:"password"`
	Address   *string `gorm:"size:255;default:null" json:"address"`
	Phone     *string `gorm:"size:20;default:null" json:"phone"`
	IsAdmin   bool    `gorm:"not null;default:false" json:"is_admin"`
	IsActive  bool    `gorm:"not null;default:true" json:"is_active"`
	RoleID    int64   `gorm:"not null;size:36" json:"role_id"`

	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Role       Role           `gorm:"foreignKey:RoleID;references:ID" json:"role"`
	PointSales []PointSale    `gorm:"many2many:member_point_sales;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"point_sales"`
}

func (u *Member) BeforeCreate(tx *gorm.DB) (err error) {
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

// func (u *Member) BeforeUpdate(tx *gorm.DB) (err error) {
// 	if tx.Statement.Changed("Password") {
// 		hashedPassword, err := utils.HashPassword(u.Password)
// 		if err != nil {
// 			return err
// 		}
// 		u.Password = hashedPassword
// 	}
// 	return
// }
