package repositories

import (
	"errors"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (t *MemberRepository) MemberGetAll() (*[]schemas.MemberDTO, error) {
	var members []schemas.Member
	if err := t.DB.Preload("Role").Find(&members).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al buscar miembros", err)
	}

	var memberDtos []schemas.MemberDTO
	for _, member := range members {
		memberDtos = append(memberDtos, schemas.MemberDTO{
			ID:        member.ID,
			FirstName: member.FirstName,
			LastName:  member.LastName,
			Username:  member.Username,
			Email:     member.Email,
			IsActive:  member.IsActive,
			CreatedAt: member.CreatedAt,
			Role:      schemas.RoleDTO{ID: member.Role.ID, Name: member.Role.Name},			
		})
	}

	return &memberDtos, nil
}

func (t *MemberRepository) MemberGetPermissionByUserID(userID string) (*schemas.Member, error) {
	var member schemas.Member
	if err := t.DB.Preload("Role").Preload("Role.Permissions").Where("user_id = ?", userID).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Miembro no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error  internoal buscar miembro", err)
	}
	return &member, nil
}

func (t *MemberRepository) MemberGetByID(id string) (*schemas.MemberResponse, error) {
	var member schemas.Member
	if err := t.DB.Preload("Role").Preload("Role.Permissions").Where("id = ?", id).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Miembro no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al buscar miembro", err)
	}

	response := schemas.MemberResponse{
		ID:        member.ID,
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Username:  member.Username,
		Email:     member.Email,
		IsActive:  member.IsActive,
		CreatedAt: member.CreatedAt,
		Role:      member.Role,
	}

	return &response, nil
}

func (t *MemberRepository) MemberCreate(memeberCreate *schemas.MemberCreate, user *schemas.AuthenticatedUser) (id string, err error) {
	newID := uuid.NewString()

	// pass, err := utils.GenerateRandomString(10)
	// if err != nil {
	// 	return "", err
	// }
	hashPass, err := utils.HashPassword("1234")
	if err != nil {
		return "", schemas.ErrorResponse(500, "Error al hashear la contraseña", err)
	}

	if err := t.DB.Create(&schemas.Member{
		ID:        newID,
		FirstName: memeberCreate.FirstName,
		LastName:  memeberCreate.LastName,
		Username:  memeberCreate.Username,
		Email:     memeberCreate.Email,
		RoleID:    memeberCreate.RoleID,
		Password:  hashPass,
	}).Error; err != nil {
		return "", schemas.ErrorResponse(500, "Error interno al crear miembro", err)
	}
	return newID, nil
}

func (t *MemberRepository) MemberDelete(id string) error {
	if err := t.DB.Delete(&schemas.Member{}, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "Miembro no encontrado", err)
		}
		return schemas.ErrorResponse(500, "Error interno al eliminar miembro", err)
	}
	return nil
}