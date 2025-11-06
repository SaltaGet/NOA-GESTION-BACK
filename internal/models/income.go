package models

import "time"

type SaleIncome struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Ticket    string    `json:"ticket"`
	Details   string    `json:"details"`
	ClientID  int64     `gorm:"not null;size:36" json:"client_id"`
	Client    Client    `gorm:"foreignKey:ClientID" json:"client"`
	SubTotal  float32   `gorm:"not null" json:"subtotal"`
	Total     float32   `gorm:"not null" json:"total"`
	Discount  float32   `gorm:"not null;default:0" json:"discount"`
	Type      string    `gorm:"not null;default:percent" json:"type_discount" validate:"oneof=amount percent"`
	IsBudget  bool      `gorm:"not null;default:false" json:"is_budget"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type ItemSaleIncome struct {
	ID           int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	SaleIncomeID int64   `gorm:"index" json:"sale_id"`
	ProductID    int64   `gorm:"index" json:"product_id"`
	Product      Product `gorm:"foreignKey:ProductID" json:"product"`
	Amount       float32 `gorm:"not null" json:"amount"`
	Price        float32 `gorm:"not null" json:"price"`
	Discount     float32 `gorm:"not null;default:0" json:"discount"`
	TypeDiscount string  `gorm:"not null;default:percent" json:"type_discount" validate:"oneof=amount percent"`
	Total        float32 `gorm:"not null" json:"total"`
}

type TypeIncome struct {
	ID   int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"not null;unique" json:"name"`
}

type OtherIncome struct { // continuar actualizando
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Total        float32    `gorm:"not null" json:"total"`
	TypeIncomeID int64      `gorm:"not null" json:"type_income_id"`
	TypeIncome   TypeIncome `gorm:"foreignKey:TypeIncomeID" json:"type_income"`
	Details      *string    `json:"details"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
