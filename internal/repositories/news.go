package repositories

import (
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func (r *MainRepository) NewsGetByID(id int64) (*schemas.NewsResponse, error) {
	var newGet models.News
	if err := r.DB.Where("id = ?", id).First(&newGet).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Noticia", schemas.Read)
	}

	var newsResponse schemas.NewsResponse
	copier.Copy(&newsResponse, &newGet)

	return &newsResponse, nil
}

func (r *MainRepository) NewsGetAll() ([]schemas.NewsResponseDTO, error) {
	var news []models.News
	if err := r.DB.Select("id", "title", "created_at").Find(&news).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Noticia", schemas.Read)
	}

	var newsResponse []schemas.NewsResponseDTO
	copier.Copy(&newsResponse, &news)

	return newsResponse, nil
}

func (r *MainRepository) NewsCreate(adminID int64, newsCreate *schemas.NewsCreate) (int64, error) {
	newNews := models.News{
		Title:   newsCreate.Title,
		Content: newsCreate.Content,
	}

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", adminID)).Error; err != nil {
			return err
		}
		if err := tx.Create(&newNews).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 0, schemas.HandlerErrorGorm(err, "Noticia", schemas.Create)
	}

	return newNews.ID, nil
}

func (r *MainRepository) NewsUpdate(adminID int64, newsUpdate *schemas.NewsUpdate) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", adminID)).Error; err != nil {
			return err
		}
		var news models.News
		if err := tx.
			Where("id = ?", newsUpdate.ID).First(&news).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Noticia", schemas.Read)
		}

		news.Title = newsUpdate.Title
		news.Content = newsUpdate.Content
		if err := tx.Save(&news).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Noticia", schemas.Update)
		}

		return nil
	})
}

func (r *MainRepository) NewsDelete(adminID int64, id int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT set_config('app.current_member_id', ?, true)", fmt.Sprintf("%d", adminID)).Error; err != nil {
			return err
		}
		var news models.News
		if err := tx.
			Where("id = ?", id).First(&news).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Noticia", schemas.Read)
		}

		if err :=	tx.Delete(&news).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Noticia", schemas.Delete)
		}

		return nil
	})
}
