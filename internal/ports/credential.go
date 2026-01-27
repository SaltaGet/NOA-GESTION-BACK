package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"

type CredentialRepository interface {
	CredentialGetMPToken(tenantID int64) (*schemas.CredentialMPTokenResponse, error)
	CredentialSetMPToken(tenantID int64, request *schemas.CredentialMPTokenRequest) error
	CredentialGetArca(tenantID int64) (*schemas.CredentialArcaResponse, error)
	CredentialSetArca(tenantID int64, request *schemas.CredentialArcaRequest) error
}

type CredentialService interface {
	CredentialGetMPToken(tenantID int64) (*schemas.CredentialMPTokenResponse, error)
	CredentialSetMPToken(tenantID int64, request *schemas.CredentialMPTokenRequest) (*string, error)
	CredentialGetArca(tenantID int64) (*schemas.CredentialArcaResponse, error)
	CredentialSetArca(tenantID int64, request *schemas.CredentialArcaRequest) error
}