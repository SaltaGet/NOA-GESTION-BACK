package repository


import (
	"time"

)

type CashRegisterOpen struct {
	OpenAmount float64 `json:"open_amount" example:"100.00"`
}

type CashRegisterClose struct {
	CloseAmount float64 `json:"close_amount" example:"100.00"`
}

type CashRegisterInformResponse struct {
	ID          int64            `json:"id"`
	MemberOpen  MemberSimpleDTO  `json:"member_open"`
	OpenAmount  float64          `json:"open_amount"`
	HourOpen    time.Time        `json:"hour_open"`
	MemberClose *MemberSimpleDTO `json:"member_close"`

	CloseAmount *float64   `json:"close_amount"`
	HourClose   *time.Time `json:"hour_close"`

	TotalIncomeCash    *float64 `json:"total_income_cash"`
	TotalIncomeOthers  *float64 `json:"total_income_others"`
	TotalExpenseCash   *float64 `json:"total_expense_cash"`
	TotalExpenseOthers *float64 `json:"total_expense_others"`

	IsClose   bool      `json:"is_close"`
	CreatedAt time.Time `json:"created_at"`
}

type CashRegisterFullResponse struct {
	ID          int64            `json:"id"`
	MemberOpen  MemberSimpleDTO  `json:"member_open"`
	OpenAmount  float64          `json:"open_amount"`
	HourOpen    time.Time        `json:"hour_open"`
	MemberClose *MemberSimpleDTO `json:"member_close"`
	CloseAmount *float64         `json:"close_amount"`
	HourClose   *time.Time       `json:"hour_close"`

	IsClose   bool      `json:"is_close"`
	CreatedAt time.Time `json:"created_at"`

	IncomeSale []*IncomeSaleSimpleResponse `json:"sale_income"`
	IncomeOther []*IncomeOtherResponse `json:"income_other"`
	// ExpenseBuy    *[]ExpenseBuyResponseSimple `json:"expense_buy"`
	ExpenseOther    []*ExpenseOtherResponse `json:"expense_others"`
}
