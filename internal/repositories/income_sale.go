package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func (i *IncomeSaleRepository) IncomeSaleGetByID(pointSaleID, id int64) (*schemas.IncomeSaleResponse, error) {
	var incomeSaleModel models.IncomeSale

	if err := i.DB.
		Preload("Member", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "username")
		}).
		Preload("Client", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "company_name", "identifier", "email", "phone")
		}).
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "income_sale_id", "product_id", "amount", "discount", "type_discount", "price", "subtotal", "total")
		}).
		Preload("Items.Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "code", "name", "price")
		}).
		Preload("Pay", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "total", "method_pay")
		}).
		First(&incomeSaleModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "ingreso no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "error al obtener los ingresos", err)
	}

	var incomeSaleSchema schemas.IncomeSaleResponse
	_ = copier.Copy(&incomeSaleSchema, &incomeSaleModel)

	return &incomeSaleSchema, nil
}

func (i *IncomeSaleRepository) IncomeSaleGetByDate(pointSaleID int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.IncomeSaleResponseDTO, int64, error) {
	offSet := (page - 1) * limit

	var incomeSaleModel []models.IncomeSale
	if err := i.DB.
		Preload("Member", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "username")
		}). 
		Preload("Client", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "company_name")
		}).
		Preload("Pay", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "total", "method_pay")
		}).
		Where("created_at BETWEEN ? AND ?", fromDate, toDate).
		Where("point_sale_id = ?", pointSaleID).
		Order("created_at DESC").
		Offset(offSet).
		Limit(limit).
		Find(&incomeSaleModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, schemas.ErrorResponse(404, "ingreso no encontrado", err)
		}
		return nil, 0, schemas.ErrorResponse(500, "error al obtener los ingresos", err)
	}

	var total int64
	if err := i.DB.Model(&models.IncomeSale{}).
		Where("created_at BETWEEN ? AND ?", fromDate, toDate).
		Where("point_sale_id = ?", pointSaleID).
		Count(&total).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "error al contar los ingresos", err)
	}

	var incomeSaleSchema []*schemas.IncomeSaleResponseDTO
	_ = copier.Copy(&incomeSaleSchema, &incomeSaleModel)

	return incomeSaleSchema, total, nil
}

func (i *IncomeSaleRepository) IncomeSaleCreate(memberID, pointSaleID int64, incomeSaleCreate *schemas.IncomeSaleCreate) (int64, error) {
	var incomeSaleID int64
	err := i.DB.Transaction(func(tx *gorm.DB) error {
		var register models.CashRegister
		if err := tx.
			Where("is_close = ? AND point_sale_id = ?", false, pointSaleID).
			Order("hour_open DESC").
			First(&register).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(400, "No hay caja abierta para este punto de venta", err)
			}
			return schemas.ErrorResponse(500, "Error al obtener la apertura de caja", err)
		}

		var isDeposit bool
		if err := tx.Model(&models.PointSale{}).
			Select("is_deposit").
			Where("id = ?", pointSaleID).
			Scan(&isDeposit).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al obtener el punto de venta", err)
		}

		var clientExist models.Client
		if err := tx.
			Select("id").
			Where("id = ?", incomeSaleCreate.ClientID).
			First(&clientExist).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(400, fmt.Sprintf("El cliente %d no existe", incomeSaleCreate.ClientID), err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el cliente", err)
		}

		var incomeSaleItems []*models.IncomeSaleItem
		subtotal := 0.0

		for _, item := range incomeSaleCreate.Items {
			var productPrice models.Product
			if err := tx.Select("price").
				Where("id = ?", item.ProductID).First(&productPrice).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return schemas.ErrorResponse(400, fmt.Sprintf("El producto %d no existe", item.ProductID), err)
				}
				return schemas.ErrorResponse(500, "Error al obtener el producto", err)
			}
			// Buscar stock del producto en el punto de venta
			if isDeposit {
				var stock models.Deposit
				if err := tx.
					Where("product_id = ?", item.ProductID).
					First(&stock).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return schemas.ErrorResponse(400, fmt.Sprintf("El producto %d no tiene stock en este punto de venta", item.ProductID), err)
					}
					return schemas.ErrorResponse(500, "Error al obtener stock", err)
				}

				// Validar stock suficiente
				if stock.Stock < float64(item.Amount) {
					return schemas.ErrorResponse(
						400,
						fmt.Sprintf("stock insuficiente para el producto %d (disponible: %.2f, requerido: %v)", item.ProductID, stock.Stock, item.Amount),
						fmt.Errorf("stock insuficiente para el producto %d (disponible: %.2f, requerido: %v)", item.ProductID, stock.Stock, item.Amount),
					)
				}

				// Restar stock
				stock.Stock -= float64(item.Amount)
				if err := tx.Save(&stock).Error; err != nil {
					return schemas.ErrorResponse(500, "Error al restar stock", err)
				}
			} else {
				var stock models.StockPointSale
				if err := tx.
					Where("point_sale_id = ? AND product_id = ?", pointSaleID, item.ProductID).
					First(&stock).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return schemas.ErrorResponse(400, fmt.Sprintf("El producto %d no tiene stock en este punto de venta", item.ProductID), err)
					}
					return schemas.ErrorResponse(500, "Error al obtener stock", err)
				}

				// Validar stock suficiente
				if stock.Stock < float64(item.Amount) {
					return schemas.ErrorResponse(
						400,
						fmt.Sprintf("stock insuficiente para el producto %d (disponible: %.2f, requerido: %v)", item.ProductID, stock.Stock, item.Amount),
						fmt.Errorf("stock insuficiente para el producto %d (disponible: %.2f, requerido: %v)", item.ProductID, stock.Stock, item.Amount),
					)
				}

				// Restar stock
				stock.Stock -= float64(item.Amount)
				if err := tx.Save(&stock).Error; err != nil {
					return schemas.ErrorResponse(500, "Error al actualizar stock", err)
				}
			}

			var priceCost models.ExpenseBuyItem
			if err := tx.
				Select("price").
				Where("product_id = ?", item.ProductID).
				Order("created_at DESC").
				First(&priceCost).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					priceCost.Price = productPrice.Price
				}
			}

			subtotalItem := item.Amount * productPrice.Price
			totalItem := 0.0

			if item.Discount > 0 {
				if item.TypeDiscount == "amount" {
					totalItem = subtotalItem - item.Discount
				} else if item.TypeDiscount == "percent" {
					totalItem = subtotalItem - (subtotalItem * item.Discount / 100)
				}
			} else {
				totalItem = subtotalItem
			}

			// Crear item en memoria
			incomeSaleItems = append(incomeSaleItems, &models.IncomeSaleItem{
				ProductID:    item.ProductID,
				Amount:       item.Amount,
				Price:        productPrice.Price,
				Price_Cost:   priceCost.Price,
				Subtotal:     subtotalItem,
				Discount:     item.Discount,
				TypeDiscount: item.TypeDiscount,
				Total:        totalItem,
			})

			subtotal += totalItem
		}

		totalIncome := 0.0
		if incomeSaleCreate.Discount > 0 {
			if incomeSaleCreate.Type == "amount" {
				totalIncome = subtotal - incomeSaleCreate.Discount
			} else if incomeSaleCreate.Type == "percent" {
				totalIncome = subtotal - (subtotal * incomeSaleCreate.Discount / 100)
			}
		} else {
			totalIncome = subtotal
		}

		income := models.IncomeSale{
			PointSaleID:    pointSaleID,
			MemberID:       memberID,
			ClientID:       incomeSaleCreate.ClientID,
			CashRegisterID: register.ID,
			Subtotal:       subtotal,
			Discount:       incomeSaleCreate.Discount,
			Type:           incomeSaleCreate.Type,
			Total:          totalIncome,
			IsBudget:       incomeSaleCreate.IsBudget,
		}

		if err := tx.Create(&income).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al crear el ingreso", err)
		}
		incomeSaleID = income.ID

		for _, item := range incomeSaleItems {
			item.IncomeSaleID = incomeSaleID
		}
		if err := tx.Create(&incomeSaleItems).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al crear items del ingreso", err)
		}

		var payModels []models.PayIncome
		totalPay := 0.0
		for _, pay := range incomeSaleCreate.Pay {
			totalPay += pay.Total
			payModels = append(payModels, models.PayIncome{
				IncomeSaleID:   incomeSaleID,
				CashRegisterID: &register.ID,
				ClientID:       &clientExist.ID,
				Total:         pay.Total,
				MethodPay:      pay.MethodPay,
			})
		}

		if totalPay != income.Total {
			message := fmt.Sprintf("la suma de los pagos (%.2f) no puede ser diferente a el total de la venta (%.2f)", totalPay, income.Total)
			return schemas.ErrorResponse(422, message, fmt.Errorf("%s", message))
		}

		if err := tx.Create(&payModels).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al crear pagos del ingreso", err)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return incomeSaleID, nil
}

func (i *IncomeSaleRepository) IncomeSaleUpdate(memberID, pointSaleID int64, incomeSaleUpdate *schemas.IncomeSaleUpdate) error {
	return i.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que la venta existe y pertenece al punto de venta
		var existingIncome models.IncomeSale
		if err := tx.
			Where("id = ? AND point_sale_id = ?", incomeSaleUpdate.ID, pointSaleID).
			First(&existingIncome).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Venta no encontrada", err)
			}
			return schemas.ErrorResponse(500, "Error al obtener la venta", err)
		}

		var isDeposit bool
		if err := tx.Model(&models.PointSale{}).
			Select("is_deposit").
			Where("id = ?", pointSaleID).
			Scan(&isDeposit).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al obtener el punto de venta", err)
		}

		// Verificar que el cliente existe
		var clientExist models.Client
		if err := tx.
			Select("id").
			Where("id = ?", incomeSaleUpdate.ClientID).
			First(&clientExist).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(400, fmt.Sprintf("El cliente %d no existe", incomeSaleUpdate.ClientID), err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el cliente", err)
		}

		// Obtener items anteriores para revertir el stock
		var oldItems []models.IncomeSaleItem
		if err := tx.Where("income_sale_id = ?", incomeSaleUpdate.ID).Find(&oldItems).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al obtener items anteriores", err)
		}

		// Revertir stock de los items anteriores
		for _, oldItem := range oldItems {
			if isDeposit {
				if err := tx.Model(&models.Deposit{}).
					Where("product_id = ?", oldItem.ProductID).
					UpdateColumn("stock", gorm.Expr("stock + ?", oldItem.Amount)).Error; err != nil {
					return schemas.ErrorResponse(500, "Error al revertir stock en depósito", err)
				}
			} else {
				if err := tx.Model(&models.StockPointSale{}).
					Where("point_sale_id = ? AND product_id = ?", pointSaleID, oldItem.ProductID).
					UpdateColumn("stock", gorm.Expr("stock + ?", oldItem.Amount)).Error; err != nil {
					return schemas.ErrorResponse(500, "Error al revertir stock", err)
				}
			}
		}

		// Eliminar items anteriores
		if err := tx.Where("income_sale_id = ?", incomeSaleUpdate.ID).Delete(&models.IncomeSaleItem{}).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al eliminar items anteriores", err)
		}

		// Procesar nuevos items
		var newIncomeSaleItems []*models.IncomeSaleItem
		subtotal := 0.0

		for _, item := range incomeSaleUpdate.Items {
			var productPrice models.Product
			if err := tx.Select("price").
				Where("id = ?", item.ProductID).First(&productPrice).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return schemas.ErrorResponse(400, fmt.Sprintf("El producto %d no existe", item.ProductID), err)
				}
				return schemas.ErrorResponse(500, "Error al obtener el producto", err)
			}

			// Validar y descontar stock
			if isDeposit {
				var stock models.Deposit
				if err := tx.
					Where("product_id = ?", item.ProductID).
					First(&stock).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return schemas.ErrorResponse(400, fmt.Sprintf("El producto %d no tiene stock en este punto de venta", item.ProductID), err)
					}
					return schemas.ErrorResponse(500, "Error al obtener stock", err)
				}

				if stock.Stock < float64(item.Amount) {
					return schemas.ErrorResponse(
						400,
						fmt.Sprintf("Stock insuficiente para el producto %d (disponible: %.2f, requerido: %.2f)", item.ProductID, stock.Stock, item.Amount),
						fmt.Errorf("stock insuficiente"),
					)
				}

				stock.Stock -= float64(item.Amount)
				if err := tx.Save(&stock).Error; err != nil {
					return schemas.ErrorResponse(500, "Error al actualizar stock", err)
				}
			} else {
				var stock models.StockPointSale
				if err := tx.
					Where("point_sale_id = ? AND product_id = ?", pointSaleID, item.ProductID).
					First(&stock).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return schemas.ErrorResponse(400, fmt.Sprintf("El producto %d no tiene stock en este punto de venta", item.ProductID), err)
					}
					return schemas.ErrorResponse(500, "Error al obtener stock", err)
				}

				if stock.Stock < float64(item.Amount) {
					return schemas.ErrorResponse(
						400,
						fmt.Sprintf("Stock insuficiente para el producto %d (disponible: %.2f, requerido: %.2f)", item.ProductID, stock.Stock, item.Amount),
						fmt.Errorf("stock insuficiente"),
					)
				}

				stock.Stock -= float64(item.Amount)
				if err := tx.Save(&stock).Error; err != nil {
					return schemas.ErrorResponse(500, "Error al actualizar stock", err)
				}
			}

			// Obtener precio de costo
			var priceCost models.ExpenseBuyItem
			if err := tx.
				Select("price").
				Where("product_id = ?", item.ProductID).
				Order("created_at DESC").
				First(&priceCost).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					priceCost.Price = productPrice.Price
				}
			}

			subtotalItem := item.Amount * productPrice.Price
			totalItem := 0.0

			if item.Discount > 0 {
				if item.TypeDiscount == "amount" {
					totalItem = subtotalItem - item.Discount
				} else if item.TypeDiscount == "percent" {
					totalItem = subtotalItem - (subtotalItem * item.Discount / 100)
				}
			} else {
				totalItem = subtotalItem
			}

			newIncomeSaleItems = append(newIncomeSaleItems, &models.IncomeSaleItem{
				IncomeSaleID: incomeSaleUpdate.ID,
				ProductID:    item.ProductID,
				Amount:       item.Amount,
				Price:        productPrice.Price,
				Price_Cost:   priceCost.Price,
				Subtotal:     subtotalItem,
				Discount:     item.Discount,
				TypeDiscount: item.TypeDiscount,
				Total:        totalItem,
			})

			subtotal += totalItem
		}

		// Calcular total
		totalIncome := 0.0
		if incomeSaleUpdate.Discount > 0 {
			if incomeSaleUpdate.Type == "amount" {
				totalIncome = subtotal - incomeSaleUpdate.Discount
			} else if incomeSaleUpdate.Type == "percent" {
				totalIncome = subtotal - (subtotal * incomeSaleUpdate.Discount / 100)
			}
		} else {
			totalIncome = subtotal
		}

		// Actualizar venta
		existingIncome.ClientID = incomeSaleUpdate.ClientID
		existingIncome.Subtotal = subtotal
		existingIncome.Discount = incomeSaleUpdate.Discount
		existingIncome.Type = incomeSaleUpdate.Type
		existingIncome.Total = totalIncome
		existingIncome.IsBudget = incomeSaleUpdate.IsBudget

		if err := tx.Save(&existingIncome).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al actualizar la venta", err)
		}

		// Crear nuevos items
		if err := tx.Create(&newIncomeSaleItems).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al crear nuevos items", err)
		}

		// Eliminar pagos anteriores
		if err := tx.Where("income_sale_id = ?", incomeSaleUpdate.ID).Delete(&models.PayIncome{}).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al eliminar pagos anteriores", err)
		}

		// Crear nuevos pagos
		var payModels []models.PayIncome
		totalPay := 0.0
		for _, pay := range incomeSaleUpdate.Pay {
			totalPay += pay.Total
			payModels = append(payModels, models.PayIncome{
				IncomeSaleID:   incomeSaleUpdate.ID,
				CashRegisterID: &existingIncome.CashRegisterID,
				ClientID:       &clientExist.ID,
				Total:         pay.Total,
				MethodPay:      pay.MethodPay,
			})
		}

		if totalPay != existingIncome.Total {
			message := fmt.Sprintf("La suma de los pagos (%.2f) no puede ser diferente al total de la venta (%.2f)", totalPay, existingIncome.Total)
			return schemas.ErrorResponse(422, message, fmt.Errorf("%s", message))
		}

		if err := tx.Create(&payModels).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al crear nuevos pagos", err)
		}

		return nil
	})
}

// IncomeSaleDelete elimina una venta y revierte el stock
func (i *IncomeSaleRepository) IncomeSaleDelete(incomeSaleID, pointSaleID int64) error {
	return i.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que la venta existe y pertenece al punto de venta
		var existingIncome models.IncomeSale
		if err := tx.
			Where("id = ? AND point_sale_id = ?", incomeSaleID, pointSaleID).
			First(&existingIncome).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Venta no encontrada", err)
			}
			return schemas.ErrorResponse(500, "Error al obtener la venta", err)
		}

		// Verificar si es depósito
		var isDeposit bool
		if err := tx.Model(&models.Deposit{}).
			Select("is_deposit").
			Where("id = ?", pointSaleID).
			Scan(&isDeposit).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al obtener el punto de venta", err)
		}

		// Obtener items para revertir stock
		var items []models.IncomeSaleItem
		if err := tx.Where("income_sale_id = ?", incomeSaleID).Find(&items).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al obtener items de la venta", err)
		}

		// Revertir stock
		for _, item := range items {
			if isDeposit {
				if err := tx.Model(&models.Deposit{}).
					Where("product_id = ?", item.ProductID).
					UpdateColumn("stock", gorm.Expr("stock + ?", item.Amount)).Error; err != nil {
					return schemas.ErrorResponse(500, "Error al revertir stock en depósito", err)
				}
			} else {
				if err := tx.Model(&models.StockPointSale{}).
					Where("point_sale_id = ? AND product_id = ?", pointSaleID, item.ProductID).
					UpdateColumn("stock", gorm.Expr("stock + ?", item.Amount)).Error; err != nil {
					return schemas.ErrorResponse(500, "Error al revertir stock", err)
				}
			}
		}

		// Eliminar pagos
		if err := tx.Where("income_sale_id = ?", incomeSaleID).Delete(&models.PayIncome{}).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al eliminar pagos", err)
		}

		// Eliminar items
		if err := tx.Where("income_sale_id = ?", incomeSaleID).Delete(&models.IncomeSaleItem{}).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al eliminar items", err)
		}

		// Eliminar venta
		if err := tx.Delete(&existingIncome).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al eliminar la venta", err)
		}

		return nil
	})
}