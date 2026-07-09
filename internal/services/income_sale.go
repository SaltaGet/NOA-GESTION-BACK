package services

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (i *IncomeSaleService) IncomeSaleGetByID(pointSaleID, id int64) (*schemas.IncomeSaleResponse, error) {
	return i.IncomeSaleRepository.IncomeSaleGetByID(pointSaleID, id)
}

func (i *IncomeSaleService) IncomeSaleGetByDate(pointSaleID int64, fromDate, toDate time.Time, page, limit int, isBudget *bool, delivered *bool) ([]*schemas.IncomeSaleResponseDTO, int64, error) {
	return i.IncomeSaleRepository.IncomeSaleGetByDate(pointSaleID, fromDate, toDate, page, limit, isBudget, delivered)
}

func (i *IncomeSaleService) IncomeSaleCreate(memberID, pointSaleID int64, incomeSaleCreate *schemas.IncomeSaleCreate) (int64, error) {
	return i.IncomeSaleRepository.IncomeSaleCreate(memberID, pointSaleID, incomeSaleCreate)
}

func (i *IncomeSaleService) IncomeSaleUpdate(memberID, pointSaleID int64, incomeSaleUpdate *schemas.IncomeSaleUpdate, maintainPrice bool) error {
	return i.IncomeSaleRepository.IncomeSaleUpdate(memberID, pointSaleID, incomeSaleUpdate, maintainPrice)
}

func (i *IncomeSaleService) IncomeSaleDelete(memberID int64, id, pointSaleID int64) error {
	return i.IncomeSaleRepository.IncomeSaleDelete(memberID, id, pointSaleID)
}