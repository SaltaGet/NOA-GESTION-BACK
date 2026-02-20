package schemas

type AuthLogin struct {
	Username string `json:"username" validate:"required,username"`
	Password string `json:"password" validate:"required"`
}

type AuthLoginAdmin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthenticatedUser struct {
	ID               int64                    `json:"id"`
	FirstName        string                   `json:"first_name"`
	LastName         string                   `json:"last_name"`
	Username         string                   `json:"username"`
	IsAdmin          bool                     `json:"is_admin"`
	RoleID           int64                    `json:"role_id"`
	RoleName         string                   `json:"role_name"`
	Permissions      []EnvironmentPermissions `json:"permissions"`
	TenantID         int64                    `json:"tenant_id"`
	TenantName       string                   `json:"tenant_name"`
	TenantIdentifier string                   `json:"tenant_identifier"`
	ListPermissions  []string                 `json:"list_permissions"`
	AcceptedTerms    bool                     `json:"accepted_terms"`
}

type EnvironmentPermissions struct {
	Environment string             `json:"environment"`
	Groups      []GroupPermissions `json:"groups"`
}

type GroupPermissions struct {
	Group       string   `json:"group"`
	Permissions []string `json:"permissions"`
}

type AuthForgotPassword struct {
	Username         string `json:"username" validate:"required"`
	TenantIdentifier string `json:"tenant_identifier" validate:"required"`
}

type AuthResetPassword struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,password" example:"Password123*"`
	ConfirmPass string `json:"confirm_pass" validate:"required,password" example:"Password123*"`
}
