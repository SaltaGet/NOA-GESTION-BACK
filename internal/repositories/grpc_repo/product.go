package grpc_repo

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strings"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *GrpcProductRepository) ProductGetByID(id int64) (*tenant.Product, error) {
	product, err := tenant.Products(
		qm.Where("id = ?", id),
		qm.Load(tenant.ProductRels.Deposits),
		qm.Load(tenant.ProductRels.Category),
	).One(context.Background(), r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("producto no encontrado")
		}
		return nil, err
	}

	return product, nil
}

func (r *GrpcProductRepository) ProductGetByCode(code string) (*tenant.Product, error) {
	product, err := tenant.Products(
		qm.Where("code = ? AND is_visible = ?", code, true),
		qm.Load(tenant.ProductRels.Deposits),
		qm.Load(tenant.ProductRels.Category),
	).One(context.Background(), r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "producto no encontrado")
		}
		return nil, err
	}

	return product, nil
}

// ProductList obtiene productos con paginación, filtros y ordenamiento
func (r *GrpcProductRepository) ProductList(req *pb.ListProductsRequest) ([]*tenant.Product, int64, error) {
	var mods []qm.QueryMod

	if req.CategoryId != nil {
		mods = append(mods, qm.Where("category_id = ?", *req.CategoryId))
	}

	if req.Search != nil {
		searchPattern := "%" + *req.Search + "%"
		mods = append(mods, qm.Where("name LIKE ?", searchPattern))
	}

	mods = append(mods, qm.Where("is_visible = ?", true))

	// Contar total antes de paginar
	total, err := tenant.Products(mods...).Count(context.Background(), r.DB)
	if err != nil {
		return nil, 0, err
	}

	if req.Sort != nil {
		sortValue := *req.Sort
		switch sortValue {
		case pb.ListProductsRequest_PRICE_LOW_TO_HIGH:
			mods = append(mods, qm.OrderBy("price ASC"))
		case pb.ListProductsRequest_PRICE_HIGH_TO_LOW:
			mods = append(mods, qm.OrderBy("price DESC"))
		case pb.ListProductsRequest_NAME_A_Z:
			mods = append(mods, qm.OrderBy("name ASC"))
		case pb.ListProductsRequest_NAME_Z_A:
			mods = append(mods, qm.OrderBy("name DESC"))
		default:
			mods = append(mods, qm.OrderBy("id DESC"))
		}
	} else {
		mods = append(mods, qm.OrderBy("id DESC"))
	}

	// Paginación
	offset := (req.Page - 1) * req.Limit
	mods = append(mods, qm.Offset(int(offset)), qm.Limit(int(req.Limit)), qm.Load(tenant.ProductRels.Deposits), qm.Load(tenant.ProductRels.Category))

	products, err := tenant.Products(mods...).All(context.Background(), r.DB)
	if err != nil {
		return nil, 0, status.Error(codes.Internal, "error interno:"+err.Error())
	}

	return products, total, nil
}

func (r *GrpcProductRepository) SaveUrlImage(req *pb.SaveImageRequest) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return status.Errorf(codes.Internal, "no se pudo iniciar transacción: %v", err)
	}

	defer tx.Rollback()

	prodExist, err := tenant.Products(qm.Where("id = ?", req.ProdId), qm.For("UPDATE")).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return status.Error(codes.NotFound, "El producto no existe")
		}
		return err
	}

	var updateCols []string

	// 1. Actualizar Imagen Principal
	if req.PrimaryImage != nil {
		prodExist.PrimaryImage = null.StringFrom(*req.PrimaryImage)
		updateCols = append(updateCols, tenant.ProductColumns.PrimaryImage)
	}

	// 2. Lógica de Imágenes Secundarias
	var currentImages []string
	if prodExist.SecondaryImages.Valid && prodExist.SecondaryImages.String != "" {
		currentImages = strings.Split(prodExist.SecondaryImages.String, ",")
	}

	var updatedList []string
	for _, keep := range req.KeepSecondaries {
		if slices.Contains(currentImages, keep) {
			updatedList = append(updatedList, keep)
		}
	}

	if len(req.SecondaryImages) > 0 {
		updatedList = append(updatedList, req.SecondaryImages...)
	}

	if len(updatedList) == 0 {
		prodExist.SecondaryImages = null.NewString("", false)
	} else {
		finalString := strings.Join(updatedList, ",")
		prodExist.SecondaryImages = null.StringFrom(finalString)
	}

	updateCols = append(updateCols, tenant.ProductColumns.SecondaryImages)

	if _, err := prodExist.Update(ctx, tx, boil.Whitelist(updateCols...)); err != nil {
		return status.Errorf(codes.Internal, "error de base de datos: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return status.Errorf(codes.Internal, "error confirmando transacción: %v", err)
	}

	return nil
}

func (r *GrpcProductRepository) ValidateProducts(ctx context.Context, req *pb.ProductValidateRequest) ([]tenant.Product, error) {
	var products []tenant.Product

	if len(req.ProductIds) == 0 {
		return products, nil
	}

	var args []interface{}
	for _, id := range req.ProductIds {
		args = append(args, id)
	}

	dbProducts, err := tenant.Products(
		qm.Select("id", "price"),
		qm.WhereIn("id IN ?", args...),
		qm.Load(tenant.ProductRels.Deposits, qm.Select("product_id", "stock")),
	).All(ctx, r.DB)

	if err != nil {
		return nil, status.Error(codes.Internal, "Error al validar productos")
	}

	for _, p := range dbProducts {
		if len(p.R.Deposits) == 0 {
			p.R.Deposits = append(p.R.Deposits, &tenant.Deposit{
				ProductID: p.ID,
			})
		}
		products = append(products, *p)
	}

	return products, nil
}
