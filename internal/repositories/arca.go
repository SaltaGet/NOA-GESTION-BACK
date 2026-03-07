package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/master"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

func (r *ArcaRepository) GetCredentialsArca(tenantID, incomeSaleID int64) (*master.Credential, error) {
	ctx := context.Background()

	exists, err := tenant.IncomeSales(
		tenant.IncomeSaleWhere.ID.EQ(incomeSaleID),
	).Exists(ctx, r.DB)
	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Ingreso de venta", schemas.Read)
	}
	if !exists {
		return nil, schemas.HandlerErrorDB(sql.ErrNoRows, "Ingreso de venta", schemas.Read)
	}

	db := database.GetMainDB()

	c, err := master.Credentials(
		qm.Select(
			master.CredentialColumns.ID, master.CredentialColumns.SocialReason, master.CredentialColumns.BusinessName,
			master.CredentialColumns.Address, master.CredentialColumns.ResponsibilityFrontIva, master.CredentialColumns.Cuit,
			master.CredentialColumns.GrossIncome, master.CredentialColumns.StartActivities, master.CredentialColumns.ArcaCertificate,
			master.CredentialColumns.ArcaKey, master.CredentialColumns.TokenArca, master.CredentialColumns.SignArca,
			master.CredentialColumns.ExpireTokenArca, master.CredentialColumns.Concept,
		),
		master.CredentialWhere.TenantID.EQ(tenantID),
	).One(ctx, db)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Credenciales", schemas.Read)
	}

	if !c.SocialReason.Valid || !c.ResponsibilityFrontIva.Valid || !c.Cuit.Valid || !c.ArcaCertificate.Valid || !c.ArcaKey.Valid {
		return nil, schemas.ErrorResponse(422, "La entidad no se puede procesar por datos incompletos, revise y complete adecuadamente las credenciales", errors.New("La entidad no se puede procesar por datos incompletos."))
	}

	credential := &master.Credential{
		ID:                     c.ID,
		SocialReason:           c.SocialReason,
		ResponsibilityFrontIva: c.ResponsibilityFrontIva,
		Cuit:                   c.Cuit,
		ArcaCertificate:        c.ArcaCertificate,
		ArcaKey:                c.ArcaKey,
	}
	if c.BusinessName.Valid {
		credential.BusinessName = c.BusinessName
	}
	if c.Address.Valid {
		credential.Address = c.Address
	}
	if c.GrossIncome.Valid {
		credential.GrossIncome = c.GrossIncome
	}
	if c.StartActivities.Valid {
		credential.StartActivities = c.StartActivities
	}
	if c.TokenArca.Valid {
		credential.TokenArca = c.TokenArca
	}
	if c.SignArca.Valid {
		credential.SignArca = c.SignArca
	}
	if c.ExpireTokenArca.Valid {
		credential.ExpireTokenArca = c.ExpireTokenArca
	}
	if c.Concept.Valid {
		credential.Concept = c.Concept
	}

	return credential, nil
}

func (r *ArcaRepository) SetTokenSignArca(v *schemas.CredentialsValidation) error {
	ctx := context.Background()
	db := database.GetMainDB()

	cuit := null.StringFrom(fmt.Sprintf("%d", v.CUIT))
	cFromDb, err := master.Credentials(master.CredentialWhere.Cuit.EQ(cuit)).One(ctx, db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schemas.ErrorResponse(404, "Credencial no encontrada", errors.New("Credencial no encontrada"))
		}
		return schemas.HandlerErrorDB(err, "Credenciales", schemas.Read)
	}

	cFromDb.TokenArca = null.StringFrom(v.Token)
	cFromDb.SignArca = null.StringFrom(v.Sign)
	cFromDb.ExpireTokenArca = null.TimeFrom(v.Expiration)

	_, err = cFromDb.Update(ctx, db, boil.Whitelist(
		master.CredentialColumns.TokenArca,
		master.CredentialColumns.SignArca,
		master.CredentialColumns.ExpireTokenArca,
	))

	if err != nil {
		return schemas.HandlerErrorDB(err, "Credenciales", schemas.Update)
	}

	return nil
}

func (r *ArcaRepository) GetLastestInvoice(w *schemas.WSFEClient, pointSale, TypeInvoice int) (int64, error) {
	req := schemas.FECompUltimoAutorizadoRequest{
		Xmlns:    "http://ar.gov.afip.dif.FEV1/",
		Auth:     w.Auth,
		PtoVta:   pointSale,
		CbteTipo: TypeInvoice,
	}

	reqXML, _ := xml.Marshal(req)
	soapEnv := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
			<soap:Body>
				%s
			</soap:Body>
		</soap:Envelope>`, string(reqXML))

	httpReq, err := http.NewRequest("POST", w.BaseURL, strings.NewReader(soapEnv))
	if err != nil {
		return 0, err
	}

	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpReq.Header.Set("SOAPAction", "http://ar.gov.afip.dif.FEV1/FECompUltimoAutorizado")

	resp, err := w.Client.Do(httpReq)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response schemas.FECompUltimoAutorizadoResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return 0, err
	}

	if len(response.Body.Response.Result.Errors.Err) > 0 {
		return 0, fmt.Errorf("AFIP Error: %s", response.Body.Response.Result.Errors.Err[0].Msg)
	}

	return response.Body.Response.Result.CbteNro, nil
}

func (r *ArcaRepository) GetInfoInvoice(w *schemas.WSFEClient, pointSale, typeInvoice int, numberInvoice int64) (*schemas.FECompConsultaResponse, error) {
	req := schemas.FECompConsultaRequest{
		Xmlns: "http://ar.gov.afip.dif.FEV1/",
		Auth:  w.Auth,
	}
	req.FeCompConsReq.CbteTipo = typeInvoice
	req.FeCompConsReq.PtoVta = pointSale
	req.FeCompConsReq.CbteNro = numberInvoice

	reqXML, _ := xml.MarshalIndent(req, "    ", "  ")
	soapEnv := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
			<soap:Body>
				%s
			</soap:Body>
		</soap:Envelope>`, string(reqXML))

	httpReq, err := http.NewRequest("POST", w.BaseURL, strings.NewReader(soapEnv))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpReq.Header.Set("SOAPAction", "http://ar.gov.afip.dif.FEV1/FECompConsultar")

	resp, err := w.Client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var response schemas.FECompConsultaResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parseando respuesta: %v", err)
	}

	if len(response.Body.Response.Result.Errors.Err) > 0 {
		return nil, fmt.Errorf("AFIP Error: %s", response.Body.Response.Result.Errors.Err[0].Msg)
	}

	return &response, nil
}

func (r *ArcaRepository) SendToWSAA(w *schemas.WSAA, cms string) ([]byte, error) {
	url := "https://wsaahomo.afip.gov.ar/ws/services/LoginCms"
	if !w.Config.Homologacion {
		url = "https://wsaa.afip.gov.ar/ws/services/LoginCms"
	}

	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
		<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:wsaa="http://wsaa.view.sua.dvadac.desein.afip.gov">
				<soapenv:Header/>
				<soapenv:Body>
						<wsaa:loginCms>
								<wsaa:in0>%s</wsaa:in0>
						</wsaa:loginCms>
				</soapenv:Body>
		</soapenv:Envelope>`, cms)

	req, err := http.NewRequest("POST", url, strings.NewReader(soapRequest))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "")

	resp, err := w.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("WSAA error %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (r *ArcaRepository) EmitInvoice(w *schemas.WSFEClient, factura *schemas.Factura) (*schemas.FECAEDetResponse, error) {
	req := schemas.FECAERequest{
		Xmlns: "http://ar.gov.afip.dif.FEV1/",
		Auth:  w.Auth,
		FeCAEReq: schemas.FeCAEReq{
			FeCabReq: schemas.FeCabReq{
				CantReg:  1,
				PtoVta:   factura.PuntoVenta,
				CbteTipo: factura.TipoComprobante,
			},
			FeDetReq: []schemas.FECAEDetRequest{factura.ToFECAEDetRequest()},
		},
	}

	soapEnv := w.BuildSOAPEnvelope(req)

	httpReq, err := http.NewRequest("POST", w.BaseURL, strings.NewReader(soapEnv))
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al realizar solicitud a ARCA", err)
	}

	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpReq.Header.Set("SOAPAction", "http://ar.gov.afip.dif.FEV1/FECAESolicitar")

	resp, err := w.Client.Do(httpReq)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al realizar solicitud a ARCA", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, schemas.ErrorResponse(resp.StatusCode, "WSFE error en la solicitud", errors.New(string(body)))
	}

	var response schemas.FECAEResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, schemas.ErrorResponse(500, "Error al parsear la respuesta de ARCA", err)
	}

	if len(response.Body.Response.Errors.Err) > 0 {
		errMsg := ""
		for _, e := range response.Body.Response.Errors.Err {
			errMsg += fmt.Sprintf("Error %d: %s\n", e.Code, e.Msg)
		}
		return nil, schemas.ErrorResponse(500, "Error en la respuesta de ARCA", fmt.Errorf("errores de AFIP:\n%s", errMsg))
	}

	if len(response.Body.Response.Results.FeDetResp.FECAEDetResponse) == 0 {
		return nil, schemas.ErrorResponse(500, "Error en la respuesta de ARCA", fmt.Errorf("sin resultados en la respuesta"))
	}

	detResp := response.Body.Response.Results.FeDetResp.FECAEDetResponse[0]

	if detResp.Resultado != "A" {
		obsMsg := ""
		for _, obs := range detResp.Observaciones.Obs {
			obsMsg += fmt.Sprintf("Obs %d: %s\n", obs.Code, obs.Msg)
		}
		return nil, schemas.ErrorResponse(500, "Factura no autorizada por ARCA", fmt.Errorf("resultado: %s\nobservaciones:\n%s", detResp.Resultado, obsMsg))
	}

	return &detResp, nil
}

func (r *ArcaRepository) SaveInvoice(factura *schemas.FacturaElectronica, incomeSaleID int64) error {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	factJson, err := json.Marshal(factura)
	if err != nil {
		return schemas.ErrorResponse(500, "Error al parsear la factura", err)
	}

	invoice := tenant.Invoice{
		InvoiceData: factJson,
	}

	if err := invoice.Insert(ctx, tx, boil.Infer()); err != nil {
		return schemas.HandlerErrorDB(err, "Factura", schemas.Create)
	}

	incomeSale, err := tenant.FindIncomeSale(ctx, tx, incomeSaleID)
	if err != nil {
		return schemas.HandlerErrorDB(err, "Ingreso de venta", schemas.Read)
	}

	incomeSale.InvoiceID = null.Int64From(invoice.ID)
	if _, err := incomeSale.Update(ctx, tx, boil.Whitelist(tenant.IncomeSaleColumns.InvoiceID, tenant.IncomeSaleColumns.UpdatedAt)); err != nil {
		return schemas.HandlerErrorDB(err, "Ingreso de venta", schemas.Update)
	}

	return tx.Commit()
}
