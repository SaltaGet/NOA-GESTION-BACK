package models

import "time"

type AuditLogAdmin struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	AdminID   int64     `gorm:"not null" json:"user_id"`
	Admin     Admin     `gorm:"foreignKey:AdminID;references:ID" json:"admin"`
	Method    string    `gorm:"not null;type:varchar(10)" json:"method"`
	Path      string    `gorm:"not null;type:varchar(255)" json:"path"`
	OldValue  *string   `gorm:"type:text" json:"old_value"`
	NewValue  *string   `gorm:"type:text" json:"new_value"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type AuditLog struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MemberID  int64     `gorm:"not null" json:"user_id"`
	Member    Member    `gorm:"foreignKey:MemberID;references:ID" json:"member"`
	Method    string    `gorm:"not null;type:varchar(10)" json:"method"`
	Path      string    `gorm:"not null;type:varchar(255)" json:"path"`
	OldValue  *string   `gorm:"type:text" json:"old_value"`
	NewValue  *string   `gorm:"type:text" json:"new_value"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
