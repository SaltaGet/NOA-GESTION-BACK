package grpc_repo

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *GrpcProductRepository) ProductGetByID(id int64) (*models.Product, error) {
	var product models.Product

	err := r.DB.Preload("StockDeposit").Preload("Category").Where("id = ?", id).First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("producto no encontrado")
		}
		return nil, err
	}

	return &product, nil
}

func (r *GrpcProductRepository) ProductGetByCode(code string) (*models.Product, error) {
	var product models.Product

	err := r.DB.Preload("StockDeposit").
		Preload("Category").
		Where("code = ? AND is_visible = ?", code, true).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "producto no encontrado")
		}
		return nil, err
	}

	return &product, nil
}

// ProductList obtiene productos con paginación, filtros y ordenamiento
func (r *GrpcProductRepository) ProductList(req *pb.ListProductsRequest) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	query := r.DB.Preload("StockDeposit").Preload("Category").Model(&models.Product{})

	if req.CategoryId != nil {
		query = query.Where("category_id = ?", *req.CategoryId)
	}

	if req.Search != nil {
		searchPattern := "%" + *req.Search + "%"
		query = query.Where("name LIKE ?", searchPattern)
	}

	query = query.Where("is_visible = ?", true)

	// Contar total antes de paginar
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if req.Sort != nil {
		sortValue := *req.Sort
		switch sortValue {
		case pb.ListProductsRequest_PRICE_LOW_TO_HIGH: // PRICE_LOW_TO_HIGH
			query = query.Order("price ASC")
		case pb.ListProductsRequest_PRICE_HIGH_TO_LOW: // PRICE_HIGH_TO_LOW
			query = query.Order("price DESC")
		case pb.ListProductsRequest_NAME_A_Z: // NAME_A_Z
			query = query.Order("name ASC")
		case pb.ListProductsRequest_NAME_Z_A: // NAME_Z_A
			query = query.Order("name DESC")
		default:
			query = query.Order("id DESC") // Por defecto
		}
	}

	// Paginación
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(int(offset)).Limit(int(req.Limit)).Find(&products).Error; err != nil {
		return nil, 0, status.Error(codes.Internal, "error interno:"+err.Error())
	}

	return products, total, nil
}

func (r *GrpcProductRepository) SaveUrlImage(req *pb.SaveImageRequest) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var prodExist models.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&prodExist, req.ProdId).Error; err != nil {
			return status.Error(codes.NotFound, "El producto no existe")
		}

		// 1. Actualizar Imagen Principal
		if req.PrimaryImage != nil {
			prodExist.PrimaryImage = req.PrimaryImage
		}

		// 2. Lógica de Imágenes Secundarias (Delta Logic)
		// Obtenemos las actuales de forma segura
		var currentImages []string
		if prodExist.SecondaryImages != nil && *prodExist.SecondaryImages != "" {
			currentImages = strings.Split(*prodExist.SecondaryImages, ",")
		}

		// Creamos la nueva lista empezando por las que el cliente quiere MANTENER
		var updatedList []string
		for _, keep := range req.KeepSecondaries {
			if slices.Contains(currentImages, keep) {
				updatedList = append(updatedList, keep)
			}
		}

		// Agregamos las NUEVAS imágenes subidas
		if len(req.SecondaryImages) > 0 {
			updatedList = append(updatedList, req.SecondaryImages...)
		}

		// 3. Guardar cambios en el campo de texto
		if len(updatedList) == 0 {
			prodExist.SecondaryImages = nil
		} else {
			finalString := strings.Join(updatedList, ",")
			prodExist.SecondaryImages = &finalString
		}

		if err := tx.Save(&prodExist).Error; err != nil {
			return status.Errorf(codes.Internal, "error de base de datos: %v", err)
		}

		return nil
	})
}

func (r *GrpcProductRepository) ValidateProducts(ctx context.Context, req *pb.ProductValidateRequest) ([]models.Product, error) {
	var products []models.Product

	err := r.DB.Preload("StockDeposit", func(db *gorm.DB) *gorm.DB {
		return db.Select("product_id", "stock")
	}).Select("id", "price").
		Where("id IN ?", req.ProductIds).
		Find(&products).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "Error al validar productos")
	}

	for i := range products {
		if products[i].StockDeposit == nil {
			products[i].StockDeposit = &models.Deposit{
				ProductID: products[i].ID,
				Stock:     0,
			}
		}
	}

	return products, nil
}
