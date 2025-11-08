package ports

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type CashRegisterRepository interface {
	CashRegisterOpen(pointSaleID uint, userID uint, amountOpen schemas.CashRegisterOpen) error
	CashRegisterGetByID(pointSaleID, id uint) (*models.CashRegister, error)
	CashRegisterClose(pointSaleID uint, userID uint, amountOpen schemas.CashRegisterClose) error
	CashRegisterInform(pointSaleID uint, userID uint, fromDate, toDate time.Time) ([]*models.CashRegister, error)
	CashRegisterExistOpen(pointSaleID uint) (bool, error)
}

type CashRegisterService interface {
	CashRegisterOpen(pointSaleID uint, userID uint, amountOpen schemas.CashRegisterOpen) error
	CashRegisterGetByID(pointSaleID, id uint) (*schemas.CashRegisterFullResponse, error)
	CashRegisterClose(pointSaleID uint, userID uint, amountOpen schemas.CashRegisterClose) error
	CashRegisterInform(pointSaleID uint, userID uint, fromDate, toDate time.Time) ([]*schemas.CashRegisterInformResponse, error)
	CashRegisterExistOpen(pointSaleID uint) (bool, error)
}
