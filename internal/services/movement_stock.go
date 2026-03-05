package services

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (m *MovementStockService) MovementStockGetByID(id int64) (*schemas.MovementStockResponse, error) {
	return m.MovementStockRepository.MovementStockGetByID(id)
}

func (m *MovementStockService) MovementStockGetByDate(page, limit int, fromDate, toDate time.Time) ([]schemas.MovementStockResponseDTO, int64, error) {
	return m.MovementStockRepository.MovementStockGetByDate(page, limit, fromDate, toDate)
}

func (m *MovementStockService) MoveStockList(userID int64, input []schemas.MovementStockList) error {
	return m.MovementStockRepository.MoveStockList(userID, input)
}
