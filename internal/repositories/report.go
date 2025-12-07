package repositories

import (
	"fmt"
	"sort"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (r *ReportRepository) ReportMovementByDatePointSale(start, end time.Time, form string) (any, error) {
	var resultados []map[string]any

	var modo string
	var dateFormat string
	if form == "month" {
		modo = "YEAR(mov.fecha), MONTH(mov.fecha)"
		dateFormat = "DATE_FORMAT(mov.fecha, '%Y-%m') as fecha"
	} else {
		modo = "DATE(mov.fecha)"
		dateFormat = "DATE(mov.fecha) as fecha"
	}

	query := fmt.Sprintf(`
	SELECT 
    ps.id as point_sale_id,
    ps.name as point_sale_name,
    %s,
    COALESCE(SUM(CASE WHEN tipo = 'ingreso_ventas' THEN total ELSE 0 END),0) as total_ingresos_ventas,
    COALESCE(SUM(CASE WHEN tipo = 'ingreso_otros' THEN total ELSE 0 END),0) as total_otros_ingresos,
    COALESCE(SUM(CASE WHEN tipo = 'egreso_otros' THEN total ELSE 0 END),0) as total_otros_egresos
	FROM (
    SELECT created_at as fecha, total, 'ingreso_ventas' as tipo, point_sale_id
			FROM income_sales
			WHERE created_at BETWEEN ? AND ?
			
			UNION ALL
			
			SELECT created_at as fecha, total, 'ingreso_otros' as tipo, point_sale_id
			FROM income_others
			WHERE created_at BETWEEN ? AND ?
			
			UNION ALL
			
			SELECT created_at as fecha, total, 'egreso_otros' as tipo, point_sale_id
			FROM expense_others
			WHERE created_at BETWEEN ? AND ?
	) AS mov
	JOIN point_sales ps ON ps.id = mov.point_sale_id
	WHERE mov.fecha BETWEEN ? AND ?
	GROUP BY ps.id, ps.name, %s
	ORDER BY ps.id
	`, dateFormat, modo)

	err := r.DB.Raw(query,
		start, end,
		start, end,
		start, end,
		start, end,
	).Scan(&resultados).Error
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al obtener los datos", err)
	}

	grouped := make(map[string][]map[string]any)
	for _, row := range resultados {
		var fecha string
		if form == "month" {
			fecha = row["fecha"].(string)
		} else {
			fecha = row["fecha"].(time.Time).Format("2006-01-02") // key simple

		}
		grouped[fecha] = append(grouped[fecha], row)
	}

	// Transformar al formato esperado
	var result []map[string]any
	for fecha, movimientos := range grouped {
		result = append(result, map[string]any{
			"fecha":      fecha,
			"movimiento": movimientos,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i]["fecha"].(string) < result[j]["fecha"].(string)
	})

	return result, nil
}

func (r *ReportRepository) ReportMovementByDate(start, end time.Time, form string) (any, error) {
	var resultados []map[string]any

	var modo string
	var dateFormat string
	if form == "month" {
		modo = "YEAR(mov.fecha), MONTH(mov.fecha)"
		dateFormat = "DATE_FORMAT(mov.fecha, '%Y-%m') as fecha"
	} else {
		modo = "DATE(mov.fecha)"
		dateFormat = "DATE(mov.fecha) as fecha"
	}

	query := fmt.Sprintf(`
	SELECT 
    %s,
    COALESCE(SUM(CASE WHEN tipo = 'egreso_compras' THEN total ELSE 0 END),0) as total_compra_egresos,
    COALESCE(SUM(CASE WHEN tipo = 'ingreso_otros' THEN total ELSE 0 END),0) as total_otros_ingresos,
    COALESCE(SUM(CASE WHEN tipo = 'egreso_otros' THEN total ELSE 0 END),0) as total_otros_egresos
	FROM (
    SELECT created_at as fecha, total, 'egreso_compras' as tipo
			FROM income_sales
			WHERE created_at BETWEEN ? AND ? AND point_sale_id IS NULL
			
			UNION ALL
			
			SELECT created_at as fecha, total, 'ingreso_otros' as tipo
			FROM income_others
			WHERE created_at BETWEEN ? AND ? AND point_sale_id IS NULL
			
			UNION ALL
			
			SELECT created_at as fecha, total, 'egreso_otros' as tipo
			FROM expense_others
			WHERE created_at BETWEEN ? AND ? AND point_sale_id IS NULL
	) AS mov
	WHERE mov.fecha BETWEEN ? AND ?
	GROUP BY %s
	`, dateFormat, modo)

	err := r.DB.Raw(query,
		start, end,
		start, end,
		start, end,
		start, end,
	).Scan(&resultados).Error
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al obtener los datos", err)
	}

	grouped := make(map[string][]map[string]any)
	for _, row := range resultados {
		var fecha string
		if form == "month" {
			fecha = row["fecha"].(string)
		} else {
			fecha = row["fecha"].(time.Time).Format("2006-01-02") // key simple

		}
		grouped[fecha] = append(grouped[fecha], row)
	}

	// Transformar al formato esperado
	var result []map[string]any
	for fecha, movimientos := range grouped {
		result = append(result, map[string]any{
			"fecha":      fecha,
			"movimiento": movimientos,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i]["fecha"].(string) < result[j]["fecha"].(string)
	})

	return result, nil
}

func (r *ReportRepository) ReportProfitableProducts(start, end time.Time) ([]schemas.ReportProfitableProducts, error) {
	var products []schemas.ReportProfitableProducts

	query := `
		SELECT 
			p.id,
			p.code,
			p.name,
			COALESCE(SUM(ii.amount), 0) AS total_quantity,
			COALESCE(SUM(ii.subtotal), 0) AS total_sales,
			COALESCE(SUM(ii.price_cost * ii.amount), 0) AS total_cost,
			COALESCE(SUM((ii.price - ii.price_cost) * ii.amount), 0) AS total_profit
		FROM products p
		LEFT JOIN income_sale_items ii 
			ON p.id = ii.product_id
			AND ii.created_at BETWEEN ? AND ?
		GROUP BY p.id, p.code, p.name
		ORDER BY total_profit DESC, total_quantity DESC
	`

	if err := r.DB.Raw(query, start, end).Scan(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ReportRepository) ReportStockProducts() ([]*models.Product, error) {
	var products []*models.Product

	if err := r.DB.
		Preload("Category").
		Preload("StockPointSales").
		Preload("StockPointSales.PointSale").
		Preload("StockDeposit").
		Find(&products).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "error al obtener productos", err)
	}

	return products, nil
}

// func (r *MainRepository) ObtenerResumenPorDiaYPuntoVenta(fechaInicio, fechaFin time.Time) ([]schemas.ResultadoPorDiaYPuntoVenta, error) {
// 	var resultados []schemas.ResultadoPorDiaYPuntoVenta

// 	query := `
// 		SELECT
// 			DATE(fecha) as fecha,
// 			point_sale_id,
// 			COALESCE(SUM(CASE WHEN tipo = 'ingreso' THEN total ELSE 0 END), 0) as total_ingresos,
// 			COALESCE(SUM(CASE WHEN tipo = 'egreso' THEN total ELSE 0 END), 0) as total_egresos,
// 			COALESCE(SUM(CASE WHEN tipo = 'cancha' THEN total ELSE 0 END), 0) as total_canchas,
// 			COALESCE(SUM(CASE WHEN tipo IN ('ingreso', 'cancha') THEN total ELSE -total END), 0) as balance
// 		FROM (
// 			SELECT created_at as fecha, total, point_sale_id, 'ingreso' as tipo
// 			FROM incomes
// 			WHERE created_at BETWEEN ? AND ?

// 			UNION ALL

// 			SELECT created_at as fecha, total, point_sale_id, 'egreso' as tipo
// 			FROM expenses
// 			WHERE created_at BETWEEN ? AND ?

// 			UNION ALL

// 			SELECT created_at as fecha, total, point_sale_id, 'cancha' as tipo
// 			FROM income_sports_courts
// 			WHERE created_at BETWEEN ? AND ?

// 			UNION ALL
// 		) as movimientos
// 		GROUP BY DATE(fecha), point_sale_id
// 		ORDER BY fecha DESC, point_sale_id
// 	`

// 	err := r.DB.Raw(query,
// 		fechaInicio, fechaFin,
// 		fechaInicio, fechaFin,
// 		fechaInicio, fechaFin,
// 	).Scan(&resultados).Error

// 	return resultados, err
// }

// ========================
// CONSULTAS POR MES
// ========================

// ObtenerResumenPorMes obtiene todos los ingresos y egresos agrupados por mes para un punto de venta
// func ObtenerResumenPorMes(pointSaleID uint, anio int) ([]schemas.ResultadoPorMes, error) {
// 	var resultados []schemas.ResultadoPorMes

// 	query := `
// 		SELECT
// 			DATE_FORMAT(fecha, '%Y-%m') as mes,
// 			YEAR(fecha) as anio,
// 			COALESCE(SUM(CASE WHEN tipo = 'ingreso' THEN total ELSE 0 END), 0) as total_ingresos,
// 			COALESCE(SUM(CASE WHEN tipo = 'egreso' THEN total ELSE 0 END), 0) as total_egresos,
// 			COALESCE(SUM(CASE WHEN tipo = 'cancha' THEN total ELSE 0 END), 0) as total_canchas,
// 			COALESCE(SUM(CASE WHEN tipo = 'compra' THEN total ELSE 0 END), 0) as total_compras,
// 			COALESCE(SUM(CASE WHEN tipo IN ('ingreso', 'cancha') THEN total ELSE -total END), 0) as balance
// 		FROM (
// 			SELECT created_at as fecha, total, 'ingreso' as tipo
// 			FROM incomes
// 			WHERE point_sale_id = ? AND YEAR(created_at) = ?

// 			UNION ALL

// 			SELECT created_at as fecha, total, 'egreso' as tipo
// 			FROM expenses
// 			WHERE point_sale_id = ? AND YEAR(created_at) = ?

// 			UNION ALL

// 			SELECT created_at as fecha, total, 'cancha' as tipo
// 			FROM income_sports_courts
// 			WHERE point_sale_id = ? AND YEAR(created_at) = ?

// 			UNION ALL

// 			SELECT created_at as fecha, total, 'compra' as tipo
// 			FROM expense_buys
// 			WHERE YEAR(created_at) = ?
// 		) as movimientos
// 		GROUP BY DATE_FORMAT(fecha, '%Y-%m'), YEAR(fecha)
// 		ORDER BY mes DESC
// 	`

// 	err := db.Raw(query,
// 		pointSaleID, anio,
// 		pointSaleID, anio,
// 		pointSaleID, anio,
// 		anio,
// 	).Scan(&resultados).Error

// 	return resultados, err
// }
