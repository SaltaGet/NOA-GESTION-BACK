package schemas

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type Client struct {
	ID        string    `gorm:"primaryKey;size:36" json:"id"`
	FirstName string    `gorm:"not null;size:30" json:"first_name"`
	LastName  string    `gorm:"not null;size:30" json:"last_name"`
	Cuil      string    `gorm:"not null;unique;size:30" json:"cuil"`
	Dni       string    `gorm:"not null;unique;size:30" json:"dni"`
	Email     string    `gorm:"not null;unique" json:"email" validate:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Vehicles  []Vehicle `gorm:"foreignKey:ClientID" json:"vehicles"`
}

type ClientCreate struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Cuil      string `json:"cuil" validate:"required"`
	Dni       string `json:"dni" validate:""`
	Email     string `json:"email" validate:"required,email"`
}

func (c *ClientCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(c)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type ClientUpdate struct {
	ID        string `json:"id" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Cuil      string `json:"cuil" validate:"required"`
	Dni       string `json:"dni" validate:""`
	Email     string `json:"email" validate:"required,email"`
}

func (c *ClientUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(c)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type ClientDTO struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ClientResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Cuil      string `json:"cuil"`
	Dni       string `json:"dni"`
	Email     string `json:"email"`
}