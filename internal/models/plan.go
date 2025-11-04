package models

type Plan struct {
	ID      int64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name    string   `gorm:"type:varchar(255)" json:"name"`
	Tenants []Tenant `gorm:"foreignKey:PlanID" json:"tenants,omitempty"`
}
