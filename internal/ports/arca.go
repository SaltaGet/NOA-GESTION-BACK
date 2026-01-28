package ports

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type ArcaRepository interface {
	GetCredentialsArca(tenantID int64) (*models.Credential, error)
	SetTokenSignArca(v *schemas.CredentialsValidation) error
	GetLastestInvoice(w *schemas.WSFEClient, pointSale, TypeInvoice int) (int64, error)
	GetInfoInvoice(w *schemas.WSFEClient, pointSale, typeInvoice int, numberInvoice int64) (*schemas.FECompConsultaResponse, error)
	SendToWSAA(w *schemas.WSAA, cms string) ([]byte, error) 
	EmitInvoice(w *schemas.WSFEClient, factura *schemas.Factura) (*schemas.FECAEDetResponse, error)
}

type ArcaService interface {
	EmitInvoice(user *schemas.AuthenticatedUser, pointSaleID int64, req *schemas.FacturaRequest, isHomo bool) (*schemas.FacturaElectronica, error)
}