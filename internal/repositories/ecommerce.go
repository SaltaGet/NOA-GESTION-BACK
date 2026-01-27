package repositories

import (
	"errors"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func (er *EcommerceRepository) GetByID(id int64) (*schemas.EcommerceResponse, error) {
	var ecommerce models.IncomeEcommerce
	if err := er.DB.Preload("Items").Preload("Items.Product").First(&ecommerce, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "compra electrónica no encontrada", err)
		}
		return nil, schemas.ErrorResponse(500, "error al obtener la compra electrónica", err)
	}

	var ecommerceResponse schemas.EcommerceResponse
	if err := copier.Copy(&ecommerceResponse, &ecommerce); err != nil {
      return nil, err
  }

  return &ecommerceResponse, nil
}

func (er *EcommerceRepository) GetByReference(reference string) (*schemas.EcommerceResponse, error) {
	var ecommerce models.IncomeEcommerce
	if err := er.DB.Preload("Items").Preload("Items.Product").Where("external_reference = ?", reference).First(&ecommerce).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "compra electrónica no encontrada", err)
		}
		return nil, schemas.ErrorResponse(500, "error al obtener la compra electrónica", err)
	}

	var ecommerceResponse schemas.EcommerceResponse
	if err := copier.Copy(&ecommerceResponse, &ecommerce); err != nil {
      return nil, err
  }

	return &ecommerceResponse, nil
}

func (er *EcommerceRepository) GetAll(page, limit int, status *string) ([]schemas.EcommerceResponseDTO, error) {
	offset := (page - 1) * limit
	var ecommerces []models.IncomeEcommerce
	query := er.DB.Offset(offset).Limit(limit)
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Select("id", "external_reference", "status", "total", "date_created", "payer_email").Find(&ecommerces).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "error al obtener las compras electrónicas", err)
	}

	var ecommerceResponses []schemas.EcommerceResponseDTO
	copier.Copy(&ecommerceResponses, &ecommerces)

	return ecommerceResponses, nil
}

func (er *EcommerceRepository) UpdateStatus(update *schemas.EcommerceStatusUpdate) error {
	result := er.DB.Model(&models.IncomeEcommerce{}).Where("id = ?", update.ID).Update("status", update.NewStatus)
	if result.Error != nil {
		return schemas.ErrorResponse(500, "error al actualizar el estado de la compra electrónica", result.Error)
	}
	if result.RowsAffected == 0 {
		return schemas.ErrorResponse(404, "compra electrónica no encontrada", nil)
	}

	return nil
}