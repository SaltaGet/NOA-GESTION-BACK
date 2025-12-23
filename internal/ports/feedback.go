package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

type FeedbackRepository interface {
	FeedbackGetAll() ([]schemas.FeedbackResponseDTO, error)
	FeedbackGetByID(id int64) (*schemas.FeedbackResponse, error)
	FeedbackCreate(news *schemas.FeedbackCreate) (int64, error)
}

type FeedbackServices interface {
	FeedbackGetAll() ([]schemas.FeedbackResponseDTO, error)
	FeedbackGetByID(id int64) (*schemas.FeedbackResponse, error)
	FeedbackCreate(news *schemas.FeedbackCreate) (int64, error)
}