package models

import "time"

type ExpenseBuy struct {
	ID             int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	MemberID       int64            `gorm:"not null" json:"member_id"`
	Member         Member           `gorm:"foreignKey:MemberID;references:ID" json:"member"`
	SupplierID     int64            `gorm:"not null" json:"supplier_id"`
	Supplier       Supplier         `gorm:"foreignKey:SupplierID;references:ID" json:"supplier"`
	Details        *string          `gorm:"size:255" json:"details"`
	ExpenseItemBuy []ExpenseBuyItem `gorm:"foreignKey:ExpenseBuyID" json:"expense_item_buys"`
	PayExpenseBuy  []PayExpenseBuy  `gorm:"foreignKey:ExpenseBuyID" json:"pay_expense"`
	Subtotal       float64          `gorm:"subtotal" json:"subtotal"`
	Discount       float64          `gorm:"not null;default:0" json:"discount"`
	TypeDiscount   string           `gorm:"not null;default:percent" json:"type_discount" validate:"oneof=amount percent"`
	Total          float64          `gorm:"total" json:"total"`
	CreatedAt      time.Time        `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt      time.Time        `gorm:"autoUpdateTime:milli" json:"updated_at"`
}

type ExpenseBuyItem struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ExpenseBuyID int64     `gorm:"not null" json:"expense_buy_id"`
	ExpenseBuy   ExpenseBuy `gorm:"foreignKey:ExpenseBuyID" json:"expense_buy"`
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
	ID             int64         `gorm:"primaryKey;autoIncrement" json:"id"`
	PointSaleID    *int64        `gorm:"" json:"point_sale_id"`
	PointSale      *PointSale    `gorm:"foreignKey:PointSaleID;references:ID" json:"point_sale"`
	MemberID       int64         `gorm:"not null" json:"member_id"`
	Member         Member        `gorm:"foreignKey:MemberID;references:ID" json:"member"`
	CashRegisterID *int64        `gorm:"" json:"register_id"`
	CashRegister   *CashRegister `gorm:"foreignKey:CashRegisterID;references:ID" json:"cash_register"`
	Details    *string       `gorm:"size:255" json:"details"`
	TypeExpenseID  int64         `gorm:"not null" json:"type_expense_id"`
	TypeExpense    TypeExpense   `gorm:"foreignKey:TypeExpenseID" json:"type_expense"`
	Total          float64       `gorm:"not null" json:"total"`
	PayMethod      string        `gorm:"size:30;default:'efectivo'" json:"pay_method" validate:"oneof=cash credit card transfer"`
	CreatedAt      time.Time     `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime:milli" json:"updated_at"`
}

type TypeExpense struct {
	ID   int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"not null;unique" json:"name"`
}
