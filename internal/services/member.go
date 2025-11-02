package services

import (
	"log"
	"runtime/debug"

	"github.com/DanielChachagua/GestionCar/pkg/models"
)

func (m *MemberService) MemberGetAll() (*[]models.MemberDTO, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic atrapado en MemberGetAll: %v", r)
			debug.PrintStack()
		}
	}()
	members, err := m.MemberRepository.MemberGetAll()
	if err != nil {
		return nil, err
	}

	return members, nil
}


func (m *MemberService) MemberGetPermissionByUserID(userID string) (*models.Member, error) {
	member, err := m.MemberRepository.MemberGetPermissionByUserID(userID)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (m *MemberService) MemberGetByID(id string) (*models.MemberResponse, error) {
	member, err := m.MemberRepository.MemberGetByID(id)
	if err != nil {
		return nil, err
	}
	
	return member, nil
}

func (m *MemberService) MemberCreate(memberCreate *models.MemberCreate, user *models.AuthenticatedUser) (string, error) {
	id, err := m.MemberRepository.MemberCreate(memberCreate, user)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (m *MemberService) MemberDelete(memberID string) error {
	if err := m.MemberRepository.MemberDelete(memberID); err != nil {
		return err
	}
	return nil
}