package services

import (
	"bytes"
	"math"
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

func (p *ProductService) ProductGetByID(id int64) (*schemas.ProductFullResponse, error) {
	product, err := p.ProductRepository.ProductGetByID(id)
	if err != nil {
		return nil, err
	}

	var productResponse schemas.ProductFullResponse

	productResponse.ID = product.ID
	productResponse.Code = product.Code
	productResponse.Name = product.Name
	productResponse.Description = product.Description
	productResponse.Category = schemas.CategoryResponse{
		ID:   product.Category.ID,
		Name: product.Category.Name,
	}
	productResponse.Price = product.Price

	if product.StockDeposit != nil {
		productResponse.StockDeposit = product.StockDeposit.Stock
	} else {
		productResponse.StockDeposit = 0
	}

	productResponse.Notifier = product.Notifier
	productResponse.MinAmount = product.MinAmount

	for _, stock := range product.StockPointSales {
		productResponse.StockPointSales = append(productResponse.StockPointSales, &schemas.PointSaleStockResponse{
			ID:        stock.PointSale.ID,
			Name:      stock.PointSale.Name,
			Stock:     stock.Stock,
			IsDeposit: stock.PointSale.IsDeposit,
		})
	}

	return &productResponse, nil
}

func (p *ProductService) ProductGetByCode(code string) (*schemas.ProductFullResponse, error) {
	product, err := p.ProductRepository.ProductGetByCode(code)
	if err != nil {
		return nil, err
	}

	var productResponse schemas.ProductFullResponse

	productResponse.ID = product.ID
	productResponse.Code = product.Code
	productResponse.Name = product.Name
	productResponse.Description = product.Description
	productResponse.Category = schemas.CategoryResponse{
		ID:   product.Category.ID,
		Name: product.Category.Name,
	}
	productResponse.Price = product.Price

	if product.StockDeposit != nil {
		productResponse.StockDeposit = product.StockDeposit.Stock
	} else {
		productResponse.StockDeposit = 0
	}

	productResponse.Notifier = product.Notifier
	productResponse.MinAmount = product.MinAmount

	for _, stock := range product.StockPointSales {
		productResponse.StockPointSales = append(productResponse.StockPointSales, &schemas.PointSaleStockResponse{
			ID:        stock.PointSale.ID,
			Name:      stock.PointSale.Name,
			Stock:     stock.Stock,
			IsDeposit: stock.PointSale.IsDeposit,
		})
	}

	return &productResponse, nil
}

func (p *ProductService) ProductGetByName(name string) ([]*schemas.ProductFullResponse, error) {
	products, err := p.ProductRepository.ProductGetByName(name)
	if err != nil {
		return nil, err
	}

	productsResponse := make([]*schemas.ProductFullResponse, len(products))

	for i, prod := range products {
		productsResponse[i] = &schemas.ProductFullResponse{
			ID:          prod.ID,
			Code:        prod.Code,
			Name:        prod.Name,
			Description: prod.Description,
			Category: schemas.CategoryResponse{
				ID:   prod.Category.ID,
				Name: prod.Category.Name,
			},
			Price:     prod.Price,
			Notifier:  prod.Notifier,
			MinAmount: prod.MinAmount,
		}
		if prod.StockDeposit != nil {
			productsResponse[i].StockDeposit = prod.StockDeposit.Stock
		} else {
			productsResponse[i].StockDeposit = 0
		}
		for _, stock := range prod.StockPointSales {
			productsResponse[i].StockPointSales = append(productsResponse[i].StockPointSales, &schemas.PointSaleStockResponse{
				ID:        stock.PointSale.ID,
				Name:      stock.PointSale.Name,
				Stock:     stock.Stock,
				IsDeposit: stock.PointSale.IsDeposit,
			})
		}
	}

	return productsResponse, nil
}

func (p *ProductService) ProductGetByCategoryID(categoryID int64) ([]*schemas.ProductFullResponse, error) {
	products, err := p.ProductRepository.ProductGetByCategoryID(categoryID)
	if err != nil {
		return nil, err
	}

	productsResponse := make([]*schemas.ProductFullResponse, len(products))

	for i, prod := range products {
		productsResponse[i] = &schemas.ProductFullResponse{
			ID:          prod.ID,
			Code:        prod.Code,
			Name:        prod.Name,
			Description: prod.Description,
			Category: schemas.CategoryResponse{
				ID:   prod.Category.ID,
				Name: prod.Category.Name,
			},
			Price:     prod.Price,
			Notifier:  prod.Notifier,
			MinAmount: prod.MinAmount,
		}
		if prod.StockDeposit != nil {
			productsResponse[i].StockDeposit = prod.StockDeposit.Stock
		} else {
			productsResponse[i].StockDeposit = 0
		}
		for _, stock := range prod.StockPointSales {
			productsResponse[i].StockPointSales = append(productsResponse[i].StockPointSales, &schemas.PointSaleStockResponse{
				ID:        stock.PointSale.ID,
				Name:      stock.PointSale.Name,
				Stock:     stock.Stock,
				IsDeposit: stock.PointSale.IsDeposit,
			})
		}
	}

	return productsResponse, nil
}

func (p *ProductService) ProductGetAll(page, limit int) ([]*schemas.ProductFullResponse, int64, error) {
	products, total, err := p.ProductRepository.ProductGetAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	productsResponse := make([]*schemas.ProductFullResponse, len(products))

	for i, prod := range products {
		productsResponse[i] = &schemas.ProductFullResponse{
			ID:          prod.ID,
			Code:        prod.Code,
			Name:        prod.Name,
			Description: prod.Description,
			Category: schemas.CategoryResponse{
				ID:   prod.Category.ID,
				Name: prod.Category.Name,
			},
			Price:     prod.Price,
			Notifier:  prod.Notifier,
			MinAmount: prod.MinAmount,
		}
		if prod.StockDeposit != nil {
			productsResponse[i].StockDeposit = prod.StockDeposit.Stock
		} else {
			productsResponse[i].StockDeposit = 0
		}
		for _, stock := range prod.StockPointSales {
			productsResponse[i].StockPointSales = append(productsResponse[i].StockPointSales, &schemas.PointSaleStockResponse{
				ID:        stock.PointSale.ID,
				Name:      stock.PointSale.Name,
				Stock:     stock.Stock,
				IsDeposit: stock.PointSale.IsDeposit,
			})
		}
	}

	return productsResponse, total, nil
}

// func (p *ProductService) ProductGenerateQR(code string) ([]byte, error) {
// 	// prod, err := p.ProductRepository.ProductGetByCode(code)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	prod := &models.Product{
// 		Code: code,
// 		Name: "coca cola 3lts",
// 		// Name: "El ají loco de la vuelta de la esqui se avecina viene diego rumbeando",
// 	}

// 	pdf := gofpdf.New("P", "mm", "A4", "")

// 	// Configurar traductor UTF-8
// 	tr := pdf.UnicodeTranslatorFromDescriptor("")

// 	pdf.AddPage()
// 	pdf.SetAutoPageBreak(false, 0)

// 	const (
// 		cols = 5
// 		rows = 10

// 		margin = 5.0

// 		pageW = 210.0
// 		pageH = 297.0
// 	)

// 	// Área útil (sin márgenes)
// 	usableW := pageW - (margin * 2)
// 	usableH := pageH - (margin * 2)

// 	cellW := usableW / cols
// 	cellH := usableH / rows

// 	// QR ocupa 65% del ancho
// 	qrSize := cellW * 0.65

// 	// Generar QR en memoria
// 	qrPNG, err := qrcode.Encode(prod.Code, qrcode.Medium, 256)
// 	if err != nil {
// 		return nil, err
// 	}

// 	qrName := "qr.png"
// 	pdf.RegisterImageOptionsReader(
// 		qrName,
// 		gofpdf.ImageOptions{ImageType: "PNG"},
// 		bytes.NewReader(qrPNG),
// 	)

// 	pdf.SetFont("Arial", "", 6.5)

// 	for r := 0; r < rows; r++ {
// 		for c := 0; c < cols; c++ {
// 			// Posición de la celda (considerando márgenes)
// 			x := margin + float64(c)*cellW
// 			y := margin + float64(r)*cellH

// 			// Centrar QR horizontalmente en la celda
// 			qrX := x + (cellW-qrSize)/2
// 			qrY := y + 1

// 			// Dibujar el QR
// 			pdf.Image(qrName, qrX, qrY, qrSize, qrSize, false, "", 0, "")

// 			// Posición del texto debajo del QR
// 			textY := qrY + qrSize + 0.5

// 			// Calcular si el texto cabe en el espacio disponible
// 			availableSpace := cellH - qrSize - 2

// 			// Convertir texto a UTF-8
// 			productName := tr(prod.Name)

// 			pdf.SetXY(x, textY)

// 			// Si el nombre es muy largo, usar dos líneas
// 			textWidth := pdf.GetStringWidth(productName)
// 			if textWidth > cellW-2 {
// 				words := strings.Fields(prod.Name)
// 				mid := len(words) / 2

// 				line1 := tr(strings.Join(words[:mid], " "))
// 				line2 := tr(strings.Join(words[mid:], " "))

// 				// Verificar que tenemos espacio para dos líneas
// 				if availableSpace >= 6 {
// 					pdf.SetXY(x, textY)
// 					pdf.CellFormat(cellW, 2.5, line1, "", 0, "C", false, 0, "")

// 					pdf.SetXY(x, textY+2.5)
// 					pdf.CellFormat(cellW, 2.5, line2, "", 0, "C", false, 0, "")
// 				}
// 			} else {
// 				// Una sola línea
// 				pdf.CellFormat(cellW, 3, productName, "", 0, "C", false, 0, "")
// 			}
// 		}
// 	}

// 	var buf bytes.Buffer
// 	err = pdf.Output(&buf)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return buf.Bytes(), nil
// }
// func (p *ProductService) ProductGenerateQR(code string) ([]byte, error) {
// 	// prod, err := p.ProductRepository.ProductGetByCode(code)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	prod := &models.Product{
// 		Code: code,
// 		Name: "El ají loco, mira lo que se avecian a la vuelta de la esquina viene diego",
// 	}

// 	pdf := gofpdf.New("P", "mm", "A4", "")
	
// 	// Configurar traductor UTF-8
// 	tr := pdf.UnicodeTranslatorFromDescriptor("")
	
// 	pdf.AddPage()
// 	pdf.SetAutoPageBreak(false, 0)
	
// 	const (
// 		cols = 10
// 		rows = 20
		
// 		margin = 5.0
		
// 		pageW = 210.0
// 		pageH = 297.0
// 	)

// 	// Área útil (sin márgenes)
// 	usableW := pageW - (margin * 2)
// 	usableH := pageH - (margin * 2)
	
// 	cellW := usableW / cols
// 	cellH := usableH / rows

// 	// QR ocupa 60% del ancho para dejar más espacio al texto
// 	qrSize := cellW * 0.60

// 	// Generar QR en memoria
// 	qrPNG, err := qrcode.Encode(prod.Code, qrcode.Medium, 256)
// 	if err != nil {
// 		return nil, err
// 	}

// 	qrName := "qr.png"
// 	pdf.RegisterImageOptionsReader(
// 		qrName,
// 		gofpdf.ImageOptions{ImageType: "PNG"},
// 		bytes.NewReader(qrPNG),
// 	)

// 	// Letra pequeña pero legible
// 	pdf.SetFont("Arial", "", 6)

// 	for r := 0; r < rows; r++ {
// 		for c := 0; c < cols; c++ {
// 			// Posición de la celda (considerando márgenes)
// 			x := margin + float64(c)*cellW
// 			y := margin + float64(r)*cellH

// 			// Centrar QR horizontalmente en la celda
// 			qrX := x + (cellW-qrSize)/2
// 			qrY := y + 1.5 // Padding superior

// 			// Dibujar el QR
// 			pdf.Image(qrName, qrX, qrY, qrSize, qrSize, false, "", 0, "")

// 			// Posición del texto debajo del QR - ESPACIO MÍNIMO
// 			textY := qrY + qrSize - 1.5 // Espacio reducido
			
// 			// Convertir texto a UTF-8
// 			productName := tr(prod.Name)
			
// 			pdf.SetXY(x, textY)
			
// 			// Si el nombre es muy largo, usar dos líneas
// 			textWidth := pdf.GetStringWidth(productName)
// 			if textWidth > cellW-1 {
// 				words := strings.Fields(prod.Name)
// 				mid := (len(words) + 1) / 2 // Mejor división
				
// 				line1 := tr(strings.Join(words[:mid], " "))
// 				line2 := tr(strings.Join(words[mid:], " "))
				
// 				pdf.SetXY(x+0.5, textY) // Pequeño margen lateral
// 				pdf.CellFormat(cellW-1, 2.5, line1, "", 0, "C", false, 0, "")
				
// 				pdf.SetXY(x+0.5, textY+2.5)
// 				pdf.CellFormat(cellW-1, 2.5, line2, "", 0, "C", false, 0, "")
// 			} else {
// 				// Una sola línea
// 				pdf.SetXY(x+0.5, textY) // Pequeño margen lateral
// 				pdf.CellFormat(cellW-1, 2.5, productName, "", 0, "C", false, 0, "")
// 			}
// 		}
// 	}

// 	var buf bytes.Buffer
// 	err = pdf.Output(&buf)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return buf.Bytes(), nil
// }
func (p *ProductService) ProductGenerateQR(code string) ([]byte, error) {
	prod, err := p.ProductRepository.ProductGetByCodeToQR(code)
	if err != nil {
		return nil, err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	
	// Configurar traductor UTF-8
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	
	pdf.AddPage()
	pdf.SetAutoPageBreak(false, 0)
	
	const (
		cols = 10
		rows = 20
		
		margin = 5.0
		
		pageW = 210.0
		pageH = 297.0
	)

	// Área útil (sin márgenes)
	usableW := pageW - (margin * 2)
	usableH := pageH - (margin * 2)
	
	cellW := usableW / cols
	cellH := usableH / rows

	// QR ocupa 60% del ancho para dejar más espacio al texto
	qrSize := cellW * 0.60

	// CALCULAR TAMAÑO DE FUENTE DINÁMICAMENTE
	// Fórmula: más filas/columnas = letra más pequeña
	// Base: 5 cols x 10 rows = tamaño 6
	baseCols := 5.0
	baseRows := 10.0
	baseFontSize := 6.0
	
	// Factor de escala basado en densidad de celdas
	scaleFactor := math.Sqrt((baseCols * baseRows) / (cols * rows))
	fontSize := baseFontSize * scaleFactor
	
	// Limitar tamaño mínimo y máximo
	if fontSize < 3 {
		fontSize = 3
	}
	if fontSize > 8 {
		fontSize = 8
	}

	// Generar QR en memoria
	qrPNG, err := qrcode.Encode(prod.Code, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	qrName := "qr.png"
	pdf.RegisterImageOptionsReader(
		qrName,
		gofpdf.ImageOptions{ImageType: "PNG"},
		bytes.NewReader(qrPNG),
	)

	// Aplicar tamaño de fuente calculado
	pdf.SetFont("Arial", "", fontSize)

	// Calcular altura de línea proporcional al tamaño de fuente
	lineHeight := fontSize * 0.5

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			// Posición de la celda (considerando márgenes)
			x := margin + float64(c)*cellW
			y := margin + float64(r)*cellH

			// Centrar QR horizontalmente en la celda
			qrX := x + (cellW-qrSize)/2
			qrY := y + 1.5 // Padding superior

			// Dibujar el QR
			pdf.Image(qrName, qrX, qrY, qrSize, qrSize, false, "", 0, "")

			// Posición del texto debajo del QR - ESPACIO MÍNIMO
			textY := qrY + qrSize - 1.5 // Espacio reducido
			
			// Convertir texto a UTF-8
			productName := tr(prod.Name)
			
			pdf.SetXY(x, textY)
			
			// Si el nombre es muy largo, usar dos líneas
			textWidth := pdf.GetStringWidth(productName)
			if textWidth > cellW-1 {
				words := strings.Fields(prod.Name)
				mid := (len(words) + 1) / 2 // Mejor división
				
				line1 := tr(strings.Join(words[:mid], " "))
				line2 := tr(strings.Join(words[mid:], " "))
				
				pdf.SetXY(x+0.5, textY) // Pequeño margen lateral
				pdf.CellFormat(cellW-1, lineHeight, line1, "", 0, "C", false, 0, "")
				
				pdf.SetXY(x+0.5, textY+lineHeight)
				pdf.CellFormat(cellW-1, lineHeight, line2, "", 0, "C", false, 0, "")
			} else {
				// Una sola línea
				pdf.SetXY(x+0.5, textY) // Pequeño margen lateral
				pdf.CellFormat(cellW-1, lineHeight, productName, "", 0, "C", false, 0, "")
			}
		}
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}


func (p *ProductService) ProductCreate(memberID int64, productCreate *schemas.ProductCreate, plan *schemas.PlanResponseDTO) (int64, error) {
	return p.ProductRepository.ProductCreate(memberID, productCreate, plan)
}

func (p *ProductService) ProductUpdate(memberID int64, productUpdate *schemas.ProductUpdate) error {
	return p.ProductRepository.ProductUpdate(memberID, productUpdate)
}
func (p *ProductService) ProductPriceUpdate(memberID int64, productUpdate *schemas.ListPriceUpdate) error {
	return p.ProductRepository.ProductPriceUpdate(memberID, productUpdate)
}

func (p *ProductService) ProductDelete(memberID int64, id int64) error {
	return p.ProductRepository.ProductDelete(memberID, id)

}
