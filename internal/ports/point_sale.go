package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

type PointSaleRepository interface {
	PointSaleCreate(memberID int64, pointSaleCreate *schemas.PointSaleCreate) (int64, error)
	PointSaleUpdate(memberID int64, pointSaleUpdate *schemas.PointSaleUpdate) error
	PointSaleUpdateMain(memberID int64, pointSaleUpdateMain *schemas.PointSaleUpdateMain) error
	PointSaleGetAllByMember(memberID int64) ([]schemas.PointSaleResponse, error)
	PointSaleGetAll() ([]schemas.PointSaleResponse, error)
	PointSaleGetByID(id int64) (*schemas.PointSaleResponse, error)
	PointSaleCount() (int64, error)
}

type PointSaleService interface {
	PointSaleCreate(memberID int64, pointSaleCreate *schemas.PointSaleCreate, plan *schemas.PlanResponseDTO) (int64, error)
	PointSaleUpdate(memberID int64, pointSaleUpdate *schemas.PointSaleUpdate) error
	PointSaleUpdateMain(memberID int64, pointSaleUpdateMain *schemas.PointSaleUpdateMain) error

	PointSaleGetAllByMember(memberID int64) ([]schemas.PointSaleResponse, error)
	PointSaleGetAll() ([]schemas.PointSaleResponse, error)
	PointSaleGetByID(id int64) (*schemas.PointSaleResponse, error)
}
