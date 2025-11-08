package services

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (r *CashRegisterService) CashRegisterExistOpen(pointSaleID uint) (bool, error) {
	return r.CashRegisterRepository.CashRegisterExistOpen(pointSaleID)
}

func (r *CashRegisterService) CashRegisterGetByID(pointSaleID, id uint) (*schemas.CashRegisterFullResponse, error) {
	register, err := r.CashRegisterRepository.CashRegisterGetByID(pointSaleID, id)
	if err != nil {
		return nil, err
	}

	return &register, nil
}

func (r *CashRegisterService) CashRegisterOpen(pointSaleID uint, userID uint, amountOpen schemas.CashRegisterOpen) error {
	return r.CashRegisterRepository.CashRegisterOpen(pointSaleID, userID, amountOpen)
}

func (r *CashRegisterService) CashRegisterClose(pointSaleID uint, userID uint, amountOpen schemas.CashRegisterClose) error {
	return r.CashRegisterRepository.CashRegisterClose(pointSaleID, userID, amountOpen)
}

func (r *CashRegisterService) CashRegisterInform(pointSaleID uint, userID uint, fromDate, toDate time.Time) ([]*schemas.CashRegisterInformResponse, error) {
	informs, err := r.CashRegisterRepository.CashRegisterInform(pointSaleID, userID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	
	return informs, nil
}