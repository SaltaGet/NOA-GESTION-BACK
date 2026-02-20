package schemas

type RoleResponseDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type RoleResponse struct {
	ID          int64                   `json:"id"`
	Name        string                  `json:"name"`
	Permissions []PermissionResponse `json:"permissions"`
}

type RolePermissionRow struct {
	RoleID      int64
	RoleName    string
	PermID      int64
	PermCode    string
	PermGroup   string
	Environment string
	Detail      string
}

type RoleCreate struct {
	Name          string  `json:"name" validate:"required"`
	PermissionsID []int64 `json:"permissions_id" validate:"required,dive"`
}

type RoleUpdate struct {
	ID            int64   `json:"id" validate:"required"`
	Name          string  `json:"name" validate:"required"`
	PermissionsID []int64 `json:"permissions_id" validate:"required,dive"`
}
