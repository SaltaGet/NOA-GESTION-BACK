package repositories

import (
	"context"
	"database/sql"
	"errors"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (r *MainRepository) UserGetByID(id int64) (*boilmodels.User, error) {
	ctx := context.Background()

	user, err := boilmodels.Users(boilmodels.UserWhere.ID.EQ(id)).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.ErrorResponse(404, "Usuario no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al buscar usuario", err)
	}

	return user, nil
}

func (r *MainRepository) UserGetByListID(ids []int64) (*[]schemas.UserDTO, error) {
	ctx := context.Background()

	var idInterfaces []interface{}
	for _, id := range ids {
		idInterfaces = append(idInterfaces, id)
	}

	users, err := boilmodels.Users(
		qm.Select(
			boilmodels.UserColumns.ID,
			boilmodels.UserColumns.FirstName,
			boilmodels.UserColumns.LastName,
			boilmodels.UserColumns.Username,
			boilmodels.UserColumns.Email,
		),
		qm.WhereIn("id IN ?", idInterfaces...),
	).All(ctx, r.DB)

	if err != nil {
		return nil, err
	}

	var response []schemas.UserDTO
	for _, user := range users {
		response = append(response, schemas.UserDTO{
			ID:        user.ID,
			FirstName: user.FirstName.String,
			LastName:  user.LastName.String,
			Username:  user.Username,
			Email:     user.Email,
		})
	}

	return &response, nil
}

func (r *MainRepository) UserGetByUsername(username string) (*boilmodels.User, error) {
	ctx := context.Background()

	user, err := boilmodels.Users(boilmodels.UserWhere.Username.EQ(username)).One(ctx, r.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, schemas.ErrorResponse(500, "Error interno al buscar usuario", err)
	}

	return user, nil
}

func (r *MainRepository) UserGetExistByUsernameEmail(username string, email string) (bool, error) {
	ctx := context.Background()

	exists, err := boilmodels.Users(
		qm.Where("email = ? OR username = ?", email, username),
	).Exists(ctx, r.DB)

	if err != nil {
		return false, schemas.ErrorResponse(500, "Error interno al buscar usuario", err)
	}

	return exists, nil
}

func (r *MainRepository) UserCreate(user *schemas.UserCreate) (int64, error) {
	ctx := context.Background()

	newUser := &boilmodels.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password, // Recordar usar hash en service!
	}

	if err := newUser.Insert(ctx, r.DB, boil.Infer()); err != nil {
		return 0, schemas.ErrorResponse(500, "Error interno al crear usuario", err)
	}

	return newUser.ID, nil
}

func (m *MainRepository) UserTenantAdd(userID, tenantID int64) error {
	ctx := context.Background()

	userTenant := &boilmodels.UserTenant{
		UserID:   userID,
		TenantID: tenantID,
	}

	if err := userTenant.Insert(ctx, m.DB, boil.Infer()); err != nil {
		return schemas.ErrorResponse(500, "Error interno al agregar usuario a tenant", err)
	}

	return nil
}

func (m *MainRepository) UserUpdate(userID int64, req *schemas.UserUpdate) error {
	ctx := context.Background()
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	user, err := boilmodels.Users(boilmodels.UserWhere.ID.EQ(userID)).One(ctx, tx)
	if err != nil {
		return schemas.ErrorResponse(404, "Usuario no encontrado", err)
	}

	user.FirstName = null.StringFrom(req.FirstName)
	user.LastName = null.StringFrom(req.LastName)
	user.Username = req.Username
	user.Email = req.Email

	if _, err := user.Update(ctx, tx, boil.Infer()); err != nil {
		return schemas.ErrorResponse(500, "Error interno al actualizar usuario", err)
	}

	return tx.Commit()
}
