package models

type Permission struct {
	ID      int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Code    string `gorm:"not null,unique" json:"code"`
	Name    string `gorm:"not null" json:"name"`
	Details string `gorm:"not null" json:"details"`
	Group   string `gorm:"not null" json:"group"`
	Environment string `gorm:"not null" json:"environment"`
	
	Roles   []Role `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"roles"`
}


