package models

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type Income struct {
	ID             string       `gorm:"primaryKey;size:36" json:"id"`
	Ticket         string       `json:"ticket"`
	Details        string       `json:"details"`
	ClientID       string       `gorm:"not null;size:36" json:"client_id"`
	VehicleID      string       `gorm:"not null;size:36" json:"vehicle_id"`
	EmployeeID     string       `gorm:"null;size:36" json:"employee_id"`
	Amount         float32      `gorm:"not null" json:"amount"`
	MovementTypeID string       `gorm:"not null;size:36" json:"movement_type_id"`
	CreatedAt      time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	Client         Client       `gorm:"foreignKey:ClientID" json:"client"`
	Vehicle        Vehicle      `gorm:"foreignKey:VehicleID" json:"vehicle"`
	Employee       Employee     `gorm:"foreignKey:EmployeeID" json:"employee"`
	MovementType   MovementType `gorm:"foreignKey:MovementTypeID;references:ID" json:"movement_type"`
	Services       []Service    `gorm:"many2many:income_services;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"services"`
}

type IncomeCreate struct {
	Ticket         string   `json:"ticket"`
	ServicesID     []string `json:"services_id" validate:"required,gt=0"`
	Details        string   `json:"details"`
	ClientID       string   `json:"client_id" validate:"required"`
	VehicleID      string   `json:"vehicle_id" validate:"required"`
	EmployeeID     string   `json:"employee_id"`
	MovementTypeID string   `json:"movement_type_id" validate:"required"`
	Amount         float32  `json:"amount" validate:"required"`
}

func (i *IncomeCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(i)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type IncomeUpdate struct {
	ID             string   `json:"id"`
	Ticket         string   `json:"ticket" validate:"required"`
	ServicesID     []string `json:"services_id" validate:"required,gt=0"`
	Details        string   `json:"details"`
	ClientID       string   `json:"client_id" validate:"required"`
	VehicleID      string   `json:"vehicle_id" validate:"required"`
	EmployeeID     string   `json:"employee_id"`
	MovementTypeID string   `json:"movement_type_id" validate:"required"`
	Amount         float32  `json:"amount" validate:"required"`
}

func (i *IncomeUpdate) Validate() error {
	validate := validator.New()
	err := validate.Struct(i)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type IncomeDTO struct {
	ID             string       `json:"id"`
	Ticket         string       `json:"ticket"`
	Amount         float32      `json:"amount"`
	CreatedAt      time.Time    `json:"created_at"`
	Client         ClientDTO       `json:"client"`
	Vehicle        VehicleDTO      `json:"vehicle"`
	Employee       EmployeeDTO     `json:"employee"`
	MovementType   MovementTypeDTO `json:"movement_type"`
	Services       []ServiceDTO    `json:"services"`
}

type IncomeResponse struct {
	ID             string       `json:"id"`
	Ticket         string       `json:"ticket"`
	Details        string       `json:"details"`
	Amount         float32      `json:"amount"`
	CreatedAt      time.Time    `json:"created_at"`
	Client         ClientResponse       `json:"client"`
	Vehicle        VehicleResponse      `json:"vehicle"`
	Employee       EmployeeResponse     `json:"employee"`
	MovementType   MovementTypeDTO `json:"movement_type"`
	Services       []ServiceDTO    `json:"services"`
}