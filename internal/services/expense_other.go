package services

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (e *ExpenseOtherService) ExpenseOtherGetByID(id int64, pointSaleID *int64) (*schemas.ExpenseOtherResponse, error) {
	return e.ExpenseOtherRepository.ExpenseOtherGetByID(id, pointSaleID)
}

func (e *ExpenseOtherService) ExpenseOtherGetByDate(pointSaleID *int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.ExpenseOtherResponse, int64, error) {
	return e.ExpenseOtherRepository.ExpenseOtherGetByDate(pointSaleID, fromDate, toDate, page, limit)
}

func (e *ExpenseOtherService) ExpenseOtherCreate(memberID int64, pointSaleID *int64, expenseOtherCreate *schemas.ExpenseOtherCreate) (int64, error) {
	return e.ExpenseOtherRepository.ExpenseOtherCreate(memberID, pointSaleID, expenseOtherCreate)
}

func (e *ExpenseOtherService) ExpenseOtherUpdate(memberID int64, pointSaleID *int64, expenseOtherUpdate *schemas.ExpenseOtherUpdate) (error) {
	return e.ExpenseOtherRepository.ExpenseOtherUpdate(memberID, pointSaleID, expenseOtherUpdate)
}

func (e *ExpenseOtherService) ExpenseOtherDelete(memberID int64, expenseOtherID int64, pointSaleID *int64) error {
	return e.ExpenseOtherRepository.ExpenseOtherDelete(memberID, expenseOtherID, pointSaleID)
}
