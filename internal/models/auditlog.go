package models

type AuditLog struct {
	ID        string `gorm:"primaryKey" json:"id"`
	UserID 	string `gorm:"not null;size:36" json:"user_id"`
	Method    string `gorm:"not null" json:"method"`
	Path      string `gorm:"not null" json:"path"`
	CreatedAt string `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt string `gorm:"autoUpdateTime" json:"updated_at"`
}