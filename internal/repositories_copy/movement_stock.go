package repositories_copy

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
	if err := r.DB.Preload("Member", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).Preload("Product").Preload("Product.Category").First(&movement, id).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Movimiento de stock", schemas.Read)
	}
	return movement, nil
}

func (r *MovementStockRepository) MovementStockGetByDate(page, limit int, fromDate, toDate time.Time) ([]models.MovementStock, int64, error) {
	offset := (page - 1) * limit

	var movements []models.MovementStock
	var total int64
	if err := r.DB.
		Preload("Member", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Preload("Product").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Find(&movements).Error; err != nil {
		return nil, 0, schemas.HandlerErrorGorm(err, "Movimiento de stock", schemas.Read)
	}

	if err := r.DB.Model(&models.MovementStock{}).Count(&total).Error; err != nil {
		return nil, 0, schemas.HandlerErrorGorm(err, "Movimiento de stock", schemas.Read)
	}

	return movements, total, nil
}

func (r *MovementStockRepository) MoveStockList(memberID int64, input []schemas.MovementStockList) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}
		// Validar que hay elementos para procesar
		if len(input) == 0 {
			return schemas.ErrorResponse(400, "no hay movimientos para procesar", fmt.Errorf("lista vacía"))
		}

		// Procesar cada producto
		for _, movementList := range input {
			// Validar producto
			var product models.Product
			if err := tx.First(&product, movementList.ProductID).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Producto", schemas.Read)
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
				var errorN error
				errorN = r.processSingleMovement(tx, memberID, movementList.ProductID, &item, idx)
				if errorN != nil { 
					return errorN
				}
			}
		}

		return nil
	})

	return err
}

func (r *MovementStockRepository) processSingleMovement(
	tx *gorm.DB,
	memberID int64,
	productID int64,
	item *schemas.MovementStockItem,
	index int,
) (error) {

	var fromID, toID int64

	ignoreStock := item.IgnoreStock != nil && *item.IgnoreStock

	// ---------------------------
	//   PROCESAR ORIGEN
	// ---------------------------
	switch item.FromType {

	case "deposit":
		fromID = 100

		// asegurar existencia
		if err := tx.Where("product_id = ?", productID).
			FirstOrCreate(&models.Deposit{ProductID: productID}).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Deposito", schemas.Create)
		}

		if !ignoreStock {
			var current float64
			if err := tx.Model(&models.Deposit{}).
				Where("product_id = ?", productID).
				Select("stock").
				Scan(&current).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Deposito", schemas.Read)
			}

			if current < item.Amount {
				return schemas.ErrorResponse(400,
					fmt.Sprintf("stock insuficiente en depósito (%d)", index+1),
					fmt.Errorf("actual %.2f < requerido %.2f", current, item.Amount))
			}
		}

		result := tx.Model(&models.Deposit{}).
			Where("product_id = ?", productID).
			UpdateColumn("stock", gorm.Expr("stock - ?", item.Amount))

		if result.Error != nil {
			return schemas.HandlerErrorGorm(result.Error, "Deposito", schemas.Update)
		}

		if result.RowsAffected == 0 {
			return schemas.ErrorResponse(404,
				fmt.Sprintf("depósito origen no encontrado (%d)", index+1),
				fmt.Errorf("RowsAffected=0"))
		}

	case "point_sale":

		var ps models.PointSale
		if err := tx.Select("id", "is_deposit").
			Where("id = ?", item.FromID).
			First(&ps).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404,
					fmt.Sprintf("punto de venta %d no encontrado (%d)", item.FromID, index+1), err)
			}
			return schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Read)
		}

		if ps.IsDeposit {
			return schemas.ErrorResponse(400,
				fmt.Sprintf("no se puede usar un punto de venta depósito como origen (%d)", index+1),
				fmt.Errorf("point_sale es depósito"))
		}

		fromID = item.FromID

		var oldPS models.StockPointSale
		tx.Where("product_id = ? AND point_sale_id = ?", productID, item.FromID).
			First(&oldPS)

		if err := tx.Where("product_id = ? AND point_sale_id = ?", productID, item.FromID).
			FirstOrCreate(&models.StockPointSale{
				ProductID:   productID,
				PointSaleID: item.FromID,
			}).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Stock Punto de venta", schemas.Create)
		}

		if !ignoreStock {
			var current float64
			if err := tx.Model(&models.StockPointSale{}).
				Where("product_id = ? AND point_sale_id = ?", productID, item.FromID).
				Select("stock").
				Scan(&current).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Stock Punto de venta", schemas.Read)
			}

			if current < item.Amount {
				return schemas.ErrorResponse(400,
					fmt.Sprintf("stock insuficiente en point_sale origen (%d)", index+1),
					fmt.Errorf("actual %.2f < requerido %.2f", current, item.Amount))
			}
		}

		result := tx.Model(&models.StockPointSale{}).
			Where("product_id = ? AND point_sale_id = ?", productID, item.FromID).
			UpdateColumn("stock", gorm.Expr("stock - ?", item.Amount))

		if result.Error != nil {
			return schemas.HandlerErrorGorm(result.Error, "Stock Punto de venta", schemas.Update)
		}

		if result.RowsAffected == 0 {
			return schemas.ErrorResponse(404,
				fmt.Sprintf("stock origen point_sale no encontrado (%d)", index+1),
				fmt.Errorf("RowsAffected=0"))
		}

	default:
		return schemas.ErrorResponse(400,
			fmt.Sprintf("tipo de origen inválido (%d)", index+1),
			fmt.Errorf("FromType='%s'"))
	}

	// ---------------------------
	//   PROCESAR DESTINO
	// ---------------------------

	switch item.ToType {

	case "deposit":
		toID = 100

		var oldDeposit models.Deposit
		tx.Where("product_id = ?", productID).First(&oldDeposit)

		if err := tx.Where("product_id = ?", productID).
			FirstOrCreate(&models.Deposit{ProductID: productID}).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Deposito", schemas.Create)
		}

		result := tx.Model(&models.Deposit{}).
			Where("product_id = ?", productID).
			UpdateColumn("stock", gorm.Expr("stock + ?", item.Amount))

		if result.Error != nil {
			return schemas.HandlerErrorGorm(result.Error, "Deposito", schemas.Update)
		}

	case "point_sale":

		var ps models.PointSale
		if err := tx.Select("id", "is_deposit").
			Where("id = ?", item.ToID).
			First(&ps).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404,
					fmt.Sprintf("punto de venta %d no encontrado (%d)", item.ToID, index+1), err)
			}
			return schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Read)
		}

		if ps.IsDeposit {
			return schemas.ErrorResponse(400,
				fmt.Sprintf("no se puede usar un punto de venta depósito como destino (%d)", index+1),
				fmt.Errorf("point_sale es depósito"))
		}

		toID = item.ToID

		var oldPS models.StockPointSale
		tx.Where("product_id = ? AND point_sale_id = ?", productID, item.ToID).
			First(&oldPS)

		if err := tx.Where("product_id = ? AND point_sale_id = ?", productID, item.ToID).
			FirstOrCreate(&models.StockPointSale{
				ProductID:   productID,
				PointSaleID: item.ToID,
			}).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Stock Punto de venta", schemas.Create)
		}

		result := tx.Model(&models.StockPointSale{}).
			Where("product_id = ? AND point_sale_id = ?", productID, item.ToID).
			UpdateColumn("stock", gorm.Expr("stock + ?", item.Amount))

		if result.Error != nil {
			return schemas.HandlerErrorGorm(result.Error, "Stock Punto de venta", schemas.Update)
		}

	default:
		return schemas.ErrorResponse(400,
			fmt.Sprintf("tipo de destino inválido (%d)", index+1),
			fmt.Errorf("ToType='%s'"))
	}

	// ---------------------------
	//   CREAR MOVIMIENTO
	// ---------------------------
	movement := models.MovementStock{
		MemberID:    memberID,
		ProductID:   productID,
		Amount:      item.Amount,
		FromID:      fromID,
		FromType:    item.FromType,
		ToID:        toID,
		ToType:      item.ToType,
		IgnoreStock: ignoreStock,
	}

	if err := tx.Create(&movement).Error; err != nil {
		return schemas.HandlerErrorGorm(err, "Movimiento", schemas.Create)
	}

	return nil
}

