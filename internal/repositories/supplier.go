package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
)

func mapToSupplierResponse(s *boilmodels.Supplier) *schemas.SupplierResponse {
	if s == nil {
		return nil
	}

	res := &schemas.SupplierResponse{
		ID:          s.ID,
		Name:        s.Name,
		CompanyName: s.CompanyName.String,
	}

	if s.Identifier.Valid {
		res.Identifier = &s.Identifier.String
	}
	if s.Address.Valid {
		res.Address = &s.Address.String
	}
	if !s.DebtLimit.IsZero() {
		limit, _ := s.DebtLimit.Big.Float64()
		res.DebtLimit = &limit
	}
	if s.Email.Valid {
		res.Email = &s.Email.String
	}
	if s.Phone.Valid {
		res.Phone = &s.Phone.String
	}
	if s.CreatedAt.Valid {
		res.CreatedAt = s.CreatedAt.Time
	}

	return res
}

func mapToSupplierResponseDTO(s *boilmodels.Supplier) *schemas.SupplierResponseDTO {
	if s == nil {
		return nil
	}

	return &schemas.SupplierResponseDTO{
		ID:          s.ID,
		Name:        s.Name,
		CompanyName: s.CompanyName.String,
	}
}

// SupplierGetByID obtiene un proveedor por ID
func (r *SupplierRepository) SupplierGetByID(id int64) (*schemas.SupplierResponse, error) {
	ctx := context.Background()

	supplier, err := boilmodels.Suppliers(
		boilmodels.SupplierWhere.ID.EQ(id),
	).One(ctx, r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Proveedor", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Proveedor", schemas.Read)
	}

	return mapToSupplierResponse(supplier), nil
}

// SupplierGetAll obtiene todos los proveedores con paginación y búsqueda
func (r *SupplierRepository) SupplierGetAll(limit, page int, search *map[string]string) ([]*schemas.SupplierResponseDTO, int64, error) {
	ctx := context.Background()
	offset := (page - 1) * limit

	var queries []qm.QueryMod

	if search != nil {
		for key, value := range *search {
			switch key {
			case "name":
				queries = append(queries, qm.Where("name ILIKE ?", "%"+value+"%"))
			case "company_name":
				queries = append(queries, qm.Where("company_name ILIKE ?", "%"+value+"%"))
			case "identifier":
				queries = append(queries, qm.Where("identifier ILIKE ?", "%"+value+"%"))
			case "email":
				queries = append(queries, qm.Where("email ILIKE ?", "%"+value+"%"))
			}
		}
	}

	total, err := boilmodels.Suppliers(queries...).Count(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Proveedor", schemas.Read)
	}

	queries = append(queries,
		qm.OrderBy("created_at DESC"),
		qm.Offset(offset),
		qm.Limit(limit),
	)

	suppliers, err := boilmodels.Suppliers(queries...).All(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Proveedor", schemas.Read)
	}

	var suppliersSchema []*schemas.SupplierResponseDTO
	for _, s := range suppliers {
		suppliersSchema = append(suppliersSchema, mapToSupplierResponseDTO(s))
	}

	return suppliersSchema, total, nil
}

// SupplierCreate crea un nuevo proveedor con auditoría
func (r *SupplierRepository) SupplierCreate(memberID int64, supplierCreate *schemas.SupplierCreate) (int64, error) {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return 0, err
	}

	supplier := &boilmodels.Supplier{
		Name:        supplierCreate.Name,
		CompanyName: null.StringFrom(supplierCreate.CompanyName),
	}

	if supplierCreate.Identifier != nil {
		supplier.Identifier = null.StringFromPtr(supplierCreate.Identifier)
	}
	if supplierCreate.Address != nil {
		supplier.Address = null.StringFromPtr(supplierCreate.Address)
	}
	if supplierCreate.DebtLimit != nil {
		supplier.DebtLimit = types.NewNullDecimal(types.NewDecimal(fmt.Sprintf("%.4f", *supplierCreate.DebtLimit)))
	}
	if supplierCreate.Email != nil {
		supplier.Email = null.StringFromPtr(supplierCreate.Email)
	}
	if supplierCreate.Phone != nil {
		supplier.Phone = null.StringFromPtr(supplierCreate.Phone)
	}

	if err := supplier.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Proveedor", schemas.Create)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return supplier.ID, nil
}

// SupplierUpdate actualiza un proveedor existente con auditoría
func (r *SupplierRepository) SupplierUpdate(memberID int64, supplierUpdate *schemas.SupplierUpdate) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	existingSupplier, err := boilmodels.Suppliers(boilmodels.SupplierWhere.ID.EQ(supplierUpdate.ID)).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schemas.HandlerErrorDB(err, "Proveedor", schemas.Read)
		}
		return schemas.HandlerErrorDB(err, "Proveedor", schemas.Read)
	}

	existingSupplier.Name = supplierUpdate.Name
	existingSupplier.CompanyName = null.StringFrom(supplierUpdate.CompanyName)

	if supplierUpdate.Identifier != nil {
		existingSupplier.Identifier = null.StringFromPtr(supplierUpdate.Identifier)
	} else {
		existingSupplier.Identifier = null.NewString("", false)
	}

	if supplierUpdate.Address != nil {
		existingSupplier.Address = null.StringFromPtr(supplierUpdate.Address)
	} else {
		existingSupplier.Address = null.NewString("", false)
	}

	if supplierUpdate.DebtLimit != nil {
		existingSupplier.DebtLimit = types.NewNullDecimal(types.NewDecimal(fmt.Sprintf("%.4f", *supplierUpdate.DebtLimit)))
	} else {
		existingSupplier.DebtLimit = types.NullDecimal{}
	}

	if supplierUpdate.Email != nil {
		existingSupplier.Email = null.StringFromPtr(supplierUpdate.Email)
	} else {
		existingSupplier.Email = null.NewString("", false)
	}

	if supplierUpdate.Phone != nil {
		existingSupplier.Phone = null.StringFromPtr(supplierUpdate.Phone)
	} else {
		existingSupplier.Phone = null.NewString("", false)
	}

	if _, err := existingSupplier.Update(ctx, tx, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Proveedor", schemas.Update)
	}

	return tx.Commit()
}

func (r *SupplierRepository) SupplierDelete(memberID int64, id int64) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	supplier, err := boilmodels.Suppliers(boilmodels.SupplierWhere.ID.EQ(id)).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schemas.HandlerErrorDB(err, "Proveedor", schemas.Read)
		}
		return schemas.HandlerErrorDB(err, "Proveedor", schemas.Read)
	}

	expenseCount, err := boilmodels.ExpenseBuys(boilmodels.ExpenseBuyWhere.SupplierID.EQ(null.Int64From(id))).Count(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Egreso de compra", schemas.Read)
	}

	if expenseCount > 0 {
		return schemas.ErrorResponse(400, "No se puede eliminar el proveedor porque tiene compras asociadas", errors.New("No se puede eliminar el proveedor porque tiene compras asociadas"))
	}

	// Soft delete
	if _, err := supplier.Delete(ctx, tx, false); err != nil {
		return schemas.HandlerErrorDB(err, "Proveedor", schemas.Delete)
	}

	return tx.Commit()
}
