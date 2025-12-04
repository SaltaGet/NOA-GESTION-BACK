package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"


type MemberRepository interface {
	MemberGetByID(id int64) (*schemas.MemberResponse, error)
	MemberGetPermissionByUserID(userID int64) (*schemas.MemberResponse, error)
	MemberGetAll(limit, page int, search *map[string]string) ([]*schemas.MemberResponseDTO, int64, error)
	MemberCreate(memeberCreate *schemas.MemberCreate) (int64, error)
	MemberUpdate(memeberUpdate *schemas.MemberUpdate) (error)
	MemberUpdatePassword(memberID int64, memeberUpdatePassword *schemas.MemberUpdatePassword) (error)
	MemberDelete(id int64) (err error)
	MemberCount() (int64, error)
}

type MemberService interface {
	MemberGetByID(id int64) (*schemas.MemberResponse, error)
	MemberGetPermissionByUserID(userID int64) (*schemas.MemberResponse, error)
	MemberGetAll(limit, page int, search *map[string]string) ([]*schemas.MemberResponseDTO, int64, error)
	MemberCreate(memeberCreate *schemas.MemberCreate, plan *schemas.PlanResponseDTO) (int64, error)
	MemberUpdate(memeberUpdate *schemas.MemberUpdate) (error)
	MemberUpdatePassword(memberID int64, memeberUpdatePassword *schemas.MemberUpdatePassword) (error)
	MemberDelete(id int64) (err error)
}