package models

import "time"

type IncomeSale struct {
	ID             int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	PointSaleID    int64            `gorm:"not null" json:"point_sale_id"`
	PointSale      PointSale        `gorm:"foreignKey:PointSaleID;references:ID" json:"point_sale"`
	MemberID       int64            `gorm:"not null" json:"member_id"`
	Member         Member           `gorm:"foreignKey:MemberID;references:ID" json:"member"`
	ClientID       int64            `gorm:"not null" json:"client_id"`
	Client         Client           `gorm:"foreignKey:ClientID;references:ID" json:"client"`
	CashRegisterID int64            `gorm:"not null" json:"register_id"`
	CashRegister   CashRegister     `gorm:"foreignKey:CashRegisterID;references:ID" json:"cash_register"`
	Items          []IncomeSaleItem `gorm:"foreignKey:IncomeSaleID" json:"items"`
	Subtotal       float64          `gorm:"not null" json:"subtotal"`
	Discount       float64          `gorm:"not null;default:0" json:"discount"`
	Type           string           `gorm:"not null;default:percent" json:"type_discount" validate:"oneof=amount percent"`
	Total          float64          `gorm:"not null" json:"total"`
	IsBudget       bool             `gorm:"not null;default:false" json:"is_budget"`
	Pay            []PayIncome      `gorm:"foreignKey:IncomeSaleID" json:"pay"`
	CreatedAt      time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
}

type IncomeSaleItem struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	IncomeSaleID int64      `gorm:"index" json:"sale_id"`
	IncomeSale   IncomeSale `gorm:"foreignKey:IncomeSaleID" json:"sale_income"`
	ProductID    int64      `gorm:"index" json:"product_id"`
	Product      Product    `gorm:"foreignKey:ProductID" json:"product"`
	Amount       float64    `gorm:"not null" json:"amount"`
	Price_Cost   float64    `gorm:"not null" json:"price_cost"`
	Price        float64    `gorm:"not null" json:"price"`
	Discount     float64    `gorm:"not null;default:0" json:"discount"`
	TypeDiscount string     `gorm:"not null;default:percent" json:"type_discount" validate:"oneof=amount percent"`
	Subtotal     float64    `gorm:"not null" json:"subtotal"`
	Total        float64    `gorm:"not null" json:"total"`
}

type IncomeOther struct {
	ID             int64         `gorm:"primaryKey;autoIncrement" json:"id"`
	PointSaleID    *int64        `gorm:"" json:"point_sale_id"`
	PointSale      *PointSale    `gorm:"foreignKey:PointSaleID;references:ID" json:"point_sale"`
	MemberID       *int64        `gorm:"" json:"member_id"`
	Member         *Member       `gorm:"foreignKey:MemberID;references:ID" json:"member"`
	CashRegisterID *int64        `gorm:"" json:"cash_register_id"`
	CashRegister   *CashRegister `gorm:"foreignKey:CashRegisterID;references:ID" json:"cash_register"`
	Total          float64       `gorm:"not null" json:"total"`
	TypeIncomeID   int64         `gorm:"not null" json:"type_income_id"`
	TypeIncome     TypeIncome    `gorm:"foreignKey:TypeIncomeID" json:"type_income"`
	Details        *string       `json:"details"`
	MethodIncome   string        `gorm:"not null;default:cash" json:"method_income" validate:"oneof=cash credit card transfer"`
	CreatedAt      time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

type TypeIncome struct {
	ID   int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"not null;unique" json:"name"`
}
