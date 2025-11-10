package ports

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type ExpenseBuyRepository interface {
	ExpenseBuyGetByID(id int64) (*schemas.ExpenseBuyResponse, error)
	ExpenseBuyGetByDate(fromDate, toDate time.Time, page, limit int) ([]*schemas.ExpenseBuyResponseSimple, int64, error)
	ExpenseBuyCreate(userID int64, incomeCreate *schemas.ExpenseBuyCreate) (int64, error)
	ExpenseBuyUpdate(userID int64, incomeCreate *schemas.ExpenseBuyUpdate) (error)
	ExpenseBuyDelete(id int64) error
}

type ExpenseBuyService interface {
	ExpenseBuyGetByID(id int64) (*schemas.ExpenseBuyResponse, error)
	ExpenseBuyGetByDate(fromDate, toDate time.Time, page, limit int) ([]*schemas.ExpenseBuyResponseSimple, int64, error)
	ExpenseBuyCreate(userID int64, incomeCreate *schemas.ExpenseBuyCreate) (int64, error)
	ExpenseBuyUpdate(userID int64, incomeCreate *schemas.ExpenseBuyUpdate) (error)
	ExpenseBuyDelete(id int64) error
}