package services

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

func (p *PointSaleService) PointSaleCreate(pointSaleCreate *schemas.PointSaleCreate) (int64, error) {
	return p.PointSaleRepository.PointSaleCreate(pointSaleCreate)
}

func (p *PointSaleService) PointSaleGetAllByMember(memberID int64) ([]schemas.PointSaleResponse, error) {
	return p.PointSaleRepository.PointSaleGetAllByMember(memberID)
}

func (p *PointSaleService) PointSaleGetAll() ([]schemas.PointSaleResponse, error) {
	return p.PointSaleRepository.PointSaleGetAll()
}

func (p *PointSaleService) PointSaleUpdate(pointSaleUpdate *schemas.PointSaleUpdate) (error) {
	return p.PointSaleRepository.PointSaleUpdate(pointSaleUpdate)
}