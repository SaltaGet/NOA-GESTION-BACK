package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"



type PointSaleService interface {
	PointSaleCreate(pointSaleCreate *schemas.PointSaleCreate) (int64, error)
	
	PointSaleGetAllByMember(memberID int64) ([]schemas.PointSaleResponse, error)
	PointSaleGetAll() ([]schemas.PointSaleResponse, error)
}

type PointSaleRepository interface {
	PointSaleCreate(pointSaleCreate *schemas.PointSaleCreate) (int64, error)
	PointSaleGetAllByMember(memberID int64) ([]schemas.PointSaleResponse, error)
	PointSaleGetAll() ([]schemas.PointSaleResponse, error)
}