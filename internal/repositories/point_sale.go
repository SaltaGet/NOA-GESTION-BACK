package repositories

import (
	"errors"
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (p *PointSaleRepository) PointSaleGetAllByMember(memberID int64) ([]schemas.PointSaleResponse, error) {
	var pointSales []schemas.PointSaleResponse
	err := p.DB.
		Model(&models.PointSale{}).
		Select("point_sales.id", "point_sales.name", "point_sales.description", "point_sales.is_deposit", "point_sales.is_main", "point_sales.number").
		Joins("JOIN member_point_sales mp ON mp.point_sale_id = point_sales.id").
		Where("mp.member_id = ?", memberID).
		Scan(&pointSales).Error
	if err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Read)
	}

	return pointSales, nil
}

func (p *PointSaleRepository) PointSaleGetAll() ([]schemas.PointSaleResponse, error) {
	var pointSales []schemas.PointSaleResponse
	err := p.DB.Model(&models.PointSale{}).Select("id", "name", "description", "is_deposit", "is_main", "number").Scan(&pointSales).Error
	if err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Read)
	}

	return pointSales, nil
}

func (p *PointSaleRepository) PointSaleGetByID(id int64) (*schemas.PointSaleResponse, error) {
	var pointSales models.PointSale
	err := p.DB.Select("id", "name", "description", "is_deposit", "is_main", "number").First(&pointSales, id).Error
	if err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Read)
	}

	var pointSaleResponse schemas.PointSaleResponse
	copier.Copy(&pointSaleResponse, &pointSales)

	return &pointSaleResponse, nil
}

func (p *PointSaleRepository) PointSaleCount() (int64, error) {
	var pointSales int64
	if err := p.DB.Model(&models.PointSale{}).Count(&pointSales).Error; err != nil {
		return 0, schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Read)
	}

	return pointSales, nil
}

func (p *PointSaleRepository) PointSaleCreate(memberID int64, pointSaleCreate *schemas.PointSaleCreate) (int64, error) {
	var pointSaleSave models.PointSale

	err := p.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var pointSaleGet []models.PointSale
		if err := tx.Where("is_main = ?", true).Find(&pointSaleGet).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Read)
		}

		pointSale := models.PointSale{
			Name:        pointSaleCreate.Name,
			Description: pointSaleCreate.Description,
			IsDeposit:   *pointSaleCreate.IsDeposit,
			Number:      pointSaleCreate.Number,
		}

		if len(pointSaleGet) == 0 {
			pointSale.IsMain = true
		}

		if err := tx.Create(&pointSale).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Create)
		}

		pointSaleSave = pointSale

		var membersAdmin []models.Member
		if err := tx.Where("is_admin = ?", true).Find(&membersAdmin).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Read)
		}

		if len(membersAdmin) > 0 {
			if err := tx.Model(&pointSale).Association("Members").Append(&membersAdmin); err != nil {
				return schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Update)
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return pointSaleSave.ID, nil
}

// PointSaleUpdate actualiza un punto de venta con auditoría
func (p *PointSaleRepository) PointSaleUpdate(memberID int64, pointSaleUpdate *schemas.PointSaleUpdate) error {
	err := p.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var pointSale models.PointSale
		if err := tx.First(&pointSale, pointSaleUpdate.ID).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Read)
		}

		pointSale.Name = pointSaleUpdate.Name
		pointSale.Description = pointSaleUpdate.Description
		pointSale.Number = pointSaleUpdate.Number

		// Si se convierte a depósito, mover el stock
		if !pointSale.IsDeposit && *pointSaleUpdate.IsDeposit {
			var stockList []models.StockPointSale
			if err := tx.Where("point_sale_id = ?", pointSale.ID).Find(&stockList).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Read)
			}

			for _, s := range stockList {

				var deposit models.Deposit
				err := tx.Where("product_id = ?", s.ProductID).First(&deposit).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					deposit = models.Deposit{
						ProductID: s.ProductID,
						Stock:     s.Stock,
					}
					if err := tx.Create(&deposit).Error; err != nil {
						return schemas.HandlerErrorGorm(err, "Déposito", schemas.Create)
					}

				} else if err == nil {
					deposit.Stock += s.Stock
					if err := tx.Save(&deposit).Error; err != nil {
						return schemas.HandlerErrorGorm(err, "Déposito", schemas.Update)
					}

				} else {
					return schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Read)
				}

				s.Stock = 0
				if err := tx.Save(&s).Error; err != nil {
					return schemas.HandlerErrorGorm(err, "Stock punto de venta", schemas.Update)
				}
			}
		}

		pointSale.IsDeposit = *pointSaleUpdate.IsDeposit

		if err := tx.Save(&pointSale).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Update)
		}

		return nil
	})

	return err
}

// PointSaleUpdateMain actualiza el punto de venta principal con auditoría
func (p *PointSaleRepository) PointSaleUpdateMain(memberID int64, pointSaleUpdateMain *schemas.PointSaleUpdateMain) error {
	var savePointSaleOld, savePointSaleNew models.PointSale

	err := p.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var pointSaleOld models.PointSale
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&pointSaleOld, pointSaleUpdateMain.ID).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Read)
		}

		savePointSaleOld = pointSaleOld

		var pointSaleNew models.PointSale
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&pointSaleNew, pointSaleUpdateMain.NewMain).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Read)
		}

		if !pointSaleOld.IsMain {
			return schemas.ErrorResponse(400, "El punto de venta indicado no es el principal actual", nil)
		}

		if pointSaleNew.IsMain {
			return schemas.ErrorResponse(400, "El nuevo punto de venta ya es el principal", nil)
		}

		// Actualizar
		if err := tx.Model(&pointSaleOld).Update("is_main", false).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Update)
		}

		if err := tx.Model(&pointSaleNew).Update("is_main", true).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Update)
		}

		// Recargar actualizados
		if err := tx.First(&savePointSaleOld, pointSaleOld.ID).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Read)
		}

		if err := tx.First(&savePointSaleNew, pointSaleNew.ID).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Punto de venta", schemas.Read)
		}

		return nil
	})

	return err
}

