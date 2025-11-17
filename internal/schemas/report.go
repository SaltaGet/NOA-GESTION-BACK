package schemas

type ReportProfitableProducts struct {
	ID           int64   `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	TotalCost    float64 `json:"total_cost"`
	TotalProfit  float64 `json:"total_profit"`
	TotalQuantity float64 `json:"total_quantity"`
	TotalSales   float64 `json:"total_sales"`
}

// type ReportMovementResponse struct {
// 	Income []*ReportCount `json:"income"`
// 	Expense []*ReportCount `json:"expense"`
// 	IncomeSportsCourts []*ReportCount `json:"income_sports_courts"`
// 	ExpenseBuy []*ReportCount `json:"expense_buy"`
// }

// type ReportCount struct {
// 	Date string `json:"date"`
// 	Group string `json:"group"`
// 	Count float64 `json:"count"`
// }

// type ReportPointSaleFull struct {
//     PointSaleID   uint    `json:"point_sale_id"`
//     PointSale     string  `json:"point_sale"`
// 		Period       time.Time `json:"period"`
//     IncomeTotal   float64 `json:"income_total"`
//     ExpenseTotal  float64 `json:"expense_total"`
//     SportsIncome  float64 `json:"sports_income"`
//     ExpenseBuy    float64 `json:"expense_buy"`
// }





// type ResultadoPorDia struct {
// 	Fecha           string  `json:"fecha"`
// 	TotalIngresos   float64 `json:"total_ingresos"`
// 	TotalEgresos    float64 `json:"total_egresos"`
// 	TotalCanchas    float64 `json:"total_canchas"`
// 	TotalCompras    float64 `json:"total_compras"`
// 	Balance         float64 `json:"balance"`
// }

// // ResultadoPorDiaYPuntoVenta representa ingresos/egresos agrupados por día y punto de venta
// type ResultadoPorDiaYPuntoVenta struct {
// 	Fecha           string  `json:"fecha"`
// 	PointSaleID     uint    `json:"point_sale_id"`
// 	TotalIngresos   float64 `json:"total_ingresos"`
// 	TotalEgresos    float64 `json:"total_egresos"`
// 	TotalCanchas    float64 `json:"total_canchas"`
// 	TotalCompras    float64 `json:"total_compras"`
// 	Balance         float64 `json:"balance"`
// }

// // ResultadoPorMes representa ingresos/egresos agrupados por mes
// type ResultadoPorMes struct {
// 	Mes             string  `json:"mes"`
// 	Anio            int     `json:"anio"`
// 	TotalIngresos   float64 `json:"total_ingresos"`
// 	TotalEgresos    float64 `json:"total_egresos"`
// 	TotalCanchas    float64 `json:"total_canchas"`
// 	TotalCompras    float64 `json:"total_compras"`
// 	Balance         float64 `json:"balance"`
// }

// // ResumenFinanciero estructura temporal para la consulta
// type ResumenFinanciero struct {
// 	Fecha           time.Time
// 	TotalIngresos   float64
// 	TotalEgresos    float64
// 	TotalCanchas    float64
// 	TotalCompras    float64
// }