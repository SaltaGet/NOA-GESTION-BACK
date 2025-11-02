package schemas

import "time"

type UserTenant struct {
	UserID   string `gorm:"primaryKey;size:36" json:"user_id"`
	TenantID string `gorm:"primaryKey;size:36" json:"tenant_id"`
	IsActive   bool   `gorm:"not null;default:true" json:"is_active"`
	IsAdmin     bool   `gorm:"not null;default:false" json:"is_admin"`
	IsDeleted bool   `gorm:"not null;default:false" json:"is_deleted"`
	DeletedAt time.Time `gorm:"default:null" json:"deleted_at"`
	 
	User   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Tenant Tenant `gorm:"foreignKey:TenantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tenant"`
}
