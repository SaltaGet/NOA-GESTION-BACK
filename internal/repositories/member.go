package repositories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// MemberGetByID obtiene un miembro por ID
func (r *MemberRepository) MemberGetByID(id int64) (*schemas.MemberResponse, error) {
	var member models.Member

	if err := r.DB.
		Preload("Role").
		Preload("Role.Permissions").
		Preload("PointSales").
		First(&member, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Miembro no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error al obtener el miembro", err)
	}

	var memberSchema schemas.MemberResponse
	copier.Copy(&memberSchema, &member)

	return &memberSchema, nil
}

// MemberGetPermissionByUserID obtiene un miembro con sus permisos por ID
func (r *MemberRepository) MemberGetPermissionByUserID(userID int64) (*schemas.MemberResponse, error) {
	var member models.Member

	if err := r.DB.
		Preload("Role.Permissions").
		Preload("PointSales").
		First(&member, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Miembro no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error al obtener el miembro", err)
	}

	var memberSchema schemas.MemberResponse
	copier.Copy(&memberSchema, &member)

	return &memberSchema, nil
}

// MemberGetAll obtiene todos los miembros con paginación y búsqueda
func (r *MemberRepository) MemberGetAll(limit, page int, search *map[string]string) ([]*schemas.MemberResponseDTO, int64, error) {
	var members []*models.Member
	var total int64

	offset := (page - 1) * limit

	query := r.DB.Model(&models.Member{})

	// Aplicar filtros de búsqueda si existen
	if search != nil {
		for key, value := range *search {
			switch key {
			case "first_name":
				query = query.Where("first_name ILIKE ?", "%"+value+"%")
			case "last_name":
				query = query.Where("last_name ILIKE ?", "%"+value+"%")
			case "username":
				query = query.Where("username ILIKE ?", "%"+value+"%")
			case "email":
				query = query.Where("email ILIKE ?", "%"+value+"%")
			case "is_active":
				isAct, err := strconv.ParseBool(value)
				if err != nil {
					return nil, 0, schemas.ErrorResponse(400, "is_active, formato no valido", fmt.Errorf("is_active formato no valido"))
				}
				query = query.Where("is_active = ?", isAct)
			}
		}
	}

	// Contar total de registros
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "Error al contar los miembros", err)
	}

	// Obtener registros con paginación
	if err := query.
		Preload("Role").
		Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&members).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "Error al obtener los miembros", err)
	}

	var membersSchema []*schemas.MemberResponseDTO
	copier.Copy(&membersSchema, &members)

	return membersSchema, total, nil
}

// MemberCreate crea un nuevo miembro
func (r *MemberRepository) MemberCreate(memberCreate *schemas.MemberCreate) (int64, error) {
	var memberID int64

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el rol existe
		var role models.Role
		if err := tx.First(&role, memberCreate.RoleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(400, fmt.Sprintf("El rol %d no existe", memberCreate.RoleID), err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el rol", err)
		}

		// Verificar que los puntos de venta existen
		var pointSales []models.PointSale
		if err := tx.Where("id IN ?", memberCreate.PointSaleIDs).Find(&pointSales).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al obtener los puntos de venta", err)
		}

		if len(pointSales) != len(memberCreate.PointSaleIDs) {
			return schemas.ErrorResponse(400, "Uno o más puntos de venta no existen", fmt.Errorf("uno o más puntos de venta no existen"))
		}

		// Crear miembro
		member := models.Member{
			FirstName: memberCreate.FirstName,
			LastName:  memberCreate.LastName,
			Username:  memberCreate.Username,
			Email:     memberCreate.Email,
			Password:  memberCreate.Password,
			Address:   memberCreate.Address,
			Phone:     memberCreate.Phone,
			RoleID:    memberCreate.RoleID,
			IsActive:  true,
			IsAdmin:   false,
		}

		if err := tx.Create(&member).Error; err != nil {
			if schemas.IsDuplicateError(err) {
				if strings.Contains(err.Error(), "username") {
					return schemas.ErrorResponse(400, "El nombre de usuario ya existe", err)
				}
				if strings.Contains(err.Error(), "email") {
					return schemas.ErrorResponse(400, "El email ya existe", err)
				}
				return schemas.ErrorResponse(400, "El miembro ya existe", err)
			}
			return schemas.ErrorResponse(500, "Error al crear el miembro", err)
		}

		memberID = member.ID

		// Asociar puntos de venta
		if err := tx.Model(&member).Association("PointSales").Append(&pointSales); err != nil {
			return schemas.ErrorResponse(500, "Error al asociar puntos de venta", err)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return memberID, nil
}

// MemberUpdate actualiza un miembro existente
func (r *MemberRepository) MemberUpdate(memberUpdate *schemas.MemberUpdate) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el miembro existe
		var existingMember models.Member
		if err := tx.First(&existingMember, memberUpdate.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Miembro no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el miembro", err)
		}

		if existingMember.IsAdmin && memberUpdate.IsActive != nil && !*memberUpdate.IsActive {
			return schemas.ErrorResponse(400, "No se puede desactivar un administrador", fmt.Errorf("no se puede desactivar un administrador"))
		}

		// Verificar que el rol existe
		var role models.Role
		if err := tx.First(&role, memberUpdate.RoleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(400, fmt.Sprintf("El rol %d no existe", memberUpdate.RoleID), err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el rol", err)
		}

		// Verificar que los puntos de venta existen
		var pointSales []models.PointSale
		if err := tx.Where("id IN ?", memberUpdate.PointSaleIDs).Find(&pointSales).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al obtener los puntos de venta", err)
		}

		if len(pointSales) != len(memberUpdate.PointSaleIDs) {
			return schemas.ErrorResponse(400, "Uno o más puntos de venta no existen", nil)
		}

		// Actualizar campos
		existingMember.FirstName = memberUpdate.FirstName
		existingMember.LastName = memberUpdate.LastName
		existingMember.Username = memberUpdate.Username
		existingMember.Email = memberUpdate.Email
		existingMember.Address = memberUpdate.Address
		existingMember.Phone = memberUpdate.Phone
		existingMember.RoleID = memberUpdate.RoleID
		if memberUpdate.IsActive != nil {
			existingMember.IsActive = *memberUpdate.IsActive
		}

		if err := tx.Save(&existingMember).Error; err != nil {
			if schemas.IsDuplicateError(err) {
				if strings.Contains(err.Error(), "username") {
					return schemas.ErrorResponse(400, "El nombre de usuario ya existe", err)
				}
				if strings.Contains(err.Error(), "email") {
					return schemas.ErrorResponse(400, "El email ya existe", err)
				}
				return schemas.ErrorResponse(400, "Error de duplicación", err)
			}
			return schemas.ErrorResponse(500, "Error al actualizar el miembro", err)
		}

		// Actualizar puntos de venta (reemplazar asociaciones existentes)
		if err := tx.Model(&existingMember).Association("PointSales").Replace(&pointSales); err != nil {
			return schemas.ErrorResponse(500, "Error al actualizar puntos de venta", err)
		}

		return nil
	})
}

// MemberUpdatePassword actualiza la contraseña de un miembro
func (r *MemberRepository) MemberUpdatePassword(memberID int64, passwordUpdate *schemas.MemberUpdatePassword) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el miembro existe
		var member models.Member
		if err := tx.First(&member, memberID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Miembro no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el miembro", err)
		}

		// Verificar contraseña actual
		if err := bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(passwordUpdate.OldPassword)); err != nil {
			return schemas.ErrorResponse(400, "La contraseña actual es incorrecta", err)
		}

		// Actualizar contraseña
		member.Password = passwordUpdate.NewPassword
		if err := tx.Save(&member).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al actualizar la contraseña", err)
		}

		return nil
	})
}

// MemberDelete elimina un miembro (soft delete)
func (r *MemberRepository) MemberDelete(id int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el miembro existe
		var member models.Member
		if err := tx.First(&member, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Miembro no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error al obtener el miembro", err)
		}

		// No permitir eliminar admin principal (opcional)
		if member.IsAdmin {
			return schemas.ErrorResponse(400, "No se puede eliminar un administrador", nil)
		}

		// Soft delete
		if err := tx.Delete(&member).Error; err != nil {
			return schemas.ErrorResponse(500, "Error al eliminar el miembro", err)
		}

		return nil
	})
}

func (r *MemberRepository) MemberCount() (int64, error) {
	var members int64
	if err := r.DB.Model(&models.Member{}).Count(&members).Error; err != nil {
		return 0, schemas.ErrorResponse(500, "error al obtner la cantidad de miembros", err)
	}

	return members, nil
}