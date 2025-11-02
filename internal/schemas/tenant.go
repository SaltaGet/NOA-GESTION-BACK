package schemas

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type Tenant struct {
	ID          string       `gorm:"primaryKey;size:36" json:"id"`
	Name        string       `gorm:"not null" json:"name"`
	Identifier  string       `gorm:"not null;unique" json:"identifier"`
	Address     string       `gorm:"not null" json:"address"`
	Phone       string       `gorm:"not null" json:"phone"`
	Email       string       `gorm:"Index;not null" json:"email"`
	CuitPdv     string       `gorm:"size:50;uniqueIndex;not null" json:"cuit_pdv"`
	IsActive    bool         `gorm:"not null;default:true" json:"is_active"`
	Connection  string       `gorm:"not null" json:"connection"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	UserTenants []UserTenant `gorm:"foreignKey:TenantID" json:"user_tenants"`
}

type TenantCreate struct {
	Name       string `json:"name" validate:"required"`
	Identifier string `json:"identifier" validate:"required"`
	Address    string `json:"address" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	CuitPdv    string `json:"cuit_pdv" validate:"required"`
}

func (t *TenantCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(t)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type TenantUpdate struct {
	ID      string `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

func (t *TenantUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(t)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type TenantResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	IsActive     bool      `json:"is_active"`
	UserIsActive bool      `json:"user_is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TenantUserCreate struct {
	TenantCreate TenantCreate `json:"tenant_create" validate:"required"`
	UserCreate   UserCreate   `json:"user_create" validate:"required"`
}

func (t *TenantUserCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(t)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}
