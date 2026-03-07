package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

func mapToClientResponseDTO(c *boilmodels.Client, debt *float64) schemas.ClientResponseDTO {
	dto := schemas.ClientResponseDTO{
		ID:        c.ID,
		FirstName: c.FirstName,
		LastName:  c.LastName,
	}
	if c.CompanyName.Valid {
		dto.CompanyName = &c.CompanyName.String
	}
	if c.Identifier.Valid {
		dto.Identifier = &c.Identifier.String
	}
	if c.Email.Valid {
		dto.Email = &c.Email.String
	}
	if c.Phone.Valid {
		dto.Phone = &c.Phone.String
	}
	if c.ResponsabilityFrontIva.Valid {
		dto.ResponsabilityFrontIVA = &c.ResponsabilityFrontIva.String
	}
	if debt != nil {
		dto.Debt = debt
	}
	return dto
}

func (r *ClientRepository) ClientGetByID(id int64) (*schemas.ClientResponse, error) {
	ctx := context.Background()

	c, err := boilmodels.Clients(
		boilmodels.ClientWhere.ID.EQ(id),
		boilmodels.ClientWhere.DeletedAt.IsNull(),
		qm.Load(boilmodels.ClientRels.MemberCreate),
		qm.Load(boilmodels.ClientRels.PayIncomes, boilmodels.PayIncomeWhere.MethodPay.EQ("credit")),
	).One(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Cliente", schemas.Read)
	}

	res := &schemas.ClientResponse{
		ID:        c.ID,
		FirstName: c.FirstName,
		LastName:  c.LastName,
	}
	if c.CompanyName.Valid {
		res.CompanyName = &c.CompanyName.String
	}
	if c.Identifier.Valid {
		res.Identifier = &c.Identifier.String
	}
	if c.Email.Valid {
		res.Email = &c.Email.String
	}
	if c.Phone.Valid {
		res.Phone = &c.Phone.String
	}
	if c.Address.Valid {
		res.Address = &c.Address.String
	}
	if c.ResponsabilityFrontIva.Valid {
		res.ResponsabilityFrontIVA = &c.ResponsabilityFrontIva.String
	}

	if c.R != nil {
		if c.R.MemberCreate != nil {
			res.MemberCreate = &schemas.MemberSimpleDTO{
				ID:        c.R.MemberCreate.ID,
				FirstName: c.R.MemberCreate.FirstName,
				LastName:  c.R.MemberCreate.LastName,
				Username:  c.R.MemberCreate.Username,
			}
		}

		if len(c.R.PayIncomes) > 0 {
			res.Pay = make([]schemas.PayDebtResponse, len(c.R.PayIncomes))
			for i, pay := range c.R.PayIncomes {
				val, _ := pay.Total.Big.Float64()
				var incomeSaleID int64
				incomeSaleID = pay.IncomeSaleID
				res.Pay[i] = schemas.PayDebtResponse{
					ID:           pay.ID,
					IncomeSaleID: incomeSaleID,
					Total:        val,
					MethodPay:    pay.MethodPay,
				}
				res.Pay[i].CreatedAt = pay.CreatedAt
			}
		}
	}

	return res, nil
}

func (r *ClientRepository) ClientGetByFilter(search string) (*[]schemas.ClientResponseDTO, error) {
	ctx := context.Background()
	searchStr := "%" + search + "%"

	boilClients, err := boilmodels.Clients(
		boilmodels.ClientWhere.DeletedAt.IsNull(),
		qm.Where("last_name ILIKE ? OR first_name ILIKE ? OR identifier ILIKE ?", searchStr, searchStr, searchStr),
		qm.Limit(10),
	).All(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Cliente", schemas.Read)
	}

	dtos := make([]schemas.ClientResponseDTO, 0, len(boilClients))
	for _, c := range boilClients {
		dtos = append(dtos, mapToClientResponseDTO(c, nil))
	}

	return &dtos, nil
}

type ClientWithDebt struct {
	boilmodels.Client `boil:",bind"`
	Debt              float64 `boil:"debt"`
}

func (r *ClientRepository) ClientGetAll(limit, page int64, search *map[string]string, filterDrbt bool) (*[]schemas.ClientResponseDTO, int64, error) {
	ctx := context.Background()

	var qms []qm.QueryMod
	qms = append(qms, boilmodels.ClientWhere.DeletedAt.IsNull())

	if filterDrbt {
		debtFormula := "COALESCE(SUM(CASE WHEN p.method_pay = 'credit' THEN p.total ELSE 0 END), 0)"
		qms = append(qms,
			qm.Select(fmt.Sprintf("clients.*, %s AS debt", debtFormula)),
			qm.LeftOuterJoin("pay_incomes p ON p.client_id = clients.id AND p.delete_at IS NULL"),
			qm.GroupBy("clients.id"),
			qm.Having(fmt.Sprintf("%s > 0", debtFormula)),
		)
	} else {
		qms = append(qms, qm.Select("clients.*"))
	}

	if search != nil {
		for key, value := range *search {
			if value == "" {
				continue
			}
			switch strings.ToLower(key) {
			case "identifier":
				qms = append(qms, qm.Where("clients.identifier ILIKE ?", "%"+value+"%"))
			case "first_name":
				qms = append(qms, qm.Where("clients.first_name ILIKE ?", "%"+value+"%"))
			case "last_name":
				qms = append(qms, qm.Where("clients.last_name ILIKE ?", "%"+value+"%"))
			case "email":
				qms = append(qms, qm.Where("clients.email ILIKE ?", "%"+value+"%"))
			}
		}
	}

	countQueryMods := make([]qm.QueryMod, len(qms))
	copy(countQueryMods, qms)

	total, err := boilmodels.Clients(countQueryMods...).Count(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Cliente", schemas.Read)
	}

	if limit > 0 {
		offset := (page - 1) * limit
		qms = append(qms, qm.Limit(int(limit)), qm.Offset(int(offset)))
	}

	var results []ClientWithDebt
	if filterDrbt {
		err := boilmodels.Clients(qms...).Bind(ctx, r.DB, &results)
		if err != nil {
			return nil, 0, schemas.HandlerErrorDB(err, "Cliente", schemas.Read)
		}
	} else {
		boilClients, err := boilmodels.Clients(qms...).All(ctx, r.DB)
		if err != nil {
			return nil, 0, schemas.HandlerErrorDB(err, "Cliente", schemas.Read)
		}
		for _, c := range boilClients {
			results = append(results, ClientWithDebt{Client: *c})
		}
	}

	dtos := make([]schemas.ClientResponseDTO, 0, len(results))
	for _, res := range results {
		var debt *float64
		if filterDrbt {
			d := res.Debt
			debt = &d
		}
		dtos = append(dtos, mapToClientResponseDTO(&res.Client, debt))
	}

	return &dtos, total, nil
}

func (r *ClientRepository) ClientCreate(memberID int64, client *schemas.ClientCreate) (int64, error) {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Cliente", schemas.Create)
	}

	newClient := boilmodels.Client{
		FirstName:              client.FirstName,
		LastName:               client.LastName,
		CompanyName:            null.StringFromPtr(client.CompanyName),
		Identifier:             null.StringFromPtr(client.Identifier),
		Email:                  null.StringFromPtr(client.Email),
		Phone:                  null.StringFromPtr(client.Phone),
		Address:                null.StringFromPtr(client.Address),
		ResponsabilityFrontIva: null.StringFromPtr(client.ResponsabilityFrontIVA),
		MemberCreateID:         memberID,
	}

	if err := newClient.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Cliente", schemas.Create)
	}

	if err := tx.Commit(); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Cliente", schemas.Create)
	}

	return newClient.ID, nil
}

func (r *ClientRepository) ClientUpdate(memberID int64, client *schemas.ClientUpdate) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	oldClient, err := boilmodels.FindClient(ctx, tx, client.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schemas.ErrorResponse(404, "Cliente no encontrado", fmt.Errorf("cliente no encontrado"))
		}
		return schemas.HandlerErrorDB(err, "Cliente", schemas.Read)
	}

	oldClient.FirstName = client.FirstName
	oldClient.LastName = client.LastName
	oldClient.CompanyName = null.StringFromPtr(client.CompanyName)
	oldClient.Identifier = null.StringFromPtr(client.Identifier)
	oldClient.Email = null.StringFromPtr(client.Email)
	oldClient.Phone = null.StringFromPtr(client.Phone)
	oldClient.Address = null.StringFromPtr(client.Address)
	oldClient.ResponsabilityFrontIva = null.StringFromPtr(client.ResponsabilityFrontIVA)

	if _, err := oldClient.Update(ctx, tx, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Cliente", schemas.Update)
	}

	if err := tx.Commit(); err != nil {
		return schemas.HandlerErrorDB(err, "Cliente", schemas.Update)
	}

	return nil
}

func (r *ClientRepository) ClientDelete(memberID, id int64) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	client, err := boilmodels.FindClient(ctx, tx, id)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Cliente", schemas.Read)
	}

	// SQLBoiler with delete_at might do a hard delete?
	// We'll manually soft delete just in case
	client.DeletedAt = null.TimeFrom(time.Now())
	if _, err := client.Update(ctx, tx, boil.Whitelist(boilmodels.ClientColumns.DeletedAt, boilmodels.ClientColumns.UpdatedAt)); err != nil {
		return schemas.HandlerErrorDB(err, "Cliente", schemas.Delete)
	}

	if err := tx.Commit(); err != nil {
		return schemas.HandlerErrorDB(err, "Cliente", schemas.Delete)
	}

	return nil
}

func (r *ClientRepository) ClientUpdateCredit(memberID, pointSaleID int64, clientUpdateCredit *schemas.ClientUpdateCredit) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	register, err := boilmodels.CashRegisters(
		boilmodels.CashRegisterWhere.IsClose.EQ(false),
		boilmodels.CashRegisterWhere.PointSaleID.EQ(pointSaleID),
		qm.OrderBy("hour_open DESC"),
	).One(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Caja", schemas.Read)
	}

	client, err := boilmodels.Clients(
		qm.Select(boilmodels.ClientColumns.ID),
		boilmodels.ClientWhere.ID.EQ(clientUpdateCredit.ID),
	).One(ctx, tx)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Cliente", schemas.Read)
	}

	total := 0.0

	for _, p := range clientUpdateCredit.PayCredit {
		payCredit, err := boilmodels.PayIncomes(
			boilmodels.PayIncomeWhere.ClientID.EQ(null.Int64From(client.ID)),
			boilmodels.PayIncomeWhere.ID.EQ(p.CreditID),
		).One(ctx, tx)

		if err != nil {
			return schemas.HandlerErrorDB(err, "Credito", schemas.Read)
		}

		payCredit.MethodPay = p.MethodPay
		payCredit.CashRegisterID = null.Int64From(register.ID)

		if _, err := payCredit.Update(ctx, tx, boil.Whitelist(boilmodels.PayIncomeColumns.MethodPay, boilmodels.PayIncomeColumns.CashRegisterID, boilmodels.PayIncomeColumns.UpdatedAt)); err != nil {
			return schemas.HandlerErrorDB(err, "Credito", schemas.Update)
		}

		val, _ := payCredit.Total.Big.Float64()
		total += val
	}

	if math.Abs(total-clientUpdateCredit.Total) > 1 {
		message := fmt.Sprintf("la diferencia entre la suma de pagos (%.2f) y el total del cliente (%.2f) no puede ser mayor que 1", total, clientUpdateCredit.Total)
		return schemas.ErrorResponse(400, message, fmt.Errorf("%s", message))
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
