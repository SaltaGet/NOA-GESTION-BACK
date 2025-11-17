package schemas

type PermissionResponse struct {
	ID      int64 `json:"id"`
	Code    string `json:"code"`
	Details string `json:"details"`
	Group   string `json:"group"`
	Environment string `json:"environment"`
}

type PermissionResponseDTO struct {
	ID      int64 `json:"id"`
	Code    string `json:"code"`
	Group   string `json:"group"`
	Environment string `json:"environment"`
}