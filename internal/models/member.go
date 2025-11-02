package models

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Member struct {
	ID        string `gorm:"primaryKey;size:36" json:"id"`
	FirstName string `gorm:"not null;size:30" json:"first_name"`
	LastName  string `gorm:"not null;size:30" json:"last_name"`
	Username  string `gorm:"unique;size:30;not null" json:"username"`
	Email     string `gorm:"unique;not null" json:"email" validate:"email"`
	Password  string `gorm:"not null" json:"password"`
	IsActive  bool   `gorm:"not null;default:true" json:"is_active"`
	RoleID    string `gorm:"not null;size:36" json:"role_id"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Role      Role      `gorm:"foreignKey:RoleID;references:ID" json:"role"`
}

type MemberCreate struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	RoleID    string `json:"role_id" validate:"required"`
}

func (m *MemberCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(m)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type MemberResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email" validate:"email"`
	IsActive  bool   `json:"is_active"`
	CreatedAt time.Time `son:"created_at"`
	Role    Role `json:"role_id"`
}

type MemberDTO struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email" validate:"email"`
	IsActive  bool   `json:"is_active"`
	CreatedAt time.Time `son:"created_at"`
	Role    RoleDTO `json:"role_id"`
}
