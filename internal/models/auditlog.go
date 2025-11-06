package models

type AuditLogAdmin struct {
	ID        int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	AdminID   int64   `gorm:"not null;size:36" json:"user_id"`
	Method    string  `gorm:"not null" json:"method"`
	Path      string  `gorm:"not null" json:"path"`
	OldValue  *string `json:"old_value"`
	NewValue  *string `json:"new_value"`
	CreatedAt string  `gorm:"autoCreateTime" json:"created_at"`
}

type AuditLog struct {
	ID        int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	MemberID  int64   `gorm:"not null;size:36" json:"user_id"`
	TenantID  int64   `gorm:"not null;size:36" json:"tenant_id"`
	Method    string  `gorm:"not null" json:"method"`
	Path      string  `gorm:"not null" json:"path"`
	OldValue  *string `json:"old_value"`
	NewValue  *string `json:"new_value"`
	CreatedAt string  `gorm:"autoCreateTime" json:"created_at"`
}
