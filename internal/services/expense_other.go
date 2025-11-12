package services

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (e *ExpenseOtherService) ExpenseOtherGetByID(id int64) (*schemas.ExpenseOtherResponse, error) {
	return e.ExpenseOtherRepository.ExpenseOtherGetByID(id)
}

func (e *ExpenseOtherService) ExpenseOtherGetByDate(pointSaleID *int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.ExpenseOtherResponseDTO, int64, error) {
	return e.ExpenseOtherRepository.ExpenseOtherGetByDate(pointSaleID, fromDate, toDate, page, limit)
}

func (e *ExpenseOtherService) ExpenseOtherCreate(memberID, pointSaleID int64, expenseOtherCreate *schemas.ExpenseOtherCreate) (int64, error) {
	return e.ExpenseOtherRepository.ExpenseOtherCreate(memberID, pointSaleID, expenseOtherCreate)
}

func (e *ExpenseOtherService) ExpenseOtherUpdate(memberID, pointSaleID int64, expenseOtherUpdate *schemas.ExpenseOtherUpdate) (error) {
	return e.ExpenseOtherRepository.ExpenseOtherUpdate(memberID, pointSaleID, expenseOtherUpdate)
}

func (e *ExpenseOtherService) ExpenseOtherDelete(expenseOtherID, pointSaleID int64) error {
	return e.ExpenseOtherRepository.ExpenseOtherDelete(expenseOtherID, pointSaleID)
}
