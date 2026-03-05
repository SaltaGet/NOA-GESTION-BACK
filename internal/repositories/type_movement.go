package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (t *TypeMovementRepository) TypeMovementGetAll(typeMovement string) ([]*schemas.TypeMovementResponse, error) {
	ctx := context.Background()
	var typeMovements []*schemas.TypeMovementResponse

	switch typeMovement {
	case "income":
		incomes, err := boilmodels.TypeIncomes(
			qm.Select(boilmodels.TypeIncomeColumns.ID, boilmodels.TypeIncomeColumns.Name),
		).All(ctx, t.DB)
		if err != nil {
			return nil, schemas.HandlerErrorDB(err, "Tipo de ingreso", schemas.Read)
		}
		for _, v := range incomes {
			typeMovements = append(typeMovements, &schemas.TypeMovementResponse{
				ID:   v.ID,
				Name: v.Name,
			})
		}
	case "expense":
		expenses, err := boilmodels.TypeExpenses(
			qm.Select(boilmodels.TypeExpenseColumns.ID, boilmodels.TypeExpenseColumns.Name),
		).All(ctx, t.DB)
		if err != nil {
			return nil, schemas.HandlerErrorDB(err, "Tipo de egreso", schemas.Read)
		}
		for _, v := range expenses {
			typeMovements = append(typeMovements, &schemas.TypeMovementResponse{
				ID:   v.ID,
				Name: v.Name,
			})
		}
	default:
		return nil, schemas.ErrorResponse(400, "tipo de movimiento no válido", fmt.Errorf("tipo de movimiento no valido: %s", typeMovement))
	}

	return typeMovements, nil
}

func (t *TypeMovementRepository) TypeMovementCreate(memberID int64, movementCreate schemas.TypeMovementCreate) error {
	ctx := context.Background()
	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	switch movementCreate.TypeMovement {
	case "income":
		typeIncome := &boilmodels.TypeIncome{Name: movementCreate.Name}
		if err := typeIncome.Insert(ctx, tx, boil.Infer()); err != nil {
			return schemas.HandlerErrorDB(err, "Tipo de ingreso", schemas.Create)
		}
	case "expense":
		typeExpense := &boilmodels.TypeExpense{Name: movementCreate.Name}
		if err := typeExpense.Insert(ctx, tx, boil.Infer()); err != nil {
			return schemas.HandlerErrorDB(err, "Tipo de egreso", schemas.Create)
		}
	default:
		return schemas.ErrorResponse(400, "tipo de movimiento no válido", fmt.Errorf("tipo de movimiento no valido: %s", movementCreate.TypeMovement))
	}

	return tx.Commit()
}

func (t *TypeMovementRepository) TypeMovementUpdate(memberID int64, movementUpdate schemas.TypeMovementUpdate) error {
	ctx := context.Background()
	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	switch movementUpdate.TypeMovement {
	case "income":
		oldIncome, err := boilmodels.TypeIncomes(boilmodels.TypeIncomeWhere.ID.EQ(movementUpdate.ID)).One(ctx, tx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return schemas.ErrorResponse(404, "tipo de ingreso no encontrado", fmt.Errorf("id %d no encontrado", movementUpdate.ID))
			}
			return schemas.HandlerErrorDB(err, "Tipo de ingreso", schemas.Read)
		}

		oldIncome.Name = movementUpdate.Name
		if _, err := oldIncome.Update(ctx, tx, boil.Infer()); err != nil {
			return schemas.HandlerErrorDB(err, "Tipo de ingreso", schemas.Update)
		}

	case "expense":
		oldExpense, err := boilmodels.TypeExpenses(boilmodels.TypeExpenseWhere.ID.EQ(movementUpdate.ID)).One(ctx, tx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return schemas.ErrorResponse(404, "tipo de egreso no encontrado", fmt.Errorf("id %d no encontrado", movementUpdate.ID))
			}
			return schemas.HandlerErrorDB(err, "Tipo de egreso", schemas.Read)
		}

		oldExpense.Name = movementUpdate.Name
		if _, err := oldExpense.Update(ctx, tx, boil.Infer()); err != nil {
			return schemas.HandlerErrorDB(err, "Tipo de egreso", schemas.Update)
		}

	default:
		return schemas.ErrorResponse(400, "tipo de movimiento no válido", fmt.Errorf("tipo de movimiento no válido: %s", movementUpdate.TypeMovement))
	}

	return tx.Commit()
}
