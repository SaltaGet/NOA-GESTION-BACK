package grpc_repo

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// func (r *GrpcCategoryRepository) CategoryGetAll() ([]*models.Category, error) {
// 	var categories []*models.Category

// 	err := r.DB.Find(&categories).Error
// 	if err != nil {
// 		return nil,  status.Errorf(codes.Internal, "error de base de datos: %v", err)
// 	}

// 	return categories, nil
// }

func (r *GrpcCategoryRepository) CategoryGetAll() ([]*models.Category, error) {
	var categories []*models.Category

	err := r.DB.
		// 1. Seleccionamos DISTINCT para evitar duplicados si hay múltiples productos válidos
		// Esto asegura que cada categoría venga una sola vez.
		Distinct("categories.*").

		// 2. Unimos con la tabla de productos
		Joins("JOIN products ON products.category_id = categories.id").

		// 3. Unimos la tabla de productos con la de depósitos (Stock)
		Joins("JOIN deposits ON deposits.product_id = products.id").

		// 4. Aplicamos los filtros requeridos
		Where("products.is_visible = ?", true).
		Where("deposits.stock > ?", 0).

		// Ejecutamos la búsqueda
		Find(&categories).Error

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error de base de datos: %v", err)
	}

	return categories, nil
}