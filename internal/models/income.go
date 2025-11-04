package models

import (
	"time"
)

type Income struct {
	ID             int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	Ticket         string       `json:"ticket"`
	Details        string       `json:"details"`
	ClientID       int64        `gorm:"not null;size:36" json:"client_id"`
	VehicleID      int64        `gorm:"not null;size:36" json:"vehicle_id"`
	EmployeeID     int64        `gorm:"null;size:36" json:"employee_id"`
	Amount         float32      `gorm:"not null" json:"amount"`
	MovementTypeID int64        `gorm:"not null;size:36" json:"movement_type_id"`
	CreatedAt      time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	Client         Client       `gorm:"foreignKey:ClientID" json:"client"`
	MovementType   MovementType `gorm:"foreignKey:MovementTypeID;references:ID" json:"movement_type"`
}
