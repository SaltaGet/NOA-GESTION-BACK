package repository


import (
	"time"
)

type FeedbackResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FeedbackResponseDTO struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type FeedbackCreate struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

