package services

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

func (es *EcommerceService) GetByID(id int64) (*schemas.EcommerceResponse, error) {
	return es.EcommerceRepository.GetByID(id)
}

func (es *EcommerceService) GetByReference(reference string) (*schemas.EcommerceResponse, error) {
	return es.EcommerceRepository.GetByReference(reference)
}

func (es *EcommerceService) GetAll(page, limit int, status *string) ([]schemas.EcommerceResponseDTO, error) {
	return es.EcommerceRepository.GetAll(page, limit, status)
}

func (es *EcommerceService) UpdateStatus(update *schemas.EcommerceStatusUpdate) error {
	return es.EcommerceRepository.UpdateStatus(update)
}