package services

import (
	"fmt"
	"sort"
	"strconv"

	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/xuri/excelize/v2"
)

func (s *ReportService) ReportExcelGet(start, end time.Time) (*excelize.File, error) {
	excel := excelize.NewFile()
	
	daysDiff := end.Sub(start).Hours() / 24
	form := "day"
	if daysDiff > 60 {
		form = "month"
	}

	productsData, profitableData, movementsData, err := s.fetchDataConcurrently(start, end, form)
	if err != nil {
		return nil, err
	}

	headerStyle := createHeaderStyle(excel)
	
	if err := s.createProductsSheet(excel, productsData, headerStyle); err != nil {
		return nil, err
	}
	
	if err := s.createProfitableSheet(excel, profitableData, headerStyle); err != nil {
		return nil, err
	}
	
	if err := s.createMovementsSheet(excel, movementsData, headerStyle); err != nil {
		return nil, err
	}

	return excel, nil
}

func (s *ReportService) fetchDataConcurrently(start, end time.Time, form string) (
	products []*models.Product,
	profitableProducts []schemas.ReportProfitableProducts,
	movements []map[string]any,
	err error,
) {
	type productResult struct {
		products []*models.Product
		err      error
	}
	type profitableResult struct {
		products []schemas.ReportProfitableProducts
		err      error
	}
	type movementResult struct {
		data any
		err  error
	}

	productsChan := make(chan productResult, 1)
	profitableChan := make(chan profitableResult, 1)
	movementsChan := make(chan movementResult, 1)

	go func() {
		prods, err := s.ReportRepository.ReportStockProducts()
		productsChan <- productResult{products: prods, err: err}
	}()

	go func() {
		profs, err := s.ReportRepository.ReportProfitableProducts(start, end)
		profitableChan <- profitableResult{products: profs, err: err}
	}()

	go func() {
		movs, err := s.ReportRepository.ReportMovementByDate(start, end, form)
		movementsChan <- movementResult{data: movs, err: err}
	}()

	productsRes := <-productsChan
	if productsRes.err != nil {
		return nil, nil, nil, productsRes.err
	}

	profitableRes := <-profitableChan
	if profitableRes.err != nil {
		return nil, nil, nil, profitableRes.err
	}

	movementsRes := <-movementsChan
	if movementsRes.err != nil {
		return nil, nil, nil, movementsRes.err
	}

	movements, ok := movementsRes.data.([]map[string]any)
	if !ok {
		return nil, nil, nil, fmt.Errorf("formato incorrecto de datos de movimientos")
	}

	return productsRes.products, profitableRes.products, movements, nil
}

func (s *ReportService) createProductsSheet(excel *excelize.File, products []*models.Product, headerStyle int) error {
	sheet := "Sheet1"
	excel.SetSheetName(sheet, "productos")
	sheet = "productos"

	points := extractPointSales(products)

	headers := []string{
		"ID", "Código", "Nombre", "Descripción", "Precio", "Categoría", "Notificador", "Cantidad Mínima",
		"Stock Depósito", "Stock Total",
	}
	for _, p := range points {
		headers = append(headers, "Stock - "+p.Name)
	}

	for i, h := range headers {
		col := colLetter(i)
		excel.SetCellValue(sheet, col+"1", h)
	}

	lastCol := colLetter(len(headers) - 1)
	excel.SetCellStyle(sheet, "A1", lastCol+"1", headerStyle)

	categoryCount := writeProductsData(excel, sheet, products, points)

	configureSheet(excel, sheet, len(headers))

	if err := addCategorySummary(excel, sheet, categoryCount, len(products), len(headers)); err != nil {
		return err
	}

	return nil
}

func (s *ReportService) createProfitableSheet(excel *excelize.File, profitableProducts []schemas.ReportProfitableProducts, headerStyle int) error {
	profitableSheet := "rentables"
	excel.NewSheet(profitableSheet)

	profitableHeaders := []string{
		"ID", "Código", "Nombre", "Cantidad Total", "Costo Total", "Ventas Totales", "Ganancia Total",
	}

	for i, h := range profitableHeaders {
		col := colLetter(i)
		excel.SetCellValue(profitableSheet, col+"1", h)
	}

	lastProfitableCol := colLetter(len(profitableHeaders) - 1)
	excel.SetCellStyle(profitableSheet, "A1", lastProfitableCol+"1", headerStyle)

	for i, product := range profitableProducts {
		row := strconv.Itoa(i + 2)
		excel.SetCellValue(profitableSheet, "A"+row, product.ID)
		excel.SetCellValue(profitableSheet, "B"+row, product.Code)
		excel.SetCellValue(profitableSheet, "C"+row, product.Name)
		excel.SetCellValue(profitableSheet, "D"+row, fmt.Sprintf("%.2f", product.TotalQuantity))
		excel.SetCellValue(profitableSheet, "E"+row, fmt.Sprintf("%.2f", product.TotalCost))
		excel.SetCellValue(profitableSheet, "F"+row, fmt.Sprintf("%.2f", product.TotalSales))
		excel.SetCellValue(profitableSheet, "G"+row, fmt.Sprintf("%.2f", product.TotalProfit))
	}

	configureSheet(excel, profitableSheet, len(profitableHeaders))

	if err := addProfitableChart(excel, profitableSheet, profitableProducts, len(profitableHeaders)); err != nil {
		return err
	}

	return nil
}

func (s *ReportService) createMovementsSheet(excel *excelize.File, movements []map[string]any, headerStyle int) error {
	movementsSheet := "movimientos"
	excel.NewSheet(movementsSheet)

	movementsHeaders := []string{
		"Fecha", "Punto de Venta", "Total Ingresos", "Total Egresos", "Total Canchas", "Balance",
	}

	for i, h := range movementsHeaders {
		col := colLetter(i)
		excel.SetCellValue(movementsSheet, col+"1", h)
	}

	lastMovCol := colLetter(len(movementsHeaders) - 1)
	excel.SetCellStyle(movementsSheet, "A1", lastMovCol+"1", headerStyle)

	currentRow := writeMovementsData(excel, movementsSheet, movements)

	configureSheet(excel, movementsSheet, len(movementsHeaders))

	if currentRow > 2 {
		if err := addMovementsSummary(excel, movementsSheet, movements, currentRow, len(movementsHeaders)); err != nil {
			return err
		}
	}

	return nil
}

func extractPointSales(products []*models.Product) []psItem {
	pointMap := make(map[int64]string)
	for _, p := range products {
		if p.StockPointSales == nil {
			continue
		}
		for _, sp := range p.StockPointSales {
			pointMap[sp.PointSaleID] = sp.PointSale.Name
		}
	}

	var points []psItem
	for id, name := range pointMap {
		points = append(points, psItem{ID: id, Name: name})
	}
	sort.Slice(points, func(i, j int) bool { return points[i].Name < points[j].Name })

	return points
}

func writeProductsData(excel *excelize.File, sheet string, products []*models.Product, points []psItem) map[string]int {
	categoryCount := make(map[string]int)

	for i, product := range products {
		row := strconv.Itoa(i + 2)

		excel.SetCellValue(sheet, "A"+row, product.ID)
		excel.SetCellValue(sheet, "B"+row, product.Code)
		excel.SetCellValue(sheet, "C"+row, product.Name)
		if product.Description != nil {
			excel.SetCellValue(sheet, "D"+row, *product.Description)
		}
		excel.SetCellValue(sheet, "E"+row, product.Price)
		excel.SetCellValue(sheet, "F"+row, product.Category.Name)
		excel.SetCellValue(sheet, "G"+row, product.Notifier)
		excel.SetCellValue(sheet, "H"+row, product.MinAmount)

		depositStock := 0.0
		if product.StockDeposit != nil {
			depositStock = product.StockDeposit.Stock
		}

		pointStockMap := make(map[int64]float64)
		if product.StockPointSales != nil {
			for _, sp := range product.StockPointSales {
				pointStockMap[sp.PointSaleID] = sp.Stock
			}
		}

		totalStock := depositStock
		excel.SetCellValue(sheet, colLetter(8)+row, fmt.Sprintf("%.2f", depositStock))

		for idx, p := range points {
			st := 0.0
			if v, ok := pointStockMap[p.ID]; ok {
				st = v
			}
			totalStock += st
			colIdx := 10 + idx
			excel.SetCellValue(sheet, colLetter(colIdx)+row, fmt.Sprintf("%.2f", st))
		}

		excel.SetCellValue(sheet, colLetter(9)+row, fmt.Sprintf("%.2f", totalStock))
		categoryCount[product.Category.Name]++
	}

	return categoryCount
}

func writeMovementsData(excel *excelize.File, sheet string, movements []map[string]any) int {
	currentRow := 2
	for _, item := range movements {
		fecha := item["fecha"].(string)
		movimientos, ok := item["movimiento"].([]map[string]any)
		if !ok {
			continue
		}

		for _, mov := range movimientos {
			row := strconv.Itoa(currentRow)

			excel.SetCellValue(sheet, "A"+row, fecha)
			excel.SetCellValue(sheet, "B"+row, mov["point_sale_name"])
			
			excel.SetCellValue(sheet, "C"+row, fmt.Sprintf("%.2f", convertToFloat64(mov["total_ingresos"])))
			excel.SetCellValue(sheet, "D"+row, fmt.Sprintf("%.2f", convertToFloat64(mov["total_egresos"])))
			excel.SetCellValue(sheet, "E"+row, fmt.Sprintf("%.2f", convertToFloat64(mov["total_canchas"])))
			excel.SetCellValue(sheet, "F"+row, fmt.Sprintf("%.2f", convertToFloat64(mov["balance"])))

			currentRow++
		}
	}
	return currentRow
}

// addCategorySummary agrega resumen y gráfico por categoría
func addCategorySummary(excel *excelize.File, sheet string, categoryCount map[string]int, productCount, headerCount int) error {
	resumeStart := productCount + 3
	resumeRow := strconv.Itoa(resumeStart)

	summaryColStartIdx := headerCount + 1
	summaryCatCol := colLetter(summaryColStartIdx)
	summaryCountCol := colLetter(summaryColStartIdx + 1)

	excel.SetCellValue(sheet, summaryCatCol+resumeRow, "Category")
	excel.SetCellValue(sheet, summaryCountCol+resumeRow, "Count")

	r := resumeStart + 1
	for cat, count := range categoryCount {
		excel.SetCellValue(sheet, summaryCatCol+strconv.Itoa(r), cat)
		excel.SetCellValue(sheet, summaryCountCol+strconv.Itoa(r), count)
		r++
	}

	categoryRange := fmt.Sprintf("%s!$%s$%d:$%s$%d", sheet, summaryCatCol, resumeStart+1, summaryCatCol, r-1)
	countRange := fmt.Sprintf("%s!$%s$%d:$%s$%d", sheet, summaryCountCol, resumeStart+1, summaryCountCol, r-1)

	return excel.AddChart(sheet, summaryCatCol+"2", &excelize.Chart{
		Type: excelize.Pie3D,
		Series: []excelize.ChartSeries{
			{
				Name:       "Productos por Categoría",
				Categories: categoryRange,
				Values:     countRange,
			},
		},
		Title: []excelize.RichTextRun{{Text: "Distribución de Productos por Categoría"}},
	})
}

// addProfitableChart agrega el gráfico de top 10 productos rentables
func addProfitableChart(excel *excelize.File, sheet string, products []schemas.ReportProfitableProducts, headerCount int) error {
	if len(products) == 0 {
		return nil
	}

	chartStartRow := len(products) + 3
	
	sortedProducts := make([]schemas.ReportProfitableProducts, len(products))
	copy(sortedProducts, products)
	sort.Slice(sortedProducts, func(i, j int) bool {
		return sortedProducts[i].TotalProfit > sortedProducts[j].TotalProfit
	})

	topCount := 10
	if len(sortedProducts) < topCount {
		topCount = len(sortedProducts)
	}

	chartCol := headerCount + 2
	nameCol := colLetter(chartCol)
	profitCol := colLetter(chartCol + 1)

	excel.SetCellValue(sheet, nameCol+strconv.Itoa(chartStartRow), "Producto")
	excel.SetCellValue(sheet, profitCol+strconv.Itoa(chartStartRow), "Ganancia")

	for i := 0; i < topCount; i++ {
		rowNum := chartStartRow + 1 + i
		excel.SetCellValue(sheet, nameCol+strconv.Itoa(rowNum), sortedProducts[i].Name)
		excel.SetCellValue(sheet, profitCol+strconv.Itoa(rowNum), sortedProducts[i].TotalProfit)
	}

	nameRange := fmt.Sprintf("%s!$%s$%d:$%s$%d", sheet, nameCol, chartStartRow+1, nameCol, chartStartRow+topCount)
	profitRange := fmt.Sprintf("%s!$%s$%d:$%s$%d", sheet, profitCol, chartStartRow+1, profitCol, chartStartRow+topCount)

	return excel.AddChart(sheet, nameCol+"2", &excelize.Chart{
		Type: excelize.Bar,
		Series: []excelize.ChartSeries{
			{
				Name:       "Ganancia",
				Categories: nameRange,
				Values:     profitRange,
			},
		},
		Title: []excelize.RichTextRun{{Text: "Top 10 Productos Más Rentables"}},
	})
}

// addMovementsSummary agrega resumen y gráfico de movimientos
func addMovementsSummary(excel *excelize.File, sheet string, movements []map[string]any, currentRow, headerCount int) error {
	chartStartRow := currentRow + 2
	
	pointSaleTotals := make(map[string]map[string]float64)
	
	for _, item := range movements {
		movimientos, ok := item["movimiento"].([]map[string]any)
		if !ok {
			continue
		}

		for _, mov := range movimientos {
			pointName := mov["point_sale_name"].(string)
			
			if pointSaleTotals[pointName] == nil {
				pointSaleTotals[pointName] = make(map[string]float64)
			}

			pointSaleTotals[pointName]["balance"] += convertToFloat64(mov["balance"])
		}
	}

	summaryCol := headerCount + 2
	pointCol := colLetter(summaryCol)
	balanceCol := colLetter(summaryCol + 1)

	excel.SetCellValue(sheet, pointCol+strconv.Itoa(chartStartRow), "Punto de Venta")
	excel.SetCellValue(sheet, balanceCol+strconv.Itoa(chartStartRow), "Balance Total")

	summaryRow := chartStartRow + 1
	for pointName, totals := range pointSaleTotals {
		excel.SetCellValue(sheet, pointCol+strconv.Itoa(summaryRow), pointName)
		excel.SetCellValue(sheet, balanceCol+strconv.Itoa(summaryRow), totals["balance"])
		summaryRow++
	}

	if summaryRow > chartStartRow+1 {
		pointRange := fmt.Sprintf("%s!$%s$%d:$%s$%d", sheet, pointCol, chartStartRow+1, pointCol, summaryRow-1)
		balanceRange := fmt.Sprintf("%s!$%s$%d:$%s$%d", sheet, balanceCol, chartStartRow+1, balanceCol, summaryRow-1)

		return excel.AddChart(sheet, pointCol+"2", &excelize.Chart{
			Type: excelize.Col,
			Series: []excelize.ChartSeries{
				{
					Name:       "Balance Total",
					Categories: pointRange,
					Values:     balanceRange,
				},
			},
			Title: []excelize.RichTextRun{{Text: "Balance Total por Punto de Venta"}},
		})
	}

	return nil
}

// configureSheet configura el ancho de columnas y congela la primera fila
func configureSheet(excel *excelize.File, sheet string, headerCount int) {
	for i := 0; i < headerCount; i++ {
		col := colLetter(i)
		excel.SetColWidth(sheet, col, col, 20)
	}

	excel.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})
}

// Helpers y tipos
type psItem struct {
	ID   int64
	Name string
}

func createHeaderStyle(excel *excelize.File) int {
	style, _ := excel.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4F81BD"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	return style
}

func colLetter(idx int) string {
	idx++ // 1-based
	col := ""
	for idx > 0 {
		rem := (idx - 1) % 26
		col = string(rune('A'+rem)) + col
		idx = (idx - 1) / 26
	}
	return col
}

func convertToFloat64(val any) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint64:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0.0
}

func (r *ReportService) ReportMovementByDate(fromDate, toDate time.Time, form string) (any, error) {
	report, err := r.ReportRepository.ReportMovementByDate(fromDate, toDate, form)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (r *ReportService) ReportProfitableProducts(start, end time.Time) ([]schemas.ReportProfitableProducts, error) {
	report, err := r.ReportRepository.ReportProfitableProducts(start, end)
	if err != nil {
		return nil, err
	}

	return report, nil
}
