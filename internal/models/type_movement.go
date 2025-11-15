package models

type TypeIncome struct {
	ID   int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"not null;unique" json:"name"`
}

type TypeExpense struct {
	ID   int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"not null;unique" json:"name"`
}
