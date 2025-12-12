package services

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (i *IncomeOtherService) IncomeOtherGetByID(id int64, pointSaleID *int64) (income *schemas.IncomeOtherResponse, err error) {
	return i.IncomeOtherRepository.IncomeOtherGetByID(id, pointSaleID)
}

func (i *IncomeOtherService) IncomeOtherGetByDate(pointSaleID *int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.IncomeOtherResponse, int64, error) {
	return i.IncomeOtherRepository.IncomeOtherGetByDate(pointSaleID, fromDate, toDate, page, limit)
}

func (i *IncomeOtherService) IncomeOtherCreate(memberID int64, pointSaleID *int64, incomeCreate *schemas.IncomeOtherCreate) (id int64, err error) {
	return i.IncomeOtherRepository.IncomeOtherCreate(memberID, pointSaleID, incomeCreate)
}

func (i *IncomeOtherService) IncomeOtherUpdate(memberID int64, pointSaleID *int64, incomeUpdate *schemas.IncomeOtherUpdate) (err error) {
	return i.IncomeOtherRepository.IncomeOtherUpdate(memberID, pointSaleID, incomeUpdate)
}

func (i *IncomeOtherService) IncomeOtherDelete(memberID int64, incomeOtherID int64, pointSaleID *int64,) error {
	return i.IncomeOtherRepository.IncomeOtherDelete(memberID, incomeOtherID, pointSaleID)
}
