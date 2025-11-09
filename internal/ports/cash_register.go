package ports

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type CashRegisterRepository interface {
	CashRegisterOpen(pointSaleID int64, userID int64, amountOpen schemas.CashRegisterOpen) error
	CashRegisterGetByID(pointSaleID, id int64) (*schemas.CashRegisterFullResponse, error)
	CashRegisterClose(pointSaleID int64, userID int64, amountOpen schemas.CashRegisterClose) error
	CashRegisterInform(pointSaleID int64, userID int64, fromDate, toDate time.Time) ([]*schemas.CashRegisterInformResponse, error)
	CashRegisterExistOpen(pointSaleID int64) (bool, error)
}

type CashRegisterService interface {
	CashRegisterOpen(pointSaleID int64, userID int64, amountOpen schemas.CashRegisterOpen) error
	CashRegisterGetByID(pointSaleID, id int64) (*schemas.CashRegisterFullResponse, error)
	CashRegisterClose(pointSaleID int64, userID int64, amountOpen schemas.CashRegisterClose) error
	CashRegisterInform(pointSaleID int64, userID int64, fromDate, toDate time.Time) ([]*schemas.CashRegisterInformResponse, error)
	CashRegisterExistOpen(pointSaleID int64) (bool, error)
}
