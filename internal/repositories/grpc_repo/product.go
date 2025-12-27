package grpc_repo

import (
	"errors"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"gorm.io/gorm"
)

func (r *GrpcProductRepository) ProductGetByCode(code string) (*models.Product, error) {
	var product models.Product
	
	err := r.DB.Preload("StockDeposit").Preload("Category").Where("code = ?", code).First(&product).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("producto no encontrado")
		}
		return nil, err
	}
	
	return &product, nil
}

// ProductList obtiene productos con paginación, filtros y ordenamiento
func (r *GrpcProductRepository) ProductList(
	page, pageSize int32,
	categoryID *int32,
	search *string,
	sortBy int32,
) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	query := r.DB.Preload("StockDeposit").Preload("Category").Model(&models.Product{})

	// Filtro por categoría
	// if categoryID != nil {
	// 	query = query.Where("category_id = ?", *categoryID)
	// }

	// Búsqueda por nombre o código
	if search != nil && *search != "" {
		searchPattern := "%" + *search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", searchPattern, searchPattern)
	}

	// Contar total antes de paginar
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Ordenamiento
	switch sortBy {
	case 0: // PRICE_LOW_TO_HIGH
		query = query.Order("price ASC")
	case 1: // PRICE_HIGH_TO_LOW
		query = query.Order("price DESC")
	case 2: // NAME_A_Z
		query = query.Order("name ASC")
	case 3: // NAME_Z_A
		query = query.Order("name DESC")
	default:
		query = query.Order("id DESC") // Por defecto
	}

	// Paginación
	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}