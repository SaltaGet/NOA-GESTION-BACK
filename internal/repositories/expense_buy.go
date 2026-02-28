package repositories

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func (r *ExpenseBuyRepository) ExpenseBuyGetByID(id int64) (*schemas.ExpenseBuyResponse, error) {
	var expenseBuy *models.ExpenseBuy

	if err := r.DB.
		Preload("Member", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "username").Unscoped()
		}).
		Preload("ExpenseBuyItem", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "expense_buy_id", "product_id", "amount", "price", "discount", "type_discount", "subtotal", "total", "created_at")
		}).
		Preload("ExpenseBuyItem.Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "code", "name", "price")
		}).
		Preload("PayExpenseBuy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "expense_buy_id", "total", "method_pay")
		}).
		Preload("Supplier", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "company_name")
		}).
		First(&expenseBuy, id).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Egreso de compras", schemas.Read)
	}

	var expenseSchema schemas.ExpenseBuyResponse
	copier.Copy(&expenseSchema, &expenseBuy)

	return &expenseSchema, nil
}

func (r *ExpenseBuyRepository) ExpenseBuyGetByDate(fromDate, toDate time.Time, page, limit int) ([]*schemas.ExpenseBuyResponseSimple, int64, error) {
	var expensesBuy []*models.ExpenseBuy

	offSet := (page - 1) * limit

	if err := r.DB.
		Preload("Supplier", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "company_name")
		}).
		Preload("PayExpenseBuy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "expense_buy_id", "total", "method_pay")
		}).
		Where("created_at BETWEEN ? AND ?", fromDate, toDate).
		Order("created_at DESC").
		Offset(offSet).
		Limit(limit).
		Find(&expensesBuy).Error; err != nil {
		return nil, 0, schemas.HandlerErrorGorm(err, "Egreso de compras", schemas.Read)
	}

	var total int64
	if err := r.DB.Model(&models.ExpenseBuy{}).
		Where("created_at BETWEEN ? AND ?", fromDate, toDate).
		Count(&total).Error; err != nil {
		return nil, 0, schemas.HandlerErrorGorm(err, "Egreso de compras", schemas.Read)
	}

	var expenseSchema []*schemas.ExpenseBuyResponseSimple
	copier.Copy(&expenseSchema, &expensesBuy)

	return expenseSchema, total, nil
}

func (r *ExpenseBuyRepository) ExpenseBuyCreate(memberID int64, expenseBuyCreate *schemas.ExpenseBuyCreate) (int64, error) {
	var expenseBuySave models.ExpenseBuy
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var supplierID models.Supplier
		if err := tx.Select("id").First(&supplierID, expenseBuyCreate.SupplierID).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Egreso de compras", schemas.Read)
		}

		var expenseItems []*models.ExpenseBuyItem
		total := 0.0

		for _, item := range expenseBuyCreate.ExpenseBuyItem {
			if item.Amount <= 0 {
				return schemas.ErrorResponse(400, fmt.Sprintf("La cantidad para el producto %d no es válida", item.ProductID), fmt.Errorf("la cantidad para el producto %d no es válida", item.ProductID))
			}

			var product models.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Egreso de compras", schemas.Read)
			}
			// Buscar stock del producto en el punto de venta
			var stock models.Deposit
			if err := tx.
				Where("product_id = ?", item.ProductID).
				FirstOrCreate(&stock, models.Deposit{ProductID: item.ProductID, Stock: 0}).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Egreso de compras", schemas.Read)
			}

			if err := tx.Model(&stock).
				Update("stock", gorm.Expr("stock + ?", item.Amount)).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Egreso de compras", schemas.Update)
			}

			subtotalItem := item.Amount * item.Price
			totalItem := 0.0
			if item.Discount > 0 {
				if item.TypeDiscount == "percent" {
					totalItem = subtotalItem - (subtotalItem * item.Discount / 100)
				} else {
					totalItem = subtotalItem - item.Discount
				}
			} else {
				totalItem = subtotalItem
			}

			expenseItems = append(expenseItems, &models.ExpenseBuyItem{
				ProductID:    item.ProductID,
				Amount:       item.Amount,
				Price:        item.Price,
				Discount:     item.Discount,
				TypeDiscount: item.TypeDiscount,
				Subtotal:     subtotalItem,
				Total:        totalItem,
			})

			total += totalItem
		}

		totalExpense := 0.0
		if expenseBuyCreate.Discount > 0 {
			if expenseBuyCreate.TypeDiscount == "percent" {
				totalExpense = total - (total * expenseBuyCreate.Discount / 100)
			} else {
				totalExpense = total - expenseBuyCreate.Discount
			}
		} else {
			totalExpense = total
		}

		expenseBuy := models.ExpenseBuy{
			MemberID:     memberID,
			SupplierID:   supplierID.ID,
			Details:      expenseBuyCreate.Details,
			Subtotal:     total,
			Discount:     expenseBuyCreate.Discount,
			TypeDiscount: expenseBuyCreate.TypeDiscount,
			Total:        totalExpense,
		}

		if err := tx.Create(&expenseBuy).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Egreso de compras", schemas.Create)
		}

		// 🔹 Asociar items
		for _, item := range expenseItems {
			item.ExpenseBuyID = expenseBuy.ID
		}
		if err := tx.Create(&expenseItems).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Items de egreso de compras", schemas.Create)
		}

		totalPay := 0.0
		var payExpenseBuy []*models.PayExpenseBuy
		for _, pay := range expenseBuyCreate.PayExpenseBuy {
			totalPay += pay.Total
			payExpenseBuy = append(payExpenseBuy, &models.PayExpenseBuy{
				ExpenseBuyID: expenseBuy.ID,
				Total:        pay.Total,
				MethodPay:    pay.MethodPay,
			})
		}

		if math.Abs(totalPay-totalExpense) > 1 {
			message := fmt.Sprintf("la diferencia entre la suma de pagos (%.2f) y el total del egreso (%.2f)", totalPay, totalExpense)
			return schemas.ErrorResponse(400, message, fmt.Errorf("%s", message))
		}

		if err := tx.Create(&payExpenseBuy).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Pagos de egreso de compras", schemas.Create)
		}

		expenseBuySave = expenseBuy

		return nil
	})

	if err != nil {
		return 0, err
	}

	return expenseBuySave.ID, nil
}

func (r *ExpenseBuyRepository) ExpenseBuyUpdate(memberID int64, expenseBuyUpdate *schemas.ExpenseBuyUpdate) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}
	
		var existingExpense models.ExpenseBuy
		if err := tx.Where("id = ?", expenseBuyUpdate.ID).First(&existingExpense).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Egreso de compras", schemas.Read)
		}

		var supplierID models.Supplier
		if err := tx.Select("id").First(&supplierID, expenseBuyUpdate.SupplierID).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Egreso de compras", schemas.Read)
		}

		// Obtener items anteriores para revertir el stock
		var oldItems []models.ExpenseBuyItem
		if err := tx.Where("expense_buy_id = ?", expenseBuyUpdate.ID).Find(&oldItems).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Egreso de compras", schemas.Read)
		}

		// Revertir stock de los items anteriores
		for _, oldItem := range oldItems {
			if err := tx.Model(&models.Deposit{}).
				Where("product_id = ?", oldItem.ProductID).
				UpdateColumn("stock", gorm.Expr("stock - ?", oldItem.Amount)).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Stock", schemas.Update)
			}
		}

		// Eliminar items anteriores
		if err := tx.Where("expense_buy_id = ?", expenseBuyUpdate.ID).Delete(&models.ExpenseBuyItem{}).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Items de egreso de compras", schemas.Delete)
		}

		// Procesar nuevos items
		var newExpenseItems []*models.ExpenseBuyItem
		total := 0.0

		for _, item := range expenseBuyUpdate.ExpenseBuyItem {
			if item.Amount <= 0 {
				return schemas.ErrorResponse(400, fmt.Sprintf("La cantidad para el producto %d no es válida", item.ProductID), fmt.Errorf("la cantidad para el producto %d no es válida", item.ProductID))
			}

			var product models.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Producto", schemas.Read)
			}

			// Buscar o crear stock del producto en el depósito
			var stock models.Deposit
			if err := tx.
				Where("product_id = ?", item.ProductID).
				FirstOrCreate(&stock, models.Deposit{ProductID: item.ProductID, Stock: 0}).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Stock Depósito", schemas.Read)
			}

			// Sumar stock
			if err := tx.Model(&stock).
				Update("stock", gorm.Expr("stock + ?", item.Amount)).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Stock Depósito", schemas.Update)
			}

			subtotalItem := item.Amount * item.Price
			totalItem := 0.0
			if item.Discount > 0 {
				if item.TypeDiscount == "percent" {
					totalItem = subtotalItem - (subtotalItem * item.Discount / 100)
				} else {
					totalItem = subtotalItem - item.Discount
				}
			} else {
				totalItem = subtotalItem
			}

			newExpenseItems = append(newExpenseItems, &models.ExpenseBuyItem{
				ExpenseBuyID: expenseBuyUpdate.ID,
				ProductID:    item.ProductID,
				Amount:       item.Amount,
				Price:        item.Price,
				Discount:     item.Discount,
				TypeDiscount: item.TypeDiscount,
				Subtotal:     subtotalItem,
				Total:        totalItem,
			})

			total += totalItem
		}

		// Calcular total con descuento general
		totalExpense := 0.0
		if expenseBuyUpdate.Discount > 0 {
			if expenseBuyUpdate.Type == "percent" {
				totalExpense = total - (total * expenseBuyUpdate.Discount / 100)
			} else {
				totalExpense = total - expenseBuyUpdate.Discount
			}
		} else {
			totalExpense = total
		}

		// Actualizar la compra
		existingExpense.SupplierID = supplierID.ID
		existingExpense.Details = expenseBuyUpdate.Details
		existingExpense.Subtotal = total
		existingExpense.Discount = expenseBuyUpdate.Discount
		existingExpense.TypeDiscount = expenseBuyUpdate.Type
		existingExpense.Total = totalExpense

		if err := tx.Save(&existingExpense).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Egreso de compras", schemas.Update)
		}

		// Crear nuevos items
		if err := tx.Create(&newExpenseItems).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Items de egreso de compras", schemas.Create)
		}

		// Eliminar pagos anteriores
		if err := tx.Where("expense_buy_id = ?", expenseBuyUpdate.ID).Delete(&models.PayExpenseBuy{}).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Pagos de egreso de compras", schemas.Delete)
		}

		// Crear nuevos pagos
		totalPay := 0.0
		var payExpenseBuy []*models.PayExpenseBuy
		for _, pay := range expenseBuyUpdate.PayExpenseBuy {
			totalPay += pay.Total
			payExpenseBuy = append(payExpenseBuy, &models.PayExpenseBuy{
				ExpenseBuyID: expenseBuyUpdate.ID,
				Total:        pay.Total,
				MethodPay:    pay.MethodPay,
			})
		}

		if math.Abs(totalPay-totalExpense) > 1 {
			message := fmt.Sprintf("la diferencia entre la suma de pagos (%.2f) y el total del egreso (%.2f)", totalPay, totalExpense)
			return schemas.ErrorResponse(400, message, fmt.Errorf("%s", message))
		}

		if err := tx.Create(&payExpenseBuy).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Pagos de egreso de compras", schemas.Create)
		}

		return nil
	})

	return err
}

func (r *ExpenseBuyRepository) ExpenseBuyDelete(memberID int64, expenseBuyID int64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var existingExpense models.ExpenseBuy
		if err := tx.Where("id = ?", expenseBuyID).First(&existingExpense).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Compra", schemas.Read)
		}

		var items []models.ExpenseBuyItem
		if err := tx.Where("expense_buy_id = ?", expenseBuyID).Find(&items).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Items de egreso de compras", schemas.Read)
		}

		// Revertir stock (restar las cantidades que se habían sumado)
		for _, item := range items {
			var stock models.Deposit
			if err := tx.Where("product_id = ?", item.ProductID).First(&stock).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Si no existe el registro de stock, continuar (caso poco probable)
					continue
				}
				return schemas.HandlerErrorGorm(err, "Stock", schemas.Read)
			}

			// Validar que hay suficiente stock para revertir
			if stock.Stock < item.Amount {
				return schemas.ErrorResponse(
					400,
					fmt.Sprintf("No se puede eliminar: stock insuficiente para el producto %d (disponible: %.2f, a revertir: %.2f)", item.ProductID, stock.Stock, item.Amount),
					fmt.Errorf("stock insuficiente para revertir"),
				)
			}

			// Restar stock
			if err := tx.Model(&stock).
				UpdateColumn("stock", gorm.Expr("stock - ?", item.Amount)).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Stock", schemas.Update)
			}
		}

		// Eliminar pagos
		if err := tx.Where("expense_buy_id = ?", expenseBuyID).Delete(&models.PayExpenseBuy{}).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Pagos de egreso de compras", schemas.Delete)
		}

		// Eliminar items
		if err := tx.Where("expense_buy_id = ?", expenseBuyID).Delete(&models.ExpenseBuyItem{}).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Items de egreso de compras", schemas.Delete)
		}

		// Eliminar la compra
		if err := tx.Delete(&existingExpense).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Egreso de compras", schemas.Delete)
		}

		return nil
	})

	return err
}
