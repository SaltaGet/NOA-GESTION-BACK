package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"


type MemberRepository interface {
	MemberGetByID(id string) (member *schemas.MemberResponse, err error)
	MemberGetPermissionByUserID(userID string) (member *schemas.Member, err error)
	MemberGetAll() (members *[]schemas.MemberDTO, err error)
	MemberCreate(memeberCreate *schemas.MemberCreate, user *schemas.AuthenticatedUser) (id string, err error)
	MemberDelete(id string) (err error)
}

type MemberService interface {
	MemberGetByID(id string) (member *schemas.MemberResponse, err error)
	MemberGetPermissionByUserID(userID string) (permission *schemas.Member, err error)
	MemberGetAll() (members *[]schemas.MemberDTO, err error)
	MemberCreate(memeberCreate *schemas.MemberCreate, user *schemas.AuthenticatedUser) (id string, err error)
	MemberDelete(id string) (err error)
}