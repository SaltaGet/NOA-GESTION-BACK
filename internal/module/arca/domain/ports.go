package domain

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/module/arca/infrastructure/repository"
	repositoryAuth "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/auth/infrastructure/repository"
	modelCredential "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/credential/domain"
)

type ArcaRepository interface {
	GetCredentialsArca(tenantID int64, incomeSaleID int64) (*modelCredential.Credential, error)
	SetTokenSignArca(v *repository.CredentialsValidation) error
	GetLastestInvoice(w *repository.WSFEClient, pointSale, TypeInvoice int) (int64, error)
	GetInfoInvoice(w *repository.WSFEClient, pointSale, typeInvoice int, numberInvoice int64) (*repository.FECompConsultaResponse, error)
	SendToWSAA(w *repository.WSAA, cms string) ([]byte, error)
	EmitInvoice(w *repository.WSFEClient, factura *repository.Factura) (*repository.FECAEDetResponse, error)
	SaveInvoice(factura *repository.FacturaElectronica, incomeSaleID int64) error
}

type ArcaService interface {
	EmitInvoice(user *repositoryAuth.AuthenticatedUser, pointSaleID int64, incomeSaleID int64, req *repository.FacturaRequest, isHomo bool) (*repository.FacturaElectronica, error)
}
