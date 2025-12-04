package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

type PointSaleRepository interface {
	PointSaleCreate(pointSaleCreate *schemas.PointSaleCreate) (int64, error)
	PointSaleUpdate(pointSaleUpdate *schemas.PointSaleUpdate) error
	PointSaleUpdateMain(pointSaleUpdateMain *schemas.PointSaleUpdateMain) error
	PointSaleGetAllByMember(memberID int64) ([]schemas.PointSaleResponse, error)
	PointSaleGetAll() ([]schemas.PointSaleResponse, error)
	PointSaleCount() (int64, error)
}

type PointSaleService interface {
	PointSaleCreate(pointSaleCreate *schemas.PointSaleCreate, plan *schemas.PlanResponseDTO) (int64, error)
	PointSaleUpdate(pointSaleUpdate *schemas.PointSaleUpdate) error
	PointSaleUpdateMain(pointSaleUpdateMain *schemas.PointSaleUpdateMain) error

	PointSaleGetAllByMember(memberID int64) ([]schemas.PointSaleResponse, error)
	PointSaleGetAll() ([]schemas.PointSaleResponse, error)
}
