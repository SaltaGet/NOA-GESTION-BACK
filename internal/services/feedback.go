package services

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

func (n *FeedbackService) FeedbackGetByID(id int64) (*schemas.FeedbackResponse, error) {
	return n.FeedbackRepository.FeedbackGetByID(id)
}

func (n *FeedbackService) FeedbackGetAll() ([]schemas.FeedbackResponseDTO, error) {
	return n.FeedbackRepository.FeedbackGetAll()
}

func (n *FeedbackService) FeedbackCreate(news *schemas.FeedbackCreate) (int64, error) {
	return n.FeedbackRepository.FeedbackCreate(news)
}

