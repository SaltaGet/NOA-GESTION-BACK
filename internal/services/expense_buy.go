package services

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (e *ExpenseBuyService) ExpenseBuyGetByID(id int64) (*schemas.ExpenseBuyResponse, error) {
	return e.ExpenseBuyRepository.ExpenseBuyGetByID(id)
}

func (e *ExpenseBuyService) ExpenseBuyGetByDate(fromDate, toDate time.Time, page, limit int) ([]*schemas.ExpenseBuyResponseSimple, int64, error) {
	return e.ExpenseBuyRepository.ExpenseBuyGetByDate(fromDate, toDate, page, limit)
}

func (e *ExpenseBuyService) ExpenseBuyCreate(memberID int64, expenseBuy *schemas.ExpenseBuyCreate) (int64, error) {
	return e.ExpenseBuyRepository.ExpenseBuyCreate(memberID, expenseBuy)
}

func (e *ExpenseBuyService) ExpenseBuyUpdate(memberID int64, expenseBuy *schemas.ExpenseBuyUpdate) error {
	return e.ExpenseBuyRepository.ExpenseBuyUpdate(memberID, expenseBuy)
}

func (e *ExpenseBuyService) ExpenseBuyDelete(memberID int64, id int64) error {
	return e.ExpenseBuyRepository.ExpenseBuyDelete(memberID,id)
}
