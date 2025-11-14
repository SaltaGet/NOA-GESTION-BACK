package ports

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type DepositRepository interface {
	DepositGetByID(id int64) (*models.Product, error)
	DepositGetByCode(code string) (*models.Product, error)
	DepositGetByName(name string) ([]*models.Product, error)
	DepositGetAll(page, limit int) ([]*models.Product, int64,error)
	DepositUpdateStock(updateStock schemas.DepositUpdateStock) (error)
}

type DepositService interface {
	DepositGetByID(id int64) (*schemas.DepositResponse, error)
	DepositGetByCode(code string) (*schemas.DepositResponse, error)
	DepositGetByName(name string) ([]*schemas.DepositResponse, error)
	DepositGetAll(page, limit int) ([]*schemas.DepositResponse, int64, error)
	DepositUpdateStock(updateStock schemas.DepositUpdateStock) (error)
}