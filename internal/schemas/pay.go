package schemas

type PayResponse struct {
	ID        int64   `json:"id"`
	Total    float64 `json:"amount"`
	MethodPay string  `json:"method_pay"`
	CreatedAt    string  `json:"created_at"`
}

type PayDebtResponse struct {
	ID           int64   `json:"id"`
	IncomeSaleID int64   `json:"income_sale_id"`
	Total       float64 `json:"amount"`
	MethodPay    string  `json:"method_pay"`
	CreatedAt    string  `json:"created_at"`
}
