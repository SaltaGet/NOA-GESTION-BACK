package repositories

import (
	"context"
	"fmt"
	"sort"
	"time"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/aarondl/sqlboiler/v4/queries"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
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
      WHERE created_at BETWEEN $1 AND $2
      
      UNION ALL
      
      SELECT created_at as fecha, total, 'ingreso_otros' as tipo, point_sale_id
      FROM income_others
      WHERE created_at BETWEEN $3 AND $4
      
      UNION ALL
      
      SELECT created_at as fecha, total, 'egreso_otros' as tipo, point_sale_id
      FROM expense_others
      WHERE created_at BETWEEN $5 AND $6
  ) AS mov
  JOIN point_sales ps ON ps.id = mov.point_sale_id
  WHERE mov.fecha BETWEEN $7 AND $8
  GROUP BY %s
  ORDER BY ps.id
  `, dateFormat, modo)

	ctx := context.Background()
	// sqlboiler uses bind and struct loading for raw queries, but we can also use dynamic slice of map
	// wait, Bind doesn't work out-of-the-box for `map[string]any` so we'll do raw SQL rows scan:
	rows, err := queries.Raw(query, start, end, start, end, start, end, start, end).QueryContext(ctx, r.DB)
	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Reporte", schemas.Read)
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	for rows.Next() {
		columns := make([]any, len(cols))
		columnPointers := make([]any, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, schemas.HandlerErrorDB(err, "Reporte", schemas.Read)
		}

		m := make(map[string]any)
		for i, colName := range cols {
			val := columnPointers[i].(*any)
			m[colName] = *val
		}
		resultados = append(resultados, m)
	}

	grouped := make(map[string][]map[string]any)
	for _, row := range resultados {
		// byte array to string mapping in standard pgx/sql queries
		var fecha string
		switch v := row["fecha"].(type) {
		case []byte:
			fecha = string(v)
		case string:
			fecha = v
		default:
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
			COALESCE(SUM(CASE WHEN tipo = 'ingreso_ventas'  THEN total ELSE 0 END), 0) AS total_ingresos,
			COALESCE(SUM(CASE WHEN tipo = 'egreso_compras'  THEN total ELSE 0 END), 0) AS total_egresos,
			COALESCE(SUM(CASE WHEN tipo = 'ingreso_otros'   THEN total ELSE 0 END), 0) AS total_canchas,
			COALESCE(SUM(CASE WHEN tipo = 'egreso_otros'    THEN total ELSE 0 END), 0) AS balance
		FROM (
			SELECT created_at AS fecha, total, 'ingreso_ventas' AS tipo
			FROM income_sales
			WHERE created_at BETWEEN $1 AND $2

			UNION ALL
			
			SELECT created_at AS fecha, total, 'egreso_compras' AS tipo
			FROM expense_buys
			WHERE created_at BETWEEN $3 AND $4 

			UNION ALL

			SELECT created_at AS fecha, total, 'ingreso_otros' AS tipo
			FROM income_others
			WHERE created_at BETWEEN $5 AND $6

			UNION ALL

			SELECT created_at AS fecha, total, 'egreso_otros' AS tipo
			FROM expense_others
			WHERE created_at BETWEEN $7 AND $8
		) AS mov
		WHERE mov.fecha BETWEEN $9 AND $10
		GROUP BY %s
		ORDER BY %s
	`, dateFormat, modo, modo)

	ctx := context.Background()
	rows, err := queries.Raw(query, start, end, start, end, start, end, start, end, start, end).QueryContext(ctx, r.DB)
	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Reporte", schemas.Read)
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	for rows.Next() {
		columns := make([]any, len(cols))
		columnPointers := make([]any, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, schemas.HandlerErrorDB(err, "Reporte", schemas.Read)
		}

		m := make(map[string]any)
		for i, colName := range cols {
			val := columnPointers[i].(*any)
			m[colName] = *val
		}
		resultados = append(resultados, m)
	}

	grouped := make(map[string][]map[string]any)
	for _, row := range resultados {
		var fecha string
		switch v := row["fecha"].(type) {
		case []byte:
			fecha = string(v)
		case string:
			fecha = v
		default:
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
			AND ii.created_at BETWEEN $1 AND $2
		GROUP BY p.id, p.code, p.name
		ORDER BY total_profit DESC, total_quantity DESC
	`

	ctx := context.Background()
	if err := queries.Raw(query, start, end).Bind(ctx, r.DB, &products); err != nil {
		return nil, schemas.HandlerErrorDB(err, "Reporte", schemas.Read)
	}

	return products, nil
}

func (r *ReportRepository) ReportStockProducts() ([]schemas.ReportStockProduct, error) {
	ctx := context.Background()

	products, err := boilmodels.Products(
		qm.Load(boilmodels.ProductRels.Category),
		qm.Load(qm.Rels(boilmodels.ProductRels.StockPointSales, boilmodels.StockPointSaleRels.PointSale)),
		qm.Load(boilmodels.ProductRels.Deposits),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Reporte", schemas.Read)
	}

	var response []schemas.ReportStockProduct
	for _, p := range products {
		price, _ := p.Price.Big.Float64()
		minAmount, _ := p.MinAmount.Big.Float64()

		prod := schemas.ReportStockProduct{
			ID:        p.ID,
			Code:      p.Code,
			Name:      p.Name,
			Price:     price,
			Notifier:  p.Notifier,
			MinAmount: minAmount,
		}

		if p.Description.Valid {
			prod.Description = &p.Description.String
		}

		if p.R.Category != nil {
			prod.Category = p.R.Category.Name
		}

		if len(p.R.Deposits) > 0 {
			stockFloat, _ := p.R.Deposits[0].Stock.Big.Float64()
			prod.DepositStock = stockFloat
		}

		for _, sp := range p.R.StockPointSales {
			pointSaleName := ""
			if sp.R.PointSale != nil {
				pointSaleName = sp.R.PointSale.Name
			}
			stockFloat, _ := sp.Stock.Big.Float64()
			prod.PointSaleStocks = append(prod.PointSaleStocks, schemas.ReportStockPointSale{
				PointSaleID:   sp.PointSaleID,
				PointSaleName: pointSaleName,
				Stock:         stockFloat,
			})
		}

		response = append(response, prod)
	}

	return response, nil
}
