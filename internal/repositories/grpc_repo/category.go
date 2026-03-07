package grpc_repo

import (
	"context"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *GrpcCategoryRepository) CategoryGetAll() ([]*tenant.Category, error) {
	ctx := context.Background()
	// SQLBoiler usa Query Mods (qm) para construir la consulta
	categories, err := tenant.Categories(
		// 1. SELECT DISTINCT para evitar duplicados por el JOIN
		qm.Select("DISTINCT categories.*"),
		
		// 2 y 3. Joins manuales (SQLBoiler requiere especificar la relación)
		qm.InnerJoin("products ON products.category_id = categories.id"),
		qm.InnerJoin("deposits ON deposits.product_id = products.id"),
		
		// 4. Filtros de visibilidad y stock
		qm.Where("products.is_visible = ?", true),
		qm.And("deposits.stock > ?", 0),
	).All(ctx, r.DB)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error de base de datos: %v", err)
	}

	return categories, nil
}

// func (r *GrpcCategoryRepository) CategoryGetAll() ([]*tenant.Category, error) {
// 	var categories []*tenant.Category

// 	err := r.DB.
// 		// 1. Seleccionamos DISTINCT para evitar duplicados si hay múltiples productos válidos
// 		// Esto asegura que cada categoría venga una sola vez.
// 		Distinct("categories.*").

// 		// 2. Unimos con la tabla de productos
// 		Joins("JOIN products ON products.category_id = categories.id").

// 		// 3. Unimos la tabla de productos con la de depósitos (Stock)
// 		Joins("JOIN deposits ON deposits.product_id = products.id").

// 		// 4. Aplicamos los filtros requeridos
// 		Where("products.is_visible = ?", true).
// 		Where("deposits.stock > ?", 0).

// 		// Ejecutamos la búsqueda
// 		Find(&categories).Error

// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "error de base de datos: %v", err)
// 	}

// 	return categories, nil
// }