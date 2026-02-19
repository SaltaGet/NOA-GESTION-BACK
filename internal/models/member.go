package models

import (
	"strings"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	"gorm.io/gorm"
)

type Member struct {
	ID        int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string  `gorm:"not null;type:varchar(30)" json:"first_name"`
	LastName  string  `gorm:"not null;type:varchar(30)" json:"last_name"`
	Username  string  `gorm:"unique;not null;type:varchar(30)" json:"username"`
	Email     string  `gorm:"unique;not null;type:varchar(100)" json:"email" validate:"email"`
	Password  string  `gorm:"not null;type:varchar(255)" json:"password"`
	Address   *string `gorm:"default:null;type:varchar(255)" json:"address"`
	Phone     *string `gorm:"default:null;type:varchar(20)" json:"phone"`
	IsAdmin   bool    `gorm:"not null;default:false" json:"is_admin"`
	IsActive  bool    `gorm:"not null;default:true" json:"is_active"`
	RoleID    int64   `gorm:"not null;type:varchar(36)" json:"role_id"`

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
