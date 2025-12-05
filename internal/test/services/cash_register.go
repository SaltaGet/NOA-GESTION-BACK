package services

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/stretchr/testify/mock"
)

type MockCashRegisterService struct {
	mock.Mock
}

func (m *MockCashRegisterService) CashRegisterOpen(pointSaleID int64, userID int64, amountOpen schemas.CashRegisterOpen) error {
	args := m.Called(pointSaleID, userID, amountOpen)
	return args.Error(0)
}

func (m *MockCashRegisterService) CashRegisterGetByID(pointSaleID, id int64) (*schemas.CashRegisterFullResponse, error) {
	args := m.Called(pointSaleID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schemas.CashRegisterFullResponse), args.Error(1)
}

func (m *MockCashRegisterService) CashRegisterClose(pointSaleID int64, userID int64, amountClose schemas.CashRegisterClose) error {
	args := m.Called(pointSaleID, userID, amountClose)
	return args.Error(0)
}

func (m *MockCashRegisterService) CashRegisterInform(pointSaleID int64, userID int64, fromDate, toDate time.Time) ([]*schemas.CashRegisterInformResponse, error) {
	args := m.Called(pointSaleID, userID, fromDate, toDate)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*schemas.CashRegisterInformResponse), args.Error(1)
}

func (m *MockCashRegisterService) CashRegisterExistOpen(pointSaleID int64) (bool, error) {
	args := m.Called(pointSaleID)
	return args.Bool(0), args.Error(1)
}
