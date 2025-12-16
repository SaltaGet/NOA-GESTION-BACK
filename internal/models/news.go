package models

type News struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string `gorm:"type:varchar(255)" json:"title"`
	Content   string `gorm:"type:text" json:"content"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at"`
}