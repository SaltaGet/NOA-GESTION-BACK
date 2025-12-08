package repositories

import (
	"errors"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func (p *PointSaleRepository) PointSaleCreate(pointSaleCreate *schemas.PointSaleCreate) (int64, error) {
	var pointSaleGet []models.PointSale
	if err := p.DB.Where("is_main = ?", true).Find(&pointSaleGet).Error; err != nil {
		return 0, schemas.ErrorResponse(500, "Error al obtener el punto de venta principal", err)
	}

	pointSale := models.PointSale{
		Name:        pointSaleCreate.Name,
		Description: pointSaleCreate.Description,
		IsDeposit:   *pointSaleCreate.IsDeposit,
	}

	if len(pointSaleGet) == 0 {
		pointSale.IsMain = true
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
		Select("point_sales.id", "point_sales.name", "point_sales.description", "point_sales.is_deposit", "point_sales.is_main").
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
	err := p.DB.Model(&models.PointSale{}).Select("id", "name", "description", "is_deposit", "is_main").Scan(&pointSales).Error
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al obtener los puntos de venta", err)
	}

	return pointSales, nil
}

func (p *PointSaleRepository) PointSaleGetByID(id int64) (*schemas.PointSaleResponse, error) {
	var pointSales models.PointSale
	err := p.DB.Select("id", "name", "description", "is_deposit", "is_main").First(&pointSales, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Punto de venta no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error al obtener los puntos de venta", err)
	}

	var pointSaleResponse schemas.PointSaleResponse
	copier.Copy(&pointSaleResponse, &pointSales)

	return &pointSaleResponse, nil
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

func (p *PointSaleRepository) PointSaleCount() (int64, error) {
	var pointSales int64
	if err := p.DB.Model(&models.PointSale{}).Count(&pointSales).Error; err != nil {
		return 0, schemas.ErrorResponse(500, "error al obtner la cantidad de puntos de ventas", err)
	}

	return pointSales, nil
}

func (p *PointSaleRepository) PointSaleUpdateMain(pointSaleUpdateMain *schemas.PointSaleUpdateMain) error {
	return p.DB.Transaction(func(tx *gorm.DB) error {

		var pointSaleOld models.PointSale
		if err := tx.First(&pointSaleOld, pointSaleUpdateMain.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Punto de venta no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el punto de venta", err)
		}

		var pointSaleNew models.PointSale
		if err := tx.First(&pointSaleNew, pointSaleUpdateMain.NewMain).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Punto de venta no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el punto de venta", err)
		}

		if !pointSaleOld.IsMain {
			return schemas.ErrorResponse(400, "El punto de venta indicado no es el principal actual", nil)
		}
		if pointSaleNew.IsMain {
			return schemas.ErrorResponse(400, "El nuevo punto de venta ya es el principal", nil)
		}

		// Actualizar usando Updates (más seguro)
		if err := tx.Model(&pointSaleOld).Update("is_main", false).Error; err != nil {
			return schemas.ErrorResponse(500, "Error actualizando el punto principal", err)
		}

		if err := tx.Model(&pointSaleNew).Update("is_main", true).Error; err != nil {
			return schemas.ErrorResponse(500, "Error actualizando el nuevo punto principal", err)
		}

		return nil
	})
}
