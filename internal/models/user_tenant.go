package models

type UserTenant struct {
	UserID    int64     `gorm:"not null" json:"user_id"`
	TenantID  int64     `gorm:"not null" json:"tenant_id"`

	User   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Tenant Tenant `gorm:"foreignKey:TenantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tenant"`
}
