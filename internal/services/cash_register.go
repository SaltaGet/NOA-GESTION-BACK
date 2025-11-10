package services

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (r *CashRegisterService) CashRegisterExistOpen(pointSaleID int64) (bool, error) {
	return r.CashRegisterRepository.CashRegisterExistOpen(pointSaleID)
}

func (r *CashRegisterService) CashRegisterGetByID(pointSaleID, id int64) (*schemas.CashRegisterFullResponse, error) {
	return r.CashRegisterRepository.CashRegisterGetByID(pointSaleID, id)
}

func (r *CashRegisterService) CashRegisterOpen(pointSaleID int64, userID int64, amountOpen schemas.CashRegisterOpen) error {
	return r.CashRegisterRepository.CashRegisterOpen(pointSaleID, userID, amountOpen)
}

func (r *CashRegisterService) CashRegisterClose(pointSaleID int64, userID int64, amountOpen schemas.CashRegisterClose) error {
	return r.CashRegisterRepository.CashRegisterClose(pointSaleID, userID, amountOpen)
}

func (r *CashRegisterService) CashRegisterInform(pointSaleID int64, userID int64, fromDate, toDate time.Time) ([]*schemas.CashRegisterInformResponse, error) {
	return r.CashRegisterRepository.CashRegisterInform(pointSaleID, userID, fromDate, toDate)
}