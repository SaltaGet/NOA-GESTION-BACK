package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"


type MemberRepository interface {
	MemberGetByID(id int64) (*schemas.MemberResponse, error)
	MemberGetPermissionByUserID(userID int64) (*schemas.MemberResponse, error)
	MemberGetAll(limit, page int, search *map[string]string) ([]*schemas.MemberResponseDTO, int64, error)
	MemberCreate(memberID int64, memeberCreate *schemas.MemberCreate) (int64, error)
	MemberUpdate(memberID int64, memeberUpdate *schemas.MemberUpdate) (error)
	MemberUpdatePassword(memberID int64, memeberUpdatePassword *schemas.MemberUpdatePassword) (error)
	MemberDelete(memberID int64, id int64) (err error)
	MemberCount() (int64, error)
}

type MemberService interface {
	MemberGetByID(id int64) (*schemas.MemberResponse, error)
	MemberGetPermissionByUserID(userID int64) (*schemas.MemberResponse, error)
	MemberGetAll(limit, page int, search *map[string]string) ([]*schemas.MemberResponseDTO, int64, error)
	MemberCreate(memberID int64, memeberCreate *schemas.MemberCreate, plan *schemas.PlanResponseDTO) (int64, error)
	MemberUpdate(memberID int64, memeberUpdate *schemas.MemberUpdate) (error)
	MemberUpdatePassword(memberID int64, memeberUpdatePassword *schemas.MemberUpdatePassword) (error)
	MemberDelete(memberID int64, id int64) (err error)
}