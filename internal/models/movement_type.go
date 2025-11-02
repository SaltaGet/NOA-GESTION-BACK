package models

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type MovementType struct {
	ID        string    `gorm:"primaryKey;size:36" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	IsIncome  bool      `gorm:"not null" json:"is_income"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type MovementTypeCreate struct {
	Name     string `json:"name"`
	IsIncome bool   `json:"is_income"`
}

func (m *MovementTypeCreate) Validate() error {
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

type MovementTypeUpdate struct {
	ID       string `gojson:"id"`
	Name     string `json:"name"`
	IsIncome bool   `json:"is_income"`
}

func (m *MovementTypeUpdate) Validate() error {
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

type MovementTypeDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	IsIncome  bool      `json:"is_income"`
}
