package models

import "time"

type ExpenseBuy struct {
	ID             int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	PointSaleID    *int64            `gorm:"" json:"point_sale_id"`
	PointSale      *PointSale        `gorm:"foreignKey:PointSaleID;references:ID" json:"point_sale"`
	MemberID       int64            `gorm:"not null" json:"member_id"`
	Member         Member           `gorm:"foreignKey:MemberID;references:ID" json:"member"`
	SupplierID     int64            `gorm:"not null" json:"supplier_id"`
	Supplier       Supplier         `gorm:"foreignKey:SupplierID;references:ID" json:"supplier"`
	RegisterID     *int64           `gorm:"" json:"register_id"`
	Register       *CashRegister    `gorm:"foreignKey:RegisterID;references:ID" json:"register"`
	Description    *string          `gorm:"size:255" json:"description"`
	ExpenseItemBuy []ExpenseBuyItem `gorm:"foreignKey:ExpenseBuyID" json:"expense_item_buys"`
	PayExpenseBuy  []PayExpenseBuy  `gorm:"foreignKey:ExpenseBuyID" json:"pay_expense"`
	Total          float64          `gorm:"total" json:"total"`
	CreatedAt      time.Time        `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt      time.Time        `gorm:"autoUpdateTime:milli" json:"updated_at"`
}

type ExpenseBuyItem struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID    int64     `gorm:"not null" json:"product_id"`
	Product      Product   `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	Amount       float64   `gorm:"not null" json:"amount"`
	Price        float64   `gorm:"not null" json:"price"`
	Discount     float64   `gorm:"not null;default:0" json:"discount"`
	TypeDiscount string    `gorm:"not null;default:percent" json:"type_discount" validate:"oneof=amount percent"`
	SubTotal     float64   `gorm:"not null" json:"subtotal"`
	Total        float64   `gorm:"not null" json:"total"`
	CreatedAt    time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`
}

type ExpenseOther struct {
	ID            int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	PointSaleID   *int64        `gorm:"" json:"point_sale_id"`
	PointSale     *PointSale    `gorm:"foreignKey:PointSaleID;references:ID" json:"point_sale"`
	MemberID      int64        `gorm:"not null" json:"member_id"`
	Member        Member       `gorm:"foreignKey:MemberID;references:ID" json:"member"`
	RegisterID    *int64        `gorm:"" json:"register_id"`
	Register      *CashRegister `gorm:"foreignKey:RegisterID;references:ID" json:"register"`
	Description   *string      `gorm:"size:255" json:"description"`
	Total         float64      `gorm:"not null" json:"total"`
	PaymentMethod string       `gorm:"size:30;default:'efectivo'" json:"payment_method" validate:"oneof=cash credit card transfer"`
	CreatedAt     time.Time    `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt     time.Time    `gorm:"autoUpdateTime:milli" json:"updated_at"`
}
