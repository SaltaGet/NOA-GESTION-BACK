package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
)

func (r *CategoryRepository) CategoryGetByID(id int64) (*tenant.Category, error) {
	ctx := context.Background()

	c, err := tenant.Categories(
		tenant.CategoryWhere.ID.EQ(id),
		tenant.CategoryWhere.DeleteAt.IsNull(),
	).One(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Categoria", schemas.Read)
	}

	return c, nil
}

func (r *CategoryRepository) CategoryGetAll() ([]*tenant.Category, error) {
	ctx := context.Background()

	boilCategories, err := tenant.Categories(
		tenant.CategoryWhere.DeleteAt.IsNull(),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Categoria", schemas.Read)
	}

	return boilCategories, nil
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

	c := tenant.Category{
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

	c, err := tenant.Categories(
		tenant.CategoryWhere.ID.EQ(categoryUpdate.ID),
		tenant.CategoryWhere.DeleteAt.IsNull(),
	).One(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Categoria", schemas.Read)
	}

	c.Name = strings.ToLower(categoryUpdate.Name)

	if _, err := c.Update(ctx, tx, boil.Whitelist(tenant.CategoryColumns.Name, tenant.CategoryColumns.UpdatedAt)); err != nil {
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

	c, err := tenant.Categories(
		tenant.CategoryWhere.ID.EQ(id),
		tenant.CategoryWhere.DeleteAt.IsNull(),
	).One(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Categoria", schemas.Read)
	}

	c.DeleteAt = null.TimeFrom(time.Now())
	if _, err := c.Update(ctx, tx, boil.Whitelist(tenant.CategoryColumns.DeleteAt, tenant.CategoryColumns.UpdatedAt)); err != nil {
		return schemas.HandlerErrorDB(err, "Categoria", schemas.Delete)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
