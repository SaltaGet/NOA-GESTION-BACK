package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

type EcommerceRepository interface {
	GetByID(id int64) (*schemas.EcommerceResponse, error)
	GetByReference(reference string) (*schemas.EcommerceResponse, error)
	GetAll(page, limit int, status *string) ([]schemas.EcommerceResponseDTO, error)
	UpdateStatus(update *schemas.EcommerceStatusUpdate) error
}

type EcommerceService interface {
	GetByID(id int64) (*schemas.EcommerceResponse, error)
	GetByReference(reference string) (*schemas.EcommerceResponse, error)
	GetAll(page, limit int, status *string) ([]schemas.EcommerceResponseDTO, error)
	UpdateStatus(update *schemas.EcommerceStatusUpdate) error
}