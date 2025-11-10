package ports

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)


type IncomeOtherRepository interface {
	IncomeOtherGetByID(id int64) (income *schemas.IncomeOtherResponse, err error)
	IncomeOtherGetByDate(pointSaleID *int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.IncomeOtherResponse, int64, error)
	IncomeOtherCreate(memberID, pointSaleID int64, incomeCreate *schemas.IncomeOtherCreate) (id int64, err error)
	IncomeOtherUpdate(memberID, pointSaleID int64, incomeUpdate *schemas.IncomeOtherUpdate) (err error)
	IncomeOtherDelete(incomeOtherID, pointSaleID int64) error
}

type IncomeOtherService interface {
	IncomeOtherGetByID(id int64) (income *schemas.IncomeOtherResponse, err error)
	IncomeOtherGetByDate(pointSaleID *int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.IncomeOtherResponse, int64, error)
	IncomeOtherCreate(memberID, pointSaleID int64, incomeCreate *schemas.IncomeOtherCreate) (id int64, err error)
	IncomeOtherUpdate(memberID, pointSaleID int64, incomeUpdate *schemas.IncomeOtherUpdate) (err error)
	IncomeOtherDelete(incomeOtherID, pointSaleID int64) error
}
