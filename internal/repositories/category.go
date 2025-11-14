package repositories

import (
	"errors"
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"gorm.io/gorm"
)

func (r *CategoryRepository) CategoryGetByID(id int64) (*models.Category, error) {
	var category *models.Category

	if err := r.DB.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "categoria no encontrada", err)
		}
		return nil, schemas.ErrorResponse(500, "error al obtener la categoria", err)
	}

	return category, nil
}

func (r *CategoryRepository) CategoryGetAll() ([]*models.Category, error) {
	var categories []*models.Category

	if err := r.DB.Find(&categories).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "error al obtener las categorias", err)
	}

	return categories, nil
}

func (r *CategoryRepository) CategoryCreate(categoryCreate *schemas.CategoryCreate) (int64, error) {
	var category models.Category

	category.Name = categoryCreate.Name

	if err := r.DB.Create(&category).Error; err != nil {
		if schemas.IsDuplicateError(err) {
			return 0, schemas.ErrorResponse(400, "la categoria "+categoryCreate.Name+" ya existe", err)
		}
		return 0, schemas.ErrorResponse(500, "error al crear la categoria", err)
	}

	return category.ID, nil
}

func (r *CategoryRepository) CategoryUpdate(categoryUpdate *schemas.CategoryUpdate) error {
	var exists bool

	if err := r.DB.Model(&models.Category{}).
		Select("count(*) > 0").
		Where("id = ?", categoryUpdate.ID).
		Find(&exists).Error; err != nil {
		return schemas.ErrorResponse(500, "error al obtener la categoria", err)
	}

	if !exists {
		return schemas.ErrorResponse(404, "categoria no encontrada", fmt.Errorf("categoria no encontrada"))
	}

	if err := r.DB.Model(&models.Category{}).
		Where("id = ?", categoryUpdate.ID).
		Updates(map[string]any{
			"name": categoryUpdate.Name,
		}).Error; err != nil {
		if schemas.IsDuplicateError(err) {
			return schemas.ErrorResponse(400, "la categoria "+categoryUpdate.Name+" ya existe", err)
		}
		return schemas.ErrorResponse(500, "error al actualizar la categoria", err)
	}

	return nil
}

func (r *CategoryRepository) CategoryDelete(id int64) error {
	if err := r.DB.Where("id = ?", id).Delete(&models.Category{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "categoria no encontrada", err)
		}
		return schemas.ErrorResponse(500, "error al eliminar la categoria", err)
	}

	return nil
}
