package repositories

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
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
	if err := r.DB.Create(&newNews).Error; err != nil {
		return 0, schemas.HandlerErrorGorm(err, "Noticia", schemas.Create)
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
		return schemas.HandlerErrorGorm(err, "Noticia", schemas.Read)
	}

	oldNews := news

	news.Title = newsUpdate.Title
	news.Content = newsUpdate.Content
	if err := r.DB.Save(&news).Error; err != nil {
		return schemas.HandlerErrorGorm(err, "Noticia", schemas.Update)
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
		return schemas.HandlerErrorGorm(err, "Noticia", schemas.Read)
	}

	if err := r.DB.Delete(&news).Error; err != nil {
		return schemas.HandlerErrorGorm(err, "Noticia", schemas.Delete)
	}

	go database.SaveAuditAdminAsync(r.DB, models.AuditLogAdmin{
		AdminID: adminID,
		Method:  "delete",
		Path:    "news",
	}, &news, nil)

	return nil
}
