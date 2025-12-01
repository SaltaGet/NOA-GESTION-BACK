package ports

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type MovementStockRepository interface {
	MovementStockGetByID(id int64) (*models.MovementStock, error)
	MovementStockGetByDate(page, limit int, fromDate, toDate time.Time) ([]*models.MovementStock, int64, error)
	MoveStockList(userID int64, input []*schemas.MovementStockList) error
}

type MovementStockService interface {
	MovementStockGetByID(id int64) (*schemas.MovementStockResponse, error)
	MovementStockGetByDate(page, limit int, fromDate, toDate time.Time) ([]*schemas.MovementStockResponseDTO, int64, error)
	MoveStockList(userID int64, input []*schemas.MovementStockList) error
}
