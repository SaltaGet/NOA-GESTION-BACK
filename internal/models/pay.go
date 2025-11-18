package models

import "time"

type PayIncome struct {
	ID             int64         `gorm:"primaryKey;autoIncrement" json:"id"`
	IncomeSaleID   int64         `gorm:"index" json:"income_sale_id"`
	IncomeSale     *IncomeSale    `gorm:"foreignKey:IncomeSaleID" json:"income_sale"`
	CashRegisterID *int64        `gorm:"index" json:"cash_register_id"`
	CashRegister   *CashRegister `gorm:"foreignKey:CashRegisterID" json:"cash_register"`
	ClientID       *int64        `gorm:"index" json:"client_id"`
	Client         *Client       `gorm:"foreignKey:ClientID" json:"client"`
	Total         float64       `gorm:"not null" json:"total"`
	MethodPay      string        `gorm:"not null;default:cash" json:"method_pay" validate:"oneof=cash credit card transfer"`
	CreatedAt      time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

type PayExpenseBuy struct {
	ID             int64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ExpenseBuyID   int64         `gorm:"index" json:"expense_buy_id"`
	ExpenseBuy     ExpenseBuy    `gorm:"foreignKey:ExpenseBuyID" json:"expense_buy"`
	CashRegisterID *int64        `gorm:"index" json:"cash_register_id"`
	CashRegister   *CashRegister `gorm:"foreignKey:CashRegisterID" json:"cash_register"`
	Total         float64       `gorm:"not null" json:"total"`
	MethodPay      string        `gorm:"not null;default:cash" json:"method_pay" validate:"oneof=cash credit card transfer"`
	CreatedAt      time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

type PayExpenseOther struct {
	ID             int64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ExpenseOtherID int64         `gorm:"index" json:"expense_other_id"`
	ExpenseOther   ExpenseOther  `gorm:"foreignKey:ExpenseOtherID" json:"expense_other"`
	CashRegisterID *int64        `gorm:"index" json:"cash_register_id"`
	CashRegister   *CashRegister `gorm:"foreignKey:CashRegisterID" json:"cash_register"`
	Total         float64       `gorm:"not null" json:"total"`
	MethodPay      string        `gorm:"not null;default:cash" json:"method_pay" validate:"oneof=cash credit card transfer"`
	CreatedAt      time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}
