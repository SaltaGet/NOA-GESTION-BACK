package repositories

import (
	"errors"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func (r *MainRepository) NewsGetByID(id int64) (*schemas.NewsResponse, error) {
	var newGet models.News
	if err := r.DB.Where("id = ?", id).First(&newGet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Noticia no encontrada", err)
		}
		return nil, schemas.ErrorResponse(500, "Error al obtener la noticia", err)
	}

	var newsResponse schemas.NewsResponse
	copier.Copy(&newsResponse, &newGet)

	return &newsResponse, nil
}

func (r *MainRepository) NewsGetAll() ([]schemas.NewsResponseDTO, error) {
	var news []models.News
	if err := r.DB.Select("id", "title", "created_at").Find(&news).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "Error al obtener las noticias", err)
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
	if err := r.DB.Create(&newNews).Error; err != nil {
		return 0, schemas.ErrorResponse(500, "Error al crear la noticia", err)
	}

	go database.SaveAuditAdminAsync(r.DB, models.AuditLogAdmin{
		AdminID: adminID,
		Method:  "create",
		Path:    "news",
	}, nil, newNews)

	return newNews.ID, nil
}

func (r *MainRepository) NewsUpdate(adminID int64, newsUpdate *schemas.NewsUpdate) error {
	var news models.News
	if err := r.DB.
		Where("id = ?", newsUpdate.ID).First(&news).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "Noticia no encontrada", err)
		}
		return schemas.ErrorResponse(500, "Error al obtener la noticia", err)
	}

	oldNews := news

	news.Title = newsUpdate.Title
	news.Content = newsUpdate.Content
	if err := r.DB.Save(&news).Error; err != nil {
		return schemas.ErrorResponse(500, "Error al actualizar la noticia", err)
	}

	go database.SaveAuditAdminAsync(r.DB, models.AuditLogAdmin{
		AdminID: adminID,
		Method:  "update",
		Path:    "news",
	}, &oldNews, &news)

	return nil
}

func (r *MainRepository) NewsDelete(adminID int64, id int64) error {
	var news models.News
	if err := r.DB.
		Where("id = ?", id).First(&news).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "Noticia no encontrada", err)
		}
		return schemas.ErrorResponse(500, "Error al obtener la noticia", err)
	}

	if err := r.DB.Delete(&news).Error; err != nil {
		return schemas.ErrorResponse(500, "Error al eliminar la noticia", err)
	}

	go database.SaveAuditAdminAsync(r.DB, models.AuditLogAdmin{
		AdminID: adminID,
		Method:  "delete",
		Path:    "news",
	}, &news, nil)

	return nil
}
