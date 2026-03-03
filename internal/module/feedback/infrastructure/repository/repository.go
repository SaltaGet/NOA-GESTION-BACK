package repository


import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
)

func (r *MainRepository) FeedbackGetByID(id int64) (*schemas.FeedbackResponse, error) {
	var newGet models.Feedback
	if err := r.DB.Where("id = ?", id).First(&newGet).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Feedback", schemas.Read)
	}

	if err := r.DB.Model(&newGet).Update("IsRead", true).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Feedback", schemas.Update)
	}

	var newsResponse schemas.FeedbackResponse
	copier.Copy(&newsResponse, &newGet)

	return &newsResponse, nil
}

func (r *MainRepository) FeedbackGetAll() ([]schemas.FeedbackResponseDTO, error) {
	var news []models.Feedback
	if err := r.DB.
		Select("id", "title", "is_read", "created_at").
		Order("created_at DESC").
		Find(&news).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Feedback", schemas.Read)
	}

	var newsResponse []schemas.FeedbackResponseDTO
	copier.Copy(&newsResponse, &news)

	return newsResponse, nil
}

func (r *MainRepository) FeedbackCreate(newsCreate *schemas.FeedbackCreate) (int64, error) {
	newFeedback := models.Feedback{
		Title:   newsCreate.Title,
		Content: newsCreate.Content,
	}
	if err := r.DB.Create(&newFeedback).Error; err != nil {
		return 0, schemas.HandlerErrorGorm(err, "Feedback", schemas.Create)
	}

	return newFeedback.ID, nil
}
