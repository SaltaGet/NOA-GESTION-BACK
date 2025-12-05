package services

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) AuthLogin(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) AuthPointSale(user *schemas.AuthenticatedUser, pointSaleID int64) (string, error) {
	args := m.Called(user, pointSaleID)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) LogoutPointSale(member *schemas.AuthenticatedUser) (string, error) {
	args := m.Called(member)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) AuthLoginAdmin(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) AuthForgotPassword(authForgotPassword *schemas.AuthForgotPassword) error {
	args := m.Called(authForgotPassword)
	return args.Error(0)
}

func (m *MockAuthService) AuthResetPassword(authResetPassword *schemas.AuthResetPassword) error {
	args := m.Called(authResetPassword)
	return args.Error(0)
}

func (m *MockAuthService) AuthAdminGetByID(id int64) (*models.Admin, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Admin), args.Error(1)
}

func (m *MockAuthService) AuthCurrentPlan(id int64) (*schemas.PlanResponseDTO, error) {
	args := m.Called(id)
	return args.Get(0).(*schemas.PlanResponseDTO), args.Error(1)
}

func (m *MockAuthService) AuthCurrentTenant(id int64) (*schemas.TenantResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*schemas.TenantResponse), args.Error(1)
}

func (m *MockAuthService) AuthCurrentUser(tenantID, memberID, pointSaleID int64) (*schemas.AuthenticatedUser, error) {
	args := m.Called(tenantID, memberID, pointSaleID)
	return args.Get(0).(*schemas.AuthenticatedUser), args.Error(1)
}
