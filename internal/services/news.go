package services

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

func (n *NewsService) NewsGetByID(id int64) (*schemas.NewsResponse, error) {
	return n.NewsRepository.NewsGetByID(id)
}

func (n *NewsService) NewsGetAll() ([]schemas.NewsResponseDTO, error) {
	return n.NewsRepository.NewsGetAll()
}

func (n *NewsService) NewsCreate(adminID int64, news *schemas.NewsCreate) (int64, error) {
	return n.NewsRepository.NewsCreate(adminID, news)
}

func (n *NewsService) NewsUpdate(adminID int64, news *schemas.NewsUpdate) (error) {
	return n.NewsRepository.NewsUpdate(adminID, news)
}

func (n *NewsService) NewsDelete(adminID int64, id int64) (error) {
	return n.NewsRepository.NewsDelete(adminID, id)
}
