package services

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
)

func (m *MovementStockService) MovementStockGetByID(id int64) (*schemas.MovementStockResponse, error) {
	movement, err := m.MovementStockRepository.MovementStockGetByID(id)
	if err != nil {
		return nil, err
	}

	var response schemas.MovementStockResponse
	_ = copier.Copy(&response, movement)

	return &response, nil
}

func (m *MovementStockService) MovementStockGetByDate(page, limit int, fromDate, toDate time.Time) ([]*schemas.MovementStockResponseDTO, int64, error) {
	movements, total, err := m.MovementStockRepository.MovementStockGetByDate(page, limit, fromDate, toDate)
	if err != nil {
		return nil, 0, err
	}

	var movementsResponse []*schemas.MovementStockResponseDTO
	_ = copier.Copy(&movementsResponse, &movements)

	return movementsResponse, total, nil
}

func (m *MovementStockService) MoveStockList(userID int64, input []*schemas.MovementStockList) error {
	return m.MovementStockRepository.MoveStockList(userID, input)
}