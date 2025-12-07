package schemas

type PointSaleStockResponse struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Stock     float64 `json:"stock"`
	IsDeposit bool    `json:"is_deposit"`
}
