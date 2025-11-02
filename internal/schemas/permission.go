package schemas

type Permission struct {
	ID    string    `gorm:"primaryKey;size:36" json:"id"`
	Code  string `gorm:"not null,unique" json:"code"`
	Details string `gorm:"not null" json:"details"`
	Group string `gorm:"not null" json:"group"`
	Roles []Role `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"roles"`
}

type PermissionResponse struct {
	ID      string `json:"id"`
	Code    string `json:"code"`
	Details string `json:"details"`
	Group   string `json:"group"`
}

type PermissionDTO struct {
	ID      string `json:"id"`
	Code    string `json:"code"`
	Details string `json:"details"`
}