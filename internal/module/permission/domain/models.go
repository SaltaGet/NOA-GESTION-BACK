package domain


type Permission struct {
	ID      int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Code    string `gorm:"type:varchar(50);not null,unique" json:"code"`
	Name    string `gorm:"type:varchar(100);not null" json:"name"`
	Details string `gorm:"type:text;not null" json:"details"`
	Group   string `gorm:"type:varchar(50);not null" json:"group"`
	Environment string `gorm:"type:varchar(20);not null" json:"environment"`
	
	Roles   []Role `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"roles"`
}


