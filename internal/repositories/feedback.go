package repositories

import (
	"context"
	"database/sql"
	"errors"

	mastermodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/master"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

func mapToFeedbackResponse(f *mastermodels.Feedback) *schemas.FeedbackResponse {
	if f == nil {
		return nil
	}
	res := &schemas.FeedbackResponse{
		ID:      f.ID,
		Title:   f.Title,
		Content: f.Content,
	}
	res.CreatedAt = f.CreatedAt
	res.UpdatedAt = f.UpdatedAt
	return res
}

func mapToFeedbackResponseDTO(f *mastermodels.Feedback) schemas.FeedbackResponseDTO {
	res := schemas.FeedbackResponseDTO{
		ID:     f.ID,
		Title:  f.Title,
		IsRead: f.IsRead,
	}
	res.CreatedAt = f.CreatedAt
	return res
}

func (r *MainRepository) FeedbackGetByID(id int64) (*schemas.FeedbackResponse, error) {
	ctx := context.Background()

	feedback, err := mastermodels.Feedbacks(mastermodels.FeedbackWhere.ID.EQ(id)).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Feedback", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Feedback", schemas.Read)
	}

	feedback.IsRead = true
	if _, err := feedback.Update(ctx, r.DB, boil.Whitelist(mastermodels.FeedbackColumns.IsRead, mastermodels.FeedbackColumns.UpdatedAt)); err != nil {
		return nil, schemas.HandlerErrorDB(err, "Feedback", schemas.Update)
	}

	return mapToFeedbackResponse(feedback), nil
}

func (r *MainRepository) FeedbackGetAll() ([]schemas.FeedbackResponseDTO, error) {
	ctx := context.Background()

	feedbacks, err := mastermodels.Feedbacks(
		qm.Select(
			mastermodels.FeedbackColumns.ID,
			mastermodels.FeedbackColumns.Title,
			mastermodels.FeedbackColumns.IsRead,
			mastermodels.FeedbackColumns.CreatedAt,
		),
		qm.OrderBy("created_at DESC"),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Feedback", schemas.Read)
	}

	var results []schemas.FeedbackResponseDTO
	for _, f := range feedbacks {
		results = append(results, mapToFeedbackResponseDTO(f))
	}

	return results, nil
}

func (r *MainRepository) FeedbackCreate(newsCreate *schemas.FeedbackCreate) (int64, error) {
	ctx := context.Background()

	newFeedback := mastermodels.Feedback{
		Title:   newsCreate.Title,
		Content: newsCreate.Content,
	}

	if err := newFeedback.Insert(ctx, r.DB, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Feedback", schemas.Create)
	}

	return newFeedback.ID, nil
}
