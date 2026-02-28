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
  
  // Cambios para compatibilidad con PostgreSQL
  if form == "month" {
    // Agrupamos por el mismo string formateado
    dateFormat = "TO_CHAR(mov.fecha, 'YYYY-MM') as fecha"
    modo = "ps.id, ps.name, TO_CHAR(mov.fecha, 'YYYY-MM')"
  } else {
    dateFormat = "TO_CHAR(mov.fecha, 'YYYY-MM-DD') as fecha"
    modo = "ps.id, ps.name, TO_CHAR(mov.fecha, 'YYYY-MM-DD')"
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
  GROUP BY %s
  ORDER BY ps.id
  `, dateFormat, modo)

  err := r.DB.Raw(query,
    start, end,
    start, end,
    start, end,
    start, end,
  ).Scan(&resultados).Error
  
  if err != nil {
    return nil, schemas.HandlerErrorGorm(err, "Reporte", schemas.Read)
  }

  grouped := make(map[string][]map[string]any)
  for _, row := range resultados {
    // IMPORTANTE: TO_CHAR siempre devuelve string.
    // Eliminamos la aserción a time.Time que causaba error.
    fecha, ok := row["fecha"].(string)
    if !ok {
        fecha = "Sin Fecha"
    }
    grouped[fecha] = append(grouped[fecha], row)
  }

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
		modo = "TO_CHAR(mov.fecha, 'YYYY-MM')"
		dateFormat = "TO_CHAR(mov.fecha, 'YYYY-MM') as fecha"
	} else {
		modo = "TO_CHAR(mov.fecha, 'YYYY-MM-DD')"
		dateFormat = "TO_CHAR(mov.fecha, 'YYYY-MM-DD') as fecha"
	}

	query := fmt.Sprintf(`
		SELECT
			%s,
			COALESCE(SUM(CASE WHEN tipo = 'ingreso_ventas'  THEN total ELSE 0 END), 0) AS total_ingreso_ventas,
			COALESCE(SUM(CASE WHEN tipo = 'egreso_compras'  THEN total ELSE 0 END), 0) AS total_compra_egresos,
			COALESCE(SUM(CASE WHEN tipo = 'ingreso_otros'   THEN total ELSE 0 END), 0) AS total_otros_ingresos,
			COALESCE(SUM(CASE WHEN tipo = 'egreso_otros'    THEN total ELSE 0 END), 0) AS total_otros_egresos
		FROM (
			SELECT created_at AS fecha, total, 'ingreso_ventas' AS tipo
			FROM income_sales
			WHERE created_at BETWEEN ? AND ?

			UNION ALL
			
			SELECT created_at AS fecha, total, 'egreso_compras' AS tipo
			FROM expense_buys
			WHERE created_at BETWEEN ? AND ? 

			UNION ALL

			SELECT created_at AS fecha, total, 'ingreso_otros' AS tipo
			FROM income_others
			WHERE created_at BETWEEN ? AND ?

			UNION ALL

			SELECT created_at AS fecha, total, 'egreso_otros' AS tipo
			FROM expense_others
			WHERE created_at BETWEEN ? AND ?
		) AS mov
		WHERE mov.fecha BETWEEN ? AND ?
		GROUP BY %s
		ORDER BY %s
	`, dateFormat, modo, modo)

	err := r.DB.Raw(query,
		start, end,
		start, end,
		start, end,
		start, end,
		start, end,
	).Scan(&resultados).Error
	if err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Reporte", schemas.Read)
	}

	grouped := make(map[string][]map[string]any)
	for _, row := range resultados {
		// Con TO_CHAR PostgreSQL siempre devuelve string, tanto para month como day
		fecha, ok := row["fecha"].(string)
		if !ok {
			return nil, fmt.Errorf("tipo inesperado en campo fecha: %T", row["fecha"])
		}
		grouped[fecha] = append(grouped[fecha], row)
	}

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
		return nil, schemas.HandlerErrorGorm(err, "Reporte", schemas.Read)
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
		return nil, schemas.HandlerErrorGorm(err, "Reporte", schemas.Read)
	}

	return products, nil
}

