package services

import (
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (p *PointSaleService) PointSaleCreate(memberID int64, pointSaleCreate *schemas.PointSaleCreate, plan *schemas.PlanResponseDTO) (int64, error) {
	amountPointSales, err := p.PointSaleRepository.PointSaleCount()
	if err != nil {
		return 0, err
	}

	if amountPointSales >= plan.AmountPointSale {
		return 0, schemas.ErrorResponse(400, "El plan no permite agregar mas puntos de venta", fmt.Errorf("el plan no permite agregar mas puntos de venta"))
	}

	return p.PointSaleRepository.PointSaleCreate(memberID, pointSaleCreate)
}

func (p *PointSaleService) PointSaleGetAllByMember(memberID int64) ([]schemas.PointSaleResponse, error) {
	return p.PointSaleRepository.PointSaleGetAllByMember(memberID)
}

func (p *PointSaleService) PointSaleGetAll() ([]schemas.PointSaleResponse, error) {
	return p.PointSaleRepository.PointSaleGetAll()
}

func (p *PointSaleService) PointSaleGetByID(id int64) (*schemas.PointSaleResponse, error) {
	return p.PointSaleRepository.PointSaleGetByID(id)
}

func (p *PointSaleService) PointSaleUpdate(memberID int64, pointSaleUpdate *schemas.PointSaleUpdate) (error) {
	return p.PointSaleRepository.PointSaleUpdate(memberID, pointSaleUpdate)
}

func (p *PointSaleService) PointSaleUpdateMain(memberID int64, pointSaleUpdate *schemas.PointSaleUpdateMain) (error) {
	return p.PointSaleRepository.PointSaleUpdateMain(memberID, pointSaleUpdate)
}