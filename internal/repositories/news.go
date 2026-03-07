package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/master"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

func mapToNewsResponse(n *boilmodels.NewsItem) *schemas.NewsResponse {
	if n == nil {
		return nil
	}

	res := &schemas.NewsResponse{
		ID:      n.ID,
		Title:   n.Title,
		Content: n.Content,
	}

	res.CreatedAt = n.CreatedAt
	res.UpdatedAt = n.UpdatedAt

	return res
}

func mapToNewsResponseDTO(n *boilmodels.NewsItem) schemas.NewsResponseDTO {
	if n == nil {
		return schemas.NewsResponseDTO{}
	}

	res := schemas.NewsResponseDTO{
		ID:    n.ID,
		Title: n.Title,
	}

	res.CreatedAt = n.CreatedAt

	return res
}

func (r *MainRepository) NewsGetByID(id int64) (*schemas.NewsResponse, error) {
	ctx := context.Background()

	news, err := boilmodels.News(boilmodels.NewsItemWhere.ID.EQ(id)).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Noticia", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Noticia", schemas.Read)
	}

	return mapToNewsResponse(news), nil
}

func (r *MainRepository) NewsGetAll() ([]schemas.NewsResponseDTO, error) {
	ctx := context.Background()

	newsList, err := boilmodels.News(
		qm.Select(
			boilmodels.NewsItemColumns.ID,
			boilmodels.NewsItemColumns.Title,
			boilmodels.NewsItemColumns.CreatedAt,
		),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Noticia", schemas.Read)
	}

	var newsResponse []schemas.NewsResponseDTO
	for _, n := range newsList {
		newsResponse = append(newsResponse, mapToNewsResponseDTO(n))
	}

	return newsResponse, nil
}

func (r *MainRepository) NewsCreate(adminID int64, newsCreate *schemas.NewsCreate) (int64, error) {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", adminID)); err != nil {
		return 0, err
	}

	newNews := &boilmodels.NewsItem{
		Title:   newsCreate.Title,
		Content: newsCreate.Content,
	}

	if err := newNews.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Noticia", schemas.Create)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newNews.ID, nil
}

func (r *MainRepository) NewsUpdate(adminID int64, newsUpdate *schemas.NewsUpdate) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", adminID)); err != nil {
		return err
	}

	news, err := boilmodels.News(boilmodels.NewsItemWhere.ID.EQ(newsUpdate.ID)).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Noticia", schemas.Read)
	}

	news.Title = newsUpdate.Title
	news.Content = newsUpdate.Content

	if _, err := news.Update(ctx, tx, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Noticia", schemas.Update)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *MainRepository) NewsDelete(adminID int64, id int64) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", adminID)); err != nil {
		return err
	}

	news, err := boilmodels.News(boilmodels.NewsItemWhere.ID.EQ(id)).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Noticia", schemas.Read)
	}

	if _, err := news.Delete(ctx, tx); err != nil {
		return schemas.HandlerErrorDB(err, "Noticia", schemas.Delete)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
