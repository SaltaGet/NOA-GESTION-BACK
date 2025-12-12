package ports

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)


type IncomeOtherRepository interface {
	IncomeOtherGetByID(id int64, pointSaleId *int64) (income *schemas.IncomeOtherResponse, err error)
	IncomeOtherGetByDate(pointSaleID *int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.IncomeOtherResponse, int64, error)
	IncomeOtherCreate(memberID int64, pointSaleID *int64, incomeCreate *schemas.IncomeOtherCreate) (id int64, err error)
	IncomeOtherUpdate(memberID int64, pointSaleID *int64, incomeUpdate *schemas.IncomeOtherUpdate) (err error)
	IncomeOtherDelete(memberID int64, incomeOtherID int64, pointSaleID *int64,) error
}

type IncomeOtherService interface {
	IncomeOtherGetByID(id int64, pointSaleId *int64) (income *schemas.IncomeOtherResponse, err error)
	IncomeOtherGetByDate(pointSaleID *int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.IncomeOtherResponse, int64, error)
	IncomeOtherCreate(memberID int64, pointSaleID *int64, incomeCreate *schemas.IncomeOtherCreate) (id int64, err error)
	IncomeOtherUpdate(memberID int64, pointSaleID *int64, incomeUpdate *schemas.IncomeOtherUpdate) (err error)
	IncomeOtherDelete(memberID int64, incomeOtherID int64, pointSaleID *int64,) error
}
