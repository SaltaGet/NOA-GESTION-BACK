package services

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (i *IncomeOtherService) IncomeOtherGetByID(id int64) (income *schemas.IncomeOtherResponse, err error) {
	return i.IncomeOtherRepository.IncomeOtherGetByID(id)
}

func (i *IncomeOtherService) IncomeOtherGetByDate(pointSaleID *int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.IncomeOtherResponse, int64, error) {
	return i.IncomeOtherRepository.IncomeOtherGetByDate(pointSaleID, fromDate, toDate, page, limit)
}

func (i *IncomeOtherService) IncomeOtherCreate(memberID, pointSaleID int64, incomeCreate *schemas.IncomeOtherCreate) (id int64, err error) {
	return i.IncomeOtherRepository.IncomeOtherCreate(memberID, pointSaleID, incomeCreate)
}

func (i *IncomeOtherService) IncomeOtherUpdate(memberID, pointSaleID int64, incomeUpdate *schemas.IncomeOtherUpdate) (err error) {
	return i.IncomeOtherRepository.IncomeOtherUpdate(memberID, pointSaleID, incomeUpdate)
}

func (i *IncomeOtherService) IncomeOtherDelete(incomeOtherID, pointSaleID int64) error {
	return i.IncomeOtherRepository.IncomeOtherDelete(incomeOtherID, pointSaleID)
}
