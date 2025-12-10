package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"gorm.io/gorm"
)

func (r *MovementStockRepository) MovementStockGetByID(id int64) (*models.MovementStock, error) {
	var movement *models.MovementStock
	if err := r.DB.Preload("Member", func(db *gorm.DB) *gorm.DB { return db.Unscoped()}).Preload("Product").Preload("Product.Category").First(&movement, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "movimiento no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "error al obtener el movimiento", err)
	}
	return movement, nil
}

func (r *MovementStockRepository) MovementStockGetByDate(page, limit int, fromDate, toDate time.Time) ([]*models.MovementStock, int64, error) {
	offset := (page - 1) * limit

	var movements []*models.MovementStock
	var total int64
	if err := r.DB.
		Preload("Member", func(db *gorm.DB) *gorm.DB { return db.Unscoped()}).
		Preload("Product").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Find(&movements).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "error al obtener movimientos", err)
	}

	if err := r.DB.Model(&models.MovementStock{}).Count(&total).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "error al contar movimientos", err)
	}

	return movements, total, nil
}

func (r *MovementStockRepository) MoveStockList(userID int64, input []*schemas.MovementStockList) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Validar que hay elementos para procesar
		if len(input) == 0 {
			return schemas.ErrorResponse(400, "no hay movimientos para procesar", fmt.Errorf("lista vacía"))
		}

		// Procesar cada producto
		for _, movementList := range input {
			// Validar producto
			var product models.Product
			if err := tx.First(&product, movementList.ProductID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return schemas.ErrorResponse(404, fmt.Sprintf("producto %d no encontrado", movementList.ProductID), err)
				}
				return schemas.ErrorResponse(500, fmt.Sprintf("error al obtener el producto %d", movementList.ProductID), err)
			}

			if product.Price <= 0.0 {
				return schemas.ErrorResponse(400, fmt.Sprintf("no se puede editar el producto %d sin precio", movementList.ProductID), 
					fmt.Errorf("no se puede editar un producto sin precio"))
			}

			// Validar que hay movimientos para el producto
			if len(movementList.MovementStockItem) == 0 {
				return schemas.ErrorResponse(400, fmt.Sprintf("no hay movimientos para el producto %d", movementList.ProductID), 
					fmt.Errorf("lista de movimientos vacía"))
			}

			// Procesar cada movimiento del producto
			for idx, item := range movementList.MovementStockItem {
				if err := r.processSingleMovement(tx, userID, movementList.ProductID, &item, idx); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *MovementStockRepository) processSingleMovement(tx *gorm.DB, userID int64, productID int64, item *schemas.MovementStockItem, index int) error {
	var fromID, toID int64

	// ===== PROCESAR ORIGEN =====
	switch item.FromType {
	case "deposit":
		fromID = 100
		
		// Asegurar que existe el registro
		if err := tx.Where("product_id = ?", productID).
			FirstOrCreate(&models.Deposit{}, &models.Deposit{
				ProductID: productID,
				Stock:     0,
			}).Error; err != nil {
			return schemas.ErrorResponse(500, fmt.Sprintf("error al inicializar el depósito (movimiento %d)", index+1), err)
		}

		// Validar stock si es necesario
		if !*item.IgnoreStock {
			var currentStock float64
			if err := tx.Model(&models.Deposit{}).
				Where("product_id = ?", productID).
				Select("stock").
				Scan(&currentStock).Error; err != nil {
				return schemas.ErrorResponse(500, fmt.Sprintf("error al verificar stock (movimiento %d)", index+1), err)
			}
			
			if currentStock < item.Amount {
				return schemas.ErrorResponse(400, fmt.Sprintf("no hay suficiente stock en depósito para transferir (movimiento %d)", index+1), 
					fmt.Errorf("stock actual: %.2f, necesario: %.2f", currentStock, item.Amount))
			}
		}

		// Actualización atómica
		result := tx.Model(&models.Deposit{}).
			Where("product_id = ?", productID).
			UpdateColumn("stock", gorm.Expr("stock - ?", item.Amount))
		
		if result.Error != nil {
			return schemas.ErrorResponse(500, fmt.Sprintf("error al actualizar stock del depósito (movimiento %d)", index+1), result.Error)
		}
		
		if result.RowsAffected == 0 {
			return schemas.ErrorResponse(404, fmt.Sprintf("no se pudo actualizar el depósito (movimiento %d)", index+1), 
				fmt.Errorf("registro no encontrado"))
		}

	case "point_sale":
		// Validar que existe el punto de venta
		var pointSale models.PointSale
		if err := tx.
			Select("id", "is_deposit").
			Where("id = ?", item.FromID).
			First(&pointSale).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return schemas.ErrorResponse(404, fmt.Sprintf("punto de venta %d no encontrado (movimiento %d)", item.FromID, index+1), err)
				}
			return schemas.ErrorResponse(500, fmt.Sprintf("error al obtener el punto de venta origen (movimiento %d)", index+1), err)
		}

		if pointSale.IsDeposit {
			return schemas.ErrorResponse(400, fmt.Sprintf("no se puede transferir stock desde un punto de venta deposito (movimiento %d), transferir desde otro punto de venta o depósito, el punto de venta es deposito", index+1), 
				fmt.Errorf("no se puede transferir stock desde un punto de venta deposito"))
		}
		
		fromID = item.FromID

		// Asegurar que existe el registro
		if err := tx.Where("product_id = ? AND point_sale_id = ?", productID, item.FromID).
			FirstOrCreate(&models.StockPointSale{}, &models.StockPointSale{
				ProductID:   productID,
				PointSaleID: item.FromID,
				Stock:       0,
			}).Error; err != nil {
			return schemas.ErrorResponse(500, fmt.Sprintf("error al inicializar stock del punto de venta origen (movimiento %d)", index+1), err)
		}

		// Validar stock si es necesario
		if !*item.IgnoreStock {
			var currentStock float64
			if err := tx.Model(&models.StockPointSale{}).
				Where("product_id = ? AND point_sale_id = ?", productID, item.FromID).
				Select("stock").
				Scan(&currentStock).Error; err != nil {
				return schemas.ErrorResponse(500, fmt.Sprintf("error al verificar stock (movimiento %d)", index+1), err)
			}
			
			if currentStock < item.Amount {
				return schemas.ErrorResponse(400, fmt.Sprintf("no hay suficiente stock en punto de venta para transferir (movimiento %d)", index+1), 
					fmt.Errorf("stock actual: %.2f, necesario: %.2f", currentStock, item.Amount))
			}
		}

		// Actualización atómica
		result := tx.Model(&models.StockPointSale{}).
			Where("product_id = ? AND point_sale_id = ?", productID, item.FromID).
			UpdateColumn("stock", gorm.Expr("stock - ?", item.Amount))
		
		if result.Error != nil {
			return schemas.ErrorResponse(500, fmt.Sprintf("error al actualizar stock del punto de venta origen (movimiento %d)", index+1), result.Error)
		}
		
		if result.RowsAffected == 0 {
			return schemas.ErrorResponse(404, fmt.Sprintf("no se pudo actualizar el punto de venta origen (movimiento %d)", index+1), 
				fmt.Errorf("registro no encontrado"))
		}

	default:
		return schemas.ErrorResponse(400, fmt.Sprintf("tipo de origen inválido (movimiento %d)", index+1), 
			fmt.Errorf("tipo de origen inválido: %s", item.FromType))
	}

	// ===== PROCESAR DESTINO =====
	switch item.ToType {
	case "deposit":
		toID = 100
		
		// Asegurar que existe el registro
		if err := tx.Where("product_id = ?", productID).
			FirstOrCreate(&models.Deposit{}, &models.Deposit{
				ProductID: productID,
				Stock:     0,
			}).Error; err != nil {
			return schemas.ErrorResponse(500, fmt.Sprintf("error al inicializar el depósito destino (movimiento %d)", index+1), err)
		}

		// Actualización atómica
		result := tx.Model(&models.Deposit{}).
			Where("product_id = ?", productID).
			UpdateColumn("stock", gorm.Expr("stock + ?", item.Amount))
		
		if result.Error != nil {
			return schemas.ErrorResponse(500, fmt.Sprintf("error al actualizar stock del depósito destino (movimiento %d)", index+1), result.Error)
		}

	case "point_sale":
		// Validar que existe el punto de venta
		var pointSale models.PointSale
		if err := tx.
			Select("id", "is_deposit").
			Where("id = ?", item.ToID).
			First(&pointSale).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return schemas.ErrorResponse(404, fmt.Sprintf("punto de venta %d no encontrado (movimiento %d)", item.FromID, index+1), err)
				}
			return schemas.ErrorResponse(500, fmt.Sprintf("error al obtener el punto de venta origen (movimiento %d)", index+1), err)
		}

		if pointSale.IsDeposit {
			return schemas.ErrorResponse(400, fmt.Sprintf("no se puede transferir stock desde un punto de venta deposito (movimiento %d), transferir desde otro punto de venta o depósito, el punto de venta es deposito", index+1), 
				fmt.Errorf("no se puede transferir stock desde un punto de venta deposito"))
		}

		toID = item.ToID

		// Asegurar que existe el registro
		if err := tx.Where("product_id = ? AND point_sale_id = ?", productID, item.ToID).
			FirstOrCreate(&models.StockPointSale{}, &models.StockPointSale{
				ProductID:   productID,
				PointSaleID: item.ToID,
				Stock:       0,
			}).Error; err != nil {
			return schemas.ErrorResponse(500, fmt.Sprintf("error al inicializar stock del punto de venta destino (movimiento %d)", index+1), err)
		}

		// Actualización atómica
		result := tx.Model(&models.StockPointSale{}).
			Where("product_id = ? AND point_sale_id = ?", productID, item.ToID).
			UpdateColumn("stock", gorm.Expr("stock + ?", item.Amount))
		
		if result.Error != nil {
			return schemas.ErrorResponse(500, fmt.Sprintf("error al actualizar stock del punto de venta destino (movimiento %d)", index+1), result.Error)
		}

	default:
		return schemas.ErrorResponse(400, fmt.Sprintf("tipo de destino inválido (movimiento %d)", index+1), 
			fmt.Errorf("tipo de destino inválido: %s", item.ToType))
	}

	// Registrar el movimiento
	movementStock := models.MovementStock{
		MemberID:      userID,
		ProductID:   productID,
		Amount:      item.Amount,
		FromID:      fromID,
		FromType:    item.FromType,
		ToID:        toID,
		ToType:      item.ToType,
		IgnoreStock: *item.IgnoreStock,
	}

	if err := tx.Create(&movementStock).Error; err != nil {
		return schemas.ErrorResponse(500, fmt.Sprintf("error al registrar el movimiento %d", index+1), err)
	}

	return nil
}