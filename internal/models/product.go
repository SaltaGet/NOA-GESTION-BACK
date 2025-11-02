package models

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID         string    `gorm:"primaryKey;size:36" json:"id"`
	Identifier string    `gorm:"not null;unique" json:"identifier"`
	Name       string    `gorm:"not null" json:"name"`
	Stock      float32   `gorm:"not null;min:0;default:0" json:"stock"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type ProductCreate struct {
	Identifier string `json:"identifier" validate:"required"`
	Name       string `json:"name" validate:"required"`
}

func (p *ProductCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type ProductUpdate struct {
	ID         string `json:"id" validate:"required"`
	Identifier string `json:"identifier"`
	Name       string `json:"name" validate:"required"`
}

func (p *ProductUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type StockUpdate struct {
	ID     string  `json:"id" validate:"required"`
	Stock  float32 `json:"stock" validate:"required"`
	Method string  `json:"method" validate:"required,oneof=add subtract update"`
}

func (p *StockUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type ProductDTO struct {
	ID         string    `json:"id"`
	Identifier string    `json:"identifier"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
}