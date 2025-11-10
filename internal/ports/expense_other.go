package ports

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type ExpenseOtherRepository interface {
	ExpenseOtherGetByID(id int64) (*schemas.ExpenseOtherResponse, error)
	ExpenseOtherGetByDate(pointSaleID *int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.ExpenseOtherResponseDTO, int64, error)
	ExpenseOtherCreate(memberID, pointSaleID, typeExpenseID int64, expenseOtherCreate *schemas.ExpenseOtherCreate) (int64, error)
	ExpenseOtherUpdate(memberID, pointSaleID int64, expenseOtherUpdate *schemas.ExpenseOtherUpdate) (error)
	ExpenseOtherDelete(expenseOtherID, pointSaleID int64) error
}

type ExpenseOtherService interface {
	ExpenseOtherGetByID(id int64) (*schemas.ExpenseOtherResponse, error)
	ExpenseOtherGetByDate(pointSaleID *int64, fromDate, toDate time.Time, page, limit int) ([]*schemas.ExpenseOtherResponseDTO, int64, error)
	ExpenseOtherCreate(memberID, pointSaleID, typeExpenseID int64, expenseOtherCreate *schemas.ExpenseOtherCreate) (int64, error)
	ExpenseOtherUpdate(memberID, pointSaleID int64, expenseOtherUpdate *schemas.ExpenseOtherUpdate) (error)
	ExpenseOtherDelete(expenseOtherID, pointSaleID int64) error
}