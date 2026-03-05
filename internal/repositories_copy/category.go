package repositories_copy

import (
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"gorm.io/gorm"
)

func (r *CategoryRepository) CategoryGetByID(id int64) (*models.Category, error) {
	var category *models.Category

	if err := r.DB.First(&category, id).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Categoria", schemas.Read)
	}

	return category, nil
}

func (r *CategoryRepository) CategoryGetAll() ([]*models.Category, error) {
	var categories []*models.Category

	if err := r.DB.Find(&categories).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Categoria", schemas.Read)
	}

	return categories, nil
}

func (r *CategoryRepository) CategoryCreate(memberID int64, categoryCreate *schemas.CategoryCreate) (int64, error) {
	category := models.Category{
		Name: categoryCreate.Name,
	}

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		if err := tx.Create(&category).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Categoria", schemas.Create)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return category.ID, nil
}

func (r *CategoryRepository) CategoryUpdate(memberID int64, categoryUpdate *schemas.CategoryUpdate) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var oldCategory models.Category
		if err := tx.First(&oldCategory, categoryUpdate.ID).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Categoria", schemas.Read)
		}
		if err := tx.Model(&models.Category{}).
			Where("id = ?", categoryUpdate.ID).
			Updates(map[string]any{
				"name": categoryUpdate.Name,
			}).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Categoria", schemas.Update)
		}

		return nil
	})
}

func (r *CategoryRepository) CategoryDelete(memberID, id int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", memberID)).Error; err != nil {
			return err
		}

		var oldCategory models.Category
		if err := tx.First(&oldCategory, id).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Categoria", schemas.Read)
		}

		res := tx.Delete(&oldCategory)
		if err := res.Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Categoria", schemas.Delete)
		}

		return nil
	})
}
