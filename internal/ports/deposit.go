package ports

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type DepositRepository interface {
	DepositGetByID(id int64) (*tenant.Product, error)
	DepositGetByCode(code string) (*tenant.Product, error)
	DepositGetByName(name string) ([]*tenant.Product, error)
	DepositGetAll(page, limit int) ([]*tenant.Product, int64,error)
	DepositUpdateStock(memberID int64, updateStock schemas.DepositUpdateStock) (error)
}

type DepositService interface {
	DepositGetByID(id int64) (*schemas.DepositResponse, error)
	DepositGetByCode(code string) (*schemas.DepositResponse, error)
	DepositGetByName(name string) ([]*schemas.DepositResponse, error)
	DepositGetAll(page, limit int) ([]*schemas.DepositResponse, int64, error)
	DepositUpdateStock(memberID int64, updateStock schemas.DepositUpdateStock) (error)
}