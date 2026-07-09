package ports

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)


type IncomeSaleRepository interface {
	IncomeSaleGetByID(pointSaleID, id int64) (*schemas.IncomeSaleResponse, error)
	IncomeSaleGetByDate(pointSaleID int64, fromDate, toDate time.Time, page, limit int, isBudget *bool, delivered *bool) ([]*schemas.IncomeSaleResponseDTO, int64, error)
	IncomeSaleCreate(memberID, pointSaleID int64, incomeSaleCreate *schemas.IncomeSaleCreate) (int64, error)
	IncomeSaleUpdate(memberID, pointSaleID int64, incomeSaleUpdate *schemas.IncomeSaleUpdate, maintainPrice bool) error
	IncomeSaleDelete(memberID int64, id, pointSaleID int64) error
}

type IncomeSaleService interface {
	IncomeSaleGetByID(pointSaleID, id int64) (*schemas.IncomeSaleResponse, error)
	IncomeSaleGetByDate(pointSaleID int64, fromDate, toDate time.Time, page, limit int, isBudget *bool, delivered *bool) ([]*schemas.IncomeSaleResponseDTO, int64, error)
	IncomeSaleCreate(memberID, pointSaleID int64, incomeSaleCreate *schemas.IncomeSaleCreate) (int64, error)
	IncomeSaleUpdate(memberID, pointSaleID int64, incomeSaleUpdate *schemas.IncomeSaleUpdate, maintainPrice bool) error
	IncomeSaleDelete(memberID int64, id, pointSaleID int64) error
}
