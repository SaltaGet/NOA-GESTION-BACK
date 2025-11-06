package repositories

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)



func (p *PointSaleRepository) PointSaleCreate(pointSaleCreate *schemas.PointSaleCreate) (int64, error) {
	pointSale := models.PointSale{
		Name:        pointSaleCreate.Name,
		Description: pointSaleCreate.Description,
		IsDeposit:   pointSaleCreate.IsDeposit,
	}

	err := p.DB.Create(&pointSale).Error
	if err != nil {
		if schemas.IsDuplicateError(err) {
			return 0, schemas.ErrorResponse(409, "El punto de venta " + pointSale.Name + " ya existe", err)
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

