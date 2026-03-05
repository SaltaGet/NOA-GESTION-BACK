package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/utils"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func mapToRoleResponse(r *boilmodels.Role) schemas.RoleResponse {
	if r == nil {
		return schemas.RoleResponse{}
	}
	res := schemas.RoleResponse{
		ID:   r.ID,
		Name: r.Name,
	}
	if r.R != nil {
		for _, p := range r.R.Permissions {
			res.Permissions = append(res.Permissions, schemas.PermissionResponse{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description.String,
			})
		}
	}
	return res
}

func mapToRoleResponseDTO(r *boilmodels.Role) schemas.RoleResponseDTO {
	if r == nil {
		return schemas.RoleResponseDTO{}
	}
	return schemas.RoleResponseDTO{
		ID:   r.ID,
		Name: r.Name,
	}
}

func mapToMemberResponse(m *boilmodels.Member) *schemas.MemberResponse {
	if m == nil {
		return nil
	}
	res := &schemas.MemberResponse{
		ID:        m.ID,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Username:  m.Username,
		Email:     m.Email,
		IsAdmin:   m.IsAdmin,
		IsActive:  m.IsActive,
	}
	if m.Address.Valid {
		res.Address = &m.Address.String
	}
	if m.Phone.Valid {
		res.Phone = &m.Phone.String
	}

	if m.R != nil {
		if m.R.Role != nil {
			res.Role = mapToRoleResponse(m.R.Role)
		}
		for _, point := range m.R.PointSales {
			p := schemas.PointSaleResponse{
				ID:        point.ID,
				Name:      point.Name,
				Number:    int64(point.Number),
				IsDeposit: point.IsDeposit,
				IsMain:    point.IsMain,
			}
			if point.Description.Valid {
				p.Description = &point.Description.String
			}
			res.PointSales = append(res.PointSales, p)
		}
	}

	return res
}

func mapToMemberResponseDTO(m *boilmodels.Member) *schemas.MemberResponseDTO {
	if m == nil {
		return nil
	}
	res := &schemas.MemberResponseDTO{
		ID:        m.ID,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Username:  m.Username,
		Email:     m.Email,
		IsAdmin:   m.IsAdmin,
		IsActive:  m.IsActive,
	}
	if m.Address.Valid {
		res.Address = &m.Address.String
	}
	if m.Phone.Valid {
		res.Phone = &m.Phone.String
	}
	if m.R != nil && m.R.Role != nil {
		res.Role = mapToRoleResponseDTO(m.R.Role)
	}

	return res
}

func (r *MemberRepository) MemberGetByID(id int64) (*schemas.MemberResponse, error) {
	ctx := context.Background()

	member, err := boilmodels.Members(
		boilmodels.MemberWhere.ID.EQ(id),
		qm.Load(boilmodels.MemberRels.Role),
		qm.Load(qm.Rels(boilmodels.MemberRels.Role, boilmodels.RoleRels.Permissions)),
		qm.Load(boilmodels.MemberRels.PointSales),
	).One(ctx, r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Miembro", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Miembro", schemas.Read)
	}

	return mapToMemberResponse(member), nil
}

func (r *MemberRepository) MemberGetPermissionByUserID(userID int64) (*schemas.MemberResponse, error) {
	ctx := context.Background()

	member, err := boilmodels.Members(
		boilmodels.MemberWhere.ID.EQ(userID),
		qm.Load(boilmodels.MemberRels.Role),
		qm.Load(qm.Rels(boilmodels.MemberRels.Role, boilmodels.RoleRels.Permissions)),
		qm.Load(boilmodels.MemberRels.PointSales),
	).One(ctx, r.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schemas.HandlerErrorDB(err, "Miembro", schemas.Read)
		}
		return nil, schemas.HandlerErrorDB(err, "Miembro", schemas.Read)
	}

	return mapToMemberResponse(member), nil
}

func (r *MemberRepository) MemberGetAll(limit, page int, search *map[string]string) ([]*schemas.MemberResponseDTO, int64, error) {
	ctx := context.Background()
	offset := (page - 1) * limit

	var qms []qm.QueryMod

	if search != nil {
		for key, value := range *search {
			switch key {
			case "first_name":
				qms = append(qms, qm.Where("first_name ILIKE ?", "%"+value+"%"))
			case "last_name":
				qms = append(qms, qm.Where("last_name ILIKE ?", "%"+value+"%"))
			case "username":
				qms = append(qms, qm.Where("username ILIKE ?", "%"+value+"%"))
			case "email":
				qms = append(qms, qm.Where("email ILIKE ?", "%"+value+"%"))
			case "is_active":
				isAct, err := strconv.ParseBool(value)
				if err != nil {
					return nil, 0, schemas.ErrorResponse(400, "is_active, formato no valido", fmt.Errorf("is_active formato no valido"))
				}
				qms = append(qms, boilmodels.MemberWhere.IsActive.EQ(isAct))
			}
		}
	}

	total, err := boilmodels.Members(qms...).Count(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Miembro", schemas.Read)
	}

	qms = append(qms,
		qm.Load(boilmodels.MemberRels.Role),
		qm.OrderBy("created_at DESC"),
		qm.Offset(offset),
		qm.Limit(limit),
	)

	members, err := boilmodels.Members(qms...).All(ctx, r.DB)
	if err != nil {
		return nil, 0, schemas.HandlerErrorDB(err, "Miembro", schemas.Read)
	}

	var membersSchema []*schemas.MemberResponseDTO
	for _, m := range members {
		membersSchema = append(membersSchema, mapToMemberResponseDTO(m))
	}

	return membersSchema, total, nil
}

func (r *MemberRepository) MemberCreate(memberID int64, memberCreate *schemas.MemberCreate) (int64, error) {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return 0, err
	}

	roleExists, err := boilmodels.Roles(boilmodels.RoleWhere.ID.EQ(memberCreate.RoleID)).Exists(ctx, tx)
	if err != nil || !roleExists {
		return 0, schemas.HandlerErrorDB(sql.ErrNoRows, "Rol", schemas.Read)
	}

	if len(memberCreate.PointSaleIDs) == 0 {
		return 0, schemas.ErrorResponse(400, "Debe seleccionar al menos un punto de venta", fmt.Errorf("debe seleccionar al menos un punto de venta"))
	}

	var pointSaleInterfaces []interface{}
	for _, id := range memberCreate.PointSaleIDs {
		pointSaleInterfaces = append(pointSaleInterfaces, id)
	}

	pointSales, err := boilmodels.PointSales(qm.WhereIn("id IN ?", pointSaleInterfaces...)).All(ctx, tx)
	if err != nil {
		return 0, schemas.HandlerErrorDB(err, "Punto de Venta", schemas.Read)
	}

	if len(pointSales) != len(memberCreate.PointSaleIDs) {
		return 0, schemas.ErrorResponse(400, "Uno o más puntos de venta no existen", fmt.Errorf("uno o más puntos de venta no existen"))
	}

	member := boilmodels.Member{
		FirstName: memberCreate.FirstName,
		LastName:  memberCreate.LastName,
		Username:  memberCreate.Username,
		Email:     memberCreate.Email,
		Password:  memberCreate.Password,
		RoleID:    memberCreate.RoleID,
		IsActive:  true,
		IsAdmin:   false,
	}

	if memberCreate.Address != nil {
		member.Address = null.StringFromPtr(memberCreate.Address)
	}
	if memberCreate.Phone != nil {
		member.Phone = null.StringFromPtr(memberCreate.Phone)
	}

	if err := member.Insert(ctx, tx, boil.Infer()); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Miembro", schemas.Create)
	}

	if err := member.AddPointSales(ctx, tx, false, pointSales...); err != nil {
		return 0, schemas.HandlerErrorDB(err, "Punto de Venta", schemas.Update)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return member.ID, nil
}

func (r *MemberRepository) MemberUpdate(memberID int64, memberUpdate *schemas.MemberUpdate) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	existingMember, err := boilmodels.Members(boilmodels.MemberWhere.ID.EQ(memberUpdate.ID)).One(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Miembro", schemas.Read)
	}

	if existingMember.IsAdmin && memberUpdate.IsActive != nil && !*memberUpdate.IsActive {
		return schemas.ErrorResponse(400, "No se puede desactivar un administrador", fmt.Errorf("no se puede desactivar un administrador"))
	}

	roleExists, err := boilmodels.Roles(boilmodels.RoleWhere.ID.EQ(memberUpdate.RoleID)).Exists(ctx, tx)
	if err != nil || !roleExists {
		return schemas.HandlerErrorDB(sql.ErrNoRows, "Rol", schemas.Read)
	}

	if len(memberUpdate.PointSaleIDs) == 0 {
		return schemas.ErrorResponse(400, "Debe seleccionar al menos un punto de venta", fmt.Errorf("debe seleccionar al menos un punto de venta"))
	}

	var pointSaleInterfaces []interface{}
	for _, id := range memberUpdate.PointSaleIDs {
		pointSaleInterfaces = append(pointSaleInterfaces, id)
	}

	pointSales, err := boilmodels.PointSales(qm.WhereIn("id IN ?", pointSaleInterfaces...)).All(ctx, tx)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Punto de Venta", schemas.Read)
	}

	if len(pointSales) != len(memberUpdate.PointSaleIDs) {
		return schemas.ErrorResponse(400, "Uno o más puntos de venta no existen", fmt.Errorf("uno o más puntos de venta no existen"))
	}

	existingMember.FirstName = memberUpdate.FirstName
	existingMember.LastName = memberUpdate.LastName
	existingMember.Username = memberUpdate.Username
	existingMember.Email = memberUpdate.Email
	existingMember.Address = null.StringFromPtr(memberUpdate.Address)
	existingMember.Phone = null.StringFromPtr(memberUpdate.Phone)
	existingMember.RoleID = memberUpdate.RoleID

	if memberUpdate.IsActive != nil {
		existingMember.IsActive = *memberUpdate.IsActive
	}

	if _, err := existingMember.Update(ctx, tx, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Miembro", schemas.Update)
	}

	if err := existingMember.SetPointSales(ctx, tx, false, pointSales...); err != nil {
		return schemas.HandlerErrorDB(err, "Punto de Venta", schemas.Update)
	}

	return tx.Commit()
}

func (r *MemberRepository) MemberUpdatePassword(memberID int64, passwordUpdate *schemas.MemberUpdatePassword) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	member, err := boilmodels.Members(boilmodels.MemberWhere.ID.EQ(memberID)).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schemas.ErrorResponse(404, "Miembro no encontrado", err)
		}
		return schemas.ErrorResponse(500, "Error al obtener el miembro", err)
	}

	if !utils.CheckPasswordHash(passwordUpdate.OldPassword, member.Password) {
		return schemas.ErrorResponse(400, "La contraseña actual es incorrecta", fmt.Errorf("la contraseña actual es incorrecta"))
	}

	passHash, err := utils.HashPassword(passwordUpdate.NewPassword)
	if err != nil {
		return schemas.ErrorResponse(500, "Error al generar la contraseña", err)
	}

	member.Password = passHash
	if _, err := member.Update(ctx, tx, boil.Whitelist(boilmodels.MemberColumns.Password, boilmodels.MemberColumns.UpdatedAt)); err != nil {
		return schemas.ErrorResponse(500, "Error al actualizar la contraseña", err)
	}

	return tx.Commit()
}

func (r *MemberRepository) MemberDelete(memberID int64, id int64) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_member_id', $1, true)", fmt.Sprintf("%d", memberID)); err != nil {
		return err
	}

	member, err := boilmodels.Members(boilmodels.MemberWhere.ID.EQ(id)).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schemas.ErrorResponse(404, "Miembro no encontrado", err)
		}
		return schemas.ErrorResponse(500, "Error al obtener el miembro", err)
	}

	if member.IsAdmin {
		return schemas.ErrorResponse(400, "No se puede eliminar un administrador", nil)
	}

	if _, err := member.Delete(ctx, tx, false); err != nil {
		return schemas.ErrorResponse(500, "Error al eliminar el miembro", err)
	}

	return tx.Commit()
}

func (r *MemberRepository) MemberCount() (int64, error) {
	ctx := context.Background()

	members, err := boilmodels.Members().Count(ctx, r.DB)
	if err != nil {
		return 0, schemas.ErrorResponse(500, "error al obtner la cantidad de miembros", err)
	}

	return members, nil
}
