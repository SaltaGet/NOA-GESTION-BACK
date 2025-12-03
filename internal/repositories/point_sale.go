package repositories

import (
	"errors"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"gorm.io/gorm"
)

func (p *PointSaleRepository) PointSaleCreate(pointSaleCreate *schemas.PointSaleCreate) (int64, error) {
	pointSale := models.PointSale{
		Name:        pointSaleCreate.Name,
		Description: pointSaleCreate.Description,
		IsDeposit:   *pointSaleCreate.IsDeposit,
	}

	err := p.DB.Create(&pointSale).Error
	if err != nil {
		if schemas.IsDuplicateError(err) {
			return 0, schemas.ErrorResponse(409, "El punto de venta "+pointSale.Name+" ya existe", err)
		}
		return 0, schemas.ErrorResponse(500, "Error al crear punto de venta", err)
	}

	var membersAdmin []models.Member
	if err := p.DB.Where("is_admin = ?", true).Find(&membersAdmin).Error; err != nil {
		return 0, schemas.ErrorResponse(500, "Error al obtener los administradores", err)
	}

	if len(membersAdmin) > 0 {
		if err := p.DB.Model(&pointSale).Association("Members").Append(&membersAdmin); err != nil {
			return 0, schemas.ErrorResponse(500, "Error al asignar punto de venta a administradores", err)
		}
	}

	return pointSale.ID, nil
}

func (p *PointSaleRepository) PointSaleGetAllByMember(memberID int64) ([]schemas.PointSaleResponse, error) {
	var pointSales []schemas.PointSaleResponse
	err := p.DB.
		Model(&models.PointSale{}).
		Select("point_sales.id", "point_sales.name", "point_sales.description", "point_sales.is_deposit").
		Joins("JOIN member_point_sales mp ON mp.point_sale_id = point_sales.id").
		Where("mp.member_id = ?", memberID).
		Scan(&pointSales).Error
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al obtener los puntos de venta", err)
	}

	return pointSales, nil
}

func (p *PointSaleRepository) PointSaleGetAll() ([]schemas.PointSaleResponse, error) {
	var pointSales []schemas.PointSaleResponse
	err := p.DB.Model(&models.PointSale{}).Select("id", "name", "description", "is_deposit").Scan(&pointSales).Error
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al obtener los puntos de venta", err)
	}

	return pointSales, nil
}

func (p *PointSaleRepository) PointSaleUpdate(pointSaleUpdate *schemas.PointSaleUpdate) error {
	return p.DB.Transaction(func(tx *gorm.DB) error {
		var pointSale models.PointSale
		if err := tx.First(&pointSale, pointSaleUpdate.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Punto de venta no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el punto de venta", err)
		}

		pointSale.Name = pointSaleUpdate.Name
		pointSale.Description = pointSaleUpdate.Description

		if !pointSale.IsDeposit && *pointSaleUpdate.IsDeposit {
			var stockList []models.StockPointSale
			if err := tx.Where("point_sale_id = ?", pointSale.ID).Find(&stockList).Error; err != nil {
				return schemas.ErrorResponse(500, "Error obteniendo el stock del punto de venta", err)
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
						return schemas.ErrorResponse(500, "Error creando registro en depósito", err)
					}
				} else if err == nil {
					deposit.Stock += s.Stock
					if err := tx.Save(&deposit).Error; err != nil {
						return schemas.ErrorResponse(500, "Error actualizando stock en depósito", err)
					}
				} else {
					return schemas.ErrorResponse(500, "Error validando stock en depósito", err)
				}

				s.Stock = 0
				if err := tx.Save(&s).Error; err != nil {
					return schemas.ErrorResponse(500, "Error limpiando stock de punto de venta", err)
				}
			}
		}

		pointSale.IsDeposit = *pointSaleUpdate.IsDeposit

		if err := tx.Save(&pointSale).Error; err != nil {
			return schemas.ErrorResponse(500, "Error actualizando punto de venta", err)
		}

		return nil
	})
}

