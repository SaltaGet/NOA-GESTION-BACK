package schemas

type RoleResponseDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type RoleResponse struct {
	ID          int64                `json:"id"`
	Name        string               `json:"name"`
	Permissions []PermissionResponseDTO `json:"permissions"`
}

type RolePermissionRow struct {
	RoleID      string
	RoleName    string
	PermID      string
	PermCode    string
	PermGroup   string
}
