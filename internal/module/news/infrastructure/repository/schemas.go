package repository


import (
	"time"

)

type NewsResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NewsResponseDTO struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type NewsCreate struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type NewsUpdate struct {
	ID      int64  `json:"id" validate:"required"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}
