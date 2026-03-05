package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func mapToModelCategory(c *boilmodels.Category) *models.Category {
	var updatedAt time.Time
	if c.UpdatedAt.Valid {
		updatedAt = c.UpdatedAt.Time
	}
	return &models.Category{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: updatedAt,
	}
}

func (r *CategoryRepository) CategoryGetByID(id int64) (*models.Category, error) {
	ctx := context.Background()

	c, err := boilmodels.Categories(
		boilmodels.CategoryWhere.ID.EQ(id),
		boilmodels.CategoryWhere.DeleteAt.IsNull(),
	).One(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Categoria", schemas.Read)
	}

	return mapToModelCategory(c), nil
}

func (r *CategoryRepository) CategoryGetAll() ([]*models.Category, error) {
	ctx := context.Background()

	boilCategories, err := boilmodels.Categories(
		boilmodels.CategoryWhere.DeleteAt.IsNull(),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Categoria", schemas.Read)
	}

	var categories []*models.Category
	for _, c := range boilCategories {
		categories = append(categories, mapToModelCategory(c))
	}

	return categories, nil
}

func (r *CategoryRepository) CategoryCreate(memberID int64, categoryCreate *schemas.CategoryCreate) (int64, error) {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Categoria", schemas.Create)
	}

	c := boilmodels.Category{
		Name: strings.ToLower(categoryCreate.Name),
	}

	if err := c.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Categoria", schemas.Create)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return c.ID, nil
}

func (r *CategoryRepository) CategoryUpdate(memberID int64, categoryUpdate *schemas.CategoryUpdate) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	c, err := boilmodels.Categories(
		boilmodels.CategoryWhere.ID.EQ(categoryUpdate.ID),
		boilmodels.CategoryWhere.DeleteAt.IsNull(),
	).One(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Categoria", schemas.Read)
	}

	c.Name = strings.ToLower(categoryUpdate.Name)

	if _, err := c.Update(ctx, tx, boil.Whitelist(boilmodels.CategoryColumns.Name, boilmodels.CategoryColumns.UpdatedAt)); err != nil {
		return schemas.HandlerErrorDB(err, "Categoria", schemas.Update)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepository) CategoryDelete(memberID, id int64) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	c, err := boilmodels.Categories(
		boilmodels.CategoryWhere.ID.EQ(id),
		boilmodels.CategoryWhere.DeleteAt.IsNull(),
	).One(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Categoria", schemas.Read)
	}

	c.DeleteAt = null.TimeFrom(time.Now())
	if _, err := c.Update(ctx, tx, boil.Whitelist(boilmodels.CategoryColumns.DeleteAt, boilmodels.CategoryColumns.UpdatedAt)); err != nil {
		return schemas.HandlerErrorDB(err, "Categoria", schemas.Delete)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
