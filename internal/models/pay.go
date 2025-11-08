package models

import "time"

type PayIncome struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	SaleIncomeID int64      `gorm:"index" json:"sale_id"`
	SaleIncome   SaleIncome `gorm:"foreignKey:SaleIncomeID" json:"sale_income"`
	Amount       float32    `gorm:"not null" json:"amount"`
	MethodPay    string     `gorm:"not null;default:cash" json:"method_pay" validate:"oneof=cash credit card transfer"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type PayExpenseBuy struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ExpenseBuyID int64      `gorm:"index" json:"sale_id"`
	ExpenseBuy   ExpenseBuy `gorm:"foreignKey:ExpenseBuyID" json:"sale_income"`
	Amount       float32    `gorm:"not null" json:"amount"`
	MethodPay    string     `gorm:"not null;default:cash" json:"method_pay" validate:"oneof=cash credit card transfer"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type PayExpenseOther struct {
	ID             int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	ExpenseOtherID int64        `gorm:"index" json:"sale_id"`
	ExpenseOther   ExpenseOther `gorm:"foreignKey:ExpenseOtherID" json:"expense_other"`
	Amount         float32      `gorm:"not null" json:"amount"`
	MethodPay      string       `gorm:"not null;default:cash" json:"method_pay" validate:"oneof=cash credit card transfer"`
	CreatedAt      time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}
